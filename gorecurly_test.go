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
	Future string `xml:"futuredate"`
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

	//create valid account with billing info
	acc5 := r.NewAccount()
	rand.Seed(int64(time.Now().Nanosecond()))
	rvalue = fmt.Sprintf("%v",rand.Intn(400000))
	acc5.AccountCode = fmt.Sprintf("%s%s","test-account-",rvalue)
	acc5.Email = "test-email-" + rvalue + "@example.com"
	acc5.FirstName = "test-fname-" + rvalue
	acc5.LastName = "test-last-" + rvalue
	acc5.B = new(BillingInfo)
	acc5.B.FirstName = "test-fname-" + rvalue
	acc5.B.LastName = "test-last-" + rvalue
	acc5.B.Number = "4111111111111111"
	acc5.B.Month = 12
	acc5.B.Year = 2015
	acc5.B.VerificationValue = "123"
	if err := acc5.Create(); err != nil {
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
	v.Set("state","active")
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
	//PLAN ADD ON TESTING

	//create 4 add ons
	addon1 := r.NewPlanAddOn()
	rand.Seed(int64(time.Now().Nanosecond()))
	rvalue = fmt.Sprintf("%v",rand.Intn(400000))
	addon1.Name = "Some updated addon"
	addon1.AddOnCode = fmt.Sprintf("%s%s","test-addon-",rvalue)
	addon1.UnitAmountInCents.SetCurrency(c.Currency,400)
	if err := addon1.Create(plan3.PlanCode); err != nil {
		t.Fatalf("Create add on failed for plan_code:%s addoncode:%s has failed: %s",plan3.PlanCode, err.Error())
	}

	addon2 := r.NewPlanAddOn()
	rand.Seed(int64(time.Now().Nanosecond()))
	rvalue = fmt.Sprintf("%v",rand.Intn(400000))
	addon2.Name = "Some updated addon"
	addon2.AddOnCode = fmt.Sprintf("%s%s","test-addon-",rvalue)
	addon2.UnitAmountInCents.SetCurrency(c.Currency,400)
	if err := addon2.Create(plan3.PlanCode); err != nil {
		t.Fatalf("Create add on failed for plan_code:%s addoncode:%s has failed: %s",plan3.PlanCode, err.Error())
	}

	addon3 := r.NewPlanAddOn()
	rand.Seed(int64(time.Now().Nanosecond()))
	rvalue = fmt.Sprintf("%v",rand.Intn(400000))
	addon3.Name = "Some updated addon"
	addon3.AddOnCode = fmt.Sprintf("%s%s","test-addon-",rvalue)
	addon3.UnitAmountInCents.SetCurrency(c.Currency,400)
	if err := addon3.Create(plan3.PlanCode); err != nil {
		t.Fatalf("Create add on failed for plan_code:%s addoncode:%s has failed: %s",plan3.PlanCode, err.Error())
	}

	addon4 := r.NewPlanAddOn()
	rand.Seed(int64(time.Now().Nanosecond()))
	rvalue = fmt.Sprintf("%v",rand.Intn(400000))
	addon4.Name = "Some updated addon"
	addon4.AddOnCode = fmt.Sprintf("%s%s","test-addon-",rvalue)
	addon4.UnitAmountInCents.SetCurrency(c.Currency,400)
	if err := addon4.Create(plan3.PlanCode); err != nil {
		t.Fatalf("Create add on failed for plan_code:%s addoncode:%s has failed: %s",plan3.PlanCode, err.Error())
	}

	//update addon
	addon1.UnitAmountInCents.SetCurrency(c.Currency,800)
	addon1.Update()
	if amt, _ := addon1.UnitAmountInCents.GetCurrency(c.Currency);amt != 800 {
		t.Fatalf("Update not successful for plan_code:%s add_on:%s", plan3.PlanCode, addon1.AddOnCode)
	}
	//delete add on
	if err := addon4.Delete(); err != nil {
		t.Fatalf("Delete Addon failed for plan_code:%s , addoncode: %s has failed: %s",plan3.PlanCode,addon4.AddOnCode, err.Error())
	}

	//END PLAN ADD ON TESTING

	//PLAN ADD ON LISTING TESTING
	v = url.Values{}
	v.Set("per_page","1")
	if addons, err :=r.GetPlanAddOns(plan3.PlanCode,v); err == nil {
		//page through
		for bcontinue := true; bcontinue; bcontinue = addons.Next() {
		}
		//page backwards
		if !addons.Prev() {
			t.Fatalf("Prev didn't work for addons")
		}
		//page start
		if !addons.Start() {
			t.Fatalf("Prev didn't work for addons")
		}
	} else {
		t.Fatal(err.Error())
	}
	//END PLAN ADD ON LISTING TESTING
	
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

	//SUBSCRIPTION TESTING
	//create 4 subs
	sub1 := r.NewSubscription()
	sub1.PlanCode = plan3.PlanCode
	sub1.Currency = c.Currency
	sub1.AttachExistingAccount(acc4)
	if err := sub1.Create(); err != nil {
		t.Fatalf("Subscription failed to be created with account:%s error:%s", acc4.AccountCode,err)
	}

	sub2 := r.NewSubscription()
	sub2.PlanCode = plan2.PlanCode
	sub2.Currency = c.Currency
	sub2.AttachExistingAccount(acc4)
	if err := sub2.Create(); err != nil {
		t.Fatalf("Subscription failed to be created with account:%s error:%s", acc4.AccountCode,err)
	}

	sub3 := r.NewSubscription()
	sub3.PlanCode = plan4.PlanCode
	sub3.Currency = c.Currency
	sub3.AttachExistingAccount(acc4)
	if err := sub3.Create(); err != nil {
		t.Fatalf("Subscription failed to be created with account:%s error:%s", acc4.AccountCode,err)
	}
	
	//create a sub w/ addons
	sub5 := r.NewSubscription()
	sub5.PlanCode = plan3.PlanCode
	sub5.Currency = c.Currency
	sub5.AttachExistingAccount(acc5)
	addons := EmbedPlanAddOn{Quantity:1,AddOnCode:addon3.AddOnCode}
	sub5.SubscriptionAddOns.UpdateAddOns(addons)
	if err := sub5.Create(); err != nil {
		t.Fatalf("Subscription failed to be created with addons account:%s err:", acc5.AccountCode,err)
	}
	//update a sub
	sub3.Quantity = "2"
	if err := sub3.Update(true);err!=nil{
		t.Fatalf("Subscription failed to be updates :%s", sub3.UUID)
	} else {
		//verify quantity update
		if sub3.Quantity != "2" {
			t.Fatalf("Subscription failed to be updates :%s qty not = 2", sub3.UUID)
		}
	}
	//update a sub w/addons
	addons.Quantity = 4
	sub5.SubscriptionAddOns.UpdateAddOns(addons)
	if err := sub5.Update(true);err!=nil{
		t.Fatalf("Subscription w/ addons failed to be updates :%s err:%s\n%v", sub5.UUID,err,sub5.SubscriptionAddOns)
	} else {
		//verify quantity update
		if e, compareaddon := sub5.SubscriptionAddOns.GetAddOn(addon3.AddOnCode); e != nil {
			t.Fatalf("Subscription w/ addons failed to be updated :%s error:%s", sub5.UUID,err)
		} else {
			if compareaddon.Quantity != 4 {
				t.Fatalf("Subscription w/ addons failed to be updates :%s qty not = 4", sub5.UUID)
			}
		}
	}
	//cancel a sub
	if err := sub1.Cancel(); err != nil {
		t.Fatalf("Subscription failed to be update :%si error:%s", sub1.UUID,err)
	} else {
		if sub1.State != "canceled" {
			t.Fatalf("Subscription failed to be cancelled :%s", sub1.UUID)
		}
	}
	//reactivate a sub
	if err := sub1.Reactivate(); err != nil {
		t.Fatalf("Subscription failed to be reactivated :%s error:%s", sub1.UUID,err)
	} else {
		if sub1.State != "active" {
			t.Fatalf("Subscription failed to be reactivated :%s", sub1.UUID)
		}
	}
	//terminate a sub
	if err := sub1.Terminate(); err != nil {
		t.Fatalf("Subscription failed to be terminated :%s error:%s", sub1.UUID,err)
	} else {
		if sub1.State != "expired" {
			t.Fatalf("Subscription failed to be terminated :%s", sub1.UUID)
		}
	}
	//postpone a sub
	threedaysfromnow := time.Now().AddDate(0,0,3)
	if err:= sub3.Postpone(threedaysfromnow); err != nil {
		t.Fatalf("Subscription failed to be postponed :%s error:%s", sub3.UUID,err)
	}
	//END SUSCRIPTION TESTING

	//SUBSCRIPTION LISTING TESTING
	v = url.Values{}
	v.Set("per_page","1")
	if subscriptions, err := r.GetSubscriptions(v); err == nil {
		//page through
		for bcontinue := true; bcontinue; bcontinue = subscriptions.Next() {
		}
		//page backwards
		if !subscriptions.Prev() {
			t.Fatalf("Prev didn't work for subscriptions")
		}
		//page start
		if !subscriptions.Start() {
			t.Fatalf("Prev didn't work for subscriptions")
		}
	} else {
		t.Fatal(err.Error())
	}
	//END SUBSCRIPTION LISTING

	//ACCOUNT SUBSCRIPTION LISTING TESTING
	v = url.Values{}
	v.Set("per_page","1")
	if subscriptions, err := r.GetAccountSubscriptions(acc4.AccountCode,v); err == nil {
		//page through
		for bcontinue := true; bcontinue; bcontinue = subscriptions.Next() {
		}
		//page backwards
		if !subscriptions.Prev() {
			t.Fatalf("Prev didn't work for account subscriptions")
		}
		//page start
		if !subscriptions.Start() {
			t.Fatalf("Prev didn't work for account subscriptions")
		}
	} else {
		t.Fatal(err.Error())
	}
	//END ACCOUNT SUBSCRIPTION LISTING TESTING

	//TRANSACTION TESTING
	//create 5 transactions
	trans1 := r.NewTransaction()
	trans1.AttachExistingAccount(acc4)
	trans1.AmountInCents = 500
	trans1.Currency = c.Currency
	if err := trans1.Create(); err != nil {
		t.Fatalf("Transaction failed to be created account:%s err:", acc4.AccountCode,err)
	}

	trans2 := r.NewTransaction()
	trans2.AttachExistingAccount(acc4)
	trans2.AmountInCents = 500
	trans2.Currency = c.Currency
	if err := trans2.Create(); err != nil {
		t.Fatalf("Transaction failed to be created account:%s err:", acc4.AccountCode,err)
	}

	trans3 := r.NewTransaction()
	trans3.AttachExistingAccount(acc4)
	trans3.AmountInCents = 500
	trans3.Currency = c.Currency
	if err := trans3.Create(); err != nil {
		t.Fatalf("Transaction failed to be created account:%s err:", acc4.AccountCode,err)
	}

	trans4 := r.NewTransaction()
	trans4.AttachExistingAccount(acc4)
	trans4.AmountInCents = 500
	trans4.Currency = c.Currency
	if err := trans4.Create(); err != nil {
		t.Fatalf("Transaction failed to be created account:%s err:", acc4.AccountCode,err)
	}

	trans5 := r.NewTransaction()
	trans5.AttachExistingAccount(acc4)
	trans5.AmountInCents = 500
	trans5.Currency = c.Currency
	if err := trans5.Create(); err != nil {
		t.Fatalf("Transaction failed to be created account:%s err:", acc4.AccountCode,err)
	}

	//full refund a transaction
	if err:= trans5.RefundAll(); err != nil {
		t.Fatalf("Transaction failed to be refunded uuid:%s err:", trans5.UUID,err)
	}
	//partial refund a transaction
	if err:= trans4.Refund(250); err != nil {
		t.Fatalf("Transaction failed to be partially refunded uuid:%s err:", trans4.UUID,err)
	}
	//END TRNSACTION TESTING

	//TRANSACTION LISTING TESTING
	v = url.Values{}
	v.Set("per_page","1")
	if transactions, err := r.GetTransactions(v); err == nil {
		//page through
		for bcontinue := true; bcontinue; bcontinue = transactions.Next() {
		}
		//page backwards
		if !transactions.Prev() {
			t.Fatalf("Prev didn't work for transactions")
		}
		//page start
		if !transactions.Start() {
			t.Fatalf("Prev didn't work for transactions")
		}
	} else {
		t.Fatal(err.Error())
	}
	//END TRANSACTION LISTING TESTING

	//ACCOUNT TRANSACTION LISTING
	v = url.Values{}
	v.Set("per_page","1")
	if transactions, err := r.GetAccountTransactions(acc4.AccountCode,v); err == nil {
		//page through
		for bcontinue := true; bcontinue; bcontinue = transactions.Next() {
		}
		//page backwards
		if !transactions.Prev() {
			t.Fatalf("Prev didn't work for account transactions")
		}
		//page start
		if !transactions.Start() {
			t.Fatalf("Prev didn't work for account transactions")
		}
	} else {
		t.Fatal(err.Error())
	}
	//END ACCOUNT TRANSACTION LISTING
}
