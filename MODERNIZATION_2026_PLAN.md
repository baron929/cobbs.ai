# cobbs.ai 2026 Modernization Plan

Date: 2026-09-24
Scope: read-only repository audit and modernization plan. No implementation changes are included in this plan. Live deployment, DNS, registry, certificate, and production-runtime claims are excluded unless directly evidenced by repository code or recorded verification output.

## Executive Summary

cobbs.ai is already a broad self-hosted AI platform, not an empty prototype. The repository contains:

- A Go 1.25 backend using Beego 1.12.12, XORM 1.2.5, MySQL/PostgreSQL/SQLite adapters, and embedded single-binary assets.
- A React 18 frontend using Ant Design 6, Create React App 5/CRACO, React Router 5, Yarn 1, ESLint 8, and Stylelint 14.
- Multiple cloud and local model providers, embeddings, RAG/vector retrieval, experience corrections, streaming, tool calls, background jobs, and task analysis.
- Built-in shell, local-file, browser, web, Office, GUI, video, and browser-use tools.
- MCP client support for stdio, SSE, and Streamable HTTP transports.
- Casbin route authorization, a separate tri-state tool guard, Prometheus metrics, and JSONL audit events.
- Docker multi-stage builds, GoReleaser, GitHub Actions, multi-architecture image publication, and Helm chart publication.

The dominant modernization problem is not missing features. It is control-plane integration. The repository has security and observability abstractions, but the highest-risk execution path still directly invokes builtin and MCP tools. The guard decision is not visibly mandatory at the call site, approval handling is not a required execution boundary, and audit delivery is best effort. This is unsafe for shell, filesystem, browser, intranet-scan, and arbitrary MCP operations.

The recommended strategy is incremental:

1. Close P0 security and provenance gaps before expanding capabilities.
2. Introduce a unified agent-run and tool-execution control plane around the existing model/tool abstractions.
3. Add explicit approval, workspace, tenant, device, and secret boundaries.
4. Modernize build/frontend/deployment surfaces in staged migrations.
5. Measure agent quality and security with repeatable evaluations before changing models or providers.

A rewrite is not justified. The provider abstraction, MCP transport abstraction, builtin tool registry, RAG/vector pipeline, route authorization, guard engine, Prometheus metrics, migration framework, and single-binary packaging should remain and be hardened behind clearer service boundaries.

## Audit Method and Evidence Rules

The audit inspected source, configuration, manifests, CI, container files, tests, and documentation. Findings are based on repository evidence, not marketing claims. The live domain and Docker repositories were previously checked and remain unverified: `cobbs.aiai.org` did not resolve, and the configured Docker Hub repositories returned `404 Not Found` during the recorded verification. The current worktree is also dirty, so a deployment commit cannot be inferred from `HEAD` alone.

Current measured repository surface:

- 37 Go test files.
- 3 frontend test/spec files.
- 32 Go files under `tool/`.
- 37 Go files under `model/`.
- 61 Go files under `controllers/`.
- One large Go module with many provider and infrastructure dependencies.

## 1. Current Architecture

### Runtime and backend

Current state:

- Go module: `github.com/baron929/cobbs.ai`.
- Go requirement/toolchain: Go 1.25.0 / toolchain 1.25.8 in [go.mod](go.mod); CI uses Go 1.25.8 in [.github/workflows/build.yml](.github/workflows/build.yml).
- HTTP framework: Beego 1.12.12.
- Persistence: XORM 1.2.5 with MySQL, PostgreSQL, and SQLite drivers in [object/adapter.go](object/adapter.go).
- Controllers and object packages contain routing and database operations directly, with global adapters and configuration state.
- Embedded frontend, skills, OCR, and PPTX worker support are wired through [embed.go](embed.go) and [embedsupport](embedsupport).

Assessment: appropriate for staged hardening and continued development; the global state and controller/object coupling are the main maintainability constraints. A Beego replacement or backend rewrite is not currently justified.

Target state:

- Keep the HTTP/runtime framework initially.
- Introduce request-scoped context, service-layer boundaries, typed configuration, explicit repositories, and dependency injection around new work.
- Isolate high-risk execution into a worker boundary rather than moving every existing controller at once.

### Frontend

Current state:

