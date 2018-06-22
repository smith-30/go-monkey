test:
	go test -v $(go list ./... | grep -v /waiig_code_1.4/)

ci-test:
	go test -cover -coverpkg $(go list ./... | grep -v /vendor/) -coverprofile go-monkey.coverprofile