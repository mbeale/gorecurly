//Recurly Client Library for go
//Works with version 2 of API only

package gorecurly

//TODO: Do all tests, do with mock server, need to fix URL so that it can be overriden
//TODO: Check all comments when finished 
//TODO: Change all paging to new request params 
//TODO: Check that state is working with lists
//TODO: Introduce stubs for all resources
//TODO: PDF Invoice
//TODO: Recurly.js signing
//TODO: transparent post (probably not)
//TODO: Double check fields and make sure no new fields were added
//TODO: Option to add no auth to header "Recurly-Skip-Authorization: true"
//TODO: Maybe some examples fetching with goroutines
//TODO: Add a variable to test if subscription is in trial
//TODO: Custom function to calculate account balance
//TODO: Custom function to calculate next billing amt
//TODO: Discount in cents for coupons not working
//TODO: Handle push notifications

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	URL               = "https://api.recurly.com/v2/"
	libversion        = "0.4"
	libname           = "Recurly-Go"
	ACCOUNTS          = "accounts"
	ADJUSTMENTS       = "adjustments"
	BILLINGINFO       = "billing_info"
	COUPONS           = "coupons"
	COUPONREDEMPTIONS = "redemption"
	INVOICES          = "invoices"
	PLANS             = "plans"
	PLANADDONS        = "add_ons"
	SUBSCRIPTIONS     = "subscriptions"
	TRANSACTIONS      = "transactions"
)

//Generic Reader
type nopCloser struct {
	io.Reader
}

//functions

//Initialize the Recurly package with your apikey and your jskey
func InitRecurly(apikey string, jskey string) *Recurly {
	r := new(Recurly)
	r.apiKey = apikey
	r.JSKey = jskey
	r.url = URL
	return r
}

//interfaces

//Paging interface to allow Next,Prev,Start
type Pager interface {
	getRawBody() []byte
}

//Recurly errors
var Error400 = errors.New("The request was invalid or could not be understood by the server. Resubmitting the request will likely result in the same error.")
var Error401 = errors.New("Your API key is missing or invalid.")
var Error402 = errors.New("Your Recurly account is in production mode but is not in good standing. Please pay any outstanding invoices.")
var Error403 = errors.New("The login is attempting to perform an action it does not have privileges to access. Verify your login credentials are for the appropriate account.")
var Error404 = errors.New("The resource was not found with the given identifier. The response body will explain which resource was not found.")
var Error405 = errors.New("The requested method is not valid at the given URL.")
var Error406 = errors.New("The request's Accept header is not set to application/xml")
var Error412 = errors.New("The request was unsuccessful because a condition was not met. For example, this message may be returned if you attempt to cancel a subscription for an account that has no subscription.")
var Error429 = errors.New("You have made too many API requests in the last hour. Future API requests will be ignored until the beginning of the next hour.")

//Recurly Generic Errors
type RecurlyError struct {
	XMLName     xml.Name `xml:"error"`
	statusCode  int
	Symbol      string `xml:"symbol"`
	Description string `xml:"description"`
	Details     string `xml:"details"`
}

//Recurly Validation Errors Array
type RecurlyValidationErrors struct {
	XMLName    xml.Name `xml:"errors"`
	statusCode int
	Errors     []RecurlyValidationError `xml:"error"`
}

//Recurly validation error
type RecurlyValidationError struct {
	XMLName     xml.Name `xml:"error"`
	FieldName   string   `xml:"field,attr"`
	Symbol      string   `xml:"symbol,attr"`
	Description string   `xml:",innerxml"`
}

//Parse Recurly XML to create a Recurly Error
func CreateRecurlyStandardError(resp *http.Response) (r RecurlyError) {
	r.statusCode = resp.StatusCode
	if xmlstring, readerr := ioutil.ReadAll(resp.Body); readerr == nil {
		if xmlerr := xml.Unmarshal(xmlstring, &r); xmlerr != nil {
			r.Description = string(xmlstring)
		}
	}
	return r
}

//Parse Recurly XML to create a Validation Error
func CreateRecurlyValidationError(resp *http.Response) (r RecurlyValidationErrors) {
	r.statusCode = resp.StatusCode
	if xmlstring, readerr := ioutil.ReadAll(resp.Body); readerr == nil {
		if xmlerr := xml.Unmarshal(xmlstring, &r); xmlerr != nil {
			//r.Description = xmlerr.Error()
			println(xmlerr.Error())
		}
	}
	return r
}

//Filter to decide which error type to create
func createRecurlyError(resp *http.Response) error {
	switch resp.StatusCode {
	case 400:
		return Error400
	case 401:
		return Error401
	case 402:
		return Error402
	case 403:
		return Error403
	case 404:
		return Error404
	case 405:
		return Error405
	case 406:
		return Error406
	case 412:
		return Error412
	case 429:
		return Error429
	case 422:
		return CreateRecurlyValidationError(resp)
	}
	return CreateRecurlyStandardError(resp)
}

