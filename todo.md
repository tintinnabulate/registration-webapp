# TODO
1. Use eury email provider to send out verification email
2. Add form fields to `/register` to match <https://eurypaa2018.com> ?
3. Use Braintree? Stripe can only receive payments in 26 countries, whereas Braintree can receive payments in 45 countries!

# I.

* Does `/register` have a consent drop down yes/no for GDPR Compliance? We are collecting/storing email/data.
* Is there a payment option other than Credit/Debit card?
* Can the databse be wiped and reused fresh again the next year? Or is the data backed up as required and the process is then restarted?
* Can one person do multiple registrations e.g.: {1 x AA, 2 x ALANON, 2 x Alateen} with a combined bill?
* Can we connect both Stripe and PayPal to the app?
* Can we add the Host's Account Number, IBAN, etc. for direct deposit?

# P.

* Does the form have validation ? E.g. needed fields?
* Is the user journey very likely to see the person has already signed up before getting to Registration? If not, having an auto redirect on initial submit might be good.
* What is the initial form field?

# S.

* Create a test suite similar to Email verifier's
	* At the very least, run a test version of the site
	* Use mock Email verifier
* Memcached
* Since IRE and EUR use the same codebase, why not just have one app, and have `/eury/signup` and `/irey/signup` select which registration form to load?
* Remove all instances of `CheckErr(err)` and replace with **PROPER ERROR HANDLING!!!**
* Add pegistration pin map showing country/city where people have registered from (Google Maps?)
* Load Customer ID based on email address (if seen before).
* Add Captcha on sign-up email entry (for rate-limiting, to stop abuse)
* Add Rate-limiting AND OAuth2 on Email verifier microservice itself
* Don't load external javascript: host it all internally.

## Status bar for:
* N users registered today, N users registered total 
* N users registered from your country
* Today’s daily reflections reading using <https://github.com/tintinnabulate/daily_readings>
