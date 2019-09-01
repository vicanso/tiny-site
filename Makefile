export GO111MODULE = on

.PHONY: default test test-cover dev

# for dev
dev:
	fresh

test: export GO_ENV=test
test:
	go test -cover ./...

test-cover: export GO_ENV=test
test-cover:
	go test -race -coverprofile=test.out ./... && go tool cover --html=test.out

list-mod:
	go list -m -u all

build:
	packr2
	go build -o tiny-site 

clean:
	packr2 clean