//Formatted General Error 
func (r RecurlyError) Error() string {
	return fmt.Sprintf("Recurly Error: %s , %s %s Status Code: %v", r.Symbol, r.Description, r.Details, r.statusCode)
}

//Formatted Validation Error
func (r RecurlyValidationErrors) Error() string {
	var rtnString string
	for _, v := range r.Errors {
		rtnString += v.FieldName + " " + v.Description + "\n"
	}
	return fmt.Sprintf("You have the following validation errors:\n%s", rtnString)
}

//Main Recurly Client
type Recurly struct {
	apiKey, JSKey, url string
	debug              bool
}

//Set verbose debugging
func (r *Recurly) EnableDebug() {
	r.debug = true
}

//Get a list of accounts
func (r *Recurly) GetAccounts(params ...url.Values) (AccountList, error) {
	accountlist := AccountList{}
	sendvars := accountlist.initParams(params)
	if err := accountlist.initList(ACCOUNTS, sendvars, r); err == nil {
		if xmlerr := xml.Unmarshal(accountlist.getRawBody(), &accountlist); xmlerr == nil {
			for k, _ := range accountlist.Account {
				accountlist.Account[k].r = r
				accountlist.Account[k].endpoint = ACCOUNTS
			}
			accountlist.r = r
			return accountlist, nil
		} else {
			if r.debug {
				println(xmlerr.Error())
			}
			return accountlist, xmlerr
		}
	} else {
		return accountlist, err
	}
	return accountlist, nil
}

//Get a list of adjustments for an account_code
func (r *Recurly) GetAdjustments(account_code string, params ...url.Values) (AdjustmentList, error) {
	adjlist := AdjustmentList{}
	sendvars := adjlist.initParams(params)
	if err := adjlist.initList(ACCOUNTS+"/"+account_code+"/"+ADJUSTMENTS, sendvars, r); err == nil {
		if xmlerr := xml.Unmarshal(adjlist.getRawBody(), &adjlist); xmlerr == nil {
			for k, _ := range adjlist.Adjustments {
				adjlist.Adjustments[k].r = r
				adjlist.Adjustments[k].endpoint = ADJUSTMENTS
			}
			adjlist.r = r
			adjlist.AccountCode = account_code
			return adjlist, nil
		} else {
			if r.debug {
				println(xmlerr.Error())
			}
			return adjlist, xmlerr
		}
	} else {
		return adjlist, err
	}
	return adjlist, nil
}

//Get a list of coupons
func (r *Recurly) GetCoupons(params ...url.Values) (CouponList, error) {
	cplist := CouponList{}
	sendvars := cplist.initParams(params)
	if err := cplist.initList(COUPONS, sendvars, r); err == nil {
		if xmlerr := xml.Unmarshal(cplist.getRawBody(), &cplist); xmlerr == nil {
			for k, _ := range cplist.Coupons {
				cplist.Coupons[k].r = r
				cplist.Coupons[k].endpoint = COUPONS
			}
			cplist.r = r
			return cplist, nil
		} else {
			if r.debug {
				println(xmlerr.Error())
			}
			return cplist, xmlerr
		}
	} else {
		return cplist, err
	}
	return cplist, nil
}

//Get a list of invoices for an account_code
func (r *Recurly) GetAccountInvoices(account_code string, params ...url.Values) (AccountInvoiceList, error) {
	invoicelist := AccountInvoiceList{}
	sendvars := invoicelist.initParams(params)
	if err := invoicelist.initList(ACCOUNTS+"/"+account_code+"/"+INVOICES, sendvars, r); err == nil {
		if xmlerr := xml.Unmarshal(invoicelist.getRawBody(), &invoicelist); xmlerr == nil {
			for k, _ := range invoicelist.Invoices {
				invoicelist.Invoices[k].r = r
				invoicelist.Invoices[k].endpoint = INVOICES
			}
			invoicelist.r = r
			invoicelist.AccountCode = account_code
			return invoicelist, nil
		} else {
			if r.debug {
				println(xmlerr.Error())
			}
			return invoicelist, xmlerr
		}
	} else {
		return invoicelist, err
	}
	return invoicelist, nil
}

//Get a list of invoices
func (r *Recurly) GetInvoices(params ...url.Values) (InvoiceList, error) {
	invoicelist := InvoiceList{}
	sendvars := invoicelist.initParams(params)
	if err := invoicelist.initList(INVOICES, sendvars, r); err == nil {
		if xmlerr := xml.Unmarshal(invoicelist.getRawBody(), &invoicelist); xmlerr == nil {
			for k, _ := range invoicelist.Invoices {
				invoicelist.Invoices[k].r = r
				invoicelist.Invoices[k].endpoint = INVOICES
			}
			invoicelist.r = r
			return invoicelist, nil
		} else {
			if r.debug {
				println(xmlerr.Error())
			}
			return invoicelist, xmlerr
		}
	} else {
		return invoicelist, err
	}
	return invoicelist, nil
}

