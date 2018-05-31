VERSION:=0.1
TAG:=v$(VERSION)

COVEROUT = cover.out
GOFMTCHECK = test -z `gofmt -l -s -w *.go | tee /dev/stderr`
GOTEST = go test -v
COVER = $(GOTEST) -coverprofile=$(COVEROUT) -covermode=atomic -race

all: get fmt test

.PHONY: fmt
fmt:
	@echo "Checking format..."
	@$(GOFMTCHECK)

.PHONY: get
get:
	go get -v

.PHONY: test
test:
	@echo "Running tests..."
	@$(COVER)

# Use the 'release' target to start a release
.PHONY: release
release: commit push
	@echo Released $(VERSION)

.PHONY: commit
commit:
	@echo Committing release $(VERSION)
	git commit -am"Release $(VERSION)"
	git tag $(TAG)

.PHONY: push
push:
	@echo Pushing release $(VERSION) to master
	git push --tags
	git push
