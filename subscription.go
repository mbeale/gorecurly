package gorecurly

import (
	"encoding/xml"
	"errors"
	"time"
)

//Subscription Struct
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

type subscriptionUpdate struct {
	XMLName            xml.Name        `xml:"subscription"`
	Timeframe          string          `xml:"timeframe,omitempty"`
	PlanCode           string          `xml:"plan_code,omitempty"`
	UnitAmountInCents  int             `xml:"unit_amount_in_cents,omitempty"`
	Quantity           string          `xml:"quantity,omitempty"`
	SubscriptionAddOns EmbedPlanAddOns `xml:"subscription_add_ons,omitempty"`
}

//Attach an existing account to the subscription before creating it
func (s *Subscription) AttachExistingAccount(a Account) (e error) {
	if s.UUID != "" {
		return RecurlyError{statusCode: 400, Description: "Subscription Already in Use and can't attach another account to it"}
	}
	s.EmbedAccount = new(Account)
	s.EmbedAccount.AccountCode = a.AccountCode
	//s.EmbedAccount.B = a.B
	return
}

//Attach a new account object to a subscription.  The account will be created along with the subscription
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

//Create an account
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
	//Hack here need to investigate why the create causes the return to append
	s.SubscriptionAddOns = EmbedPlanAddOns{}
	if err := s.r.doCreateReturn(sc, &s, s.endpoint); err == nil {
		return nil
	} else {
		return err
	}
	return nil
}

//Update an account
func (s *Subscription) Update(now bool) error {
	sub := subscriptionUpdate{
		PlanCode:           s.PlanCode,
		UnitAmountInCents:  s.UnitAmountInCents,
		Quantity:           s.Quantity,
		SubscriptionAddOns: s.SubscriptionAddOns,
	}
	if now {
		sub.Timeframe = "now"
	} else {
		sub.Timeframe = "renewal"
	}
	//Hack here need to investigate why the update causes the return to append
	s.SubscriptionAddOns = EmbedPlanAddOns{}
	return s.r.doUpdateReturn(sub, &s, s.endpoint+"/"+s.UUID)
}

//Reactivate a cancelled account
func (s *Subscription) Reactivate() error {
	return s.r.doUpdateReturn(nil, &s, s.endpoint+"/"+s.UUID+"/reactivate")
}

//Cancel an account
func (s *Subscription) Cancel() error {
	return s.r.doUpdateReturn(nil, &s, s.endpoint+"/"+s.UUID+"/cancel")
}

//Postpone an accounts renewal datetime
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


//A struct to have embedded plan add ons
type EmbedPlanAddOns struct {
	PlanAddOns []*EmbedPlanAddOn `xml:"subscription_add_on"`
}

//Either Insert or Update a new AddOn.  Can only have One AddOnCode in slice once.
func (e *EmbedPlanAddOns) UpdateAddOns(a EmbedPlanAddOn) {
	var found bool
	for k,v := range e.PlanAddOns {
		if v.AddOnCode == a.AddOnCode {
			e.PlanAddOns[k].UnitAmountInCents = a.UnitAmountInCents
			e.PlanAddOns[k].Quantity = a.Quantity
			found = true
			break
		}
	}
	if !found {
		newa := EmbedPlanAddOn{
			Quantity:          a.Quantity,
			AddOnCode:         a.AddOnCode,
			UnitAmountInCents: a.UnitAmountInCents,
		}
		e.PlanAddOns = append(e.PlanAddOns, &newa)
	}
	/*if err, finder := e.GetAddOn(a.AddOnCode); err == nil {
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
	}*/
}

//Delete AddOn by AddOn Code.
func (e *EmbedPlanAddOns) DeleteAddOn(code string) {
	if len(e.PlanAddOns) > 0 {
		addons := []*EmbedPlanAddOn{}
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
			a = e.PlanAddOns[k]
			return
		}
	}
	err = errors.New("Code does not exist in current array")
	return
}
//An embedded plan add on
type EmbedPlanAddOn struct {
	AddOnCode         string `xml:"add_on_code,omitempty"`
	Quantity          int    `xml:"quantity,omitempty"`
	UnitAmountInCents int    `xml:"unit_amount_in_cents,omitempty"`
}