//Get a list of Plans
func (r *Recurly) GetPlans(params ...url.Values) (PlanList, error) {
	planlist := PlanList{}
	sendvars := planlist.initParams(params)
	if err := planlist.initList(PLANS, sendvars, r); err == nil {
		if xmlerr := xml.Unmarshal(planlist.getRawBody(), &planlist); xmlerr == nil {
			for k, _ := range planlist.Plans {
				planlist.Plans[k].r = r
				planlist.Plans[k].endpoint = PLANS
			}
			planlist.r = r
			return planlist, nil
		} else {
			if r.debug {
				println(xmlerr.Error())
			}
			return planlist, xmlerr
		}
	} else {
		return planlist, err
	}
	return planlist, nil
}

//Get a list of add ons for a plan_code
func (r *Recurly) GetPlanAddOns(plan_code string, params ...url.Values) (planaddonlist PlanAddOnList, e error) {
	sendvars := planaddonlist.initParams(params)
	if err := planaddonlist.initList(PLANS+"/"+plan_code+"/add_ons", sendvars, r); err == nil {
		if xmlerr := xml.Unmarshal(planaddonlist.getRawBody(), &planaddonlist); xmlerr == nil {
			for k, _ := range planaddonlist.AddOns {
				planaddonlist.AddOns[k].r = r
			}
			planaddonlist.r = r
			planaddonlist.PlanCode = plan_code
			return
		} else {
			if r.debug {
				println(xmlerr.Error())
			}
			return planaddonlist, xmlerr
		}
	} else {
		return planaddonlist, err
	}
	return
}

//Get a list of subscriptions
func (r *Recurly) GetSubscriptions(params ...url.Values) (SubscriptionList, error) {
	subs := SubscriptionList{}
	sendvars := subs.initParams(params)
	if err := subs.initList(SUBSCRIPTIONS, sendvars, r); err == nil {
		if xmlerr := xml.Unmarshal(subs.getRawBody(), &subs); xmlerr == nil {
			for k, _ := range subs.Subscriptions {
				subs.Subscriptions[k].r = r
				subs.Subscriptions[k].endpoint = SUBSCRIPTIONS
			}
			subs.r = r
			return subs, nil
		} else {
			if r.debug {
				println(xmlerr.Error())
			}
			return subs, xmlerr
		}
	} else {
		return subs, err
	}
	return subs, nil
}

//Get a list of subscriptions for an account_code
func (r *Recurly) GetAccountSubscriptions(account_code string, params ...url.Values) (AccountSubscriptionList, error) {
	subs := AccountSubscriptionList{}
	sendvars := subs.initParams(params)
	if err := subs.initList(ACCOUNTS+"/"+account_code+"/"+SUBSCRIPTIONS, sendvars, r); err == nil {
		if xmlerr := xml.Unmarshal(subs.getRawBody(), &subs); xmlerr == nil {
			for k, _ := range subs.Subscriptions {
				subs.Subscriptions[k].r = r
				subs.Subscriptions[k].endpoint = SUBSCRIPTIONS
			}
			subs.r = r
			subs.AccountCode = account_code
			return subs, nil
		} else {
			if r.debug {
				println(xmlerr.Error())
			}
			return subs, xmlerr
		}
	} else {
		return subs, err
	}
	return subs, nil
}

//Get a list of transactions
func (r *Recurly) GetTransactions(params ...url.Values) (TransactionList, error) {
	subs := TransactionList{}
	sendvars := subs.initParams(params)
	if err := subs.initList(TRANSACTIONS, sendvars, r); err == nil {
		if xmlerr := xml.Unmarshal(subs.getRawBody(), &subs); xmlerr == nil {
			for k, _ := range subs.Transactions {
				subs.Transactions[k].r = r
				subs.Transactions[k].endpoint = TRANSACTIONS
			}
			subs.r = r
			return subs, nil
		} else {
			if r.debug {
				println(xmlerr.Error())
			}
			return subs, xmlerr
		}
	} else {
		return subs, err
	}
	return subs, nil
}

//Get a list of transactions for an account_code
func (r *Recurly) GetAccountTransactions(account_code string, params ...url.Values) (AccountTransactionList, error) {
	subs := AccountTransactionList{}
	sendvars := subs.initParams(params)
	if err := subs.initList(ACCOUNTS+"/"+account_code+"/"+TRANSACTIONS, sendvars, r); err == nil {
		if xmlerr := xml.Unmarshal(subs.getRawBody(), &subs); xmlerr == nil {
			for k, _ := range subs.Transactions {
				subs.Transactions[k].r = r
				subs.Transactions[k].endpoint = TRANSACTIONS
			}
			subs.r = r
			subs.AccountCode = account_code
			return subs, nil
		} else {
			if r.debug {
				println(xmlerr.Error())
			}
			return subs, xmlerr
		}
	} else {
		return subs, err
	}
	return subs, nil
}

