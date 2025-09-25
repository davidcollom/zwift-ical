VERSION ?= $(shell cat VERSION)

build:
	docker build --build-arg VERSION=$(VERSION) -t davidcollom/zwiftcal:$(VERSION)  .

debug: build
	docker run --rm -ti -v $(PWD):/app/ -v $(PWD)/.cache:/app/cache davidcollom/zwiftcal:$(VERSION) bash

publish: build
	docker push davidcollom/zwiftcal:$(VERSION)
	
test:
	go test -v ./...

coverage:
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html

coverage-report:
	go tool cover -func=coverage.out
