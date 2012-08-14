package gorecurly

import (
	"encoding/xml"
	"errors"
	"time"
)

//Adjustment struct
type Adjustment struct {
	XMLName           xml.Name `xml:"adjustment"`
	endpoint          string
	r                 *Recurly
	Type              string       `xml:"type,attr"`
	AccountCode       string       `xml:"-"`
	Account           *AccountStub `xml:"account,omitempty"`
	UUID              string       `xml:"uuid,omitempty"`
	Description       string       `xml:"description,omitempty"`
	AccountingCode    string       `xml:"accounting_code,omitempty"`
	Origin            string       `xml:"origin,omitempty"`
	UnitAmountInCents int          `xml:"unit_amount_in_cents,omitempty"`
	Quantity          int          `xml:"quantity,omitempty"`
	DiscountInCents   int          `xml:"discount_in_cents,omitempty"`
	TaxInCents        int          `xml:"tax_in_cents,omitempty"`
	Currency          string       `xml:"currency,omitempty"`
	Taxable           bool         `xml:"taxable,omitempty"`
	StartDate         *time.Time   `xml:"start_date,omitempty"`
	EndDate           RecurlyDate  `xml:"end_date,omitempty"`
	CreatedAt         *time.Time   `xml:"created_at,omitempty"`
}

//Create a new adjustment and load updated fields
func (a *Adjustment) Create() error {
	if a.UUID != "" {
		return RecurlyError{statusCode: 400, Description: "Adjustment Already created"}
	}
	return a.r.doCreate(&a, ACCOUNTS+"/"+a.AccountCode+"/"+a.endpoint)
}

//Delete an adjustment
func (a *Adjustment) Delete() error {
	return a.r.doDelete(a.endpoint + "/" + a.UUID)
}

func (a *Adjustment) GetAccount() (Account, error) {
	if a.Account == nil {
		return Account{}, errors.New("Account Stub is nil")
	}
	return a.r.GetAccount(a.Account.GetCode())
}
