# Lambda Code
test-verbose:
	go test -v ./...

test:
	go test ./...

build:
	GOOS=linux GOARCH=arm64 go build -o bootstrap
	zip main.zip bootstrap
	mv main.zip ./infrastructure/service

# Terraform
init:
	cd infrastructure && terraform init

plan:
	cd infrastructure && terraform plan --var-file=vars.tfvars -out tfplan

apply:
	cd infrastructure && terraform apply tfplan

destroy:
	cd infrastructure && terraform destroy --var-file=vars.tfvars
