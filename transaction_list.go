package gorecurly

import (
	"encoding/xml"
)


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
		*t, _ = t.r.GetTransactions(t.NextParams())
	} else {
		return false
	}
	return true
}

//Get previous set of transactions
func (t *TransactionList) Prev() bool {
	if t.prev != "" {
		*t, _ = t.r.GetTransactions(t.PrevParams())
	} else {
		return false
	}
	return true
}

//Go to start set of transactions
func (t *TransactionList) Start() bool {
	if t.prev != "" {
		*t, _ = t.r.GetTransactions(t.StartParams())
	} else {
		return false
	}
	return true
}


//A listing of transactions by account
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
		*a,_ = a.r.GetAccountTransactions(a.AccountCode,a.NextParams())
	} else {
		return false
	}
	return true
}

//Get previous set of transactions
func (a *AccountTransactionList) Prev() ( bool) {
	if a.prev != "" {
		*a,_ = a.r.GetAccountTransactions(a.AccountCode,a.PrevParams())
	} else {
		return false
	}
	return true
}

//Go to start set of transactions
func (a *AccountTransactionList) Start() ( bool) {
	if a.prev != "" {
		*a,_ = a.r.GetAccountTransactions(a.AccountCode,a.StartParams())
	} else {
		return false
	}
	return true
}

