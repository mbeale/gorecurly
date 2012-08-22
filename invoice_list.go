package gorecurly

import (
	"encoding/xml"
	"net/url"
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
		v := url.Values{}
		v.Set("cursor",i.next)
		v.Set("per_page",i.perPage)
		*i,_ = i.r.GetInvoices(v)
	} else {
		return false
	}
	return true
}

//Get previous set of invoices
func (i *InvoiceList) Prev() ( bool) {
	if i.prev != "" {
		v := url.Values{}
		v.Set("cursor",i.prev)
		v.Set("per_page",i.perPage)
		*i,_ = i.r.GetInvoices(v)
	} else {
		return false
	}
	return true
}

//Go to start set of invoices
func (i *InvoiceList) Start() ( bool) {
	if i.prev != "" {
		v := url.Values{}
		v.Set("per_page",i.perPage)
		*i,_ = i.r.GetInvoices(v)
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
		v := url.Values{}
		v.Set("cursor",a.next)
		v.Set("per_page",a.perPage)
		*a,_ = a.r.GetAccountInvoices(a.AccountCode,v)
	} else {
		return false
	}
	return true
}

//Get previous set of invoices by account
func (a *AccountInvoiceList) Prev() ( bool) {
	if a.prev != "" {
		v := url.Values{}
		v.Set("cursor",a.prev)
		v.Set("per_page",a.perPage)
		*a,_ = a.r.GetAccountInvoices(a.AccountCode,v)
	} else {
		return false
	}
	return true
}

//Go to start set of invoices by account
func (a *AccountInvoiceList) Start() ( bool) {
	if a.prev != "" {
		v := url.Values{}
		v.Set("per_page",a.perPage)
		*a,_ = a.r.GetAccountInvoices(a.AccountCode,v)
	} else {
		return false
	}
	return true
}
