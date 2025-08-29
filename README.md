# Blog Site

A small Go-powered static site generator for my blog, with an optional CMS helper and SQLite-backed content/migrations.

## Quick start

```bash
# Prereqs: Go 1.24+, SQLite, entr (for watch)
make site-build-draft   # build draft site to ./dist
make serve              # serve ./dist at http://localhost:3000
```

## Common tasks

```bash
make site-build-release  # production build
make site-watch          # rebuild on changes (assets/, templates/, cmd/, internal/)
make cms-build           # build CMS binary to ./cms
make site-debug          # debug site builder with dlv
make cms-debug           # debug CMS with dlv
```

## Lint/format/tests

```bash
make goimports        # format + organize imports
go vet ./...          # static checks
go test ./...         # run all tests
# single test by regex
go test ./path/to/pkg -run '^TestName$'
```

## Migrations

```bash
# Apply SQLite migrations to ./database/content.sqlite3
# (reads env from migrations/goose.env)
. ./migrations/goose.env; go tool goose up
```

## Project layout

- cmd/site: site builder entrypoint
- cmd/cms: CMS CLI entrypoint
- internal: application packages (no external imports)
- templates: HTML templates
- assets: CSS/JS/images; build outputs to ./dist
- migrations: SQL migrations managed by goose

## Config

- config.json: site configuration

## License

MIT (see repository for details).
