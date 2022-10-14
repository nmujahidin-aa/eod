BINARY=eod
test: clean generate
	go test -v -cover -covermode=atomic ./...

coverage: clean generate
	bash coverage.sh --html

run: generate
	go run .

build:
	go build -o ${BINARY} .
	
generate:
	go generate ./...

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	@find . -name *mock* -delete
	@rm -rf .cover
	
.PHONY: test coverage clean build run