//Get a single account by account_code
func (r *Recurly) GetAccount(account_code string) (account Account, err error) {
	account = r.NewAccount()
	if resp, err := r.createRequest(ACCOUNTS+"/"+account_code, "GET", nil, nil); err == nil {
		if resp.StatusCode == 200 {
			if body, readerr := ioutil.ReadAll(resp.Body); readerr == nil {
				if r.debug {
					println(resp.Status)
					for k, _ := range resp.Header {
						println(k + ":" + resp.Header[k][0])
					}
					fmt.Printf("%s\n", body)
					fmt.Printf("Content-Length:%v\n", resp.ContentLength)
				}
				//load object xml
				if xmlerr := xml.Unmarshal(body, &account); xmlerr != nil {
					account.B = nil
					return account, xmlerr
				}
				//everything went fine
				return account, nil
			} else {
				//return read error
				return account, readerr
			}
			return account, nil
		} else {
			return account, createRecurlyError(resp)
		}
	} else {
		return account, err
	}
	return account, nil
}

//Get a single adjustment by uuid
func (r *Recurly) GetAdjustment(uuid string) (adj Adjustment, err error) {
	adj = r.NewAdjustment()
	if resp, err := r.createRequest(ADJUSTMENTS+"/"+uuid, "GET", nil, nil); err == nil {
		if resp.StatusCode == 200 {
			if body, readerr := ioutil.ReadAll(resp.Body); readerr == nil {
				if r.debug {
					println(resp.Status)
					for k, _ := range resp.Header {
						println(k + ":" + resp.Header[k][0])
					}
					fmt.Printf("%s\n", body)
					fmt.Printf("Content-Length:%v\n", resp.ContentLength)
				}
				//load object xml
				if xmlerr := xml.Unmarshal(body, &adj); xmlerr != nil {
					return adj, xmlerr
				}
				//everything went fine
				return adj, nil
			} else {
				//return read error
				return adj, readerr
			}
			return adj, nil
		} else {
			return adj, createRecurlyError(resp)
		}
	} else {
		return adj, err
	}
	return adj, nil
}

//Get a single coupon redemption by account_code
func (r *Recurly) GetCouponRedemption(account_code string) (red Redemption, err error) {
	red.r = r
	red.AccountCode = account_code
	if resp, err := r.createRequest(ACCOUNTS+"/"+account_code+"/redemption", "GET", nil, nil); err == nil {
		if resp.StatusCode == 200 {
			if body, readerr := ioutil.ReadAll(resp.Body); readerr == nil {
				if r.debug {
					println(resp.Status)
					for k, _ := range resp.Header {
						println(k + ":" + resp.Header[k][0])
					}
					fmt.Printf("%s\n", body)
					fmt.Printf("Content-Length:%v\n", resp.ContentLength)
				}
				//load object xml
				if xmlerr := xml.Unmarshal(body, &red); xmlerr != nil {
					return red, xmlerr
				}
				//everything went fine
				return red, nil
			} else {
				//return read error
				return red, readerr
			}
			return red, nil
		} else {
			return red, createRecurlyError(resp)
		}
	} else {
		return red, err
	}
	return red, nil
}

//Get a single coupon by uuid
func (r *Recurly) GetCoupon(uuid string) (coupon Coupon, err error) {
	coupon = r.NewCoupon()
	if resp, err := r.createRequest(COUPONS+"/"+uuid, "GET", nil, nil); err == nil {
		if resp.StatusCode == 200 {
			if body, readerr := ioutil.ReadAll(resp.Body); readerr == nil {
				if r.debug {
					println(resp.Status)
					for k, _ := range resp.Header {
						println(k + ":" + resp.Header[k][0])
					}
					fmt.Printf("%s\n", body)
					fmt.Printf("Content-Length:%v\n", resp.ContentLength)
				}
				//load object xml
				if xmlerr := xml.Unmarshal(body, &coupon); xmlerr != nil {
					return coupon, xmlerr
				}
				//everything went fine
				return coupon, nil
			} else {
				//return read error
				return coupon, readerr
			}
			return coupon, nil
		} else {
			return coupon, createRecurlyError(resp)
		}
	} else {
		return coupon, err
	}
	return coupon, nil
}

