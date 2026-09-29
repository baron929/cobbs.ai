<div align="center">

<img src="https://cdn.cobbs.aiai.org/img/cobbs.ai-logo_1900x450.png" alt="cobbs.ai" width="480">

<br/>
<br/>

**Next-generation personal AI assistant powered by LLM, RAG and agent loops — ships as a single binary, no installation needed**

*Supporting computer-use, browser-use and coding agent*

<br/>

[![Build](https://github.com/baron929/cobbs.ai/actions/workflows/build.yml)](https://github.com/baron929/cobbs.ai/actions/workflows/build.yml)
[![Release](https://img.shields.io/github/v/release/baron929/cobbs.ai?style=flat-square&color=4f46e5)](https://github.com/baron929/cobbs.ai/releases/latest)
[![Docker Pulls](https://img.shields.io/docker/pulls/casbin/cobbs.ai?style=flat-square&color=0ea5e9)](https://hub.docker.com/r/casbin/cobbs.ai)
[![Go Report](https://goreportcard.com/badge/github.com/baron929/cobbs.ai?style=flat-square)](https://goreportcard.com/report/github.com/baron929/cobbs.ai)
[![License](https://img.shields.io/github/license/baron929/cobbs.ai?style=flat-square&color=22c55e)](https://github.com/baron929/cobbs.ai/blob/master/LICENSE)
[![Discord](https://img.shields.io/discord/1022748306096537660?logo=discord&label=discord&color=5865F2&style=flat-square)](https://discord.gg/5rPsrAzK7S)

<br/>

[**Live Demo**](https://demo.cobbs.aiai.org) · [**Playground**](https://try.cobbs.aiai.org) · [**Docs**](https://www.cobbs.aiai.org) · [**Discord**](https://discord.gg/5rPsrAzK7S)

</div>

---

English | [中文](./README_zh.md)

---

## Overview

cobbs.ai is an open-source, self-hosted AI assistant for connecting language models, searching your own knowledge, and running agent workflows. It supports multiple model providers, retrieval-augmented generation (RAG), MCP integrations, and browser, shell, and file tools. These capabilities can access sensitive data or perform consequential actions; review [Security](#security) and [SECURITY.md](SECURITY.md) before enabling them.

Start with the [Local Development](#local-development) guide, browse the [API documentation](#api-documentation), or try the [online demo](https://demo.cobbs.aiai.org).

## Architecture

- **Web client:** React single-page application in `web/`; in development it runs on port `13001` and proxies API requests to the backend.
- **HTTP service:** Go and Beego serve the REST API under `/api`, authentication and authorization filters, and the production web assets.
- **Application and persistence:** `model/` coordinates model and tool workflows; `object/` contains XORM-backed entities and database access.
- **Integrations:** Provider adapters connect to hosted and local models, embedding services, storage backends, and MCP servers.
- **Distribution:** Docker and GoReleaser workflows package the service and frontend; static assets are also supported for single-binary deployments.

---

## Sponsors

<table>
  <tr>
    <td width="300" align="center">
      <a href="https://go.apimart.ai/gh-cobbs.ai" target="_blank"><img src="https://cdn.cobbs.aiai.org/img/sponsor_apimart.png" alt="APIMart" width="280"></a>
    </td>
    <td>
      Thanks to APIMart for sponsoring this project! APIMart is a low-cost API platform for AI image &amp; video generation &mdash; GPT-Image-2 from $0.006/image, 160+ images per dollar. One async API covers both image and video: submit a task, get an ID, fetch results via polling or callback. Batch tens of thousands of images without timeouts, switch models without changing code. Pay-as-you-go with no monthly fee &mdash; <a href="https://go.apimart.ai/gh-cobbs.ai" target="_blank">sign up here</a> to get started.
    </td>
  </tr>
</table>

---

## Features

### 🤖 30+ Model Providers

Connect every major LLM provider and switch between them per conversation — no code changes required.

<div align="center">

`OpenAI` · `Azure OpenAI` · `Anthropic Claude` · `Google Gemini` · `DeepSeek` · `Mistral` · `Grok` · `Qwen` · `Doubao` · `Moonshot` · `ChatGLM` · `Baichuan` · `Ernie` · `iFlytek` · `HuggingFace` · `Cohere` · `Amazon Bedrock` · `OpenRouter` · `Ollama` · `APIMart` · `and more`

</div>

---

### 🔄 Autonomous Agent Loops

| Capability                 | Description                                                                                        |
|:---------------------------|:---------------------------------------------------------------------------------------------------|
| **Browser-Use**            | Drive a real browser — navigate, click, fill forms, scrape, and screenshot pages                   |
| **Web Search & Fetch**     | Search the web and pull live page content into the agent's context                                 |
| **Shell Execution**        | Run shell commands and scripts directly from the agent loop                                        |
| **Office Automation**      | Read and write Word, Excel, and PowerPoint files                                                   |
| **MCP Integration**        | Plug in any MCP-compatible server (SSE / Stdio / StreamableHTTP) and expose its tools to the agent |
| **Transparent Tool Calls** | See every tool invocation, its arguments, and its return value — step by step                      |

---

### 📚 RAG & Knowledge Base

| Capability               | Description                                                                                   |
|:-------------------------|:----------------------------------------------------------------------------------------------|
| **Document Ingestion**   | Upload PDFs, Word docs, Excel sheets, and more — chunked, embedded, and indexed automatically |
| **Semantic Search**      | Retrieves the most relevant passages from your knowledge base before each LLM response        |
| **Pluggable Embeddings** | OpenAI, Azure, Gemini, Qwen, Cohere, Jina, HuggingFace, local models, and more                |
| **Isolated Stores**      | Organise knowledge into separate stores; assign them per chat or per application              |

---

### ⚡ Workflow Automation

| Capability                           | Description                                                         |
|:-------------------------------------|:--------------------------------------------------------------------|
| **Visual Workflow Builder**          | Compose multi-step pipelines with a BPMN-style drag-and-drop editor |
| **Conditional & Parallel Execution** | Branch on gateway conditions; run independent tasks concurrently    |
| **Task Scheduling**                  | Trigger workflows or agent jobs on a recurring schedule             |
| **Usage Analytics**                  | Track token consumption and cost per provider, model, and user      |

---

### 🏗️ Platform Features

| Capability                  | Description                                                                                                       |
|:----------------------------|:------------------------------------------------------------------------------------------------------------------|
| **Single Binary**           | One executable file — no installer, no runtime dependencies. Download and run instantly on any supported platform |
| **Native Windows Support**  | Runs directly on Windows — no WSL, no Docker, no Linux subsystem required                                         |
| **Single Sign-On**          | OIDC / OAuth2 / LDAP / SAML via the built-in auth layer                                                           |
| **Multi-tenancy**           | Isolated workspaces per user or organisation                                                                      |
| **REST API + Swagger UI**   | Every feature is accessible programmatically                                                                      |
| **Audit Logs**              | Full activity history for every action                                                                            |
| **File & Media Management** | Built-in storage for files, images, and video content                                                             |

---

### 📊 Admin Dashboard

| Panel                   | What you get                                                                                |
|:------------------------|:--------------------------------------------------------------------------------------------|
| **Usage Statistics**    | Token & cost metrics per app, user, and model — with interactive charts and heatmaps        |
| **Activity Monitoring** | Real-time system operations with success/error rates, operation-type breakdowns, and trends |
| **Tool Management**     | Centralised CRUD for all agent tools: browser, shell, office, web search, and more          |
| **Request Logs**        | Full request/response payloads with JSON formatting, filtering, and debugging               |

---

## Tech Stack

| Area | Technology |
|:---|:---|
| Backend | Go 1.25+, Beego 1.12, XORM |
| Frontend | React 18, Ant Design, React Router, CRACO / Create React App |
| Database | MySQL by default; PostgreSQL and SQLite adapters are included |
| Authorization | Casbin-backed route policies and application-level tool authorization |
| Observability | Prometheus metrics and application logs |
| Packaging | Docker, embedded web assets, GoReleaser, GitHub Actions |

## Project Structure

| Path | Purpose |
|:---|:---|
| `controllers/` | HTTP handlers and API behavior |
| `routers/` | Route registration and request filters |
| `model/` | Model providers, agent loops, and orchestration |
| `tool/`, `toolauth/`, `guard/` | Built-in tools, tool authorization, and policy decisions |
| `object/`, `storage/` | Database entities, persistence, and storage adapters |
| `conf/` | Runtime configuration and production validation |
| `web/` | React application, frontend tests, and localization |
| `swagger/` | Generated API documentation assets |
| `deploy/`, `scripts/`, `.github/workflows/` | Deployment assets, helper scripts, and CI/release workflows |

## Local Development

### Prerequisites

- Go 1.25.0 or newer (the module toolchain is Go 1.25.8).
- Node.js 20 or newer and Yarn 1.x for the frontend.
- A configured database and authentication provider. The default database driver is MySQL.

Configure `conf/app.conf` (or provide equivalent environment overrides) with a valid database `dataSourceName` and the required authentication settings before starting the backend. Production mode validates required settings and exits when they are missing. Keep credentials out of source control.

Run the backend and frontend in separate terminals:

```bash
# Repository root: start the API on port 14000
go run main.go
```

```bash
# Frontend terminal
cd web
yarn install --frozen-lockfile
yarn start
```

The frontend development server listens on port `13001`; open [http://localhost:13001](http://localhost:13001). For the packaged application, use port `14000`.

### Docker Compose

The root Compose file requires these environment variables: `MYSQL_ROOT_PASSWORD`, `COBBSAI_ISSUER`, `COBBSAI_CLIENT_ID`, `COBBSAI_CLIENT_SECRET`, `COBBSAI_ORGANIZATION`, and `COBBSAI_APPLICATION`. Supply them through your local environment or a protected secrets manager, then run:

```bash
docker compose up --build
```

The database is bound to loopback for local administration; the application remains available on port `14000`. Do not expose the database port publicly.

## API Documentation

The generated Swagger UI is served at [`/swagger`](http://localhost:14000/swagger) when the backend is running. API routes use the `/api/` prefix; authentication and permissions depend on the route. For example, sign-in is `POST /api/signin`. Prometheus metrics are available at `GET /api/metrics` to administrators only. See `routers/router.go` and the generated Swagger assets for the current endpoint definitions.

## Testing

Run backend tests and static analysis from the repository root:

```bash
go test ./...
go vet ./...
```

Run frontend tests and lint from `web/`:

```bash
CI=true yarn test --watchAll=false --runInBand
yarn eslint src/ --ext .js
yarn build
```

Some integration tests need external services or database configuration. A skipped integration test is not evidence that the external integration works. The production frontend build can require substantial memory; CI is the reference environment if a local build runs out of memory.

## Security

Read [SECURITY.md](SECURITY.md) for vulnerability reporting, security controls, known limitations, and verification status. In particular:

- Run the service with least privilege and place it behind a TLS-terminating reverse proxy for network deployments.
- Treat model output and MCP tool requests as untrusted. Shell, filesystem, browser, and process tools are not fully isolated by an OS sandbox; restrict their host and network access externally.
- Use strong, private database and authentication credentials. Do not commit secrets or put credentials in URLs.
- Keep the database private; only expose the application through an appropriately authenticated and monitored ingress.
- Audit logs are best-effort JSONL sidecars and are not a substitute for durable centralized security logging.

No repository-level checklist can guarantee that a deployment is secure. Review the complete security policy and deployment environment before exposing the service.

The frontend dependency audit on 2026-09-29 reported 231 high, 115 moderate, and 32 low advisory instances, with zero critical instances. Counts can include multiple dependency paths to the same advisory. The legacy spreadsheet package `xlsx` remains affected by advisories for which the npm registry reports no patched version; it is currently used for spreadsheet export. Treat dependency auditing as outstanding security work, not as a completed hardening claim.

## Deployment

For a local or self-hosted installation, Docker Compose uses the configuration above. Pre-built installers are also available for Linux, macOS, and Windows:

```bash
curl -fsSL https://raw.githubusercontent.com/baron929/cobbs.ai/master/scripts/install.sh | bash
```

```powershell
irm https://raw.githubusercontent.com/baron929/cobbs.ai/master/scripts/install.ps1 | iex
```

The installers download and run the latest release; inspect scripts before executing them in sensitive environments. For production, configure secrets outside the repository, terminate TLS at a trusted proxy, restrict network access, and use immutable image digests. See `Dockerfile`, `docker-compose.yml`, and `.github/workflows/build.yml` for the current build and release pipeline.

## Monitoring

- `GET /api/metrics` exposes Prometheus-format metrics and requires administrator authorization.
- `GET /api/get-prometheus-info` returns the application metrics summary and also requires administrator authorization.
- Application logs use the configured Beego logger. Audit events are written as best-effort JSONL sidecars; configure permissions, rotation, retention, and external collection for your deployment.
- The admin dashboard includes usage and activity views. Do not expose monitoring endpoints to untrusted networks.

## Screenshots / Demo

<div align="center">

| Usage Analytics | Activity Monitoring |
|:---:|:---:|
| ![Usage Analytics](https://raw.githubusercontent.com/the-open-agent/static/master/img/screenshot-usages.png) | ![Activity Monitoring](https://raw.githubusercontent.com/the-open-agent/static/master/img/screenshot-activities.png) |
| **Tool Management** | **Detailed Logs** |
| ![Tool Management](https://raw.githubusercontent.com/the-open-agent/static/master/img/screenshot-tools.png) | ![Detailed Logs](https://raw.githubusercontent.com/the-open-agent/static/master/img/screenshot-logs.png) |

</div>

| Demo | URL | Notes |
|:---|:---|:---|
| Live Preview | [demo.cobbs.aiai.org](https://demo.cobbs.aiai.org) | Read-only tour; no account needed |
| Playground | [try.cobbs.aiai.org](https://try.cobbs.aiai.org) | Data resets every five minutes |

Full product documentation: [cobbs.aiai.org](https://www.cobbs.aiai.org).

## Future Improvements

The roadmap is tracked in [MODERNIZATION_2026_PLAN.md](MODERNIZATION_2026_PLAN.md), [MODERNIZATION_BACKLOG.md](MODERNIZATION_BACKLOG.md), and [PERFORMANCE_OPTIMIZATION_PLAN.md](PERFORMANCE_OPTIMIZATION_PLAN.md). Key areas include stronger OS-level isolation for consequential tools, durable audit and approval workflows, broader database-backed security tests, and frontend CI gates for tests and lint. These are planned work, not claims that the current release already provides them.

## Community

- **Discord:** [discord.gg/5rPsrAzK7S](https://discord.gg/5rPsrAzK7S)
- **Issues and pull requests:** welcome; open an issue first to discuss larger changes.
- **Sponsors:** thanks to [APIMart](https://go.apimart.ai/gh-cobbs.ai) for supporting the project.

| Environment      | URL                          | Notes                                             |
|:-----------------|:-----------------------------|:--------------------------------------------------|
| **Live Preview** | https://demo.cobbs.aiai.org | Read-only tour — no account needed                |
| **Playground**   | https://try.cobbs.aiai.org  | Make changes freely — data resets every 5 minutes |

---

## License

[Apache 2.0](https://github.com/baron929/cobbs.ai/blob/master/LICENSE)
