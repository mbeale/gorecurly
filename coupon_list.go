package gorecurly

import (
	"encoding/xml"
)

//Coupon List Struct
type CouponList struct {
	Paging
	r       *Recurly
	XMLName xml.Name `xml:"coupons"`
	Coupons []Coupon `xml:"coupon"`
}

//Get next set of Coupons
func (c *CouponList) Next() bool {
	if c.next != "" {
		*c, _ = c.r.GetCoupons(c.NextParams())
	} else {
		return false
	}
	return true
}

//Get previous set of coupons
func (c *CouponList) Prev() bool {
	if c.prev != "" {
		*c, _ = c.r.GetCoupons(c.PrevParams())
	} else {
		return false
	}
	return true
}

//Go to start set of coupons
func (c *CouponList) Start() bool {
	if c.prev != "" {
		*c, _ = c.r.GetCoupons(c.StartParams())
	} else {
		return false
	}
	return true
}
