build:
	@(cd service && make build)

test:
	go test -v ./...

pull-upstream:
	git pull upstream main