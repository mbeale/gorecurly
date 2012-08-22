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

	//create valid account with billing info for invoicing
	acc4 := r.NewAccount()
	rand.Seed(int64(time.Now().Nanosecond()))
	rvalue = fmt.Sprintf("%v",rand.Intn(400000))
	acc4.AccountCode = fmt.Sprintf("%s%s","test-invoice-account-",rvalue)
	acc4.Email = "test-email-" + rvalue + "@example.com"
	acc4.FirstName = "test-fname-" + rvalue
	acc4.LastName = "test-last-" + rvalue
	acc4.B = new(BillingInfo)
	acc4.B.FirstName = "test-fname-" + rvalue
	acc4.B.LastName = "test-last-" + rvalue
	acc4.B.Number = "4111111111111111"
	acc4.B.Month = 12
	acc4.B.Year = 2015
	acc4.B.VerificationValue = "123"
	if err := acc4.Create(); err != nil {
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
	adj.AccountCode = acc3.AccountCode
	adj.Description = "some extra credit"
	adj.UnitAmountInCents = -2000
	adj.Currency = c.Currency
	adj.Quantity = 1
	if err := adj.Create(); err != nil {
		t.Fatalf("Create Adjustment for account_code:%s has failed: %s",acc3.AccountCode, err.Error())
	} else {
		if _,err := r.GetAdjustment(adj.UUID); err != nil {
			t.Fatalf("Couldn't find Adjustment uuid:%s has failed: %s",adj.UUID, err.Error())
		}
	}
	//create invalid charge
	adj = r.NewAdjustment()
	adj.Description = "some extra charge"
	adj.UnitAmountInCents = 2000
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
		t.Fatalf("Create Adjustment for uuid:%s has failed: %s",adj.UUID, err.Error())
	} else {
		if d,err := r.GetAdjustment(adj.UUID); err == nil {
			if err := d.Delete(); err != nil {
				t.Fatalf("Delete of adjustment failed :%s msg:%s",adj.UUID, err.Error())
			}
		} else {
			t.Fatalf("Couldn't find Adjustment uuid:%s has failed: %s",adj.UUID, err.Error())
		}
	}
	//END ADJUSTMENT TESTING

	//ADJUSTMENT LISTING TESTING
	v = url.Values{}
	v.Set("per_page","1")
	if adjs, err :=r.GetAdjustments(acc1.AccountCode,v); err == nil {
		//page through
		for bcontinue := true; bcontinue; bcontinue = adjs.Next() {
		}
		//page backwards
		if !adjs.Prev() {
			t.Fatalf("Prev didn't work for adjustments")
		}
		//page start
		if !adjs.Start() {
			t.Fatalf("Prev didn't work for adjustments")
		}
	} else {
		t.Fatal(err.Error())
	}
	//END ADJUSTMENT LISTING TESTING

	//PLAN TESTING
	//Create 4 plans
	//plan1
	rand.Seed(int64(time.Now().Nanosecond()))
	rvalue = fmt.Sprintf("%v",rand.Intn(400000))
	plan1 := r.NewPlan();
	plan1.Name = "Some Plan"
	plan1.PlanCode = fmt.Sprintf("%s%s","test-plan-",rvalue)
	plan1.SetupFeeInCents.SetCurrency(c.Currency,3000)
	plan1.UnitAmountInCents.SetCurrency(c.Currency,7000)
	if err := plan1.Create(); err != nil {
		t.Fatalf("Create Plan failed for plan_code:%s has failed: %s",plan1.PlanCode, err.Error())
	} else {
		if _,err := r.GetPlan(plan1.PlanCode); err != nil  {
			t.Fatalf("Couldn't find plan :%s has failed: %s",plan1.PlanCode, err.Error())
		}
	}
	//plan2
	rand.Seed(int64(time.Now().Nanosecond()))
	rvalue = fmt.Sprintf("%v",rand.Intn(400000))
	plan2 := r.NewPlan();
	plan2.Name = "Some Plan"
	plan2.PlanCode = fmt.Sprintf("%s%s","test-plan-",rvalue)
	plan2.SetupFeeInCents.SetCurrency(c.Currency,3000)
	plan2.UnitAmountInCents.SetCurrency(c.Currency,7000)
	if err := plan2.Create(); err != nil {
		t.Fatalf("Create Plan failed for plan_code:%s has failed: %s",plan2.PlanCode, err.Error())
	} else {
		if _,err := r.GetPlan(plan2.PlanCode); err != nil  {
			t.Fatalf("Couldn't find plan :%s has failed: %s",plan2.PlanCode, err.Error())
		}
	}
	//plan3
	rand.Seed(int64(time.Now().Nanosecond()))
	rvalue = fmt.Sprintf("%v",rand.Intn(400000))
	plan3 := r.NewPlan();
	plan3.Name = "Some Plan"
	plan3.PlanCode = fmt.Sprintf("%s%s","test-plan-",rvalue)
	plan3.SetupFeeInCents.SetCurrency(c.Currency,3000)
	plan3.UnitAmountInCents.SetCurrency(c.Currency,7000)
	if err := plan3.Create(); err != nil {
		t.Fatalf("Create Plan failed for plan_code:%s has failed: %s",plan3.PlanCode, err.Error())
	} else {
		if _,err := r.GetPlan(plan3.PlanCode); err != nil  {
			t.Fatalf("Couldn't find plan :%s has failed: %s",plan3.PlanCode, err.Error())
		}
	}
	//plan4
	rand.Seed(int64(time.Now().Nanosecond()))
	rvalue = fmt.Sprintf("%v",rand.Intn(400000))
	plan4 := r.NewPlan();
	plan4.Name = "Some Plan"
	plan4.PlanCode = fmt.Sprintf("%s%s","test-plan-",rvalue)
	plan4.SetupFeeInCents.SetCurrency(c.Currency,3000)
	plan4.UnitAmountInCents.SetCurrency(c.Currency,7000)
	if err := plan4.Create(); err != nil {
		t.Fatalf("Create Plan failed for plan_code:%s has failed: %s",plan4.PlanCode, err.Error())
	} else {
		if _,err := r.GetPlan(plan4.PlanCode); err != nil  {
			t.Fatalf("Couldn't find plan :%s has failed: %s",plan4.PlanCode, err.Error())
		}
	}
	//create plan error
	plan5 := r.NewPlan();
	if err := plan5.Create(); err == nil {
		t.Fatalf("Plan creation should have failed")
	}
	//update plan
	plan1.SetupFeeInCents.SetCurrency(c.Currency,0)
	if err := plan1.Update(); err != nil {
		t.Fatalf("Update Plan failed for plan_code:%s has failed: %s",plan1.PlanCode, err.Error())
	}
	//delete plan
	if err := plan1.Delete(); err != nil {
		t.Fatalf("Delete Plan failed for plan_code:%s has failed: %s",plan1.PlanCode, err.Error())
	}

	//END PLAN TESTING

	//LIST PLAN TESTING
	v = url.Values{}
	v.Set("per_page","1")
	if plans, err :=r.GetPlans(v); err == nil {
		//page through
		for bcontinue := true; bcontinue; bcontinue = plans.Next() {
		}
		//page backwards
		if !plans.Prev() {
			t.Fatalf("Prev didn't work for plans")
		}
		//page start
		if !plans.Start() {
			t.Fatalf("Prev didn't work for plans")
		}
	} else {
		t.Fatal(err.Error())
	}
	//END LIST PLAN TESTING
	
	//COUPON TESTING
	//create 4 coupons
	PlanCodes := PlanCode{}
	PlanCodes.PlanCode = append(PlanCodes.PlanCode,plan2.PlanCode, plan3.PlanCode)
	cp1 := r.NewCoupon()
	rand.Seed(int64(time.Now().Nanosecond()))
	rvalue = fmt.Sprintf("%v",rand.Intn(400000))
	cp1.CouponCode = fmt.Sprintf("%s%s","test-coupon-",rvalue)
	cp1.Name = "Coupon for API"
	cp1.DiscountType = "percent"
	cp1.DiscountPercent = 10
	cp1.SingleUse = false
	cp1.MaxRedemptions = "10"
	cp1.AppliesToAllPlans = false
	cp1.PlanCodes = &PlanCodes
	if err := cp1.Create(); err != nil {
		t.Fatalf("Create Coupon failed for coupon_code:%s has failed: %s",cp1.CouponCode, err.Error())
	} else {
		if _,err := r.GetCoupon(cp1.CouponCode); err != nil  {
			t.Fatalf("Couldn't find Coupon code:%s has failed: %s",cp1.CouponCode, err.Error())
		}
	}
	cp2 := r.NewCoupon()
	rand.Seed(int64(time.Now().Nanosecond()))
	rvalue = fmt.Sprintf("%v",rand.Intn(400000))
	cp2.CouponCode = fmt.Sprintf("%s%s","test-coupon-",rvalue)
	cp2.Name = "Coupon for API"
	cp2.DiscountType = "percent"
	cp2.DiscountPercent = 10
	cp2.SingleUse = false
	cp2.MaxRedemptions = "10"
	cp2.AppliesToAllPlans = false
	cp2.PlanCodes = &PlanCodes
	if err := cp2.Create(); err != nil {
		t.Fatalf("Create Coupon failed for coupon_code:%s has failed: %s",cp2.CouponCode, err.Error())
	} else {
		if _,err := r.GetCoupon(cp2.CouponCode); err != nil  {
			t.Fatalf("Couldn't find Coupon code:%s has failed: %s",cp2.CouponCode, err.Error())
		}
	}
	cp3 := r.NewCoupon()
	rand.Seed(int64(time.Now().Nanosecond()))
	rvalue = fmt.Sprintf("%v",rand.Intn(400000))
	cp3.CouponCode = fmt.Sprintf("%s%s","test-coupon-",rvalue)
	cp3.Name = "Coupon for API"
	cp3.DiscountType = "percent"
	cp3.DiscountPercent = 10
	cp3.SingleUse = false
	cp3.MaxRedemptions = "10"
	cp3.AppliesToAllPlans = false
	cp3.PlanCodes = &PlanCodes
	if err := cp3.Create(); err != nil {
		t.Fatalf("Create Coupon failed for coupon_code:%s has failed: %s",cp3.CouponCode, err.Error())
	} else {
		if _,err := r.GetCoupon(cp3.CouponCode); err != nil  {
			t.Fatalf("Couldn't find Coupon code:%s has failed: %s",cp3.CouponCode, err.Error())
		}
	}
	cp4 := r.NewCoupon()
	rand.Seed(int64(time.Now().Nanosecond()))
	rvalue = fmt.Sprintf("%v",rand.Intn(400000))
	cp4.CouponCode = fmt.Sprintf("%s%s","test-coupon-",rvalue)
	cp4.Name = "Coupon for API"
	cp4.DiscountType = "percent"
	cp4.DiscountPercent = 10
	cp4.SingleUse = false
	cp4.MaxRedemptions = "10"
	cp4.AppliesToAllPlans = false
	cp4.PlanCodes = &PlanCodes
	if err := cp4.Create(); err != nil {
		t.Fatalf("Create Coupon failed for coupon_code:%s has failed: %s",cp4.CouponCode, err.Error())
	} else {
		if _,err := r.GetCoupon(cp4.CouponCode); err != nil  {
			t.Fatalf("Couldn't find Coupon code:%s has failed: %s",cp4.CouponCode, err.Error())
		}
	}
	//create invalid coupon
	cp5 := r.NewCoupon()
	if err := cp5.Create(); err == nil {
		t.Fatalf("Coupon creation should have failed")
	}
	//deactivate coupon
	if err := cp1.Deactivate(); err != nil {
		t.Fatalf("Couldn't deactivate Coupon code:%s failed: %s",cp1.CouponCode, err.Error())
	}
	//END COUPON TESTING

	//COUPON REDEMPTION TESTING
	//redeem coupon 
	if err := cp4.Redeem(acc1.AccountCode,c.Currency); err != nil {
		t.Fatalf("Redemption failed for account_code:%s message:%s", acc1.AccountCode, err.Error())
	} else {
		//check if successful
		if red, err := r.GetCouponRedemption(acc1.AccountCode); err != nil {
			t.Fatalf("Error retreiving coupon redemption for account_code:%s err:%s",acc1.AccountCode,err.Error())
		} else {
			if red.Coupon.GetCode() != cp4.CouponCode {
				t.Fatalf("Coupon codes do not match")
			} else {
				//remove redemption
				if err = red.Delete();err != nil {
					t.Fatalf("Deletion error:%s",err.Error())
				}
			}
		}
	}
	//END COUPON REDEMPTION TESTING
	
	//COUPON LISTING TESTING 
	v = url.Values{}
	v.Set("per_page","1")
	v.Set("state","redeemable")
	if coupons, err := r.GetCoupons(v); err == nil {
		//page through
		for bcontinue := true; bcontinue; bcontinue = coupons.Next() {
		}
		//page backwards
		if !coupons.Prev() {
			t.Fatalf("Prev didn't work for coupons")
		}
		//page start
		if !coupons.Start() {
			t.Fatalf("Prev didn't work for coupons")
		}
	} else {
		t.Fatal(err.Error())
	}
	//END COUPON LISTING TESTING

	//INVOICE TESTING
	//generate an invoice from pending charges on account without billing info *2
	adj = r.NewAdjustment()
	adj.AccountCode = acc1.AccountCode
	adj.Description = "some extra charge"
	adj.UnitAmountInCents = 2000
	adj.Currency = c.Currency
	adj.Quantity = 1
	if err := adj.Create(); err != nil {
		t.Fatalf("Create Adjustment for account_code:%s has failed: %s",acc1.AccountCode, err.Error())
	}
	inv4 := r.NewInvoice()
	if err := inv4.InvoicePendingCharges(acc1.AccountCode); err != nil {
		t.Fatalf("Invoice Pending charges failed for account_code:%s message:%s", acc1.AccountCode, err.Error())
	}
	adj = r.NewAdjustment()
	adj.AccountCode = acc1.AccountCode
	adj.Description = "some extra charge"
	adj.UnitAmountInCents = 2000
	adj.Currency = c.Currency
	adj.Quantity = 1
	if err := adj.Create(); err != nil {
		t.Fatalf("Create Adjustment for account_code:%s has failed: %s",acc1.AccountCode, err.Error())
	}
	inv5 := r.NewInvoice()
	if err := inv5.InvoicePendingCharges(acc1.AccountCode); err != nil {
		t.Fatalf("Invoice Pending charges failed for account_code:%s message:%s", acc1.AccountCode, err.Error())
	}
	///generate an invoice from pending charges 3 times
	//generate inv1
	adj = r.NewAdjustment()
	adj.AccountCode = acc4.AccountCode
	adj.Description = "some extra charge"
	adj.UnitAmountInCents = 2000
	adj.Currency = c.Currency
	adj.Quantity = 1
	if err := adj.Create(); err != nil {
		t.Fatalf("Create Adjustment for account_code:%s has failed: %s",acc4.AccountCode, err.Error())
	}
	inv1 := r.NewInvoice()
	if err := inv1.InvoicePendingCharges(acc4.AccountCode); err != nil {
		t.Fatalf("Invoice Pending charges failed for account_code:%s message:%s", acc4.AccountCode, err.Error())
	}
	//create charge
	adj = r.NewAdjustment()
	adj.AccountCode = acc4.AccountCode
	adj.Description = "some extra charge"
	adj.UnitAmountInCents = 2000
	adj.Currency = c.Currency
	adj.Quantity = 1
	if err := adj.Create(); err != nil {
		t.Fatalf("Create Adjustment for account_code:%s has failed: %s",acc4.AccountCode, err.Error())
	}
	inv2 := r.NewInvoice()
	if err := inv2.InvoicePendingCharges(acc4.AccountCode); err != nil {
		t.Fatalf("Invoice Pending charges failed for account_code:%s message:%s", acc4.AccountCode, err.Error())
	}
	adj = r.NewAdjustment()
	adj.AccountCode = acc4.AccountCode
	adj.Description = "some extra charge"
	adj.UnitAmountInCents = 2000
	adj.Currency = c.Currency
	adj.Quantity = 1
	if err := adj.Create(); err != nil {
		t.Fatalf("Create Adjustment for account_code:%s has failed: %s",acc4.AccountCode, err.Error())
	}
	inv3 := r.NewInvoice()
	if err := inv3.InvoicePendingCharges(acc4.AccountCode); err != nil {
		t.Fatalf("Invoice Pending charges failed for account_code:%s message:%s", acc4.AccountCode, err.Error())
	}
	//get invoice for acc1
	if invoices, err := r.GetAccountInvoices(acc1.AccountCode,v); err == nil {
		marksuccesful := false
		for _, invoice := range invoices.Invoices {
			if marksuccesful {
				//mark invoice as failed
				if err = invoice.MarkFailed(); err != nil {
					t.Fatalf("Marking failed failed inv num:%s error:%s", invoice.InvoiceNumber, err.Error())
				}
			} else {
				//mark invoice as successful
				marksuccesful = true
				if err = invoice.MarkSuccessful(); err != nil {
					t.Fatalf("Marking successful failed inv num:%s error:%s", invoice.InvoiceNumber, err.Error())
				}
			}
		}
	}
	//END INVOICE TESTING

	//ACCOUNT INVOICE LISTING
	v = url.Values{}
	v.Set("per_page","1")
	if invoices, err := r.GetAccountInvoices(acc4.AccountCode,v); err == nil {
		//page through
		for bcontinue := true; bcontinue; bcontinue = invoices.Next() {
		}
		//page backwards
		if !invoices.Prev() {
			t.Fatalf("Prev didn't work for account invoices")
		}
		//page start
		if !invoices.Start() {
			t.Fatalf("Prev didn't work for account invoices")
		}
	} else {
		t.Fatal(err.Error())
	}
	//END ACCOUNT INVOICE LISTING

	//INVOICE LISTING
	v = url.Values{}
	v.Set("per_page","1")
	if invoices, err := r.GetInvoices(v); err == nil {
		//page through
		for bcontinue := true; bcontinue; bcontinue = invoices.Next() {
		}
		//page backwards
		if !invoices.Prev() {
			t.Fatalf("Prev didn't work for invoices")
		}
		//page start
		if !invoices.Start() {
			t.Fatalf("Prev didn't work for invoices")
		}
	} else {
		t.Fatal(err.Error())
	}
	//END INVOICE LISTING
}
