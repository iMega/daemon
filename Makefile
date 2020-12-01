REPO = github.com/imega/daemon
CWD = /go/src/$(REPO)
GO_IMG = golang:alpine

test: lint unit acceptance

lint:
	@-docker run --rm -t -v $(CURDIR):$(CWD) -w $(CWD) \
		golangci/golangci-lint golangci-lint run

unit:
	@docker run --rm -w $(CWD) -v $(CURDIR):$(CWD) \
		$(GO_IMG) sh -c "go list ./... | grep -v 'tests' | xargs go test -vet=off -coverprofile cover.out"

acceptance: down
	GO_IMG=$(GO_IMG) CWD=$(CWD) docker-compose up -d --build --scale acceptance=0
	GO_IMG=$(GO_IMG) CWD=$(CWD) docker-compose up --abort-on-container-exit acceptance

down:
	GO_IMG=$(GO_IMG) CWD=$(CWD) docker-compose down -v --remove-orphans
