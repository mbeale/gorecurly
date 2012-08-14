package gorecurly

import (
	"encoding/xml"
	"net/url"
	"time"
)

type Invoice struct {
	XMLName xml.Name `xml:"invoice"`
	endpoint string
	r *Recurly
	Account *AccountStub `xml:"account,omitempty"`
	UUID string `xml:"uuid,omitempty"`
	State string `xml:"state,omitempty"`
	InvoiceNumber string `xml:"invoice_number,omitempty"`
	PONumber string `xml:"po_number,omitempty"`
	VATNumber string `xml:"vat_number,omitempty"`
	SubtotalInCents int `xml:"subtotal_in_cents,omitempty"`
	TaxInCents int `xml:"tax_in_cents,omitempty"`
	TotalInCents int `xml:"total_in_cents,omitempty"`
	Currency string `xml:"currency,omitempty"`
	CreatedAt *time.Time `xml:"created_at,omitempty"`
	LineItems []LineItems `xml:"line_items,omitempty"`
	//Transactions Transaction `xml:"transactions,omitempty"`
}

type InvoiceList struct {
	Paging
	r *Recurly
	XMLName xml.Name `xml:"invoices"`
	Invoices []Invoice `xml:"invoice"`
}


//Get next set of Coupons
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

//Get previous set of accounts
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

//Go to start set of accounts
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

type AccountInvoiceList struct {
	Paging
	r *Recurly
	XMLName xml.Name `xml:"invoices"`
	AccountCode string `xml:"-"`
	Invoices []Invoice `xml:"invoice"`
}


//Get next set of Coupons
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

//Get previous set of accounts
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

//Go to start set of accounts
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
