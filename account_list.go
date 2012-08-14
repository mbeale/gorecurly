package gorecurly

import (
	"encoding/xml"
	"net/url"
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
		v := url.Values{}
		v.Set("cursor", a.next)
		v.Set("per_page", a.perPage)
		*a, _ = a.r.GetAccounts(v)
	} else {
		return false
	}
	return true
}

//Get previous set of accounts, will return false if no previous accounts
func (a *AccountList) Prev() bool {
	if a.prev != "" {
		v := url.Values{}
		v.Set("cursor", a.prev)
		v.Set("per_page", a.perPage)
		*a, _ = a.r.GetAccounts(v)
	} else {
		return false
	}
	return true
}

//Go to start set of accounts, returns false if no valid records
func (a *AccountList) Start() bool {
	if a.prev != "" {
		v := url.Values{}
		v.Set("per_page", a.perPage)
		*a, _ = a.r.GetAccounts(v)
	} else {
		return false
	}
	return true
}
