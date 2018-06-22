GOFILES = $(shell go list ./... | grep -v /vendor/)

vet:
	go vet -v $(GOFILES)

test:
	go test -v $(go list ./... | grep -v /waiig_code_1.4/)

ci-test:
	go test -v -cover -coverpkg $(GOFILES) -coverprofile go-monkey.coverage.out