package gorecurly

import (
	"encoding/xml"
	"time"
)

//Coupon object
type Coupon struct {
	XMLName            xml.Name `xml:"coupon"`
	endpoint           string
	r                  *Recurly
	AccountCode        string      `xml:"-"`
	CouponCode         string      `xml:"coupon_code"`
	Name               string      `xml:"name"`
	State              string      `xml:"state,omitempty"`
	HostedDescription  string      `xml:"hosted_description,omitempty"`
	InvoiceDescription string      `xml:"invoice_description,omitempty"`
	DiscountType       string      `xml:"discount_type,omitempty"`
	DiscountPercent    int         `xml:"discount_percent,omitempty"`
	//DiscountInCents    int         `xml:"discount_in_cents,omitempty"`
	RedeemByDate       RecurlyDate `xml:"redeem_by_date,omitempty"`
	SingleUse          bool        `xml:"single_use,omitempty"`
	AppliesForMonths   string      `xml:"applies_for_months,omitempty"`
	MaxRedemptions     string      `xml:"max_redemptions,omitempty"`
	AppliesToAllPlans  bool        `xml:"applies_to_all_plans,omitempty"`
	CreatedAt          *time.Time  `xml:"created_at,omitempty"`
	PlanCodes          *PlanCode   `xml:"plan_codes,omitempty"`
}

type createCoupon struct {
	XMLName            xml.Name   `xml:"coupon"`
	CouponCode         string     `xml:"coupon_code"`
	Name               string     `xml:"name"`
	HostedDescription  string     `xml:"hosted_description,omitempty"`
	InvoiceDescription string     `xml:"invoice_description,omitempty"`
	RedeemByDate       *time.Time `xml:"redeem_by_date,omitempty"`
	SingleUse          bool       `xml:"single_use,omitempty"`
	AppliesForMonths   string     `xml:"applies_for_months,omitempty"`
	MaxRedemptions     string     `xml:"max_redemptions,omitempty"`
	AppliesToAllPlans  bool       `xml:"applies_to_all_plans,omitempty"`
	DiscountType       string     `xml:"discount_type,omitempty"`
	DiscountPercent    int        `xml:"discount_percent,omitempty"`
	//DiscountInCents    int        `xml:"discount_in_cents,omitempty"`
	PlanCodes          *PlanCode  `xml:"plan_codes,omitempty"`
}

//Create a new coupon
func (c *Coupon) Create() error {
	if c.CreatedAt != nil {
		return RecurlyError{statusCode: 400, Description: "Coupon Already created"}
	}
	//return c.r.doCreate(&c, c.endpoint)
	cc := createCoupon{
		CouponCode:         c.CouponCode,
		Name:               c.Name,
		HostedDescription:  c.HostedDescription,
		InvoiceDescription: c.InvoiceDescription,
		SingleUse:          c.SingleUse,
		AppliesForMonths:   c.AppliesForMonths,
		MaxRedemptions:     c.MaxRedemptions,
		AppliesToAllPlans:  c.AppliesToAllPlans,
		DiscountType:       c.DiscountType,
		DiscountPercent:    c.DiscountPercent,
		//DiscountInCents:    c.DiscountInCents,
		PlanCodes:          c.PlanCodes,
	}
	gd, err := c.RedeemByDate.GetDate()
	if err == nil {
		cc.RedeemByDate = &gd
	}
	return c.r.doCreateReturn(cc, &c, c.endpoint)
}

//Redeem a coupon on an account
func (c *Coupon) Redeem(account_code string, currency string) error {
	redemption := Redemption{AccountCode: account_code, Currency: currency}
	redemption.r = c.r
	return redemption.r.doCreate(&redemption, c.endpoint+"/"+c.CouponCode+"/redeem")
}

//Deactivate a coupon
func (c *Coupon) Deactivate() error {
	return c.r.doDelete(c.endpoint + "/" + c.CouponCode)
}

//Coupon Stub struct
type CouponStub struct {
	XMLName xml.Name `xml:"coupon"`
	stub
}

