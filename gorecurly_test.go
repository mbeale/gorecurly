package gorecurly

import (
	"os"
	"net/url"
	"fmt"
	"math/rand"
	"time"
	"testing"
	"encoding/xml"
	"io/ioutil"
)

type Config struct {
	XMLName xml.Name `xml:"settings"`
	APIKey string `xml:"apikey"`
	JSKey string `xml:"jskey"`
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

func TestA(t *testing.T) {
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
	rvalue := fmt.Sprintf("%v",rand.Intn(400))
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
	rvalue = fmt.Sprintf("%v",rand.Intn(400))
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
	rvalue = fmt.Sprintf("%v",rand.Intn(400))
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
	rvalue = fmt.Sprintf("%v",rand.Intn(400))
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
}
