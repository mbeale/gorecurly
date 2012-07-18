package gorecurly

import (
	"encoding/xml"
	"net/url"
	"time"
)

//adjustment struct
type Adjustment struct{
	XMLName xml.Name `xml:"adjustment"`
	endpoint string
	r *Recurly
	Type string `xml:"type,attr"`
	AccountCode string `xml:"-"`
	UUID string `xml:"uuid,omitempty"`
	Description string `xml:"description,omitempty"`
	AccountingCode string `xml:"accounting_code,omitempty"`
	Origin string `xml:"origin,omitempty"`
	UnitAmountInCents int `xml:"unit_amount_in_cents,omitempty"`
	Quantity int `xml:"quantity,omitempty"`
	DiscountInCents int `xml:"discount_in_cents,omitempty"`
	TaxInCents int `xml:"tax_in_cents,omitempty"`
	Currency string `xml:"currency,omitempty"`
	Taxable bool `xml:"taxable,omitempty"`
	StartDate *time.Time `xml:"start_date,omitempty"`
	EndDate *time.Time `xml:"end_date,omitempty"`
	CreatedAt *time.Time `xml:"created_at,omitempty"`
}

//Create a new adjustment and load updated fields
func (a *Adjustment) Create() (error) {
	if a.UUID != "" {
		return RecurlyError{statusCode:400,Description:"Adjustment Already created"}
	}
	return a.r.doCreate(&a,ACCOUNTS + "/" + a.AccountCode + "/" + a.endpoint)
}

//delete and adjustment
func (a *Adjustment) Delete() (error) {
	return a.r.doDelete(a.endpoint + "/" + a.UUID)
}
type AdjustmentList struct {
	Paging
	r *Recurly
	AccountCode string
	XMLName xml.Name `xml:"adjustments"`
	Adjustments []Adjustment `xml:"adjustment"`
}

//Get next set of adjustments
func (a *AdjustmentList) Next() (bool) {
	if a.next != "" {
		v := url.Values{}
		v.Set("cursor",a.next)
		v.Set("per_page",a.perPage)
		*a,_ = a.r.GetAdjustments(a.AccountCode,v)
	} else {
		return false
	}
	return true
}

//Get previous set of accounts
func (a *AdjustmentList) Prev() ( bool) {
	if a.prev != "" {
		v := url.Values{}
		v.Set("cursor",a.prev)
		v.Set("per_page",a.perPage)
		*a,_ = a.r.GetAdjustments(a.AccountCode,v)
	} else {
		return false
	}
	return true
}

//Go to start set of accounts
func (a *AdjustmentList) Start() ( bool) {
	if a.prev != "" {
		v := url.Values{}
		v.Set("per_page",a.perPage)
		*a,_ = a.r.GetAdjustments(a.AccountCode,v)
	} else {
		return false
	}
	return true
}

