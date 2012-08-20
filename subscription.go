package gorecurly

import (
	"encoding/xml"
	"errors"
	"net/url"
	"time"
)

type Subscription struct {
	XMLName                xml.Name `xml:"subscription"`
	endpoint               string
	r                      *Recurly
	Timeframe              string          `xml:"timeframe,omitempty"`
	Account                *AccountStub    `xml:"account,omitempty"`
	EmbedAccount           *Account        `xml:"-"`
	Plan                   *PlanStub       `xml:"plan,omitempty"`
	PlanCode               string          `xml:"-"`
	UUID                   string          `xml:"uuid,omitempty"`
	State                  string          `xml:"state,omitempty"`
	UnitAmountInCents      int             `xml:"unit_amount_in_cents,omitempty"`
	CouponCode             string          `xml:"-"`
	Currency               string          `xml:"currency,omitempty"`
	Quantity               string          `xml:"quantity,omitempty"`
	ActivatedAt            *time.Time      `xml:"activated_at,omitempty"`
	CanceledAt             RecurlyDate     `xml:"canceled_at,omitempty"`
	ExpiresAt              RecurlyDate     `xml:"expires_at,omitempty"`
	CurrentPeriodStartedAt *time.Time      `xml:"current_period_starts_at,omitempty"`
	CurrentPeriodEndsAt    *time.Time      `xml:"currenct_period_ends_at,omitempty"`
	RemainingBillingCycles string          `xml:"remaining_billing_cycles,omitempty"`
	TrialStartedAt         RecurlyDate     `xml:"trial_started_at,omitempty"`
	TrialEndsAt            RecurlyDate     `xml:"trial_ends_at,omitempty"`
	SubscriptionAddOns     EmbedPlanAddOns `xml:"subscription_add_ons,omitempty"`
	StartsAt               *time.Time      `xml:"-"`
	FirstRenewalDate       *time.Time      `xml:"-"`
	TotalBillingCycles     string          `xml:"total_billing_cycles,omitempty"`
}

type subscriptionCreate struct {
	XMLName            xml.Name        `xml:"subscription"`
	PlanCode           string          `xml:"plan_code,omitempty"`
	CouponCode         string          `xml:"coupon_code,omitempty"`
	Account            *Account        `xml:"account,omitempty"`
	UnitAmountInCents  int             `xml:"unit_amount_in_cents,omitempty"`
	Currency           string          `xml:"currency,omitempty"`
	Quantity           string          `xml:"quantity,omitempty"`
	SubscriptionAddOns EmbedPlanAddOns `xml:"subscription_add_ons,omitempty"`
	TrialEndsAt        *time.Time      `xml:"trial_ends_at,omitempty"`
	StartsAt           *time.Time      `xml:"starts_at,omitempty"`
	FirstRenewalDate   *time.Time      `xml:"first_renewal_date,omitempty"`
	TotalBillingCycles string          `xml:"total_billing_cycles,omitempty"`
}

type SubscriptionUpdate struct {
	XMLName            xml.Name        `xml:"subscription"`
	Timeframe          string          `xml:"timeframe,omitempty"`
	PlanCode           string          `xml:"plan_code,omitempty"`
	UnitAmountInCents  int             `xml:"unit_amount_in_cents,omitempty"`
	Quantity           string          `xml:"quantity,omitempty"`
	SubscriptionAddOns EmbedPlanAddOns `xml:"subscription_add_ons,omitempty"`
}

func (s *Subscription) AttachExistingAccount(a Account) (e error) {
	if s.UUID != "" {
		return RecurlyError{statusCode: 400, Description: "Subscription Already in Use and can't attach another account to it"}
	}
	s.EmbedAccount = new(Account)
	s.EmbedAccount.AccountCode = a.AccountCode
	//s.EmbedAccount.B = a.B
	return
}
func (s *Subscription) AttachAccount(a Account) (e error) {
	if s.UUID != "" {
		return RecurlyError{statusCode: 400, Description: "Subscription Already in Use and can't attach another account to it"}
	}
	s.EmbedAccount = new(Account)
	a.CreatedAt = nil
	a.State = ""
	//some more may need to be blanked out
	a.HostedLoginToken = ""
	s.EmbedAccount = &a
	return
}

func (s *Subscription) Create() error {
	if s.UUID != "" {
		return RecurlyError{statusCode: 400, Description: "Subscription Already in Use"}
	}
	t := new(time.Time)
	decode, err := s.TrialEndsAt.GetDate()
	if err == nil {
		t = &decode
	}
	sc := subscriptionCreate{
		PlanCode:           s.PlanCode,
		Account:            s.EmbedAccount,
		Currency:           s.Currency,
		UnitAmountInCents:  s.UnitAmountInCents,
		CouponCode:         s.CouponCode,
		Quantity:           s.Quantity,
		SubscriptionAddOns: s.SubscriptionAddOns,
		StartsAt:           s.StartsAt,
		TrialEndsAt:        t,
		FirstRenewalDate:   s.FirstRenewalDate,
		TotalBillingCycles: s.TotalBillingCycles,
	}
	if err := s.r.doCreateReturn(sc, &s, s.endpoint); err == nil {
		return nil
	} else {
		return err
	}
	return nil
}

