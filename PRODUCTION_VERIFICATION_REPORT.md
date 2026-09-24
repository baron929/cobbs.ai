# cobbs.ai Production Verification Report

Date: 2026-09-24
Scope: production verification and approved foundation remediation. No new features or broad architectural changes were added.

## Executive Summary

Production verification remains partially blocked. The configured production domain does not resolve, the configured Docker Hub repositories return `404 Not Found`, and no local backend is listening on port `14000`. The Go toolchain is available in the current environment and the required targeted backend verification now passes.

The repository-level configuration risks have been remediated: `conf/app.conf` now uses `runmode = prod`, production Compose requires injected authentication and database settings, and database credentials are no longer embedded in source configuration. Runtime deployment and external ownership findings remain unverified.

No deployed image digest, deployed health response, HTTPS certificate, live CORS behavior, live authentication flow, or deployed-commit match can be claimed as verified from this environment.

## Remediation Implemented

The approved foundation changes are now present:

- [conf/app.conf](conf/app.conf) uses production mode and contains no database password.
- [conf/conf.go](conf/conf.go) validates required production database and authentication settings.
- [main.go](main.go) fails before database initialization when production settings are missing.
- [docker-compose.yml](docker-compose.yml) requires database and Casdoor values from the environment.
- [docker-entrypoint.sh](docker-entrypoint.sh) fails when `MYSQL_ROOT_PASSWORD` is absent.
- [.gitignore](.gitignore) ignores local environment and secret files.
- [.env.example](.env.example) contains placeholders only.
- CI uses a run-scoped database password instead of a fixed committed password.
- [Dockerfile](Dockerfile) and [riscv64.Dockerfile](riscv64.Dockerfile) include OCI source, version, revision, and build-time labels.
- The prior report and contributor example no longer disclose the fixed credential.

Validation evidence:

- `docker compose config -q` passes with synthetic placeholder values.
- `docker compose config -q` fails with exit code `1` when required values are absent.
- The credential scan reports no non-placeholder database credential matches.
- Secret-path ignore checks pass with `git check-ignore --no-index`.
- `git diff --check` passes.

The Docker daemon and Casdoor service are unavailable locally. The required targeted Go tests and repository build were executed successfully; Casdoor-backed tests reported explicit skips because their external service and fixtures were unavailable.

## 1. Live Domain: `cobbs.aiai.org`

**Files/configuration involved**

- [README.md](README.md) and [README_zh.md](README_zh.md): advertised production URLs.
- [conf/conf.go](conf/conf.go): default static asset/domain configuration.
- [routers/router.go](routers/router.go): documented external API URL.

**Exact procedure**

```powershell
Resolve-DnsName cobbs.aiai.org
Invoke-WebRequest -Uri https://cobbs.aiai.org -Method Head -MaximumRedirection 5
```

**Expected result**

DNS returns one or more valid records, and HTTPS returns a successful response, normally `200`, with a valid certificate and an expected final URL.

**Actual result**

```text
DNS_ERROR_RCODE_NAME_ERROR
HTTPS_ERROR=The remote name could not be resolved: 'cobbs.aiai.org'
```

**Status: FAILED**

**Root cause**: the domain has no resolvable DNS record from the current network. Deployment, DNS, or the advertised domain name must be checked before application-level verification can proceed.

## 2. Docker Image Ownership, Tags, and Digests

**Files/configuration involved**

- [.github/workflows/build.yml](.github/workflows/build.yml): publishes `casbin/cobbs.ai` and `casbin/cobbs.ai-all-in-one`.
- [Dockerfile](Dockerfile) and [riscv64.Dockerfile](riscv64.Dockerfile): image build targets.
- [docker-compose.yml](docker-compose.yml): local image/build configuration.

**Exact procedure**

```powershell
$repos = @('casbin/cobbs.ai','casbin/cobbs.ai-all-in-one')
foreach ($repo in $repos) {
  Invoke-RestMethod "https://hub.docker.com/v2/repositories/$repo/tags?page_size=100"
}
```

For an authenticated or deployment-capable environment, verify multi-architecture manifests and digests with:

```bash
docker buildx imagetools inspect casbin/cobbs.ai:<tag>
docker buildx imagetools inspect casbin/cobbs.ai-all-in-one:<tag>
docker image inspect casbin/cobbs.ai:<tag> --format '{{json .RepoDigests}}'
```

**Expected result**

The repositories exist under the intended owner, the expected release and `latest` tags exist, and the deployed runtime digest exactly matches the approved immutable digest.

**Actual result**

Docker Hub API returned:

```text
casbin/cobbs.ai: 404 Not Found
casbin/cobbs.ai-all-in-one: 404 Not Found
```

The lower-level registry request also returned `401 Unauthorized`.

**Status: FAILED / BLOCKED**

**Root cause**: either the repositories do not exist under `casbin`, are private, or require credentials not available here. Ownership and deployed digest cannot be verified until the registry owner supplies authenticated access or confirms the correct repository names.

