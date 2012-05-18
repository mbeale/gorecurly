//Main GoRecurly Package
package gorecurly

//TODO: Do all account tests

import (
	"net/http"
	"io"
	"bytes"
	"io/ioutil"
	"fmt"
	"strings"
	"encoding/xml"
	"net/url"
)

const (
	URL = "https://api.recurly.com/v2/"
	libversion = "0.1"
	libname = "Recurly-Go"
	ACCOUNTS = "accounts"

)
//Generic Reader
type nopCloser struct {
	io.Reader
}
//functions

//Initialize the Recurly package with your apikey and your jskey
func InitRecurly(apikey string,jskey string) (*Recurly){
	r := new (Recurly)
	r.apiKey = apikey
	r.JSKey = jskey
	return r
}

//interfaces


//Paging interface to allow Next,Prev,Start
type Pager interface {
	getRawBody() []byte
}

//Recurly Generic Errors
type RecurlyError struct {
	XMLName xml.Name `xml:"error"`
	statusCode int
	Symbol string `xml:"symbol"`
	Description string `xml:"description"`
	Details string `xml:"details"`
}

//Recurly Validation Errors Array
type RecurlyValidationErrors struct {
	XMLName xml.Name `xml:"errors"`
	statusCode int
	Errors []RecurlyValidationError `xml:"error"`
}

//Recurly validation error
type RecurlyValidationError struct {
	XMLName xml.Name `xml:"error"`
	FieldName string `xml:"field,attr"`
	Symbol string `xml:"symbol,attr"`
	Description string `xml:",innerxml"`
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
		println(string(xmlstring))
		if xmlerr := xml.Unmarshal(xmlstring, &r); xmlerr != nil {
			//r.Description = xmlerr.Error()
			println(xmlerr.Error())
		}
	}
	return r
}

//Filter to decide which error type to create
func createRecurlyError(resp *http.Response) ( error) {
	if resp.StatusCode != 422 {
		return CreateRecurlyStandardError(resp)
	} else {
		return CreateRecurlyValidationError(resp)
	}
	return nil
}

//Formatted General Error 
func (r RecurlyError) Error() string {
	return fmt.Sprintf("Recurly Error: %s , %s %s Status Code: %v", r.Symbol,r.Description, r.Details,r.statusCode)
}

//Formatted Validation Error
func (r RecurlyValidationErrors) Error() string {
	var rtnString string
	for _,v := range r.Errors {
		rtnString += v.FieldName + " " + v.Description + "\n"
	}
	return fmt.Sprintf("You have the following validation errors:\n%s", rtnString)
}

//Main Recurly Client
type Recurly struct {
	apiKey, JSKey  string
	debug bool
}

//Set verbose debugging
func (r *Recurly) EnableDebug() {
	r.debug = true
}