//Get invoice by uuid
func (r *Recurly) GetInvoice(uuid string) (invoice Invoice, err error) {
	invoice = r.NewInvoice()
	if resp, err := r.createRequest(INVOICES+"/"+uuid, "GET", nil, nil); err == nil {
		if resp.StatusCode == 200 {
			if body, readerr := ioutil.ReadAll(resp.Body); readerr == nil {
				if r.debug {
					println(resp.Status)
					for k, _ := range resp.Header {
						println(k + ":" + resp.Header[k][0])
					}
					fmt.Printf("%s\n", body)
					fmt.Printf("Content-Length:%v\n", resp.ContentLength)
				}
				//load object xml
				if xmlerr := xml.Unmarshal(body, &invoice); xmlerr != nil {
					return invoice, xmlerr
				}
				//everything went fine
				return invoice, nil
			} else {
				//return read error
				return invoice, readerr
			}
			return invoice, nil
		} else {
			return invoice, createRecurlyError(resp)
		}
	} else {
		return invoice, err
	}
	return invoice, nil
}

//Get a single plan by plan_code
func (r *Recurly) GetPlan(plan_code string) (plan Plan, err error) {
	plan = r.NewPlan()
	if resp, err := r.createRequest(PLANS+"/"+plan_code, "GET", nil, nil); err == nil {
		if resp.StatusCode == 200 {
			if body, readerr := ioutil.ReadAll(resp.Body); readerr == nil {
				if r.debug {
					println(resp.Status)
					for k, _ := range resp.Header {
						println(k + ":" + resp.Header[k][0])
					}
					fmt.Printf("%s\n", body)
					fmt.Printf("Content-Length:%v\n", resp.ContentLength)
				}
				//load object xml
				if xmlerr := xml.Unmarshal(body, &plan); xmlerr != nil {
					return plan, xmlerr
				}
				//everything went fine
				return plan, err
			} else {
				//return read error
				return plan, readerr
			}
			return plan, nil
		} else {
			return plan, createRecurlyError(resp)
		}
	} else {
		return plan, err
	}
	return plan, nil
}

//Get a single plan add on by plan_code and add_on_code
func (r *Recurly) GetPlanAddOn(plan_code, add_on_code string) (plan PlanAddOn, err error) {
	plan = r.NewPlanAddOn()
	if resp, err := r.createRequest(PLANS+"/"+plan_code+"/add_ons/"+add_on_code, "GET", nil, nil); err == nil {
		if resp.StatusCode == 200 {
			if body, readerr := ioutil.ReadAll(resp.Body); readerr == nil {
				if r.debug {
					println(resp.Status)
					for k, _ := range resp.Header {
						println(k + ":" + resp.Header[k][0])
					}
					fmt.Printf("%s\n", body)
					fmt.Printf("Content-Length:%v\n", resp.ContentLength)
				}
				//load object xml
				if xmlerr := xml.Unmarshal(body, &plan); xmlerr != nil {
					return plan, xmlerr
				}
				//everything went fine
				return plan, err
			} else {
				//return read error
				return plan, readerr
			}
			return plan, nil
		} else {
			return plan, createRecurlyError(resp)
		}
	} else {
		return plan, err
	}
	return plan, nil
}

//Get a single subscription by uuid
func (r *Recurly) GetSubscription(uuid string) (sub Subscription, err error) {
	sub = r.NewSubscription()
	if resp, err := r.createRequest(SUBSCRIPTIONS+"/"+uuid, "GET", nil, nil); err == nil {
		if resp.StatusCode == 200 {
			if body, readerr := ioutil.ReadAll(resp.Body); readerr == nil {
				if r.debug {
					println(resp.Status)
					for k, _ := range resp.Header {
						println(k + ":" + resp.Header[k][0])
					}
					fmt.Printf("%s\n", body)
					fmt.Printf("Content-Length:%v\n", resp.ContentLength)
				}
				//load object xml
				if xmlerr := xml.Unmarshal(body, &sub); xmlerr != nil {
					return sub, xmlerr
				}
				//everything went fine
				return sub, nil
			} else {
				//return read error
				return sub, readerr
			}
			return sub, nil
		} else {
			return sub, createRecurlyError(resp)
		}
	} else {
		return sub, err
	}
	return sub, nil
}

//Get a single transaction by uuid
func (r *Recurly) GetTransaction(uuid string) (tran Transaction, err error) {
	tran = r.NewTransaction()
	if resp, err := r.createRequest(TRANSACTIONS+"/"+uuid, "GET", nil, nil); err == nil {
		if resp.StatusCode == 200 {
			if body, readerr := ioutil.ReadAll(resp.Body); readerr == nil {
				if r.debug {
					println(resp.Status)
					for k, _ := range resp.Header {
						println(k + ":" + resp.Header[k][0])
					}
					fmt.Printf("%s\n", body)
					fmt.Printf("Content-Length:%v\n", resp.ContentLength)
				}
				//load object xml
				if xmlerr := xml.Unmarshal(body, &tran); xmlerr != nil {
					return tran, xmlerr
				}
				//everything went fine
				return tran, nil
			} else {
				//return read error
				return tran, readerr
			}
			return tran, nil
		} else {
			return tran, createRecurlyError(resp)
		}
	} else {
		return tran, err
	}
	return tran, nil
}