## 3. Complete Backend and Frontend Together

**Files/configuration involved**

- [.github/workflows/build.yml](.github/workflows/build.yml): CI Go tests, backend build, and frontend build.
- [go.mod](go.mod): Go module and toolchain requirements.
- [web/package.json](web/package.json): frontend scripts.
- [embed.go](embed.go): embedded frontend/runtime assets.
- [conf/app.conf](conf/app.conf): backend port and runtime configuration.

**Exact repository procedure**

```powershell
go version
go test ./...
Push-Location web
yarn install --frozen-lockfile
yarn lint:js
yarn build
Pop-Location
go run main.go
```

Then, from another terminal:

```powershell
Invoke-WebRequest http://localhost:14000/api/health
```

**Expected result**

Go tests pass, frontend lint and build pass, the backend starts, and `/api/health` returns HTTP `200`.

**Actual result**

- `go` is not installed or not on `PATH`; `go version` failed with `The term 'go' is not recognized`.
- `web/node_modules` was absent.
- Frontend dependency installation was attempted but remained network-bound while fetching approximately `2121` packages and was stopped.
- No local backend was running; `http://localhost:14000/api/health` returned `Unable to connect to the remote server`.

**Status: BLOCKED**

**Root cause**: missing Go toolchain and incomplete frontend dependency installation. This is an environment limitation, not evidence that the code passes or fails.

## 4. Frontend-to-Backend API Communication

**Files/configuration involved**

- [web/src/backend/FetchFilter.js](web/src/backend/FetchFilter.js): frontend fetch interception and API handling.
- [web/src/Conf.js](web/src/Conf.js): runtime API/static configuration.
- [routers/router.go](routers/router.go): backend API routes.
- [routers/cors_filter.go](routers/cors_filter.go): CORS response and origin checks.

**Exact procedure after deployment**

```powershell
$base = 'https://cobbs.aiai.org'
Invoke-WebRequest "$base/api/health" -Method Get
Invoke-WebRequest "$base/api/get-signin-options" -Method Get
```

In browser DevTools, load the frontend and verify that API requests target the same production origin, return expected JSON, and do not show failed preflight requests or mixed-content errors.

For a CORS preflight:

```bash
curl -i -X OPTIONS https://cobbs.aiai.org/api/get-signin-options \
  -H 'Origin: https://cobbs.aiai.org' \
  -H 'Access-Control-Request-Method: GET'
```

**Expected result**

The frontend loads, `/api/health` and `/api/get-signin-options` return successful JSON responses, and allowed-origin preflight responses include the expected CORS headers.

**Actual result**

Not executable against production because `cobbs.aiai.org` does not resolve. Local verification also failed because no backend was listening on port `14000`.

**Status: BLOCKED**

**Root cause**: upstream domain and local runtime are unavailable. The source code has a same-origin API path and a CORS filter, but source inspection is not proof of deployed communication.

## 5. Production Environment Variables

**Files/configuration involved**

- [conf/app.conf](conf/app.conf)
- [docker-compose.yml](docker-compose.yml)
- [.github/workflows/build.yml](.github/workflows/build.yml)
- [web/src/Conf.js](web/src/Conf.js)

**Exact procedure**

On the deployment host, inspect non-secret names and validate secret presence without printing values:

```bash
docker compose config
printenv | cut -d= -f1 | sort
for name in CASDOOR_ENDPOINT CASDOOR_CLIENT_ID CASDOOR_CLIENT_SECRET JWT_PUBLIC_KEY MYSQL_ROOT_PASSWORD; do
  test -n "${!name:-}" && echo "$name=SET" || echo "$name=MISSING"
done
```

**Expected result**

Production uses non-development mode, secrets come from a secret manager or protected environment, authentication settings are populated, database credentials are not committed plaintext, and frontend runtime configuration points to the deployed API/static origins.

**Actual repository result after remediation**

- [conf/app.conf](conf/app.conf) now contains `runmode = prod` and an empty DSN placeholder.
- Production startup validation requires issuer/endpoint, client ID, client secret, organization, application, and database settings.
- [docker-compose.yml](docker-compose.yml) requires those values through environment interpolation.
- [docker-entrypoint.sh](docker-entrypoint.sh) fails when `MYSQL_ROOT_PASSWORD` is missing.
- [web/src/Conf.js](web/src/Conf.js) still defaults frontend runtime values to empty strings; production runtime injection remains an external deployment responsibility.

**Status: REPOSITORY FOUNDATION REMEDIATED; DEPLOYMENT UNVERIFIED**

**Remaining root cause**: the actual deployment environment was not accessible, so runtime overrides, secret-manager integration, and live authentication settings remain unverified.

## 6. HTTPS, CORS, Authentication, and Health Checks

**Files/configuration involved**