- React 18.2, Ant Design 6.3.5, React Router 5.3.3, CRA 5.0.1, CRACO 6.4.5.
- Yarn 1 lockfile and class-heavy legacy page patterns coexist with newer components.
- Build/lint scripts exist in [web/package.json](web/package.json), but CI builds the frontend and does not run frontend tests or lint as mandatory gates.

Assessment: the UI framework can remain. CRA/CRACO and Router 5 are the main modernization liabilities. Migrate in stages after the backend execution contract is stable.

Target state:

- Vite or another maintained build tool, Router 6+, modern test runner, and strict TypeScript adoption for new platform surfaces.
- Preserve Ant Design and existing page workflows.
- Add explicit agent run, tool status, approval, cancellation, and verification states without exposing private chain-of-thought.

### Agent and model runtime

Current state:

- Model providers are implemented behind a common model abstraction under [model](model), with many cloud and local providers.
- [model/mcp.go](model/mcp.go) loops over model-generated tool calls and invokes builtin or MCP tools.
- Streaming, tool-call accumulation, vision handling, recovery prompts, and session messages exist.
- Context is primarily recent message history, per-store memory limits, vector retrieval, and experience/correction data.

Assessment: broad capability exists, but planning, durable task state, approval state, retries, cancellation, and structured execution state are not centralized as a durable agent runtime.

Target state:

- An `AgentRun` state machine with run/task/tool-call IDs, explicit phases, cancellation, retry policy, approval transitions, bounded context, and durable status.
- Provider capability metadata for streaming, tools, structured output, vision, context limits, usage, cost, and retry behavior.
- Specialist-agent/handoff support only after the single-agent control plane is safe and observable.

## 2. Component Decisions

### Components to retain

- Existing model/provider abstraction and provider-specific adapters.
- MCP transport abstraction and existing stdio/SSE/Streamable HTTP support.
- Builtin tool registry in [tool/tool.go](tool/tool.go) and [tool/builtin_tool](tool/builtin_tool).
- RAG/vector and experience pipelines.
- Casbin route authorization in [authz/authz.go](authz/authz.go).
- Tri-state guard engine in [guard/casbin_guard.go](guard/casbin_guard.go), after mandatory integration.
- Prometheus instrumentation and embedded single-binary distribution.
- Existing migration, object, audit, and provider test foundations.

### Components to upgrade

- Tool execution boundary and policy evaluation.
- Shell/filesystem/browser/MCP security controls.
- Multi-user/tenant isolation and device ownership.
- Secret storage, logging redaction, audit durability, and tracing.
- Beego/XORM boundaries and global state around new functionality.
- CRA/CRACO, Router 5, Yarn 1, lint/test CI coverage.
- Docker base-image pinning, action pinning, SBOM/signing/provenance, health checks, and immutable deployment references.
- Provider timeout/retry/cancellation/usage/cost/capability contracts.

### Components to replace only with technical evidence

- Direct unrestricted host execution for production agent tools: replace with policy-mediated worker execution.
- Best-effort audit as the only compliance trail: replace or supplement with durable security-event storage.
- Mutable `latest` deployment identity: replace with digest-based deployment.
- Passwordless sudo in the application container: remove; use a dedicated worker or narrowly scoped capabilities.
- CRA/CRACO: replace after a compatibility spike proves the build can migrate without disrupting embedded assets.
- Beego: do not replace immediately; evaluate only after service boundaries and API compatibility tests exist.

## 3. 2026 AI Capability Target

### Agent runtime

Current gap: model/tool loops exist, but execution state is not a first-class durable object.

Target:

- `AgentRun`, `Task`, `ToolCall`, `ApprovalRequest`, and `ExecutionAttempt` records with owner/tenant/device scope.
- State transitions: `queued -> planning -> awaiting_approval -> executing -> verifying -> completed/failed/cancelled`.
- Context budgets, compaction, retry/recovery, idempotency keys, cancellation propagation, and resumability.
- Structured model outputs validated against schemas before execution.
- Agent handoffs and specialist agents represented as child runs with inherited but narrowed permissions.

Files likely affected: `model/`, `controllers/message*`, `controllers/task*`, `object/message*`, `object/task*`, new service packages, frontend chat/task surfaces.

Dependencies: existing model provider contract, database migration framework, approval/guard APIs, event/audit schema.

Risk: high API and persistence risk. Roll out behind existing chat behavior and maintain compatibility with current message storage.

