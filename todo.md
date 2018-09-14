# TODO

## HIGH

* Create a test suite similar to Email verifier's

	* At the very least, run a test version of the site

	* Use mock Email verifier - https://medium.com/@tech_phil/how-to-stub-external-services-in-go-8885704e8c53

* Remove all instances of `CheckErr(err)` and replace with **PROPER ERROR
  HANDLING!!!**

* Add Captcha on sign-up email entry (for rate-limiting, to stop abuse)

* Does `/register` have a consent drop down yes/no for GDPR Compliance? We are
  collecting/storing email/data.

* Don't load external javascript: host it all internally.

* Somehow expose User table to Host Registration Chairperson

## MEDIUM

* Add form fields to `/register` to match <https://eurypaa2018.com> ?

* Load Customer ID based on email address (if seen before). ?

* Use Gorilla SecureCookie instead of having RegistrationForm table? https://www.calhoun.io/securing-cookies-in-go/

* Can one person do multiple registrations e.g.: {1 x AA, 2 x ALANON, 2 x
  Alateen} with a combined bill?

* Is there a payment option other than Credit/Debit card?
	* Can we connect both Stripe and PayPal to the app?
	* Can we add the Host's Account Number, IBAN, etc. for direct deposit?


## LOW

* Does the form have validation ? E.g. needed fields?

* Can the databse be wiped and reused fresh again the next year? Or is the data
  backed up as required and the process is then restarted?

* What is the initial form field?

* Since IRE and EUR use the same codebase, why not just have one app, and have
  `/eury/signup` and `/irey/signup` select which registration form to load?

* Fix `go lint` warnings.

* Use eury email provider to send out verification email

* Use Braintree? Stripe can only receive payments in 26 countries, whereas
  Braintree can receive payments in 45 countries!

* Add registration pin map showing country/city where people have registered
from (Google Maps?)

* Status Bar: N users registered today, N users registered total 

* Status Bar: N users registered from your country

* Status Bar: Today’s daily reflections reading using
  <https://github.com/tintinnabulate/daily_readings>
