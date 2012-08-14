package gorecurly

import (
	"os"
	"testing"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/url"
	"time"
)

type Config struct {
	XMLName xml.Name `xml:"settings"`
	APIKey string `xml:"apikey"`
	JSKey string `xml:"jskey"`
	Currency string `xml:"currency"`
}

func (c *Config) LoadConfig() error {
	file, err := os.Open("config.xml") // For read access.
	if err != nil {
		return err
	}
	if body, readerr := ioutil.ReadAll(file); readerr == nil {
		if xmlerr := xml.Unmarshal(body, &c); xmlerr != nil {
			return xmlerr
		}
	} else {
		return readerr
	}
	return nil
}

func TestLive(t *testing.T) {
	//live testing
	c := Config{}
	//load config
	if err := c.LoadConfig(); err != nil{
		t.Fatalf("Configuration failed to load: %s", err.Error())
	}
	//init recurly
	r := InitRecurly(c.APIKey,c.JSKey)

	//ACCOUNT TESTS
	//create invalid account
	acc0 := r.NewAccount()
	if err := acc0.Create(); err == nil {
		t.Fatal("Should have failed with blank create")
	}

	//create valid account
	rand.Seed(int64(time.Now().Nanosecond()))
	rvalue := fmt.Sprintf("%v",rand.Intn(400000))
	acc0.AccountCode = fmt.Sprintf("%s%s","test-account-",rvalue)
	acc0.Email = "test-email-" + rvalue + "@example.com"
	acc0.FirstName = "test-fname-" + rvalue
	acc0.LastName = "test-last-" + rvalue
	if err := acc0.Create(); err != nil {
		t.Fatal(err.Error())
	}

	//create valid account
	acc1 := r.NewAccount()
	rand.Seed(int64(time.Now().Nanosecond()))
	rvalue = fmt.Sprintf("%v",rand.Intn(400000))
	acc1.AccountCode = fmt.Sprintf("%s%s","test-account-",rvalue)
	acc1.Email = "test-email-" + rvalue + "@example.com"
	acc1.FirstName = "test-fname-" + rvalue
	acc1.LastName = "test-last-" + rvalue
	if err := acc1.Create(); err != nil {
		t.Fatal(err.Error())
	}

	//create valid account with billing info
	acc2 := r.NewAccount()
	rand.Seed(int64(time.Now().Nanosecond()))
	rvalue = fmt.Sprintf("%v",rand.Intn(400000))
	acc2.AccountCode = fmt.Sprintf("%s%s","test-account-",rvalue)
	acc2.Email = "test-email-" + rvalue + "@example.com"
	acc2.FirstName = "test-fname-" + rvalue
	acc2.LastName = "test-last-" + rvalue
	acc2.B = new(BillingInfo)
	acc2.B.FirstName = "test-fname-" + rvalue
	acc2.B.LastName = "test-last-" + rvalue
	acc2.B.Number = "4111111111111111"
	acc2.B.Month = 12
	acc2.B.Year = 2015
	acc2.B.VerificationValue = "123"
	if err := acc2.Create(); err != nil {
		t.Fatal(err.Error())
	}

	//create valid account with billing info
	acc3 := r.NewAccount()
	rand.Seed(int64(time.Now().Nanosecond()))
	rvalue = fmt.Sprintf("%v",rand.Intn(400000))
	acc3.AccountCode = fmt.Sprintf("%s%s","test-account-",rvalue)
	acc3.Email = "test-email-" + rvalue + "@example.com"
	acc3.FirstName = "test-fname-" + rvalue
	acc3.LastName = "test-last-" + rvalue
	acc3.B = new(BillingInfo)
	acc3.B.FirstName = "test-fname-" + rvalue
	acc3.B.LastName = "test-last-" + rvalue
	acc3.B.Number = "4111111111111111"
	acc3.B.Month = 12
	acc3.B.Year = 2015
	acc3.B.VerificationValue = "123"
	if err := acc3.Create(); err != nil {
		t.Fatal(err.Error())
	}

	//get account
	getacc, err := r.GetAccount(acc2.AccountCode)
	if err != nil {
		t.Fatal(err.Error())
	} 

	//update account
	getacc.Email = "NewEmail@example.com"
	if updateerr := getacc.Update(); updateerr != nil {
		t.Fatal(err.Error())
	}

	//close account
	if closeerr := getacc.Close(); closeerr != nil {
		t.Fatal(err.Error())
	}
	//Get closed account and check state
	getacc, err = r.GetAccount(acc2.AccountCode)
	if err != nil {
		t.Fatal(err.Error())
	} 
	if getacc.State != "closed" {
		t.Fatal("Account state was not = closed, instead it was %s", getacc.State)
	} else {
		//reopen account
		if roerr := getacc.Reopen(); roerr != nil {
			t.Fatal(err.Error())
		} 
	}

	//ACCOUNT LIST TESTS
	//list accounts
	v := url.Values{}
	v.Set("per_page","2")
	if accounts, err :=r.GetAccounts(v); err == nil {
		//page through
		for bcontinue := true; bcontinue; bcontinue = accounts.Next() {
		}
		//page backwards
		accounts.Prev() 
		//page start
		accounts.Start() 
	} else {
		t.Fatal(err.Error())
	}

	//BILLINGINFO TESTING
	//update with invalid card
	if bi, err := r.GetBillingInfo(acc3.AccountCode); err != nil {
		t.Fatalf("Get billing information failed: %s account_code:%s", err.Error(),acc3.AccountCode)
	} else {
		bi.Number = "4000-0000-0000-0002"
		bi.VerificationValue = "123"
		bi.Month = 12
		bi.Year = time.Now().Year() + 1
		if err = bi.Update(); err == nil {
			t.Fatalf("Credit card update should have failed")
		}
	}
	//update with valid card
	if bi, err := r.GetBillingInfo(acc3.AccountCode); err != nil {
		t.Fatalf("Get billing information failed: %s account_code:%s", err.Error(),acc3.AccountCode)
	} else {
		bi.Number = "4111-1111-1111-1111"
		bi.VerificationValue = "123"
		bi.Month = 12
		bi.Year = time.Now().Year() + 1
		if err = bi.Update(); err != nil {
			t.Fatalf("Credit card update failed:%s",err.Error())
		}
	}
	//delete billing info off an account
	if bi, err := r.GetBillingInfo(acc3.AccountCode); err != nil {
		t.Fatalf("Get billing information failed: %s", err.Error())
	} else {
		if err = bi.Delete(); err != nil {
			t.Fatalf("Delete billing information for account_code:%s has failed: %s",acc3.AccountCode, err.Error())
		} else {
			if _, err := r.GetBillingInfo(acc3.AccountCode); err == nil {
				t.Fatalf("Delete billing information failed because billing info still exists: %s", err.Error())
			}
		}
	}
	//END BILLING INFO TESTING

	//ADJUSTMENT TESTING
	//create charge
	adj := r.NewAdjustment()
	adj.AccountCode = acc1.AccountCode
	adj.Description = "some extra charge"
	adj.UnitAmountInCents = 2000
	adj.Currency = c.Currency
	adj.Quantity = 1
	if err := adj.Create(); err != nil {
		t.Fatalf("Create Adjustment for account_code:%s has failed: %s",acc1.AccountCode, err.Error())
	} else {
		if _,err := r.GetAdjustment(adj.UUID); err != nil {
			t.Fatalf("Couldn't find Adjustment uuid:%s has failed: %s",adj.UUID, err.Error())
		}
	}

	//create credit
	adj = r.NewAdjustment()
	adj.AccountCode = acc1.AccountCode
	adj.Description = "some extra credit"
	adj.UnitAmountInCents = -2000
	adj.Currency = c.Currency
	adj.Quantity = 1
	if err := adj.Create(); err != nil {
		t.Fatalf("Create Adjustment for account_code:%s has failed: %s",acc1.AccountCode, err.Error())
	} else {
		if _,err := r.GetAdjustment(adj.UUID); err != nil {
			t.Fatalf("Couldn't find Adjustment uuid:%s has failed: %s",adj.UUID, err.Error())
		}
	}
	//create invalid charge
	adj = r.NewAdjustment()
	adj.Description = "some extra credit"
	adj.UnitAmountInCents = -2000
	adj.Currency = c.Currency
	adj.Quantity = 1
	if err := adj.Create(); err == nil {
		t.Fatalf("Adjustment should have failed")
	}

	//create and delete charge
	adj = r.NewAdjustment()
	adj.AccountCode = acc1.AccountCode
	adj.Description = "some extra charge"
	adj.UnitAmountInCents = 2000
	adj.Currency = c.Currency
	adj.Quantity = 1
	if err := adj.Create(); err != nil {
		t.Fatalf("Create Adjustment for account_code:%s has failed: %s",acc1.AccountCode, err.Error())
	} else {
		if d,err := r.GetAdjustment(adj.UUID); err == nil {
			if err := d.Delete(); err == nil {
				println("Success Delete:" + d.UUID)
			} else {
				println(err.Error())
			}
		} else {
			t.Fatalf("Couldn't find Adjustment uuid:%s has failed: %s",adj.UUID, err.Error())
		}
	}
	//END ADJUSTMENT TESTING

	//ADJUSTMENT LISTING TESTING
	v = url.Values{}
	v.Set("per_page","1")
	if accounts, err :=r.GetAdjustments(acc1.AccountCode,v); err == nil {
		//page through
		for bcontinue := true; bcontinue; bcontinue = accounts.Next() {
		}
		//page backwards
		if !accounts.Prev() {
			t.Fatalf("Prev didn't work for adjustments")
		}
		//page start
		if !accounts.Start() {
			t.Fatalf("Prev didn't work for adjustments")
		}
	} else {
		t.Fatal(err.Error())
	}
	//END ADJUSTMENT LISTING TESTING
}
