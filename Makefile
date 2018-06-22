GOFILES = $(shell go list ./... | grep -v /vendor/ | grep -v /waiig_code_1.4/)

vet:
	go vet -v $(GOFILES)

test:
	go test -v $(GOFILES)

ci-test:
	go test -v -cover -coverpkg $(GOFILES) -coverprofile go-monkey.coverage.out