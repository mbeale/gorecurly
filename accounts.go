package gorecurly

import (
	"encoding/xml"
	"net/url"
	"time"
)

//Account struct
type Account struct{
	XMLName xml.Name `xml:"account"`
	endpoint string
	r *Recurly
	AccountCode string `xml:"account_code"`
	Username string `xml:"username"`
	Email string `xml:"email"`
	State string `xml:"state,omitempty"`
	FirstName string `xml:"first_name"`
	LastName string `xml:"last_name"`
	CompanyName string `xml:"company_name"`
	AcceptLanguage string `xml:"accept_language"`
	HostedLoginToken string `xml:"hosted_login_token,omitempty"`
	CreatedAt *time.Time `xml:"created_at,omitempty"`
	B *BillingInfo `xml:"billing_info,omitempty"` 
}

//Create a new account and load updated fields
func (a *Account) Create() (error) {
	if a.CreatedAt != nil || a.HostedLoginToken != "" || a.State != "" {
		return RecurlyError{statusCode:400,Description:"Account Code Already in Use"}
	}
	err := a.r.doCreate(&a,a.endpoint)
	if err == nil {
		a.B = nil
	}
	return err
}

//Update an account 
func (a *Account) Update() (error) {
	newaccount := new(Account)
	*newaccount = *a
	newaccount.State = ""
	newaccount.HostedLoginToken = ""
	newaccount.CreatedAt = nil
	newaccount.B = nil
	return a.r.doUpdate(newaccount,a.endpoint + "/" + a.AccountCode)
}

//Close an account
func (a *Account) Close() (error) {
	return a.r.doDelete(a.endpoint + "/" + a.AccountCode)
}

//Reopen a closed account
func (a *Account) Reopen() (error) {
	newaccount := new(Account)
	return a.r.doUpdate(newaccount,a.endpoint + "/" + a.AccountCode + "/reopen")
}

//Account Stub struct
type AccountStub struct {
	XMLName xml.Name `xml:"account"`
	stub
}


//Account pager
type AccountList struct {
	Paging
	r *Recurly
	XMLName xml.Name `xml:"accounts"`
	Account []Account `xml:"account"`
}

//Get next set of accounts
func (a *AccountList) Next() (bool) {
	if a.next != "" {
		v := url.Values{}
		v.Set("cursor",a.next)
		v.Set("per_page",a.perPage)
		*a,_ = a.r.GetAccounts(v)
	} else {
		return false
	}
	return true
}

//Get previous set of accounts
func (a *AccountList) Prev() ( bool) {
	if a.prev != "" {
		v := url.Values{}
		v.Set("cursor",a.prev)
		v.Set("per_page",a.perPage)
		*a,_ = a.r.GetAccounts(v)
	} else {
		return false
	}
	return true
}

//Go to start set of accounts
func (a *AccountList) Start() ( bool) {
	if a.prev != "" {
		v := url.Values{}
		v.Set("per_page",a.perPage)
		*a,_ = a.r.GetAccounts(v)
	} else {
		return false
	}
	return true
}

