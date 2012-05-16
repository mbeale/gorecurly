package gorecurly

import (
	"net/http"
	"io"
	"bytes"
	"io/ioutil"
	"fmt"
	"encoding/xml"
	"net/url"
)

const (
	URL = "https://api.recurly.com/v2/"
	libversion = "0.1"
	libname = "Recurly-Go"

)
type nopCloser struct {
	io.Reader
}
//functions

func InitRecurly(apikey string,jskey string) (*Recurly){
	r := new (Recurly)
	r.apiKey = apikey
	r.JSKey = jskey
	return r
}

//interfaces
type Pager interface {
	Next() (*Pager,error)
	//SetData([]byte) 
	getRawBody() []byte
}

//Recurly Errors
type RecurlyError struct {
	XMLName xml.Name `xml:"error"`
	statusCode int
	Symbol string `xml:"symbol"`
	Description string `xml:"description"`
	Details string `xml:"details"`
}

type RecurlyValidationErrors struct {
	XMLName xml.Name `xml:"errors"`
	statusCode int
	Errors []RecurlyValidationError `xml:"error"`
}

type RecurlyValidationError struct {
	XMLName xml.Name `xml:"error"`
	FieldName string `xml:"field,attr"`
	Symbol string `xml:"symbol,attr"`
	Description string `xml:",innerxml"`
}

func CreateRecurlyStandardError(resp *http.Response) (r RecurlyError) {
	r.statusCode = resp.StatusCode
	if xmlstring, readerr := ioutil.ReadAll(resp.Body); readerr == nil {
		if xmlerr := xml.Unmarshal(xmlstring, &r); xmlerr != nil {
			r.Description = string(xmlstring)
		}
	}
	return r
}

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

func createRecurlyError(resp *http.Response) ( error) {
	if resp.StatusCode != 422 {
		return CreateRecurlyStandardError(resp)
	} else {
		return CreateRecurlyValidationError(resp)
	}
	return nil
}

func (r RecurlyError) Error() string {
	return fmt.Sprintf("Recurly Error: %s , %s %s Status Code: %v", r.Symbol,r.Description, r.Details,r.statusCode)
}

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

func (r *Recurly) EnableDebug() {
	r.debug = true
}

func (r *Recurly) GetAccounts(params ...url.Values) (AccountList, error){
	sendvars := url.Values{}
	if params != nil {
		sendvars = params[0] 
	} 
	accountlist := AccountList{}
	if err := accountlist.initList("accounts",sendvars,r); err == nil {
		if xmlerr := xml.Unmarshal(accountlist.getRawBody(), &accountlist); xmlerr == nil {
			for k,_ := range accountlist.Account {
				accountlist.Account[k].r = r
				accountlist.Account[k].endpoint = "accounts"
			}
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

func (r *Recurly) GetAccount(account_code string) (account Account, err error) {
	account = r.NewAccount()
	if resp,err := r.createRequest("accounts/" + account_code,"GET", nil, nil); err == nil {
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

func (r *Recurly) NewAccount() (account Account) {
	account.r = r
	account.endpoint = "accounts"
	return
}

func (r *Recurly) createRequest(endpoint string, method string, params url.Values, msgbody []byte) (*http.Response, error) { 
	client := &http.Client{}

	u, err := url.Parse(URL + endpoint)
	if err != nil {
		return nil,err
	}
	u.RawQuery = params.Encode()
	body := nopCloser{bytes.NewBufferString(string(msgbody))}
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

type Paging struct {
	rawBody []byte
	count, next, prev, start string
}

func (p Paging) getRawBody() ([]byte) {
	return p.rawBody
}
func (p *Paging) SetData(rb []byte, count string, link string) {
	p.rawBody = rb
	p.count = count
	println(link)
	println(count)
	p.next = link
	p.prev = link
	p.start = link
}

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

type BillingInfo struct {
	XMLName xml.Name `xml:"billing_info"`
	endpoint string
	r *Recurly
	FirstName string `xml:"first_name,omitempty"`
	LastName string `xml:"last_name,omitempty"`
}

type AccountList struct {
	Paging
	XMLName xml.Name `xml:"accounts"`
	Account []Account `xml:"account"`
}

func (a AccountList) Next() (*Pager, error) {
	accountlist := new(Pager)
	return accountlist,nil
}

func (a AccountList) Prev() (AccountList, error) {
	accountlist := AccountList{}
	return accountlist,nil
}

func (a AccountList) Start() (AccountList, error) {
	accountlist := AccountList{}
	return accountlist,nil
}

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

func (a *Account) Create() (error) {
	if a.CreatedAt != "" || a.HostedLoginToken != "" || a.State != "" {
		return RecurlyError{statusCode:400,Description:"Account Code Already in Use"}
	}
	return a.r.doCreate(&a,a.endpoint)
}

func (a *Account) Update() (error) {
	newaccount := new(Account)
	*newaccount = *a
	newaccount.State = ""
	newaccount.HostedLoginToken = ""
	newaccount.CreatedAt = ""
	newaccount.B = nil
	return a.r.doUpdate(newaccount,a.endpoint + "/" + a.AccountCode)
}

func (a *Account) Close() (error) {
	return a.r.doDelete(a.endpoint + "/" + a.AccountCode)
}

func (a *Account) Reopen() (error) {
	newaccount := new(Account)
	return a.r.doUpdate(newaccount,a.endpoint + "/" + a.AccountCode + "/reopen")
}

/* end resource objects */


