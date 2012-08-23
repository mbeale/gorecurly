package gorecurly

import (
	"encoding/xml"
)

//The invoice list struct
type InvoiceList struct {
	Paging
	r *Recurly
	XMLName xml.Name `xml:"invoices"`
	Invoices []Invoice `xml:"invoice"`
}


//Get next set of invoices
func (i *InvoiceList) Next() (bool) {
	if i.next != "" {
		*i,_ = i.r.GetInvoices(i.NextParams())
	} else {
		return false
	}
	return true
}

//Get previous set of invoices
func (i *InvoiceList) Prev() ( bool) {
	if i.prev != "" {
		*i,_ = i.r.GetInvoices(i.PrevParams())
	} else {
		return false
	}
	return true
}

//Go to start set of invoices
func (i *InvoiceList) Start() ( bool) {
	if i.prev != "" {
		*i,_ = i.r.GetInvoices(i.StartParams())
	} else {
		return false
	}
	return true
}

//Get the list of invoices by account
type AccountInvoiceList struct {
	InvoiceList
	AccountCode string `xml:"-"`
}


//Get next set of invoices by account
func (a *AccountInvoiceList) Next() (bool) {
	if a.next != "" {
		*a,_ = a.r.GetAccountInvoices(a.AccountCode,a.NextParams())
	} else {
		return false
	}
	return true
}

//Get previous set of invoices by account
func (a *AccountInvoiceList) Prev() ( bool) {
	if a.prev != "" {
		*a,_ = a.r.GetAccountInvoices(a.AccountCode,a.PrevParams())
	} else {
		return false
	}
	return true
}

//Go to start set of invoices by account
func (a *AccountInvoiceList) Start() ( bool) {
	if a.prev != "" {
		*a,_ = a.r.GetAccountInvoices(a.AccountCode,a.StartParams())
	} else {
		return false
	}
	return true
}