//Create a new Account
func (r *Recurly) NewAccount() (account Account) {
	account.r = r
	account.endpoint = ACCOUNTS
	return
}

//Create a new Adjustment
func (r *Recurly) NewAdjustment() (adj Adjustment) {
	adj.r = r
	adj.endpoint = ADJUSTMENTS
	return
}

//Create new Billing Info
func (r *Recurly) NewBillingInfo() (bi BillingInfo) {
	bi.r = r
	bi.endpoint = BILLINGINFO
	return
}

//Create a new Coupon
func (r *Recurly) NewCoupon() (c Coupon) {
	c.r = r
	c.endpoint = COUPONS
	return
}

//Create a new Plan
func (r *Recurly) NewPlan() (plan Plan) {
	plan.r = r
	plan.SetupFeeInCents = new(CurrencyArray)
	plan.UnitAmountInCents = new(CurrencyArray)
	plan.endpoint = PLANS
	return
}

//Create a new plan Add On
func (r *Recurly) NewPlanAddOn() (planAddOn PlanAddOn) {
	planAddOn.r = r
	planAddOn.UnitAmountInCents = new(CurrencyArray)
	return
}

//Create a new Subscription
func (r *Recurly) NewSubscription() (subscription Subscription) {
	subscription.r = r
	subscription.endpoint = SUBSCRIPTIONS
	return
}

//Create a new transaction
func (r *Recurly) NewTransaction() (transaction Transaction) {
	transaction.r = r
	transaction.endpoint = TRANSACTIONS
	return
}

//Invoice Pending Charges on an account
func (r *Recurly) InvoicePendingCharges(account_code string) (invoice Invoice, e error) {
	invoice.r = r
	e = invoice.r.doCreate(&invoice, ACCOUNTS+"/"+account_code+"/invoices")
	return
}

//Create a new Invoice
func (r *Recurly) NewInvoice() (invoice Invoice) {
	invoice.r = r
	invoice.endpoint = INVOICES
	return
}

//Get a single accounts billing info by account_code
func (r *Recurly) GetBillingInfo(account_code string) (bi BillingInfo, err error) {
	bi = r.NewBillingInfo()
	if resp, err := r.createRequest(ACCOUNTS+"/"+account_code+"/"+BILLINGINFO, "GET", nil, nil); err == nil {
		if resp.StatusCode == 200 {
			if body, readerr := ioutil.ReadAll(resp.Body); readerr == nil {
				if r.debug {
					println(resp.Status)
					for k, _ := range resp.Header {
						println(k + ":" + resp.Header[k][0])
					}
					fmt.Printf("%s\n", body)
					fmt.Printf("Content-Length:%v\n", resp.ContentLength)
				}
				//load object xml
				if xmlerr := xml.Unmarshal(body, &bi); xmlerr != nil {
					return bi, xmlerr
				}
				//everything went fine
				bi.Account.endpoint = ACCOUNTS
				return bi, nil
			} else {
				//return read error
				return bi, readerr
			}
			return bi, nil
		} else {
			return bi, createRecurlyError(resp)
		}
	} else {
		return bi, err
	}
	return bi, nil
}

//Create a request to Recurly and return that response object
func (r *Recurly) createRequest(endpoint string, method string, params url.Values, msgbody []byte) (*http.Response, error) {
	client := &http.Client{}
	u, err := url.Parse(r.url + endpoint)
	if err != nil {
		return nil, err
	}
	u.RawQuery = u.RawQuery + params.Encode()
	body := nopCloser{bytes.NewBufferString(string(msgbody))}
	if r.debug {
		fmt.Printf("Endpoint Requested: %s Method: %s Body: %s\n", u.String(), method, string(msgbody))
	}
	if req, err := http.NewRequest(method, u.String(), body); err != nil {
		return nil, err
	} else {
		req.Header.Add("Accept", "application/xml")
		req.Header.Add("Accept-Language", "en-US")
		req.Header.Add("User-Agent", libname+" version="+libversion)
		req.Header.Add("Content-Type", "application/xml; charset=utf-8")
		req.ContentLength = int64(len(string(msgbody)))
		req.SetBasicAuth(r.apiKey, "")
		if resp, resperr := client.Do(req); resperr == nil {
			return resp, nil
		} else {
			return nil, resperr
		}
	}
	return nil, nil
}

