GOFILES = $(shell go list ./... | grep -v /vendor/ | grep -v /waiig_code_1.4/)

.PHONY: \
	vet \
	test \
	ci-test \

vet:
	go vet -v $(GOFILES)

test:
	go test -v $(GOFILES)

ci-test:
	go test -v -cover -race -coverpkg=./... $(GOFILES) -coverprofile=coverage.out