install-test:
	go get -u golang.org/x/lint/golint

install-build:
	go get github.com/golang/dep/cmd/dep
	dep ensure

build:
	go build -o pararius-scraper main.go

lint:
	golint -set_exit_status $(call get_go_packages,vendor)

test:
	go test -v -race $(call get_go_packages,vendor)

coverage:
	@jexia-static/golang/coverage.sh
	@echo "HTML report: .cover/coverage.html"


run-service:
	go run main.go

# Return list of Go packages found in current directory
# #1 (space-separated) list of package names (grep/regex patterns) to ignore
define get_go_packages
	$(shell go list ./... $(foreach pattern,$1,| grep -v /$(pattern)/))
endef
