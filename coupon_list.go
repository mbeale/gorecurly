package gorecurly

import (
	"encoding/xml"
	"net/url"
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
		v := url.Values{}
		v.Set("cursor", c.next)
		v.Set("per_page", c.perPage)
		*c, _ = c.r.GetCoupons(v)
	} else {
		return false
	}
	return true
}

//Get previous set of coupons
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

//Go to start set of coupons
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
