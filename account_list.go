package gorecurly

import (
	"encoding/xml"
)

//Account pager
type AccountList struct {
	Paging
	r       *Recurly
	XMLName xml.Name  `xml:"accounts"`
	Account []Account `xml:"account"`
}

//Get next set of accounts, will return false if no more accounts
func (a *AccountList) Next() bool {
	if a.next != "" {
		*a, _ = a.r.GetAccounts(a.NextParams())
	} else {
		return false
	}
	return true
}

//Get previous set of accounts, will return false if no previous accounts
func (a *AccountList) Prev() bool {
	if a.prev != "" {
		*a, _ = a.r.GetAccounts(a.PrevParams())
	} else {
		return false
	}
	return true
}

//Go to start set of accounts, returns false if no valid records
func (a *AccountList) Start() bool {
	if a.prev != "" {
		*a, _ = a.r.GetAccounts(a.StartParams())
	} else {
		return false
	}
	return true
}
