init:
	terraform init

plan:
	terraform plan --var-file=vars.tfvars -out tfplan

apply:
	terraform apply tfplan

destroy:
	terraform destroy --var-file=vars.tfvars