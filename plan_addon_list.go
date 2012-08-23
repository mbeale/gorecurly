package gorecurly

import (
	"encoding/xml"
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
		*p, _ = p.r.GetPlanAddOns(p.PlanCode,p.NextParams())
	} else {
		return false
	}
	return true
}

//Get previous set of accounts
func (p *PlanAddOnList) Prev() bool {
	if p.prev != "" {
		*p, _ = p.r.GetPlanAddOns(p.PlanCode,p.PrevParams())
	} else {
		return false
	}
	return true
}

//Go to start set of accounts
func (p *PlanAddOnList) Start() bool {
	if p.prev != "" {
		*p, _ = p.r.GetPlanAddOns(p.PlanCode,p.StartParams())
	} else {
		return false
	}
	return true
}
