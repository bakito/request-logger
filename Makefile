# Run golangci-lint
lint:
	golangci-lint run

# Run golangci-lint --fix
fix-lint:
	golangci-lint run --fix

# Run go vet against code
vet:
	go vet ./...

# Run go mod tidy
tidy:
	go mod tidy

# Run tests
test: tidy lint
	go test ./...  -coverprofile=coverage.out
	go tool cover -func=coverage.out

release: tools
	@version=$$(semver); \
	git tag -s $$version -m"Release $$version"
	goreleaser --rm-dist

test-release: tools
	goreleaser --skip-publish --snapshot --rm-dist

tools:
ifeq (, $(shell which goreleaser))
 $(shell go get github.com/goreleaser/goreleaser)
endif
