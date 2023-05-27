build:
	@(cd service && make build)

test:
	go test -v ./...