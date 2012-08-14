package gorecurly

import (
	"encoding/xml"
	"net/url"
)

//Adjustment Paging Struct
type AdjustmentList struct {
	Paging
	r           *Recurly
	AccountCode string
	XMLName     xml.Name     `xml:"adjustments"`
	Adjustments []Adjustment `xml:"adjustment"`
}

//Get next set of adjustments
func (a *AdjustmentList) Next() bool {
	if a.next != "" {
		v := url.Values{}
		v.Set("cursor", a.next)
		v.Set("per_page", a.perPage)
		*a, _ = a.r.GetAdjustments(a.AccountCode, v)
	} else {
		return false
	}
	return true
}

//Get previous set of accounts
func (a *AdjustmentList) Prev() bool {
	if a.prev != "" {
		v := url.Values{}
		v.Set("cursor", a.prev)
		v.Set("per_page", a.perPage)
		*a, _ = a.r.GetAdjustments(a.AccountCode, v)
	} else {
		return false
	}
	return true
}

//Go to start set of accounts
func (a *AdjustmentList) Start() bool {
	if a.prev != "" {
		v := url.Values{}
		v.Set("per_page", a.perPage)
		*a, _ = a.r.GetAdjustments(a.AccountCode, v)
	} else {
		return false
	}
	return true
}