//process create request and return the updated interface
func (r *Recurly) doCreateReturn(v, ret interface{}, endpoint string) (e error) {
	if xmlstring, err := xml.MarshalIndent(v, "", "    "); err == nil {
		xmlstring = []byte(xml.Header + string(xmlstring))
		if r.debug {
			fmt.Printf("%s\n", xmlstring)
		}
		if resp, reqerr := r.createRequest(endpoint, "POST", nil, xmlstring); reqerr == nil {
			if resp.StatusCode < 400 {
				if body, readerr := ioutil.ReadAll(resp.Body); readerr == nil {
					if r.debug {
						println(resp.Status)
						for k, _ := range resp.Header {
							println(k + ":" + resp.Header[k][0])
						}
						fmt.Printf("%s\n", body)
						fmt.Printf("Content-Length:%v\n", resp.ContentLength)
					}
					//load object xml
					if xmlerr := xml.Unmarshal(body, ret); xmlerr != nil {
						return xmlerr
					}
					//everything went fine
					return nil
				} else {
					//return read error
					return readerr
				}
				return nil
			} else {
				return createRecurlyError(resp)
			}
		} else {
			return reqerr
		}
	} else {
		return err
	}
	return nil
}

//Create a resource from struct, uses POST method
func (r *Recurly) doCreate(v interface{}, endpoint string) error {
	if xmlstring, err := xml.MarshalIndent(v, "", "    "); err == nil {
		xmlstring = []byte(xml.Header + string(xmlstring))
		if r.debug {
			fmt.Printf("%s\n", xmlstring)
		}
		if resp, reqerr := r.createRequest(endpoint, "POST", nil, xmlstring); reqerr == nil {
			if resp.StatusCode < 400 {
				if body, readerr := ioutil.ReadAll(resp.Body); readerr == nil {
					if r.debug {
						println(resp.Status)
						for k, _ := range resp.Header {
							println(k + ":" + resp.Header[k][0])
						}
						fmt.Printf("%s\n", body)
						fmt.Printf("Content-Length:%v\n", resp.ContentLength)
					}
					//load object xml
					if xmlerr := xml.Unmarshal(body, v); xmlerr != nil {
						return xmlerr
					}
					//everything went fine
					return nil
				} else {
					//return read error
					return readerr
				}
				return nil
			} else {
				return createRecurlyError(resp)
			}
		} else {
			return reqerr
		}
	} else {
		return err
	}
	return nil
}

//Update a resource from Struct, then return the updated object uses PUT method
func (r *Recurly) doUpdateReturn(v, ret interface{}, endpoint string) error {
	if xmlstring, err := xml.MarshalIndent(v, "", "    "); err == nil {
		if v != nil {
			xmlstring = []byte(xml.Header + string(xmlstring))
		}
		if r.debug {
			fmt.Printf("%s\n", xmlstring)
		}
		if resp, reqerr := r.createRequest(endpoint, "PUT", nil, xmlstring); reqerr == nil {
			if resp.StatusCode < 400 {
				if body, readerr := ioutil.ReadAll(resp.Body); readerr == nil {
					if r.debug {
						println(resp.Status)
						for k, _ := range resp.Header {
							println(k + ":" + resp.Header[k][0])
						}
						fmt.Printf("%s\n", body)
						fmt.Printf("Content-Length:%v\n", resp.ContentLength)
					}
					//load object xml
					if xmlerr := xml.Unmarshal(body, ret); xmlerr != nil {
						return xmlerr
					}
					//everything went fine
					return nil
				} else {
					//return read error
					return readerr
				}
				return nil
			} else {
				return createRecurlyError(resp)
			}
		} else {
			return reqerr
		}
	} else {
		return err
	}
	return nil
}

//Update a resource from Struct, uses PUT method
func (r *Recurly) doUpdate(v interface{}, endpoint string) error {
	if xmlstring, err := xml.MarshalIndent(v, "", "    "); err == nil {
		xmlstring = []byte(xml.Header + string(xmlstring))
		if r.debug {
			fmt.Printf("%s\n", xmlstring)
		}
		if resp, reqerr := r.createRequest(endpoint, "PUT", nil, xmlstring); reqerr == nil {
			if resp.StatusCode < 400 {
				return nil
			} else {
				return createRecurlyError(resp)
			}
		} else {
			return reqerr
		}
	} else {
		return err
	}
	return nil
}

//Delete a resource, uses DELETE method
func (r *Recurly) doDelete(endpoint string) error {
	if resp, reqerr := r.createRequest(endpoint, "DELETE", nil, nil); reqerr == nil {
		if resp.StatusCode < 400 {
			return nil
		} else {
			return createRecurlyError(resp)
		}
	} else {
		return reqerr
	}
	return nil
}

/* paging struct */

//A struct to assist in paging result sets
type Paging struct {
	rawBody                    []byte
	count, next, prev, perPage string
	UrlVars url.Values
}

//Return the rawBody Var
func (p Paging) getRawBody() []byte {
	return p.rawBody
}
func (p *Paging) initParams(params []url.Values) (sendvars url.Values){
	if params != nil {
		sendvars = params[0]
		if sendvars.Get("per_page") == "" {
			p.perPage = "50"
		} else {
			p.perPage = sendvars.Get("per_page")
		}
	}
	return
}

