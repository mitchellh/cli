TEST?=./...

default: test

# test runs the test suite and vets the code
test: generate
	go list $(TEST) | xargs -n1 go test -timeout=60s -parallel=10 $(TESTARGS)

# testrace runs the race checker
testrace: generate
	go list $(TEST) | xargs -n1 go test -race $(TESTARGS)

# updatedeps installs all the dependencies to run and build
updatedeps:
	go list ./... \
		| xargs go list -f '{{ join .Deps "\n" }}{{ printf "\n" }}{{ join .TestImports "\n" }}' \
		| grep -v github.com/mitchellh/cli \
		| xargs go get -f -u -v

# generate runs `go generate` to build the dynamically generated source files
generate:
	find . -type f -name '.DS_Store' -delete
	go generate ./...

.PHONY: test testrace updatedeps generate
