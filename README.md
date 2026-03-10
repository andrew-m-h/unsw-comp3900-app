# UNSW COMP3900 App

Full-stack guestbook application: Go backend (Chi, DynamoDB) and Vue 3 frontend (Vite), with AWS deployment (App Runner, CloudFront, S3).

---

## GitHub Actions

Workflows live in `.github/workflows/`.

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| **Pre-build checks** | Push/PR to `main` | Runs unit tests, frontend tests, builds frontend, starts the full stack (`make local-up`), runs integration and E2E tests, then tears down. Must pass before deploy. |
| **Build and deploy** | After Pre-build checks succeeds on `main`, or manual `workflow_dispatch` | Deploys BaseStack (CDK) → builds and pushes Docker image to ECR → uploads frontend to S3 and deploys AppStack (App Runner + CloudFront). |
| **Build and push to ECR** | Manual or `workflow_call` | Builds the backend Docker image and pushes it to AWS ECR. Uses OIDC (no long-lived keys). |
| **Upload frontend to S3** | Manual or `workflow_call` | Builds the Vue frontend and syncs `frontend/dist/` to the static assets S3 bucket from BaseStack. |

**Secrets / variables (for deploy):**

- `AWS_ROLE_ARN` – IAM role for ECR push (OIDC).
- `AWS_CDK_DEPLOY_ROLE_ARN` – IAM role for CDK deploy and S3.
- `ECR_REPOSITORY_URI` – ECR repository URI (from BaseStack output).
- `AWS_REGION` (optional) – e.g. `ap-southeast-2` or `us-east-1`.

---

## How to use the Makefile

Run from the **repository root**.

### Build and run

| Target | Description |
|--------|-------------|
| `make build` | Build backend binary into `bin/server`. |
| `make run` | Run backend on host; uses LocalStack env (`AWS_ENDPOINT_URL=http://127.0.0.1:4566`). Use after `make local-resources-up`. |
| `make tidy` | Run `go mod tidy` in `backend/`. |
| `make clean` | Remove `bin/`. |

### Docker (image only)

| Target | Description |
|--------|-------------|
| `make docker-build` | Build Docker image from `backend/` (default image name: `unsw-comp3900-app`). |
| `make docker-push` | Build and push image (override with `IMAGE=...`). |

### Frontend (Vue / Vite)

| Target | Description |
|--------|-------------|
| `make frontend-install` | `npm install` in `frontend/`. |
| `make frontend-build` | Install deps and build; output in `frontend/dist/`. |
| `make frontend-dev` | Start Vite dev server (proxies `/api` to `http://localhost:8080`). |
| `make frontend-preview` | Serve production build with `vite preview`. |
| `make frontend-clean` | Remove `frontend/dist` and `frontend/node_modules`. |
| `make frontend-test` | Run Vitest tests (`npm run test:run`). Requires `frontend-install` first. |

### Local stack (Docker Compose)

| Target | Description |
|--------|-------------|
| `make local-up` | Start LocalStack, dynamodb-init, backend, and frontend (nginx). **Build frontend first:** `make frontend-build && make local-up`. |
| `make local-down` | Stop and remove containers. |
| `make local-resources-up` | Start only LocalStack + dynamodb-init (no backend/frontend). Use when running backend on host with `make run`. |
| `make local-resources-down` | Stop LocalStack. |

### Tests

| Target | Description |
|--------|-------------|
| `make unittest` | Backend unit tests (no external deps; uses mocks). |
| `make integration` | Backend integration tests (in-process server + real LocalStack/DynamoDB). Prerequisite: `make local-resources-up`. |
| `make e2e` | E2E tests (HTTP client vs running backend). Prerequisite: `make local-up`; tests hit `http://localhost:3000`. |

### Version

| Target | Description |
|--------|-------------|
| `make version-bump` | Bump patch version in `backend/version.go`. |
| `make version-minor` | Bump minor version. |

---

## Backend Code

- **Language / runtime:** Go 1.24+
- **Location:** `backend/`
- **Entrypoint:** `backend/main.go` – creates DynamoDB guestbook client, wires Chi router, listens on `:8080`.

**Layout:**

- `main.go` – bootstrap; `server.NewHandler` returns the HTTP handler.
- `internal/server/` – Chi router, middleware, route registration.
- `internal/handlers/` – HTTP handlers (health, guestbook CRUD).
- `internal/guestbook/` – DynamoDB client (create, get, list, delete), entry types, and test mock.
- `internal/errors/` – error types and handler helpers.
- `internal/middleware/` – logging, etc.
- `tests/` – `e2e/` (external HTTP), `integration/` (in-process + LocalStack), `testclient/`.

**Environment (local with LocalStack):**

- `AWS_ACCESS_KEY_ID=test`, `AWS_SECRET_ACCESS_KEY=test`
- `AWS_DEFAULT_REGION=us-east-1`
- `AWS_ENDPOINT_URL=http://127.0.0.1:4566`
- `GUESTBOOK_TABLE_NAME=Guestbook`

**Commands:**

```bash
make build          # build binary
make run            # run on host (after local-resources-up)
make unittest       # unit tests
make integration    # integration tests (after local-resources-up)
make e2e            # E2E (after local-up)
```

---

## Frontend Code

- **Stack:** Vue 3, Vite 6, Vitest.
- **Location:** `frontend/`
- **Entry:** `index.html` → `src/main.js` → `src/App.vue`.

**Layout:**

- `src/main.js` – app bootstrap.
- `src/App.vue` – root component and routing (if any).
- `src/components/` – e.g. `WeddingHome.vue`, `ExamplePage.vue`.
- `src/api/` – API client (e.g. `guestbook.js`).
- `vite.config.js` – base `/static/`, dev proxy `/api` → `http://localhost:8080`.

**Commands:**

```bash
make frontend-install   # npm install
make frontend-build     # production build → frontend/dist
make frontend-dev       # dev server (API proxy to backend)
make frontend-preview   # preview production build
make frontend-test      # Vitest
```

---

## Local development

### Option A: Full stack in Docker (backend + frontend + LocalStack)

1. **Build frontend** (nginx serves from `frontend/dist`):
   ```bash
   make frontend-build && make local-up
   ```
2. **Use:**
   - App (frontend + API): **http://localhost:3000**
   - Backend only: **http://localhost:8080**
   - LocalStack: **http://localhost:4566**
3. **Stop:** `make local-down`

Run E2E against this stack: `make e2e`.

### Option B: Backend on host, frontend dev server

1. **Start only LocalStack + DynamoDB table:**
   ```bash
   make local-resources-up
   ```
2. **Run backend on host:**
   ```bash
   make run
   ```
   (Backend uses `AWS_ENDPOINT_URL=http://127.0.0.1:4566` and `GUESTBOOK_TABLE_NAME=Guestbook`.)
3. **Run frontend dev server** (Vite proxies `/api` to `http://localhost:8080`):
   ```bash
   make frontend-dev
   ```
4. Open the Vite dev URL (e.g. http://localhost:5173). Stop LocalStack when done: `make local-resources-down`.

### Integration tests (backend + LocalStack, no browser)

```bash
make local-resources-up
make integration
make local-resources-down
```

### Prerequisites

- **Go** 1.24+
- **Node.js** 20+ (for frontend and CDK)
- **Docker** and **Docker Compose** (for `local-up`, `local-resources-up`, and E2E)

Optional: `.env.example` is used by docker-compose for LocalStack/backend env; copy to `.env` if you need to override.