//Return params for next request
func (p *Paging) NextParams() url.Values{
	v := p.UrlVars
	v.Set("cursor", p.next)
	v.Set("per_page", p.perPage)
	return v
}

//Return params for next request
func (p *Paging) StartParams() url.Values{
	v := p.UrlVars
	v.Del("cursor")
	v.Set("per_page", p.perPage)
	return v
}

//Return params for next request
func (p *Paging) PrevParams() url.Values{
	v := p.UrlVars
	v.Set("cursor", p.prev)
	v.Set("per_page", p.perPage)
	return v
}

//Set header data for paging
func (p *Paging) SetData(rb []byte, count string, links string, params url.Values) {
	p.rawBody = rb
	p.count = count
	p.next = ""
	p.prev = ""
	p.UrlVars = params
	for _, v := range strings.SplitN(links, ",", -1) {
		link := strings.SplitN(v, ";", -1)
		link[0] = strings.Replace(link[0], "<", "", -1)
		link[0] = strings.Replace(link[0], ">", "", -1)
		if u, err := url.Parse(link[0]); err == nil {
			values := u.Query()
			switch link[1] {
			case " rel=\"next\"":
				p.next = values.Get("cursor")
			case " rel=\"prev\"":
				p.prev = values.Get("cursor")
			}
		}
	}
}

//Initialize the paging list values
func (p *Paging) initList(endpoint string, params url.Values, r *Recurly) error {
	if resp, err := r.createRequest(endpoint, "GET", params, make([]byte, 0)); err == nil {
		if resp.StatusCode < 400 {
			defer resp.Body.Close()
			if body, readerr := ioutil.ReadAll(resp.Body); readerr == nil {
				if r.debug {
					println(resp.Status)
					for k, _ := range resp.Header {
						println(k + ":" + resp.Header[k][0])
					}
					fmt.Printf("%s\n", body)
					fmt.Printf("Content-Length:%v\n", resp.ContentLength)
				}
				if x := len(resp.Header["Link"]); x > 0 {
					p.SetData(body, resp.Header["X-Records"][0], resp.Header["Link"][0],params)
				} else {
					p.SetData(body, resp.Header["X-Records"][0], "",params)
				}
				//everything went fine
				return nil
			} else {
				//return read error
				return readerr
			}
		} else {
			return createRecurlyError(resp)
		}
	} else {
		//return error message
		return err
	}
	return nil
}

/*resource objects */

//A struct to help with marshalling currency
type CurrencyMarshalArray struct {
	CurrencyList []*Currency `xml:""`
}

//A struct to help with marshalling currency
type CurrencyArray struct {
	CurrencyList []Currency `xml:"unit_amount_in_cents"`
}

//This helps you set an amount for a 3-Digit Currency.
//This function will update the amount if the currency already
//Exists in the CurrencyArray
func (c *CurrencyArray) SetCurrency(currency string, amount int) {
	if k := c.findCurrency(currency); k >= 0 {
		//update instead of insert
		c.CurrencyList[k].Amount = fmt.Sprintf("%v", amount)
	} else {
		newc := Currency{Amount: fmt.Sprintf("%v", amount)}
		newc.XMLName.Local = currency
		c.CurrencyList = append(c.CurrencyList, newc)
	}
}

//find a currency
func (c *CurrencyArray) findCurrency(currency string) (key int) {
	if c == nil {
		return -1
	}
	for k, v := range c.CurrencyList {
		if v.XMLName.Local == currency {
			return k
		}
	}
	return -1
}

//Given a 3-Digit currency, return the value or an error if currency not found
func (c *CurrencyArray) GetCurrency(currency string) (value int, e error) {
	if k := c.findCurrency(currency); k >= 0 {
		value, e = strconv.Atoi(c.CurrencyList[k].Amount)
		return
	}
	e = errors.New(fmt.Sprintf("%s not found", currency))
	return
}

//A Currency Struct
type Currency struct {
	XMLName xml.Name `xml:""`
	Amount  string   `xml:",chardata"`
}

/* Stub */
//A stub place holder
type stub struct {
	HREF     string `xml:"href,attr"`
	endpoint string `xml:",-"`
}

//Get the code/uuid contained in the stub
func (s stub) GetCode() (code string) {
	code = "invalidcode"
	if s.HREF != "" {
		//code = strings.Replace(s.HREF, URL, "", -1)
		codes := strings.SplitN(s.HREF, "/", -1)
		code = codes[len(codes)-1]
	}
	return
}

//A recurly date object
type RecurlyDate struct {
	Raw string `xml:",innerxml"`
}

//Convert the date to RFC3339 format
func (r RecurlyDate) GetDate() (time.Time, error) {
	if r.Raw == "" {
		return time.Now(), errors.New("Datetime is blank")
	}
	t, err := time.Parse(time.RFC3339, r.Raw)
	return t, err
}

/* end resource objects */
