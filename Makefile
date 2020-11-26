REPO = github.com/imega/daemon
IMG = imega/daemon
TAG = latest
CWD = /go/src/$(REPO)
GO_IMG = golang:alpine

test: lint

lint:
	@docker run --rm -t -v $(CURDIR):$(CWD) -w $(CWD) \
		golangci/golangci-lint golangci-lint run