Tests: state-machine unit tests, cancellation/retry integration tests, idempotency tests, owner isolation tests, provider contract tests, frontend state rendering tests.

Rollback: keep legacy synchronous chat path behind a feature flag until the new run state machine passes evaluation and production shadow checks.

### Human-in-the-loop

Current gap: a tri-state guard exists, but approval is not a mandatory execution boundary in the observed `model/mcp.go` path.

Target:

- Every consequential builtin/MCP operation produces a policy decision.
- `allow` executes, `deny` returns a safe structured error, `ask` creates an approval request and pauses the run.
- Approval is bound to user, tenant, device, tool, resource, arguments hash, expiry, and one-time use.
- Shell, file write/delete/move, package installation, Git push, Docker, deployment, browser evaluation, downloads, and external MCP calls default to ask or deny.

Files: [guard/casbin_guard.go](guard/casbin_guard.go), [object/tool_policy.go](object/tool_policy.go), [model/mcp.go](model/mcp.go), tool implementations, controllers, frontend approval components.

Security implication: the model never becomes the authorization authority; application code enforces the decision.

### Tool architecture

Target common contract:

```text
ToolDescriptor
  name, category, version, input schema, output schema, capabilities
ToolRequest
  run ID, task ID, user, tenant, device, resource, arguments, deadline
ToolDecision
  allow/ask/deny, matched rule, reason, expiry
ToolResult
  success, structured data, safe summary, error class, duration
AuditEvent
  actor, run, task, tool call, decision, action, outcome, resource hash
```

Each tool must validate inputs, enforce resource boundaries, propagate context cancellation, apply a timeout, return structured output, and emit redacted audit metadata.

## 4. Security Modernization

### Highest-risk current behavior

- [model/mcp.go](model/mcp.go) directly calls `BuiltinTools.ExecuteTool` and `conn.CallTool`; the guard is not visibly mandatory there.
- [tool/shell.go](tool/shell.go) executes arbitrary commands, working directories, background sessions, and PTYs.
- [tool/local_file.go](tool/local_file.go) accepts absolute paths and supports read/write/move operations without an evident configured workspace root.
- [tool/browser.go](tool/browser.go) enables arbitrary navigation and JavaScript-capable browser operations and launches Chrome with `--no-sandbox`.
- Browser/web fetch and MCP URL paths require SSRF/private-network controls.
- MCP intranet scanning intentionally probes private address ranges and requires explicit authorization.
- [routers/authz.go](routers/authz.go) policy includes many anonymous mutating API routes; this requires a route-by-route security review, not a blanket trust assumption.
- Provider/MCP credentials are persisted in models and masked in some API responses; storage encryption and lifecycle controls are not established.
- Audit is asynchronous and explicitly drops events when its queue is full in [audit/audit.go](audit/audit.go).
- Container/application hardening needs removal of passwordless sudo and mutable base images.

### Security target

- Mandatory tool guard enforcement with fail-closed defaults.
- Workspace roots and path canonicalization using secure join logic.
- Shell command policy, executable allowlists, environment filtering, resource quotas, timeout/cancellation, and isolated worker execution.
- Browser URL policy, private/link-local/metadata IP blocking, download isolation, credential partitioning, page-content trust boundaries, and confirmation for external actions.
- MCP allowlists, server identity, transport-specific authentication, tool capability declarations, health checks, per-server scopes, and audit records.
- Secret redaction at logs/audit/model context boundaries; encrypted-at-rest provider secrets or an external secret manager.
- Tenant/user/device authorization at service boundaries, not only route filters.
- Security event durability policy with backpressure and loss alerts.
- Revocable device identity and explicit cloud-to-local permission scopes.

## 5. MCP Modernization

Current state:

- [mcp/client.go](mcp/client.go) supports stdio, SSE, and Streamable HTTP with bearer headers.
- [mcp/scan.go](mcp/scan.go) scans intranet CIDRs for MCP endpoints.
- MCP tools are converted into model tool schemas and executed from [model/mcp.go](model/mcp.go).

Assessment:

- Transport interoperability exists, but server trust, capability approval, authorization, health state, and execution audit need to be first-class.
- The exact current MCP SDK/spec compatibility must be verified against the supported 2026 protocol before changing the dependency; the repository currently pins `github.com/ThinkInAIXYZ/go-mcp v0.2.24`.

Target:

