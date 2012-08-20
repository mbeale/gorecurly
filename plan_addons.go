package gorecurly

import (
	"encoding/xml"
	"errors"
	"time"
)

//Plan Add on fields struct
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

//Plan add on
type PlanAddOn struct {
	XMLName xml.Name `xml:"add_on"`
	PlanAddOnFields
	UnitAmountInCents *CurrencyArray `xml:"unit_amount_in_cents,omitempty"`
}

type tempPlanAddOn struct {
	XMLName xml.Name `xml:"add_on"`
	PlanAddOnFields
	UnitAmountInCents *CurrencyMarshalArray `xml:"unit_amount_in_cents,omitempty"`
}

//Create plan add on given a plan code
func (p *PlanAddOn) Create(plan_code string) error {
	if p.CreatedAt != nil {
		return RecurlyError{statusCode: 400, Description: "Add on Code Already in Use"}
	}
	return p.r.doCreate(&p, PLANS+"/"+plan_code+"/add_ons")
}

//Update a plan add on
func (p *PlanAddOn) Update() error {
	newaddon := new(tempPlanAddOn)
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

//Delete plan add on
func (p *PlanAddOn) Delete() error {
	return p.r.doDelete(PLANS + "/" + p.Plan.GetCode() + "/add_ons/" + p.AddOnCode)
}


