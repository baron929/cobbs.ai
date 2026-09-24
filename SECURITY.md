# Security Policy

## Reporting a Vulnerability

We appreciate security researchers and users who report vulnerabilities responsibly. To ensure your report is handled in a timely manner and we can keep users safe, please follow the guidelines below.

**Please do NOT report security vulnerabilities through public GitHub issues.**

Instead, please report them by sending an email to [baron@cobbs.ai](mailto:baron@cobbs.ai).

We will endeavor to respond as quickly as possible and work with you to understand and resolve the issue promptly.

## Current Security Boundaries

The server enforces API authorization through the Casbin-backed filter in `routers/authz_filter.go`, with controller-level signed-in, admin, global-admin, and store-admin checks. Tool authorization is centralized in `toolauth/authorization.go` and is invoked by the MCP execution path in `model/mcp.go`. Missing authorization policy fails closed; unknown and consequential tools require approval or are denied, and tool execution does not trust the model's decision.

Approval grants are stored in the `ApprovalRecord` table. Raw approval tokens are returned only to the authenticated grant caller and are not persisted. The persisted token hash is unique, scope-bound, expiry-bound, and consumed with a conditional database update that requires exactly one affected row. Consumption occurs before tool execution, so failed or timed-out executions require a new approval. The implementation is covered by focused unit tests, but database-specific concurrency behavior still requires supported-database integration verification.

The approval endpoint derives the grant owner and subject from the authenticated session. Client-supplied owner, subject, and device identities are not accepted. Device identity is currently unavailable in the observed execution path and is therefore not treated as verified.

## Known Limitations

- Local shell, filesystem, browser, and process tools are application-authorized but are not fully isolated by a container, OS sandbox, seccomp policy, or dedicated worker boundary. Run the service under a least-privilege operating-system account and restrict its filesystem and network access externally.
- Approval persistence and concurrent-consumption integration tests require a supported database. The targeted object package currently compiles and passes; database-specific concurrency behavior remains unverified in this environment.
- Audit logging is a best-effort JSONL sidecar. It records safe metadata and never gates authorization. Use filesystem permissions, rotation, retention, and centralized collection appropriate to the deployment.
- Application auto-sign-in requires Basic or Bearer credentials in headers. Query-string credentials are rejected because URLs are routinely logged and retained.
- The application bearer digest now uses SHA-256 instead of the former MD5-derived value; clients using the legacy derived token must migrate to the new header token.
- Prometheus metric names use underscore-separated names because dots are invalid in Prometheus identifiers; dashboards using the old names must be updated.

## Verification Status (2026-09-24)

The following repository checks currently pass:

```text
go test ./model ./txt ./object ./audio ./storage ./split ./audit
go build ./...
git diff --check
```

Casdoor-backed migration and email tests skip when Casdoor is not configured or reachable. Those skips prevent nil SDK-client panics, but they are not evidence of a live authentication or email integration pass. The full backend suite, frontend checks, deployed HTTPS/CORS behavior, OS-level tool isolation, and database-specific approval contention remain separate verification tasks.

## Verification Commands

Focused authorization tests:

```text
go test ./toolauth
```

Production package compilation and broader checks should be run before deployment:

```text
go build ./...
go test ./...
go vet ./...
go test -race ./...
```

Do not treat a green focused authorization test as evidence that the operating-system worker boundary, all route object authorization, or database-specific concurrency behavior has been fully verified.
