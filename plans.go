package gorecurly

import (
	"encoding/xml"
	"errors"
	"time"
)

type PlanFields struct {
	endpoint string
	r        *Recurly
	//AddOns *AddOnsStub `xml:"add_ons,omitempty"`
	Name                     string     `xml:"name,omitempty"`
	PlanCode                 string     `xml:"plan_code,omitempty"`
	Description              string     `xml:"description,omitempty"`
	SuccessUrl               string     `xml:"success_url,omitempty"`
	CancelUrl                string     `xml:"cancel_url,omitempty"`
	DisplayDonationAmounts   bool       `xml:"display_donation_amounts,omitempty"`
	DisplayQuantity          bool       `xml:"display_quantity,omitempty"`
	DisplayPhoneNumber       bool       `xml:"display_phone_number,omitempty"`
	BypassHostedConfirmation bool       `xml:"bypass_hosted_confirmation,omitempty"`
	UnitName                 string     `xml:"unit_name,omitempty"`
	PaymentPageTOSLink       string     `xml:"payment_page_tos_link,omitempty"`
	PlanIntervalLength       int        `xml:"plan_interval_length,omitempty"`
	PlanIntervalUnit         string     `xml:"plan_interval_unit,omitempty"`
	AccountingCode           string     `xml:"accounting_code,omitempty"`
	CreatedAt                *time.Time `xml:"created_at,omitempty"`
}

type TempPlan struct {
	XMLName xml.Name `xml:"plan"`
	PlanFields
	SetupFeeInCents   *CurrencyMarshalArray `xml:"setup_fee_in_cents,omitempty"`
	UnitAmountInCents *CurrencyMarshalArray `xml:"unit_amount_in_cents,omitempty"`
}

type PlanList struct {
	Paging
	r       *Recurly
	XMLName xml.Name `xml:"plans"`
	Plans   []Plan   `xml:"plan"`
}
type Plan struct {
	XMLName xml.Name `xml:"plan"`
	PlanFields
	SetupFeeInCents   *CurrencyArray `xml:"setup_fee_in_cents,omitempty"`
	UnitAmountInCents *CurrencyArray `xml:"unit_amount_in_cents,omitempty"`
}

func (p *Plan) Create() error {
	if p.CreatedAt != nil {
		return RecurlyError{statusCode: 400, Description: "Plan Code Already in Use"}
	}
	return p.r.doCreate(&p, p.endpoint)
}

func (p *Plan) Update() error {
	newplan := new(TempPlan)
	newplan.Name = p.Name
	newplan.PlanCode = p.PlanCode
	newplan.UnitName = p.UnitName
	newplan.PlanIntervalUnit = p.PlanIntervalUnit
	newplan.CreatedAt = nil
	//Total hack job 
	//due to limitation of XML.marshal not recognizing "any" tag
	//could be fixed in future go releases
	setupFeeInCents := make([]*Currency, len(p.SetupFeeInCents.CurrencyList))
	unitAmountInCents := make([]*Currency, len(p.UnitAmountInCents.CurrencyList))
	newplan.SetupFeeInCents = &CurrencyMarshalArray{setupFeeInCents}
	newplan.UnitAmountInCents = &CurrencyMarshalArray{unitAmountInCents}
	for k, _ := range p.SetupFeeInCents.CurrencyList {
		newplan.SetupFeeInCents.CurrencyList[k] = &p.SetupFeeInCents.CurrencyList[k]
	}
	for k, _ := range p.UnitAmountInCents.CurrencyList {
		newplan.UnitAmountInCents.CurrencyList[k] = &p.UnitAmountInCents.CurrencyList[k]
	}
	//end hack job
	if len(newplan.SetupFeeInCents.CurrencyList) <= 0 {
		newplan.SetupFeeInCents = nil
	}
	if len(newplan.UnitAmountInCents.CurrencyList) <= 0 {
		newplan.UnitAmountInCents = nil
	}

	return p.r.doUpdate(newplan, p.endpoint+"/"+p.PlanCode)
}

func (p *Plan) Delete() error {
	return p.r.doDelete(p.endpoint + "/" + p.PlanCode)
}

//Account Stub struct
type PlanStub struct {
	XMLName xml.Name `xml:"plan"`
	stub
}

type PlanAddOnFields struct {
	endpoint                    string
	r                           *Recurly
	Plan                        *PlanStub  `xml:"plan,omitempty"`
	Name                        string     `xml:"name,omitempty"`
	AddOnCode                   string     `xml:"add_on_code,omitempty"`
	DisplayQuantityOnHostedPage bool       `xml:"display_quantity_on_hosted_page,omitempty"`
	DefaultQuantity             int        `xml:"default_quantity,omitempty"`
	CreatedAt                   *time.Time `xml:"created_at,omitempty"`
}

type PlanAddOn struct {
	XMLName xml.Name `xml:"add_on"`
	PlanAddOnFields
	UnitAmountInCents *CurrencyArray `xml:"unit_amount_in_cents,omitempty"`
}

type TempPlanAddOn struct {
	XMLName xml.Name `xml:"add_on"`
	PlanAddOnFields
	UnitAmountInCents *CurrencyMarshalArray `xml:"unit_amount_in_cents,omitempty"`
}

func (p *PlanAddOn) Create(plan_code string) error {
	if p.CreatedAt != nil {
		return RecurlyError{statusCode: 400, Description: "Add on Code Already in Use"}
	}
	return p.r.doCreate(&p, PLANS+"/"+plan_code+"/add_ons")
}

func (p *PlanAddOn) Update() error {
	newaddon := new(TempPlanAddOn)
	newaddon.Name = p.Name
	newaddon.DisplayQuantityOnHostedPage = p.DisplayQuantityOnHostedPage
	newaddon.DefaultQuantity = p.DefaultQuantity
	newaddon.CreatedAt = nil
	//Total hack job 
	//due to limitation of XML.marshal not recognizing "any" tag
	//could be fixed in future go releases
	unitAmountInCents := make([]*Currency, len(p.UnitAmountInCents.CurrencyList))
	newaddon.UnitAmountInCents = &CurrencyMarshalArray{unitAmountInCents}
	for k, _ := range p.UnitAmountInCents.CurrencyList {
		newaddon.UnitAmountInCents.CurrencyList[k] = &p.UnitAmountInCents.CurrencyList[k]
	}
	//end hack job
	if len(newaddon.UnitAmountInCents.CurrencyList) <= 0 {
		newaddon.UnitAmountInCents = nil
	}

	if p.Plan != nil {
		return p.r.doUpdate(newaddon, PLANS+"/"+p.Plan.GetCode()+"/add_ons/"+p.AddOnCode)
	}
	return errors.New("Plan Does not exist")
}

func (p *PlanAddOn) Delete() error {
	return p.r.doDelete(PLANS + "/" + p.Plan.GetCode() + "/add_ons/" + p.AddOnCode)
}

type PlanAddOnList struct {
	Paging
	r       *Recurly
	XMLName xml.Name    `xml:"add_ons"`
	AddOns  []PlanAddOn `xml:"add_on"`
}
