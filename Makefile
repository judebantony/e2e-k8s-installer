# Sample Makefile for infrastructure provisioning

.PHONY: init plan apply destroy validate clean format health-check

init:
	@echo "Initializing infrastructure..."
	@terraform init

plan:
	@echo "Planning infrastructure changes..."
	@terraform plan

apply:
	@echo "Applying infrastructure changes..."
	@terraform apply

destroy:
	@echo "Destroying infrastructure..."
	@terraform destroy

validate:
	@echo "Validating infrastructure configuration..."
	@terraform validate

clean:
	@echo "Cleaning up temporary files..."
	@rm -rf .terraform/

format:
	@echo "Formatting configuration files..."
	@terraform fmt

health-check:
	@echo "Running infrastructure health checks..."
	@curl -f http://localhost:8080/health || echo "Health check failed"