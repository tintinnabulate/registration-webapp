all:
	# launch dev version of app on localhost
	dev_appserver.py .
test:
	# verbose mode, get code coverage, check for race conditions, on all *_test.go files in this package
	go test -v -cover -race ./...
deploy:
	# deploy to live site, creating a new instance (do this before overwrite!)
	rm -r vendor
	gcloud app deploy
	git reset --hard HEAD
overwrite:
	# deploy to live site overwriting version 0 (do this only after you've tested deploy!)
	rm -r vendor
	gcloud app deploy --version 0
	git reset --hard HEAD
