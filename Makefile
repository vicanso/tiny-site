.PHONY: default test test-cover dev

defalt: dev

# for dev
dev: export CONFIG=./configs
dev:
	fresh

# for test
test: export VIPER_INIT_TEST=true
test: export GO_ENV=test
test:
	go test -race -cover ./...

test-cover: export VIPER_INIT_TEST=true
test-cover: export GO_ENV=test
test-cover:
	go test -race -coverprofile=test.out ./... && go tool cover --html=test.out
