package gorecurly

import (
	"testing"
)

type Config struct {
	apikey, jskey string
}

func (c *Config) LoadConfig() bool {
	return false
}

func TestA(t *testing.T) {
	c := Config{}
	//load config
	if !c.LoadConfig() {
		t.Fatalf("Configuration failed to load.")
	}
	//init recurly
	//r := InitRecurly(c.apikey,c.jskey)

	//ACCOUNT TESTS
	//create invalid account
	//create valid account
	//create valid account
	//create valid account with billing info
	//create valid account with billing info
	//get account
	//update account
	//close account
	//reopen account
	//list accounts
	//page through
	//page backwards
	//page start
}
