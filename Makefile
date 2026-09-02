BUF_IMAGE := bufbuild/buf:1.47.2
GO_IMAGE := golang:1.25.1
TERRAFORM_IMAGE := hashicorp/terraform:1.13.3

.PHONY: generate generate-check test build run docker-up docker-down terraform-validate
generate:
	docker run --rm -u "$$(id -u):$$(id -g)" -e HOME=/tmp -v "$$(pwd):/workspace" -w /workspace $(BUF_IMAGE) generate

generate-check:
	@tmp="$$(mktemp -d)"; mkdir -p "$$tmp/project"; cp -R buf.yaml buf.gen.yaml proto "$$tmp/project/"; docker run --rm -u "$$(id -u):$$(id -g)" -e HOME=/tmp -v "$$tmp/project:/workspace" -w /workspace $(BUF_IMAGE) generate >/dev/null; diff -ru backend/generated "$$tmp/project/backend/generated"; diff -ru frontend/generated "$$tmp/project/frontend/generated"; rm -rf "$$tmp"

test:
	docker run --rm -v "$$(pwd)/backend:/app" -w /app $(GO_IMAGE) sh -c 'go test ./... && go vet ./...'
	cd frontend && npm test

build:
	docker compose build
	cd frontend && npm run build

run:
	docker compose up --build backend

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down --remove-orphans

terraform-validate:
	docker run --rm -v "$$(pwd)/infra/terraform:/infra" -w /infra $(TERRAFORM_IMAGE) fmt -check
	docker run --rm -v "$$(pwd)/infra/terraform:/infra" -w /infra $(TERRAFORM_IMAGE) init -backend=false
	docker run --rm -v "$$(pwd)/infra/terraform:/infra" -w /infra $(TERRAFORM_IMAGE) validate
