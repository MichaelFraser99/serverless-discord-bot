build:
	@(GOARCH=arm64 GOOS=linux go build -mod=vendor -o bootstrap main.go && zip main.zip bootstrap)