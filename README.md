# registration-webapp

How it works...

## User Signup

Uses vmail microservice! <https://github.com/tintinnabulate/vmail>

1. `POST vmail_microservice_url/signup/{email_address}`
2. `GET vmail_microservice_url/verify/{code}`

## User Registration & Payment

Also uses vmail microservice!

1. When user clicks "Register", if `GET vmail_microservice_url/signup/{email_address}` response JSON = `{"Address": "email_address", "Success": true, "Note": ""}` GOTO (2), else GOTO (3)
2. Take payment. If successful payment, GOTO (4), else GOTO (1).
3. Redirect to signup page
4. Store user object in Registrations database.

## Notes

* See <https://tutorialedge.net/golang/consuming-restful-api-with-go/>
* See <https://stackoverflow.com/questions/26035816/testing-post-method-with-csrf>
* See <https://gowebexamples.com/password-hashing/>
