package gorecurly

import (
	"encoding/xml"
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
		*p, _ = p.r.GetPlans(p.NextParams())
	} else {
		return false
	}
	return true
}

//Get previous set of accounts
func (p *PlanList) Prev() bool {
	if p.prev != "" {
		*p, _ = p.r.GetPlans(p.PrevParams())
	} else {
		return false
	}
	return true
}

//Go to start set of accounts
func (p *PlanList) Start() bool {
	if p.prev != "" {
		*p, _ = p.r.GetPlans(p.StartParams())
	} else {
		return false
	}
	return true
}