- Explicit MCP server registration and approval state.
- Authenticated transport profiles and identity pinning where applicable.
- Tool schema validation and capability/policy mapping.
- Health and reconnect state with bounded timeouts.
- MCP integration tests for stdio, HTTP, auth failure, malformed schemas, denied tools, timeouts, and audit events.

Replacement is not justified until a compatibility test demonstrates the existing SDK cannot meet required protocol/auth behavior.

## 6. Model Provider Architecture

Current state:

- Many provider adapters under [model](model) and [embedding](embedding), including OpenAI-compatible, Anthropic, Gemini, Bedrock, Cohere, Hugging Face, local, and regional providers.
- The abstraction already supports streaming and tool-related behavior, but provider-specific retries, context limits, pricing, and error semantics are duplicated or manually maintained.

Target contract:

- `Capabilities()` for tools, structured output, vision, streaming, context, embeddings, and reasoning.
- `Generate(ctx, request)` with deadline/cancellation.
- Typed usage and cost metadata.
- Normalized error classes: auth, quota, rate limit, transient, invalid request, context overflow, safety refusal.
- Retry policy owned by the runtime, with provider hints and idempotency.
- Fallback selection based on capability and policy, never on silent data/tenant changes.
- Contract tests shared by every provider adapter.

Upgrade path: retain existing adapters and add the contract incrementally. Do not replace all SDKs together. Dependency upgrades require provider-specific compatibility tests and recorded API behavior.

## 7. Local PC Agent Architecture

Current state: local shell, filesystem, GUI, browser, and Docker-adjacent capabilities are present as server-side/builtin tools, but device identity, device registration, revocation, scoped grants, and cloud-to-device approval are not established as a unified subsystem.

Target:

- Device registration with asymmetric device identity and rotating credentials.
- Explicit scopes: directories, repositories, shell classes, Docker, browser profiles, network destinations.
- Device heartbeat, health, ownership, revocation, and last-seen state.
- Cloud sends signed task requests; local agent validates owner, scope, approval, expiry, and nonce before execution.
- Local agent remains the final authority for local permissions and can disconnect independently.
- No cloud feature can obtain unrestricted host control by default.

Files likely affected: new `device/` and `task/` services, auth/session models, tool registry, frontend device management, audit/event schema, deployment docs.

Risks: key lifecycle, offline behavior, replay, user consent, and backward compatibility. Start with one local agent and explicit manual approval; do not build a broad remote-control protocol first.

## 8. Multi-User Cloud Architecture

Current state:

- Casdoor authentication, sessions, owner fields, store restrictions, route Casbin roles, and masking helpers exist.
- Isolation is inconsistent at the tool/system-entity boundary; some resources intentionally use `admin`, and public route policy includes mutating operations.

Target:

- Explicit tenant and user identity on every task, run, tool call, resource, device, secret, and audit event.
- Central service authorization with deny-by-default tenant checks.
- Row-level ownership tests and negative tests for every resource family.
- Per-tenant rate, cost, storage, concurrency, and model quotas.
- Device ownership and credential scope tied to tenant/user.

Database migration: add tenant/user/device foreign-key or serialized identity fields to new run/task/event tables; backfill existing owner fields; preserve compatibility for system/admin resources with explicit system tenancy.

## 9. Memory and Context

Current state:

- Per-store `MemoryLimit`, recent message retrieval, vector/RAG retrieval, and experience/correction data exist.
- A durable, user-controlled, provenance-aware memory lifecycle is not established.

Target layers:

1. Current conversation: bounded and user-visible.
2. Task state: structured, resumable, and separate from chat text.
3. Project/workspace context: repository and files with scope/provenance.
4. User preferences: explicit opt-in and editable.
5. Long-term memory: typed records with retention, deletion, consent, and provenance.
6. Retrieved documents/tool results: source, timestamp, trust, and sensitivity metadata.

Required controls: compaction, relevance retrieval, sensitive-data removal, retention/deletion, provenance, tenant isolation, and context-budget accounting.

## 10. Browser and Computer Use

Current state:

- Chromedp browser navigation, screenshots, evaluation, clicks, downloads/browser-use capabilities exist.
- Browser launches with `--no-sandbox`; URL and page-origin restrictions are not visibly centralized.

Target:

