package gorecurly

import (
	"encoding/xml"
)

//Subscription pager
type SubscriptionList struct {
	Paging
	r       *Recurly
	XMLName xml.Name  `xml:"subscriptions"`
	Subscriptions []Subscription `xml:"subscription"`
}

//Get next set of subscriptions
func (s *SubscriptionList) Next() bool {
	if s.next != "" {
		*s, _ = s.r.GetSubscriptions(s.NextParams())
	} else {
		return false
	}
	return true
}

//Get previous set of subscriptions
func (s *SubscriptionList) Prev() bool {
	if s.prev != "" {
		*s, _ = s.r.GetSubscriptions(s.PrevParams())
	} else {
		return false
	}
	return true
}

//Go to start set of subscriptions
func (s *SubscriptionList) Start() bool {
	if s.prev != "" {
		*s, _ = s.r.GetSubscriptions(s.StartParams())
	} else {
		return false
	}
	return true
}

//List of subscriptions for an account
type AccountSubscriptionList struct {
	Paging
	r *Recurly
	XMLName xml.Name `xml:"subscriptions"`
	AccountCode string `xml:"-"`
	Subscriptions []Subscription `xml:"subscriptions"`
}


//Get next set of subscriptions
func (a *AccountSubscriptionList) Next() (bool) {
	if a.next != "" {
		*a,_ = a.r.GetAccountSubscriptions(a.AccountCode,a.NextParams())
	} else {
		return false
	}
	return true
}

//Get previous set of subscriptions
func (a *AccountSubscriptionList) Prev() ( bool) {
	if a.prev != "" {
		*a,_ = a.r.GetAccountSubscriptions(a.AccountCode,a.PrevParams())
	} else {
		return false
	}
	return true
}

//Go to start set of subscriptions
func (a *AccountSubscriptionList) Start() ( bool) {
	if a.prev != "" {
		*a,_ = a.r.GetAccountSubscriptions(a.AccountCode,a.StartParams())
	} else {
		return false
	}
	return true
}
