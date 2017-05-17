GO ?= go
COVERAGEDIR = coverage
ifdef CIRCLE_ARTIFACTS
  COVERAGEDIR = $(CIRCLE_ARTIFACTS)
endif

all: test cover
fmt:
	$(GO) fmt ./...
test:
	if [ ! -d coverage ]; then mkdir coverage; fi
	$(GO) test -v ./proto3 -race -cover -coverprofile=$(COVERAGEDIR)/proto3.coverprofile
	$(GO) test -v ./proto3/generator -race -cover -coverprofile=$(COVERAGEDIR)/proto3_generator.coverprofile
cover:
	$(GO) tool cover -html=$(COVERAGEDIR)/proto3.coverprofile -o $(COVERAGEDIR)/proto3.html
	$(GO) tool cover -html=$(COVERAGEDIR)/proto3_generator.coverprofile -o $(COVERAGEDIR)/proto3_generator.html
tc: test cover
coveralls:
	gover $(COVERAGEDIR) $(COVERAGEDIR)/coveralls.coverprofile
	goveralls -coverprofile=$(COVERAGEDIR)/coveralls.coverprofile -service=circle-ci -repotoken=$(COVERALLS_TOKEN)
clean:
	$(GO) clean
	rm -rf coverage/