- [routers/cors_filter.go](routers/cors_filter.go): CORS origin validation and preflight headers.
- [routers/authz_filter.go](routers/authz_filter.go): API authorization filter.
- [controllers/system_info.go](controllers/system_info.go): `/api/health` and version endpoints.
- [main.go](main.go): filter registration, including CORS and HSTS.
- [auth/auth.go](auth/auth.go) and [controllers/account.go](controllers/account.go): Casdoor authentication setup.

**Exact procedure**

```bash
curl -Iv https://cobbs.aiai.org/api/health
curl -i https://cobbs.aiai.org/api/health
curl -i -X OPTIONS https://cobbs.aiai.org/api/get-account \
  -H 'Origin: https://cobbs.aiai.org' \
  -H 'Access-Control-Request-Method: GET' \
  -H 'Access-Control-Request-Headers: Authorization'
curl -i https://cobbs.aiai.org/api/get-account
```

Verify that:

- TLS certificate hostname matches the production domain.
- HTTP redirects to HTTPS.
- `/api/health` is `200` without authentication.
- protected endpoints reject unauthenticated access with the expected `401`/authorization response.
- allowed CORS origins receive explicit headers and disallowed origins are rejected.

**Expected result**

All four controls behave as above in the deployed environment.

**Actual result**

HTTPS and all application-level checks are blocked because DNS resolution fails. No deployed certificate, CORS response, authentication response, or health response was observed.

**Status: BLOCKED**

## 7. Full Test Suite

**Files/configuration involved**

- [.github/workflows/build.yml](.github/workflows/build.yml): authoritative CI test command.
- Go `*_test.go` files throughout the repository.
- [web/package.json](web/package.json): frontend test script.

**Exact commands**

Backend, matching CI:

```bash
go test -v -vet=asmdecl,assign,atomic,bools,buildtag,cgocall,composites,copylocks,defers,directive,errorsas,framepointer,httpresponse,ifaceassert,loopclosure,lostcancel,nilfunc,shift,sigchanyzer,slog,stdmethods,stringintconv,structtag,testinggoroutine,tests,timeformat,unmarshal,unreachable,unsafeptr,unusedresult $(go list ./...) -tags skipCi
```

Frontend:

```bash
cd web
yarn install --frozen-lockfile
CI=true yarn test --watchAll=false
yarn lint
CI=false yarn build
```

**Expected result**

Every backend test, frontend test, lint check, and production build exits with code `0`.

**Actual result**

The required repository-level backend gate completed successfully:

```text
go test ./model ./txt ./object ./audio ./storage ./split ./audit  PASS
go build ./...                                                    PASS
git diff --check                                                   PASS
```

The focused Casdoor migration test also exits successfully but is explicitly skipped when Casdoor is unavailable. The full `go test ./...`, `go vet ./...`, race test, and frontend checks were not claimed as complete in this targeted verification pass.

**Status: PARTIALLY VERIFIED**

**Remaining limitation**: live Casdoor integration, full backend-suite coverage, race testing, frontend dependencies, Docker runtime checks, and deployed verification remain environment-dependent.

## 8. Deployed Version Versus Intended Git Commit

**Files/configuration involved**

- [internal/cli/version.go](internal/cli/version.go): build version/commit metadata.
- [.goreleaser.yaml](.goreleaser.yaml): injects version and commit at release build time.
- [.github/workflows/build.yml](.github/workflows/build.yml): passes `COMMIT=${{ github.sha }}` into Docker builds.

**Exact repository procedure**

```powershell
git rev-parse HEAD
git status --short
git log -1 --format='%H %s'
```

Exact deployed-image procedure:

```bash
docker pull casbin/cobbs.ai:<approved-tag>
docker inspect --format '{{index .Config.Labels "org.opencontainers.image.revision"}}' casbin/cobbs.ai:<approved-tag>
docker run --rm casbin/cobbs.ai:<approved-tag> version
curl -s https://cobbs.aiai.org/api/get-version-info
```

Compare all returned commit IDs to the approved immutable Git SHA. Also compare the running container's `RepoDigests` with the digest recorded during release approval.

**Actual repository result**

The current repository `HEAD` is:

```text
726a263d66e2442b1323ea70e5793533f89e1039
```

The worktree is dirty with the rebrand and verification-related changes. No deployed container or live version endpoint was reachable, so no deployed commit or digest could be compared.

**Status: BLOCKED**

**Root cause**: production domain and Docker repositories are unavailable/unverified, and the current worktree is not a clean release commit.

## Approval Gate

Do not fix or promote anything based only on this report. Approval should first confirm:

1. The canonical production domain and DNS ownership.
2. The correct Docker Hub owner/repositories and registry credentials.
3. The approved release Git SHA.
4. The production secret/environment management mechanism.
5. Access to the deployment host or CI environment needed to run the blocked checks.

Once those are confirmed, rerun every blocked procedure and attach the raw command output, image digests, HTTP responses, and test logs to the release record.
