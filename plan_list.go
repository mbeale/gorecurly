package gorecurly

import (
	"encoding/xml"
	"net/url"
)

//Listing of plans
type PlanList struct {
	Paging
	r       *Recurly
	XMLName xml.Name `xml:"plans"`
	Plans   []Plan   `xml:"plan"`
}

//Get next set of Coupons
func (p *PlanList) Next() bool {
	if p.next != "" {
		v := url.Values{}
		v.Set("cursor", p.next)
		v.Set("per_page", p.perPage)
		*p, _ = p.r.GetPlans(v)
	} else {
		return false
	}
	return true
}

//Get previous set of accounts
func (p *PlanList) Prev() bool {
	if p.prev != "" {
		v := url.Values{}
		v.Set("cursor", p.prev)
		v.Set("per_page", p.perPage)
		*p, _ = p.r.GetPlans(v)
	} else {
		return false
	}
	return true
}

//Go to start set of accounts
func (p *PlanList) Start() bool {
	if p.prev != "" {
		v := url.Values{}
		v.Set("per_page", p.perPage)
		*p, _ = p.r.GetPlans(v)
	} else {
		return false
	}
	return true
}
