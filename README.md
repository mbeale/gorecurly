gorecurly
=========
This package is meant to be used with Recurly.com's services.  Only works with version 2 of their api.

Sort of ready for production.  There are live tests to do most CRUD type operations agains resources but there could still be some edge case issues

Installing
==========

	go get github.com/mbeale/gorecurly

Examples
=======

	
	package main

	import (
		"github.com/mbeale/gorecurly"
	)

	func main() {
		r := gorecurly.InitRecurly("fad7d9622a9a49489393d4139609f804", "e44b36f13c92465eb519d70e24b4054c")
		r.EnableDebug() //this will print out more verbose messaging
		acc := r.NewAccount()
		acc.AccountCode = 'test-account'
		acc.Email = "muemail@example.com"
		acc.B = new(gorecurly.BillingInfo)
		acc.B.FirstName = "First"
		acc.B.LastName = "Last"
		acc.B.Number = "4111111111111111"
		acc.B.Month = "12"
		acc.B.Year = "2015"
		acc.B.VerificationValue = "123"
		if err := acc.Create(); err == nil {
			println("Status:" + acc.State)
		} else {
			println(err.Error())
		}
	}

More examples in test

Documentation
=============

Main documentation at [GoPkgDoc](http://go.pkgdoc.org/github.com/mbeale/gorecurly)

Why Live Testing
================

I thought that live testing against an actual connection would be more beneficial then testing against XML fixtures.  If Recurly changes their XML layout, that could
provide problems that testing against the fixtures wouldn't show.  I guess I am of the belief that any library that is considered a client library should actually run
tests against the server.  

Also, the live testing function is a great place for finding examples on using the library.

There are some issues with live testing.  If you already have a lot of test data in your sandbox account, the tests could time out.  I would suggest only testing 
after you have cleared out the test data and never test against a production account, only sandbox.  There is use of random numbers for some object creation, there could be conflicts even though I have never ran into any.

There is a config.xml which you need to alter to suit your account to complete live testing.

TODO
====

* PDF Invoice
* Recurly.js signing
* transparent post (probably not)
* Option to add no auth to header "Recurly-Skip-Authorization: true"
* Discount in cents for coupons not working

Recurly.com
===========

They will be unable to provide support for this library so open issues in this github repo. 

Contributing
============

I welcome all that are willing to contribute.  If you submit a pull request, please make sure that proper tests are submitted if adding new functionality or that
all the tests pass if submitting a bug fix.