//Get a list of accounts
func (r *Recurly) GetAccounts(params ...url.Values) (AccountList, error){
	accountlist := AccountList{}
	sendvars := url.Values{}
	if params != nil {
		sendvars = params[0] 
		accountlist.perPage = sendvars.Get("per_page")
	} 
	if err := accountlist.initList(ACCOUNTS,sendvars,r); err == nil {
		if xmlerr := xml.Unmarshal(accountlist.getRawBody(), &accountlist); xmlerr == nil {
			for k,_ := range accountlist.Account {
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

//Get a single account by key
func (r *Recurly) GetAccount(account_code string) (account Account, err error) {
	account = r.NewAccount()
	if resp,err := r.createRequest(ACCOUNTS + "/" + account_code,"GET", nil, nil); err == nil {
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
					return account,xmlerr
				}
				//everything went fine
				return  account,nil
			} else {
				//return read error
				return account,readerr
			}
			return account,nil
		} else {
			return account,createRecurlyError(resp)
		}
	} else {
		return account, err
	}
	return account, nil
}

//Create a new Account
func (r *Recurly) NewAccount() (account Account) {
	account.r = r
	account.endpoint = ACCOUNTS
	return
}

//Create a request to Recurly and return that response object
func (r *Recurly) createRequest(endpoint string, method string, params url.Values, msgbody []byte) (*http.Response, error) { 
	client := &http.Client{}

	u, err := url.Parse(URL + endpoint)
	if err != nil {
		return nil,err
	}
	u.RawQuery = params.Encode()
	body := nopCloser{bytes.NewBufferString(string(msgbody))}
	if r.debug {
		fmt.Printf("Endpoint Requested: %s Method: %s Body: %s\n", u.String(), method, string(msgbody))
	}
	if req, err := http.NewRequest(method, u.String(), body); err != nil {
		return nil,err
	} else {
		req.Header.Add("Accept", "application/xml")
		req.Header.Add("Accept-Language", "en-US")
		req.Header.Add("User-Agent", libname + " version=" + libversion)
		req.Header.Add("Content-Type","application/xml; charset=utf-8")
		req.ContentLength = int64(len(string(msgbody)))
		req.SetBasicAuth(r.apiKey,"")
		if resp, resperr := client.Do(req); resperr == nil {
			return resp, nil
		} else {
			return nil,resperr
		}
	}
	return nil, nil
}

//Create a resource from struct, uses POST method
func (r *Recurly) doCreate(v interface{}, endpoint string) (error) {
	if xmlstring, err := xml.MarshalIndent(v, "", "    "); err == nil {
		xmlstring = []byte(xml.Header + string(xmlstring))
		if r.debug {
			fmt.Printf("%s\n",xmlstring)
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
					return  nil
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
func (r *Recurly) doUpdate(v interface{}, endpoint string) (error) {
	if xmlstring, err := xml.MarshalIndent(v, "", "    "); err == nil {
		xmlstring = []byte(xml.Header + string(xmlstring))
		if r.debug {
			fmt.Printf("%s\n",xmlstring)
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
func (r *Recurly) doDelete(endpoint string) (error) {
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
	rawBody []byte
	count, next, prev, perPage string
}

//Return the rawBody Var
func (p Paging) getRawBody() ([]byte) {
	return p.rawBody
}

//Set header data for paging
func (p *Paging) SetData(rb []byte, count string, links string) {
	p.rawBody = rb
	p.count = count
	p.next = ""
	p.prev = ""
	for _,v := range strings.SplitN(links,",",-1) {
		println(v)
		link := strings.SplitN(v,";",-1)
		link[0] = strings.Replace(link[0],"<","",-1)
		link[0] = strings.Replace(link[0],">","",-1)
		if u, err := url.Parse(link[0]); err == nil {
			values := u.Query() 
			switch link[1] {
			case " rel=\"next\"" :
				p.next = values.Get("cursor")
			case " rel=\"prev\"" :
				p.prev = values.Get("cursor")
			}
		} 
	}
}

//Initialize the paging list values
func (p *Paging) initList(endpoint string, params url.Values, r *Recurly) ( error) { 
	if resp, err := r.createRequest(endpoint,"GET",params, make([]byte,0)); err == nil {
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
			p.SetData(body,resp.Header["X-Records"][0],resp.Header["Link"][0])
			//everything went fine
			return  nil
		} else {
			//return read error
			return readerr
		}
	} else {
		//return error message
		return err
	}
	return nil
}

/*resource objects */

//Billing Info struct
type BillingInfo struct {
	XMLName xml.Name `xml:"billing_info"`
	endpoint string
	r *Recurly
	FirstName string `xml:"first_name,omitempty"`
	LastName string `xml:"last_name,omitempty"`
	Address1 string `xml:"address1,omitempty"`
	Address2 string `xml:"address2,omitempty"`
	City string `xml:"city,omitempty"`
	State string `xml:"state,omitempty"`
	Zip string `xml:"zip,omitempty"`
	Country string `xml:"country,omitempty"`
	Phone string `xml:"phone,omitempty"`
	VatNumber string `xml:"vat_number,omitempty"`
	IPAddress string `xml:"ip_address,omitempty"`
	IPAddressCountry string `xml:"ip_address_country,omitempty"`
	Number string `xml:"number,omitempty"`
	FirstSix string `xml:"first_six,omitempty"`
	LastFour string `xml:"last_four,omitempty"`
	VerificationValue string `xml:"verification_value,omitempty"`
	CardType string `xml:"card_type,omitempty"`
	Month string `xml:"month,omitempty"`
	Year string `xml:"year,omitempty"`
	BillingAgreementID string `xml:"billing_agreement_id,omitempty"`
}

//Account pager
type AccountList struct {
	Paging
	r *Recurly
	XMLName xml.Name `xml:"accounts"`
	Account []Account `xml:"account"`
}

//Get next set of accounts
func (a *AccountList) Next() (bool) {
	if a.next != "" {
		v := url.Values{}
		v.Set("cursor",a.next)
		v.Set("per_page",a.perPage)
		*a,_ = a.r.GetAccounts(v)
	} else {
		return false
	}
	return true
}

//Get previous set of accounts
func (a *AccountList) Prev() ( bool) {
	if a.prev != "" {
		v := url.Values{}
		v.Set("cursor",a.prev)
		v.Set("per_page",a.perPage)
		*a,_ = a.r.GetAccounts(v)
	} else {
		return false
	}
	return true
}

//Go to start set of accounts
func (a *AccountList) Start() ( bool) {
	if a.prev != "" {
		v := url.Values{}
		v.Set("per_page",a.perPage)
		*a,_ = a.r.GetAccounts(v)
	} else {
		return false
	}
	return true
}

//Account struct
type Account struct{
	XMLName xml.Name `xml:"account"`
	endpoint string
	r *Recurly
	AccountCode string `xml:"account_code"`
	Username string `xml:"username"`
	Email string `xml:"email"`
	State string `xml:"state,omitempty"`
	FirstName string `xml:"first_name"`
	LastName string `xml:"last_name"`
	CompanyName string `xml:"company_name"`
	AcceptLanguage string `xml:"accept_language"`
	HostedLoginToken string `xml:"hosted_login_token,omitempty"`
	CreatedAt string `xml:"created_at,omitempty"`
	B *BillingInfo `xml:"billing_info,omitempty"` 
}

//Create a new account and load updated fields
func (a *Account) Create() (error) {
	if a.CreatedAt != "" || a.HostedLoginToken != "" || a.State != "" {
		return RecurlyError{statusCode:400,Description:"Account Code Already in Use"}
	}
	err := a.r.doCreate(&a,a.endpoint)
	if err == nil {
		a.B = nil
	}
	return err
}

//Update an account 
func (a *Account) Update() (error) {
	newaccount := new(Account)
	*newaccount = *a
	newaccount.State = ""
	newaccount.HostedLoginToken = ""
	newaccount.CreatedAt = ""
	newaccount.B = nil
	return a.r.doUpdate(newaccount,a.endpoint + "/" + a.AccountCode)
}

//Close an account
func (a *Account) Close() (error) {
	return a.r.doDelete(a.endpoint + "/" + a.AccountCode)
}

//Reopen a closed account
func (a *Account) Reopen() (error) {
	newaccount := new(Account)
	return a.r.doUpdate(newaccount,a.endpoint + "/" + a.AccountCode + "/reopen")
}

/* end resource objects */


