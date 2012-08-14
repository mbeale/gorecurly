package gorecurly

import (
	"encoding/xml"
	"net/url"
	"time"
)

type Redemption struct {
	XMLName                xml.Name `xml:"redemption"`
	r                      *Recurly
	AccountCode            string     `xml:"account_code,omitempty"`
	SingleUse              bool       `xml:"single_use,omitempty"`
	TotalDiscountedInCents int        `xml:"total_discounted_in_cents,omitempty"`
	Currency               string     `xml:"currency,omitempty"`
	CreatedAt              *time.Time `xml:"created_at,omitempty"`
}

//delete and redemption
func (r *Redemption) Delete() error {
	return r.r.doDelete(ACCOUNTS + "/" + r.AccountCode + "/redemption")
}

type Coupon struct {
	XMLName           xml.Name `xml:"coupon"`
	endpoint          string
	r                 *Recurly
	AccountCode       string     `xml:"-"`
	CouponCode        string     `xml:"coupon_code"`
	Name              string     `xml:"name"`
	State             string     `xml:"state,omitempty"`
	DiscountType      string     `xml:"discount_type,omitempty"`
	DiscountPercent   int        `xml:"discount_percent,omitempty"`
	RedeemByDate      *time.Time `xml:"redeem_by_date,omitempty"`
	SingleUse         bool       `xml:"single_use,omitempty"`
	AppliesForMonths  string     `xml:"applies_for_months,omitempty"`
	MaxRedemptions    int        `xml:"max_redemptions,omitempty"`
	AppliesToAllPlans bool       `xml:"applies_to_all_plans,omitempty"`
	CreatedAt         *time.Time `xml:"created_at,omitempty"`
	PlanCodes         PlanCode   `xml:"plan_codes,omitempty"`
}

//Create a new adjustment and load updated fields
func (c *Coupon) Create() error {
	if c.CreatedAt != nil {
		return RecurlyError{statusCode: 400, Description: "Coupon Already created"}
	}
	return c.r.doCreate(&c, c.endpoint)
}

//Redeem a coupon on an account
func (c *Coupon) Redeem(account_code string, currency string) error {
	redemption := Redemption{AccountCode: account_code, Currency: currency}
	redemption.r = c.r
	return redemption.r.doCreate(&redemption, c.endpoint+"/"+c.CouponCode+"/redeem")
}

//delete and adjustment
func (c *Coupon) Deactivate() error {
	return c.r.doDelete(c.endpoint + "/" + c.CouponCode)
}

type CouponList struct {
	Paging
	r       *Recurly
	XMLName xml.Name `xml:"coupons"`
	Coupons []Coupon `xml:"coupon"`
}

//Get next set of Coupons
func (c *CouponList) Next() bool {
	if c.next != "" {
		v := url.Values{}
		v.Set("cursor", c.next)
		v.Set("per_page", c.perPage)
		*c, _ = c.r.GetCoupons(v)
	} else {
		return false
	}
	return true
}

//Get previous set of accounts
func (c *CouponList) Prev() bool {
	if c.prev != "" {
		v := url.Values{}
		v.Set("cursor", c.prev)
		v.Set("per_page", c.perPage)
		*c, _ = c.r.GetCoupons(v)
	} else {
		return false
	}
	return true
}

//Go to start set of accounts
func (c *CouponList) Start() bool {
	if c.prev != "" {
		v := url.Values{}
		v.Set("per_page", c.perPage)
		*c, _ = c.r.GetCoupons(v)
	} else {
		return false
	}
	return true
}
