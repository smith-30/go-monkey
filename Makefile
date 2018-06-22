test:
	go test -v $(go list ./... | grep -v /waiig_code_1.4/)