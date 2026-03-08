.PHONY: build run tidy clean docker-build docker-push
.PHONY: frontend-install frontend-build frontend-dev frontend-preview frontend-clean
.PHONY: local-up local-down local-resources-up local-resources-down localstack-init e2e

IMAGE ?= unsw-comp3900-app
FRONTEND_DIR = frontend
BACKEND_DIR = backend

build:
	cd $(BACKEND_DIR) && go build -o ../bin/server .

# Run backend on host. Exports LocalStack env (127.0.0.1). Use with: make local-resources-up (then make run).
run:
	export AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test AWS_DEFAULT_REGION=us-east-1 \
	       AWS_ENDPOINT_URL=http://127.0.0.1:4566 GUESTBOOK_TABLE_NAME=Guestbook && \
	cd $(BACKEND_DIR) && go run .

tidy:
	cd $(BACKEND_DIR) && go mod tidy

clean:
	rm -rf bin/

docker-build:
	docker build -f Dockerfile -t $(IMAGE) $(BACKEND_DIR)

docker-push: docker-build
	docker push $(IMAGE)

# Frontend (Vue.js) — generates HTML/JS/CSS in frontend/dist
frontend-install:
	cd $(FRONTEND_DIR) && npm install

frontend-build: frontend-install
	cd $(FRONTEND_DIR) && npm run build

frontend-dev:
	cd $(FRONTEND_DIR) && npm run dev

frontend-preview:
	cd $(FRONTEND_DIR) && npm run preview

frontend-clean:
	rm -rf $(FRONTEND_DIR)/dist $(FRONTEND_DIR)/node_modules

version-bump:
	cd $(BACKEND_DIR) && go tool goversion -version-file=version.go patch

version-minor:
	cd $(BACKEND_DIR) && go tool goversion -version-file=version.go minor


# Local integration tests: LocalStack + backend + frontend (docker-compose)
# Build frontend first so nginx can serve it: make frontend-build && make local-up
local-up:
	docker compose up -d --build

local-down:
	docker compose down

# Resource-only stack: LocalStack + dynamodb-init (no backend/frontend). Run backend on host with AWS_ENDPOINT_URL=http://127.0.0.1:4566.
local-resources-up:
	docker compose up -d localstack dynamodb-init

local-resources-down:
	docker compose stop localstack

# Optional: create DynamoDB table in LocalStack manually (normally done by dynamodb-init in docker-compose)
localstack-init:
	@echo "Table is created automatically by the dynamodb-init service when you run: make local-up"

# E2E tests: run against the real backend (external HTTP client). Prerequisite: start stack first with make local-up, then make e2e.
# Issues HTTP requests to E2E_BASE_URL (default http://localhost:3000). Set E2E_BASE_URL when running tests inside Docker.
e2e:
	cd $(BACKEND_DIR) && E2E_BASE_URL=http://localhost:3000 go test -tags=e2e ./e2e -v -count=1
