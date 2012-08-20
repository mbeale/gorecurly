package gorecurly

import (
	"encoding/xml"
	"time"
)

//Standar Plan Fields
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

type tempPlan struct {
	XMLName xml.Name `xml:"plan"`
	PlanFields
	SetupFeeInCents   *CurrencyMarshalArray `xml:"setup_fee_in_cents,omitempty"`
	UnitAmountInCents *CurrencyMarshalArray `xml:"unit_amount_in_cents,omitempty"`
}

//Plan Struct
type Plan struct {
	XMLName xml.Name `xml:"plan"`
	PlanFields
	SetupFeeInCents   *CurrencyArray `xml:"setup_fee_in_cents,omitempty"`
	UnitAmountInCents *CurrencyArray `xml:"unit_amount_in_cents,omitempty"`
}

//Create a plan
func (p *Plan) Create() error {
	if p.CreatedAt != nil {
		return RecurlyError{statusCode: 400, Description: "Plan Code Already in Use"}
	}
	return p.r.doCreate(&p, p.endpoint)
}

//Update a plan
func (p *Plan) Update() error {
	newplan := new(tempPlan)
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

//Delete a plan
func (p *Plan) Delete() error {
	return p.r.doDelete(p.endpoint + "/" + p.PlanCode)
}

//Plan Stub struct
type PlanStub struct {
	XMLName xml.Name `xml:"plan"`
	stub
}

//A struct to be embedded for plan_code
type PlanCode struct {
	XMLName  xml.Name `xml:"plan_codes"`
	PlanCode []string `xml:"plan_code"`
}

