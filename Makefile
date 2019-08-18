all:
	# launch dev version of app on localhost
	go run .
test:
	# verbose mode, get code coverage, check for race conditions, on all *_test.go files in this package
	go test -v -cover -race ./...
coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
deploy:
	# deploy to live site, creating a new instance (do this before overwrite!).
	go generate
	gcloud app deploy --project 000000
overwrite:
	# deploy to live site overwriting version 0 (do this only after you've tested deploy!)
	go generate
	gcloud app deploy --project 000000 --version 0
