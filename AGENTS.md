Agent quickstart for this repo

Build, run, watch
- Prereqs: Go 1.24+, SQLite. Tools auto-installed via `go mod` tools section.
- Dev site (draft): `make site-build-draft` (writes dist/)
- Release build: `make site-build-release`
- Debug: `make site-debug` (dlv) or `make cms-debug`
- Watch rebuild: `make site-watch` (requires entr)
- Static serve built site: `make serve`
- CMS binary: `make cms-build` (outputs ./cms)

Lint/format/tests
- Format/imports: `make goimports` (golang.org/x/tools/cmd/goimports). Keep std imports first, third‑party next, local last; grouped with blank lines.
- Vet/static checks: `go vet ./...`
- Tests (all): `go test ./...`
- Single test file: `go test ./internal/site -run TestName`
- Single test (regex): `go test ./path/to/pkg -run '^TestExact$'`

Code style
- Language: Go. Use context propagation, avoid globals; prefer dependency injection via funcs/structs in `internal/...`.
- Errors: return wrapped errors with context (`fmt.Errorf("...: %w", err)`), handle at boundaries; no panics in library code; log minimally in CLI entrypoints under `cmd/...`.
- Types: prefer explicit types; zero‑value safe structs; avoid interface{}—define small interfaces at consumer side.
- Naming: packages lower_snake? (Go convention: all lowercase, short); exported symbols have doc comments; constructors `NewType`.
- Imports: no side‑effects; keep internal packages under `internal/...`, shared under `pkg/...`.
- Formatting: run goimports; keep 120 col soft limit.

Migrations/content
- DB migrations via goose (tool in go.mod). Env: `migrations/goose.env`. Apply: `goose -dir migrations sqlite3 database/content.sqlite3 up`.

Notes
- No Cursor/Copilot rules present. If added later (`.cursor/rules/` or `.github/copilot-instructions.md`), mirror here.
- Config lives in `config.json`; builder entrypoint `cmd/site/main.go`. Keep dist/ as build artifact; don’t commit.
