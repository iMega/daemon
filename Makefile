REPO = github.com/imega/daemon
IMG = imega/daemon
TAG = latest
CWD = /go/src/$(REPO)
GO_IMG = golang:alpine

test: lint unit

lint:
	@docker run --rm -t -v $(CURDIR):$(CWD) -w $(CWD) \
		golangci/golangci-lint golangci-lint run

unit:
	@docker run --rm -w $(CWD) -v $(CURDIR):$(CWD) \
		$(GO_IMG) sh -c "go list ./... | xargs go test -vet=off -coverprofile cover.out"
