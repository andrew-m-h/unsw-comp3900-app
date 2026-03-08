.PHONY: build run tidy clean docker-build docker-push
.PHONY: frontend-install frontend-build frontend-dev frontend-preview frontend-clean

IMAGE ?= unsw-comp3900-app
FRONTEND_DIR = frontend
BACKEND_DIR = backend

build:
	cd $(BACKEND_DIR) && go build -o ../bin/server .

run:
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