- Browser worker isolation with a disposable profile per run.
- Network egress policy and SSRF blocking.
- Download quarantine and malware/content scanning.
- Page content marked untrusted; prompt injection cannot alter authorization.
- Browser evaluate/click/submit/external side effects require explicit policy and approval.
- Credential access restricted to named domains and user-approved profiles.
- Screenshots and extracted content redacted or retained under data policy.

Tests: malicious page fixtures, private-IP redirects, credential exfiltration prompts, download restrictions, approval transitions, and browser worker cleanup.

## 11. Observability and Evaluation

### Observability target

Current state: Prometheus metrics and JSONL audit exist, but there is no clearly integrated OpenTelemetry trace context across HTTP -> task -> model -> tool -> MCP, and audit can drop events.

Target fields:

- request ID, run ID, task ID, tool-call ID, user/tenant/device IDs, model/provider, attempt, policy decision, latency, token usage, cost, outcome, and error class.
- Structured logs with secrets and sensitive content redacted.
- OpenTelemetry traces and metrics with sampling and retention controls.
- Durable security-event path separate from best-effort operational logs.
- SLOs and alerts for availability, tool failures, approval latency, model latency, queue depth, dropped audit events, and cost spikes.

### Evaluation target

Add repeatable fixtures and scorecards for:

- Tool selection and argument validity.
- Structured output adherence.
- Multi-step coding and test tasks.
- Retry/recovery and cancellation.
- Prompt injection resistance.
- Shell/filesystem policy enforcement.
- MCP identity and permission enforcement.
- Cross-user isolation.
- Browser malicious-content handling.
- Regression across provider/model versions.

Evaluations must separate model quality from authorization correctness. A model choosing the wrong tool is a quality failure; a denied tool executing is a security failure.

## 12. Dependency Modernization Plan

The repository currently pins many dependencies, but a current-version claim requires network-backed vulnerability and compatibility checks that were not available in this audit. The baseline must be recorded before upgrades.

### Keep and audit first

- Go 1.25.8 toolchain: keep aligned across `go.mod`, CI, Dockerfile, and RISC-V build. The RISC-V Dockerfile currently uses Go 1.23.11 and must be aligned or explicitly isolated.
- Casbin 2.135.0, Prometheus client 1.15.0, OAuth/Casdoor, database drivers, and provider SDKs: upgrade only after tests and security scans.
- MCP SDK `github.com/ThinkInAIXYZ/go-mcp v0.2.24`: verify current protocol support before changing.

### Technical debt to migrate

- Beego 1.12.12 and XORM 1.2.5: supported short-term by compatibility boundary; evaluate upgrade/replacement only with API and migration tests.
- Chromedp 0.9.5 and its pinned cdproto snapshot: browser security and Chrome compatibility review required.
- Mixed provider SDKs and `replace` directives: generate a dependency/SBOM inventory and contract-test adapters.
- Frontend CRA 5, CRACO 6, Router 5, ESLint 8, Stylelint 14, Yarn 1: staged migration to Vite, Router 6+, current lint/test tooling.

Required supply-chain work:

- `go mod verify`, `govulncheck`, OSV/dependency scanner, npm/Yarn audit equivalent, SBOM, license report, lockfile integrity, action SHA pinning, signed release artifacts, and provenance attestations.
- No blanket latest-version upgrade. Every upgrade records current version, target version, reason, compatibility result, breaking changes, and rollback.

## 13. Database and API Migration Strategy

Database:

- Use additive migrations for agent runs, task state, approvals, devices, tool-call audit metadata, tenant identity, and secret references.
- Backfill ownership/tenant fields before enforcing non-null constraints.
- Add indexes for owner/tenant/run/status/time queries.
- Keep migration rollback or forward-fix procedures and test against MySQL, PostgreSQL, and SQLite where supported.
- Do not store raw credentials or unbounded model/tool payloads in new tables.

API:

- Preserve existing endpoints during migration.
- Add versioned run/task/approval APIs rather than changing chat semantics in place.
- Define typed error classes and idempotency behavior.
- Deprecate endpoints only after frontend and external API compatibility tests pass.
- Treat public mutating routes and OpenAI-compatible endpoints as explicit security review surfaces.

## 14. Frontend Modernization Plan

Phase 1: retain Ant Design and current pages; add agent run status, tool-call summary, approval, cancellation, retry, and verification components using existing API conventions.

Phase 2: add device/project/workspace and security-event views.

