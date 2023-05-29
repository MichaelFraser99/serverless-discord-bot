build:
	@(cd service && make build)

test-verbose:
	go test -v ./...

test:
	go test ./...

pull-upstream:
	git pull upstream main