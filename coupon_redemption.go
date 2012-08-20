package gorecurly

import (
	"encoding/xml"
	"time"
)

//Coupon Redemption
type Redemption struct {
	XMLName                xml.Name `xml:"redemption"`
	r                      *Recurly
	Account			*AccountStub `xml:"account,omitempty"`
	Coupon			*CouponStub `xml:"coupon,omitempty"`
	AccountCode            string     `xml:"account_code,omitempty"`
	SingleUse              bool       `xml:"single_use,omitempty"`
	TotalDiscountedInCents int        `xml:"total_discounted_in_cents,omitempty"`
	Currency               string     `xml:"currency,omitempty"`
	CreatedAt              *time.Time `xml:"created_at,omitempty"`
}

//Remove a coupon from an account
func (r *Redemption) Delete() error {
	return r.r.doDelete(ACCOUNTS + "/" + r.Account.GetCode() + "/redemption")
}

