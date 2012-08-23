package gorecurly

import (
	"encoding/xml"
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
		*a, _ = a.r.GetAdjustments(a.AccountCode, a.NextParams())
	} else {
		return false
	}
	return true
}

//Get previous set of accounts
func (a *AdjustmentList) Prev() bool {
	if a.prev != "" {
		*a, _ = a.r.GetAdjustments(a.AccountCode, a.PrevParams())
	} else {
		return false
	}
	return true
}

//Go to start set of accounts
func (a *AdjustmentList) Start() bool {
	if a.prev != "" {
		*a, _ = a.r.GetAdjustments(a.AccountCode, a.StartParams())
	} else {
		return false
	}
	return true
}

