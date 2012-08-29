package gorecurly

import (
	"encoding/xml"
	"fmt"
	"time"
)

//Transaction Object
type Transaction struct {
	XMLName  xml.Name `xml:"transaction"`
	endpoint string
	r        *Recurly
	Account  *AccountStub `xml:"account,omitempty"`
	//Invoice *InvoiceStub `xml:"invoice,omitempty"`
	//Subscription *SubscriptionStub `xml:"subscription,omitempty"`
	EmbedAccount    *Account `xml:"-"`
	UUID            string   `xml:"uuid,omitempty"`
	Action          string   `xml:"action,omitempty"`
	State           string   `xml:"state,omitempty"`
	AmountInCents   int      `xml:"amount_in_cents,omitempty"`
	TaxInCents      int      `xml:"tax_in_cents,omitempty"`
	Currency        string   `xml:"currency,omitempty"`
	Status          string   `xml:"status,omitempty"`
	Reference       string   `xml:"reference,omitempty"`
	Test            bool     `xml:"test,omitempty"`
	Voidable        bool     `xml:"voidable,omitempty"`
	Refundable      bool     `xml:"refundable,omitempty"`
	CVVResult       string   `xml:"cvv_result,omitempty"`
	AVSResult       string   `xml:"avs_result,omitempty"`
	AVSResultStreet string   `xml:"avs_result_street,omitempty"`
	AVSResultPostal string   `xml:"avs_result_postal,omitempty"`
	//Details not implemented
	CreatedAt *time.Time `xml:"created_at,omitempty"`
}

type transactionCreate struct {
	XMLName       xml.Name `xml:"transaction"`
	Account       *Account `xml:"account,omitempty"`
	AmountInCents int      `xml:"amount_in_cents,omitempty"`
	Currency      string   `xml:"currency,omitempty"`
}

//Attach an existing account to a transaction
func (t *Transaction) AttachExistingAccount(a Account) (e error) {
	if t.UUID != "" {
		return RecurlyError{statusCode: 400, Description: "Subscription Already in Use and can't attach another account to it"}
	}
	t.EmbedAccount = new(Account)
	t.EmbedAccount.AccountCode = a.AccountCode
	return
}

//Attach a new account to a transaction
func (t *Transaction) AttachAccount(a Account) (e error) {
	if t.UUID != "" {
		return RecurlyError{statusCode: 400, Description: "Subscription Already in Use and can't attach another account to it"}
	}
	t.EmbedAccount = new(Account)
	a.CreatedAt = nil
	a.State = ""
	//some more may need to be blanked out
	a.HostedLoginToken = ""
	t.EmbedAccount = &a
	return
}

//Create a transaction
func (t *Transaction) Create() error {
	if t.UUID != "" {
		return RecurlyError{statusCode: 400, Description: "Subscription Already in Use"}
	}
	tc := transactionCreate{
		Account:       t.EmbedAccount,
		Currency:      t.Currency,
		AmountInCents: t.AmountInCents,
	}
	if err := t.r.doCreateReturn(tc, &t, t.endpoint); err == nil {
		return nil
	} else {
		return err
	}
	return nil
}

//Refund a partial amount from a transaction 
func (t *Transaction) Refund(amount int) error {
	return t.r.doDelete(t.endpoint + "/" + t.UUID + "?amount_in_cents=" + fmt.Sprintf("%v",amount))
}
//Completely refund a transaction
func (t *Transaction) RefundAll() error {
	return t.r.doDelete(t.endpoint + "/" + t.UUID)
}