Phase 3: migrate build/test tooling from CRA/CRACO to Vite or an equivalent maintained tool, then Router 5 to Router 6+.

Phase 4: introduce TypeScript for new task/device/security surfaces and progressively type shared API models.

Accessibility, responsive behavior, loading/error/empty states, streamed status updates, and no private chain-of-thought display are acceptance requirements.

## 15. Docker, CI, and Deployment Strategy

Current risks:

- Mutable `latest` base images in [Dockerfile](Dockerfile) and [riscv64.Dockerfile](riscv64.Dockerfile).
- RISC-V Go version mismatch.
- Older/mixed GitHub Action major versions.
- Frontend build is gated, but frontend lint/tests are not clearly mandatory.
- Compose exposes MySQL port 3306 and lacks a database health condition.
- Docker images publish mutable `latest` tags; recorded verification could not confirm repository ownership or deployed digest.
- OCI provenance labels now exist in the current worktree, but no deployed image evidence exists.

Target release chain:

```text
approved Git SHA
  -> reproducible CI build
  -> signed binary/image/SBOM/provenance
  -> immutable image digest
  -> Helm/Compose deployment by digest
  -> runtime version/commit and digest check
```

Required controls: pinned base/action digests, non-root runtime, no passwordless sudo, health checks, resource limits, private database network, secret injection, image signing, SBOM/vulnerability gates, deployment admission by digest, rollback, and migration compatibility checks.

## 16. Phased Roadmap

### Phase A: P0 security gate

- Wire guard and approval into every builtin/MCP execution.
- Enforce workspace/path, shell, browser, network, and MCP policies.
- Remove sensitive logging and add redaction tests.
- Define tenant/user scope and cross-user negative tests.
- Remove passwordless sudo, pin runtime images, and close production secret/config gaps.

Exit criteria: denied/ask operations cannot execute; security regression suite passes; no raw credentials appear in logs/audit/model context.

### Phase B: P1 production control plane

- Introduce durable agent runs, tasks, approvals, cancellation, retries, and audit events.
- Add provider capability/usage/error contracts.
- Add OpenTelemetry correlation and durable security-event policy.
- Add device registration/scopes/revocation foundation.
- Make CI run frontend lint/tests, security scans, SBOM, and provenance checks.

Exit criteria: a task can be reconstructed and safely resumed/cancelled, with tenant/device ownership enforced.

### Phase C: P2 maintainability and UX

- Add service/repository boundaries around new work.
- Migrate frontend build/test tooling in a compatibility branch.
- Add context/memory lifecycle and provenance.
- Add browser worker isolation and evaluation fixtures.
- Add provider contract/evaluation matrix.

Exit criteria: new features do not require controller/object cross-coupling; UI exposes safe execution status and approvals.

### Phase D: P3 advanced platform capabilities

- Specialist-agent handoffs.
- Multi-device task routing.
- Durable long-term memory with user controls.
- Advanced evaluation dashboards and model routing.
- Optional Beego/XORM replacement only if measured maintenance/security cost justifies it.

## 17. Global Risk and Rollback Rules

- Do not combine backend framework, database, frontend build, and agent-runtime rewrites in one release.
- Gate every P0/P1 change with negative security tests.
- Use feature flags for new agent execution and task APIs.
- Keep database migrations additive and reversible or forward-fixable.
- Preserve legacy API responses until compatibility clients are migrated.
- Keep the existing provider/tool implementations behind adapters while the control plane changes.
- Roll back by disabling the new execution path, deploying the prior immutable image digest, and applying only tested database forward fixes.

## 18. Definition of Modernization Done

The repository is ready for continued 2026 platform development when:

- Every consequential tool call is policy-authorized and, when required, approved by a human.
- User, tenant, device, task, run, and secret ownership is enforced in service code.
- Local execution is isolated, scoped, revocable, and never cloud-unrestricted.
- MCP servers are registered, authenticated, health-checked, permissioned, and audited.
- Providers share capability, retry, usage, cost, and cancellation contracts.
- Context and memory have provenance, retention, deletion, and sensitivity controls.
- Logs/metrics/traces reconstruct safe execution without secrets.
- Frontend tests/lint and backend/security/integration/evaluation suites gate CI.
- Images and deployments are signed, scanned, reproducible, and addressed by digest.
- Documentation describes implemented behavior rather than aspirational capabilities.
