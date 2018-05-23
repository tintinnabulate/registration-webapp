# Register

## Signup

Use signup microservice!

1. `POST signup_microservice/signup/{email_address}`
2. `GET signup_microservice/verify/{code}`

## Registration

Again, use signup microservice!

1. When they click "Register", if `GET signup_microservice/signup/{email_address}` response JSON = `{"Address": "email_address", "Success": true, "Note": ""}` GOTO (2), else GOTO (3)
2. Take payment. If successful payment, GOTO (4), else GOTO (1).
3. Redirect to signup page
4. Store user in Registrations database.

## Notes

See [https://tutorialedge.net/golang/consuming-restful-api-with-go/]()