func (s *Subscription) Update(now bool) error {
	sub := SubscriptionUpdate{
		PlanCode:           s.PlanCode,
		UnitAmountInCents:  s.UnitAmountInCents,
		Quantity:           s.Quantity,
		SubscriptionAddOns: s.SubscriptionAddOns,
	}
	if len(sub.SubscriptionAddOns.PlanAddOns) <= 0 {
		sub.SubscriptionAddOns.PlanAddOns = nil
	}
	if now {
		sub.Timeframe = "now"
	} else {
		sub.Timeframe = "renewal"
	}
	return s.r.doUpdateReturn(sub, &s, s.endpoint+"/"+s.UUID)
}

func (s *Subscription) Reactivate() error {
	return s.r.doUpdateReturn(nil, &s, s.endpoint+"/"+s.UUID+"/reactivate")
}

func (s *Subscription) Cancel() error {
	return s.r.doUpdateReturn(nil, &s, s.endpoint+"/"+s.UUID+"/cancel")
}

func (s *Subscription) Postpone(renewal time.Time) error {
	return s.r.doUpdateReturn(nil, &s, s.endpoint+"/"+s.UUID+"/postpone?next_renewal_date=" + renewal.Format(time.RFC3339))
}

func (s *Subscription) terminate(refund string) error {
	return s.r.doUpdateReturn(nil, &s, s.endpoint+"/"+s.UUID+"/terminate?refund="+refund)
}

//Terminate with full refund of the last charge for the current subscription term.
func (s *Subscription) TerminateWithFullRefund() error {
	return s.terminate("full")
}

//Terminate and Prorates a refund based on the amount of time remaining in the current bill cycle.
func (s *Subscription) TerminateWithPartialRefund() error {
	return s.terminate("partial")
}

//Terminate without refund
func (s *Subscription) Terminate() error {
	return s.terminate("none")
}

func (s *Subscription) Delete() error {
	return s.r.doDelete(s.endpoint + "/" + s.UUID)
}

type EmbedPlanAddOns struct {
	PlanAddOns []EmbedPlanAddOn `xml:"subscription_add_on"`
}

//Either Insert or Update a new AddOn.  Can only have One AddOnCode in slice once.
func (e *EmbedPlanAddOns) UpdateAddOns(a *EmbedPlanAddOn) {
	if err, finder := e.GetAddOn(a.AddOnCode); err == nil {
		finder.AddOnCode = a.AddOnCode
		finder.Quantity = a.Quantity
		finder.UnitAmountInCents = a.UnitAmountInCents
	} else {
		newa := EmbedPlanAddOn{
			Quantity:          a.Quantity,
			AddOnCode:         a.AddOnCode,
			UnitAmountInCents: a.UnitAmountInCents,
		}
		e.PlanAddOns = append(e.PlanAddOns, newa)
	}
}

//Delete AddOn by AddOn Code.
func (e *EmbedPlanAddOns) DeleteAddOn(code string) {
	if len(e.PlanAddOns) > 0 {
		addons := []EmbedPlanAddOn{}
		for _, addon := range e.PlanAddOns {
			if addon.AddOnCode != code {
				addons = append(addons, addon)
			}
		}
		e.PlanAddOns = addons
	}
}

//Return a pointer to an add on for manipulation.
func (e *EmbedPlanAddOns) GetAddOn(code string) (err error, a *EmbedPlanAddOn) {
	for k, addon := range e.PlanAddOns {
		if code == addon.AddOnCode {
			a = &e.PlanAddOns[k]
			return
		}
	}
	err = errors.New("Code does not exist in current array")
	return
}

type EmbedPlanAddOn struct {
	AddOnCode         string `xml:"add_on_code,omitempty"`
	Quantity          int    `xml:"quantity,omitempty"`
	UnitAmountInCents int    `xml:"unit_amount_in_cents,omitempty"`
}
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
		v := url.Values{}
		v.Set("cursor", s.next)
		v.Set("per_page", s.perPage)
		*s, _ = s.r.GetSubscriptions(v)
	} else {
		return false
	}
	return true
}

//Get previous set of subscriptions
func (s *SubscriptionList) Prev() bool {
	if s.prev != "" {
		v := url.Values{}
		v.Set("cursor", s.prev)
		v.Set("per_page", s.perPage)
		*s, _ = s.r.GetSubscriptions(v)
	} else {
		return false
	}
	return true
}

//Go to start set of subscriptions
func (s *SubscriptionList) Start() bool {
	if s.prev != "" {
		v := url.Values{}
		v.Set("per_page", s.perPage)
		*s, _ = s.r.GetSubscriptions(v)
	} else {
		return false
	}
	return true
}

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
		v := url.Values{}
		v.Set("cursor",a.next)
		v.Set("per_page",a.perPage)
		*a,_ = a.r.GetAccountSubscriptions(a.AccountCode,v)
	} else {
		return false
	}
	return true
}

//Get previous set of subscriptions
func (a *AccountSubscriptionList) Prev() ( bool) {
	if a.prev != "" {
		v := url.Values{}
		v.Set("cursor",a.prev)
		v.Set("per_page",a.perPage)
		*a,_ = a.r.GetAccountSubscriptions(a.AccountCode,v)
	} else {
		return false
	}
	return true
}

//Go to start set of subscriptions
func (a *AccountSubscriptionList) Start() ( bool) {
	if a.prev != "" {
		v := url.Values{}
		v.Set("per_page",a.perPage)
		*a,_ = a.r.GetAccountSubscriptions(a.AccountCode,v)
	} else {
		return false
	}
	return true
}
