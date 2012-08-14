package gorecurly

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"time"
)

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

type TransactionCreate struct {
	XMLName       xml.Name `xml:"transaction"`
	Account       *Account `xml:"account,omitempty"`
	AmountInCents int      `xml:"amount_in_cents,omitempty"`
	Currency      string   `xml:"currency,omitempty"`
}

func (t *Transaction) AttachExistingAccount(a Account) (e error) {
	if t.UUID != "" {
		return RecurlyError{statusCode: 400, Description: "Subscription Already in Use and can't attach another account to it"}
	}
	t.EmbedAccount = new(Account)
	t.EmbedAccount.AccountCode = a.AccountCode
	return
}
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

func (t *Transaction) Create() error {
	if t.UUID != "" {
		return RecurlyError{statusCode: 400, Description: "Subscription Already in Use"}
	}
	tc := TransactionCreate{
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

//Refund a transaction 
func (t *Transaction) Refund(amount int) error {
	return t.r.doDelete(t.endpoint + "/" + t.UUID + "?amount_in_cents=" + fmt.Sprintf("%v",amount))
}
//Completely refund a transaction
func (t *Transaction) RefundAll() error {
	return t.r.doDelete(t.endpoint + "/" + t.UUID)
}

//Transaction pager
type TransactionList struct {
	Paging
	r       *Recurly
	XMLName xml.Name  `xml:"transactions"`
	Transactions []Transaction `xml:"transaction"`
}

//Get next set of transactions
func (t *TransactionList) Next() bool {
	if t.next != "" {
		v := url.Values{}
		v.Set("cursor", t.next)
		v.Set("per_page", t.perPage)
		*t, _ = t.r.GetTransactions(v)
	} else {
		return false
	}
	return true
}

//Get previous set of transactions
func (t *TransactionList) Prev() bool {
	if t.prev != "" {
		v := url.Values{}
		v.Set("cursor", t.prev)
		v.Set("per_page", t.perPage)
		*t, _ = t.r.GetTransactions(v)
	} else {
		return false
	}
	return true
}

//Go to start set of transactions
func (t *TransactionList) Start() bool {
	if t.prev != "" {
		v := url.Values{}
		v.Set("per_page", t.perPage)
		*t, _ = t.r.GetTransactions(v)
	} else {
		return false
	}
	return true
}

type AccountTransactionList struct {
	Paging
	r *Recurly
	XMLName xml.Name `xml:"transactions"`
	AccountCode string `xml:"-"`
	Transactions []Transaction `xml:"transaction"`
}


//Get next set of transactions
func (a *AccountTransactionList) Next() (bool) {
	if a.next != "" {
		v := url.Values{}
		v.Set("cursor",a.next)
		v.Set("per_page",a.perPage)
		*a,_ = a.r.GetAccountTransactions(a.AccountCode,v)
	} else {
		return false
	}
	return true
}

//Get previous set of transactions
func (a *AccountTransactionList) Prev() ( bool) {
	if a.prev != "" {
		v := url.Values{}
		v.Set("cursor",a.prev)
		v.Set("per_page",a.perPage)
		*a,_ = a.r.GetAccountTransactions(a.AccountCode,v)
	} else {
		return false
	}
	return true
}

//Go to start set of transactions
func (a *AccountTransactionList) Start() ( bool) {
	if a.prev != "" {
		v := url.Values{}
		v.Set("per_page",a.perPage)
		*a,_ = a.r.GetAccountTransactions(a.AccountCode,v)
	} else {
		return false
	}
	return true
}
