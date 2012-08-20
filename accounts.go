package gorecurly

import (
	"encoding/xml"
	"time"
)

//Account struct
type Account struct {
	XMLName          xml.Name `xml:"account"`
	endpoint         string
	r                *Recurly
	AccountCode      string       `xml:"account_code"`
	Username         string       `xml:"username,omitempty"`
	Email            string       `xml:"email,omitempty"`
	State            string       `xml:"state,omitempty"`
	FirstName        string       `xml:"first_name,omitempty"`
	LastName         string       `xml:"last_name,omitempty"`
	CompanyName      string       `xml:"company_name,omitempty"`
	AcceptLanguage   string       `xml:"accept_language,omitempty"`
	HostedLoginToken string       `xml:"hosted_login_token,omitempty"`
	CreatedAt        *time.Time   `xml:"created_at,omitempty"`
	B                *BillingInfo `xml:"billing_info,omitempty"`
}

//Load the billing information for this account
func (a *Account) LoadBilling() error {
	bi, err := a.r.GetBillingInfo(a.AccountCode)
	if a.B == nil {
		a.B = new(BillingInfo)
	}
	a.B = &bi
	return err
}

//Return adjustments for this account
func (a *Account) GetAdjustments() (AdjustmentList, error) {
	return a.r.GetAdjustments(a.AccountCode)
}

//Return invoices for this account
func (a *Account) GetInvoices() (AccountInvoiceList, error) {
	return a.r.GetAccountInvoices(a.AccountCode)
}

//Return subscriptions for this account
func (a *Account) GetSubscriptions() (AccountSubscriptionList, error) {
	return a.r.GetAccountSubscriptions(a.AccountCode)
}

//Return transactions for this account
func (a *Account) GetTransactions() (AccountTransactionList, error) {
	return a.r.GetAccountTransactions(a.AccountCode)
}

//Create a new account and load updated fields
func (a *Account) Create() error {
	if a.CreatedAt != nil || a.HostedLoginToken != "" || a.State != "" {
		return RecurlyError{statusCode: 400, Description: "Account Code Already in Use"}
	}
	err := a.r.doCreate(&a, a.endpoint)
	if err == nil {
		a.B = nil
	}
	return err
}

//Update an account 
func (a *Account) Update() error {
	newaccount := new(Account)
	*newaccount = *a
	newaccount.State = ""
	newaccount.HostedLoginToken = ""
	newaccount.CreatedAt = nil
	newaccount.B = nil
	return a.r.doUpdate(newaccount, a.endpoint+"/"+a.AccountCode)
}

//Close an account
func (a *Account) Close() error {
	return a.r.doDelete(a.endpoint + "/" + a.AccountCode)
}

//Close an account
func (a *Account) RemoveRedemption() error {
	return a.r.doDelete(a.endpoint + "/" + a.AccountCode + "/redemption")
}

//Reopen a closed account
func (a *Account) Reopen() error {
	newaccount := new(Account)
	return a.r.doUpdate(newaccount, a.endpoint+"/"+a.AccountCode+"/reopen")
}

//Account Stub struct
type AccountStub struct {
	XMLName xml.Name `xml:"account"`
	stub
}

