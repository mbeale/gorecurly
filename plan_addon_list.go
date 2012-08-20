package gorecurly

import (
	"encoding/xml"
	"net/url"
)

//Plan add on list
type PlanAddOnList struct {
	Paging
	r        *Recurly
	XMLName  xml.Name    `xml:"add_ons"`
	PlanCode string      `xml:"-"`
	AddOns   []PlanAddOn `xml:"add_on"`
}

//Get next set of Coupons
func (p *PlanAddOnList) Next() bool {
	if p.next != "" {
		v := url.Values{}
		v.Set("cursor", p.next)
		v.Set("per_page", p.perPage)
		*p, _ = p.r.GetPlanAddOns(p.PlanCode,v)
	} else {
		return false
	}
	return true
}

//Get previous set of accounts
func (p *PlanAddOnList) Prev() bool {
	if p.prev != "" {
		v := url.Values{}
		v.Set("cursor", p.prev)
		v.Set("per_page", p.perPage)
		*p, _ = p.r.GetPlanAddOns(p.PlanCode,v)
	} else {
		return false
	}
	return true
}

//Go to start set of accounts
func (p *PlanAddOnList) Start() bool {
	if p.prev != "" {
		v := url.Values{}
		v.Set("per_page", p.perPage)
		*p, _ = p.r.GetPlanAddOns(p.PlanCode,v)
	} else {
		return false
	}
	return true
}
