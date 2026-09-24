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

## What is cobbs.ai?

cobbs.ai is an open-source personal AI assistant that brings together powerful LLMs, your own knowledge base, and autonomous agent loops — all in one self-hostable platform. Connect any model provider, build a RAG knowledge base from your documents, and let agents browse the web, run code, and call any MCP-compatible tool on your behalf.

<div align="center">
<br/>

|                                               Usage Analytics                                                |                                                 Activity Monitoring                                                  |
|:------------------------------------------------------------------------------------------------------------:|:--------------------------------------------------------------------------------------------------------------------:|
| ![Usage Analytics](https://raw.githubusercontent.com/the-open-agent/static/master/img/screenshot-usages.png) | ![Activity Monitoring](https://raw.githubusercontent.com/the-open-agent/static/master/img/screenshot-activities.png) |
|                                             **Tool Management**                                              |                                                  **Detailed Logs**                                                   |
| ![Tool Management](https://raw.githubusercontent.com/the-open-agent/static/master/img/screenshot-tools.png)  |       ![Detailed Logs](https://raw.githubusercontent.com/the-open-agent/static/master/img/screenshot-logs.png)       |

<br/>
</div>

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

## Quick Start

Pre-built binaries for **Linux**, **macOS**, and **Windows** (`x86_64` / `arm64`). The installer downloads the latest release and starts cobbs.ai on **port 14000**.

**macOS / Linux / WSL**
```bash
curl -fsSL https://raw.githubusercontent.com/baron929/cobbs.ai/master/scripts/install.sh | bash
```

**Windows (PowerShell)**
```powershell
irm https://raw.githubusercontent.com/baron929/cobbs.ai/master/scripts/install.ps1 | iex
```

> **Windows runs natively** — no WSL, no Docker required.

Open [http://localhost:14000](http://localhost:14000) and you're in.

> Optional env vars: `COBBSAI_VERSION`, `INSTALL_DIR`, `BIN_DIR`

**Build from source**

Prerequisites:
- **Backend**: [Go](https://golang.org/dl/) 1.25.0+
- **Frontend**: [Node.js](https://nodejs.org/) 20+ and [Yarn](https://classic.yarnpkg.com/) 1.x

```bash
# Backend
go build

# Frontend
cd web && yarn install && yarn start
```

**Docker**
```bash
docker-compose up
```

Open [http://localhost:14000](http://localhost:14000) once the containers are running.

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

## Online Demo

| Environment      | URL                          | Notes                                             |
|:-----------------|:-----------------------------|:--------------------------------------------------|
| **Live Preview** | https://demo.cobbs.aiai.org | Read-only tour — no account needed                |
| **Playground**   | https://try.cobbs.aiai.org  | Make changes freely — data resets every 5 minutes |

---

## Documentation

Full docs at **[https://www.cobbs.aiai.org](https://www.cobbs.aiai.org)**

---

## Community

- **Discord** — [discord.gg/5rPsrAzK7S](https://discord.gg/5rPsrAzK7S) · chat with maintainers and other users
- **Issues & PRs** — welcome! Please open an issue first to discuss larger changes

---

## License

[Apache 2.0](https://github.com/baron929/cobbs.ai/blob/master/LICENSE)
