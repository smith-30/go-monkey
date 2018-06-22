vet:
	go vet -v $(go list ./... | grep -v /vendor/)

test:
	go test -v $(go list ./... | grep -v /waiig_code_1.4/)

ci-test:
	go test -v -cover -race $(go list ./... | grep -v /vendor/) -coverprofile go-monkey.coverage.out