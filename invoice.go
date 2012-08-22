package gorecurly

import (
	"encoding/xml"
	"errors"
	"time"
)

//Invoice object
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
	//Transactions []Transaction `xml:"transactions,omitempty"`
}

//Invoice any pending charges given an acount code
func (i *Invoice) InvoicePendingCharges(account_code string) error {
	_, err := i.r.createRequest(ACCOUNTS + "/" + account_code + "/invoices", "POST", nil, nil)
	return err
}

//Mark an invoice as successfully paid
func (i *Invoice) MarkSuccessful() error {
	if i.UUID == "" {
		return errors.New("Not a valid invoice")
	}
	_, err := i.r.createRequest(INVOICES + "/" + i.InvoiceNumber + "/mark_successful", "PUT", nil, nil)
	return err
}

//Mark an invoice as failed
func (i *Invoice) MarkFailed() error {
	if i.UUID == "" {
		return errors.New("Not a valid invoice")
	}
	_, err := i.r.createRequest(INVOICES + "/" + i.InvoiceNumber + "/mark_failed", "PUT", nil, nil)
	return err
}


//Listing of line items in a transaction
type LineItems struct {
	XMLName    xml.Name `xml:"line_items"`
	Adjustment []Adjustment
}

