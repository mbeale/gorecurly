package gorecurly

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var accountGet,accountCreate string

func init(){
	accountGet = `
		<?xml version="1.0" encoding="UTF-8"?>
		<account href="https://api.recurly.com/v2/accounts/test21">
      			<adjustments href="https://api.recurly.com/v2/accounts/test21/adjustments"/>
    			<billing_info href="https://api.recurly.com/v2/accounts/test21/billing_info"/>
  			<invoices href="https://api.recurly.com/v2/accounts/test21/invoices"/>
		        <subscriptions href="https://api.recurly.com/v2/accounts/test21/subscriptions"/>
			<transactions href="https://api.recurly.com/v2/accounts/test21/transactions"/>
			<account_code>test21</account_code>
  			<state>active</state>
			<username nil="nil"></username>
			<email>verena@example.com</email>
			<first_name>Verena</first_name>
			<last_name>Example</last_name>
			<company_name nil="nil"></company_name>
			<accept_language>en-US,en;q=0.8</accept_language>
			<hosted_login_token>1781d57cd7cfacc216314349e286ff4a</hosted_login_token>
			<created_at type="datetime">2012-03-14T21:08:20Z</created_at>
		</account>
		`
	accountCreate = `
		<?xml version="1.0" encoding="UTF-8"?>
		<account href="https://api.recurly.com/v2/accounts/abcdef1234567890">
			<adjustments href="https://api.recurly.com/v2/accounts/abcdef1234567890/adjustments"/>
			<billing_info href="https://api.recurly.com/v2/accounts/abcdef1234567890/billing_info"/>
		  	<invoices href="https://api.recurly.com/v2/accounts/abcdef1234567890/invoices"/>
			<subscriptions href="https://api.recurly.com/v2/accounts/abcdef1234567890/subscriptions"/>
      			<transactions href="https://api.recurly.com/v2/accounts/abcdef1234567890/transactions"/>
    			<account_code>abcdef1234567890</account_code>
  			<username>shmohawk58</username>
		        <email>larry.david@example.com</email>
			<first_name>Larry</first_name>
			<last_name>David</last_name>
			<company_name>Home Box Office</company_name>
			<accept_language>en-US</accept_language>
			<created_at type="datetime">2011-04-30T12:00:00Z</created_at>
		</account>
	`
}

func TestGetAccount(t *testing.T) {
	http.HandleFunc("/accounts/test21", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s",accountGet)
	})
	ts := httptest.NewServer(nil)
	defer ts.Close()

	r := InitRecurly("","")
	r.url = "http://" + ts.Listener.Addr().String() + "/"
	if acc, e := r.GetAccount("test21"); e != nil {
		t.Fatal(e.Error())
	} else {
		if acc.AccountCode != "test21" {
			t.Fatal("Could not parse the account object correctly")
		}
	}
}

func TestCreate(t *testing.T) {
}
