# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Development Commands

### Building and Testing
- **Build the project**: `make` or `make cms-build`
- **Run tests**: `go test ./...`
- **Format code**: `make goimports` (uses `go tool goimports -w -l .`)

### Site Generation
- **Build site (production)**: `make site-build-release` or `go run ./cmd/site/main.go`
- **Build site with drafts**: `make site-build-draft` or `go run ./cmd/site/main.go -d`
- **Watch for changes**: `make site-watch` (uses `entr` to rebuild on file changes)
- **Clean build directory**: `make site-clean`
- **Serve locally**: `make serve` (serves on port 3000)
- **Deploy to production**: `make site-deploy`

### Content Management System (CMS)
- **Build CMS**: `make cms-build` (creates `./cms` binary)
- **Run CMS**: `./cms --help` for available commands
- **CMS operations**: 
  - `./cms entries [create|update|list|delete]`
  - `./cms categories [create|update|delete]`
  - `./cms tags [create|update|delete]`
  - `./cms authors [create|update|delete]`

### Database Management
- **Apply migrations**: `source migrations/goose.env && goose up`
- **Create migration**: `goose create <migration_name> sql`
- **Database location**: `./database/content.sqlite3`

### Debugging
- **Debug CMS**: `make cms-debug` (uses `dlv`)
- **Debug site generator**: `make site-debug`

## Architecture Overview

This is a Go-based blog platform with two main applications:

### 1. Static Site Generator (`cmd/site/`)
- Generates static HTML from SQLite database content
- Uses Go templates in `templates/` directory
- Outputs to `dist/` directory
- Supports draft mode with `-d` flag
- Uses Goldmark for Markdown rendering

### 2. Terminal-based CMS (`cmd/cms/`)
- TUI application built with `tview` for managing content
- Direct SQLite database interaction
- Supports CRUD operations for entries, categories, tags, and authors
- Interactive forms and selection menus

### Core Components

#### Data Layer (`internal/store/`)
- SQLite database connection and queries
- Database location: `./database/content.sqlite3`
- Uses `sqlx` for database operations
- Foreign keys enabled

#### Models (`internal/model/`)
- Core data structures: Entry, Category, Tag, Author, Stats
- Represents database schema in Go structs

#### Site Generation (`internal/site/`)
- `MySite` struct contains all site configuration and data
- Template rendering and static file generation
- Asset copying from `assets/` to `dist/`

#### Template System (`internal/templates/`)
- HTML template loading and parsing
- Template functions and helpers

#### Utilities (`internal/util/`)
- Helper functions for slugification, truncation, etc.

### Database Schema
- Entries: blog posts with markdown content, categories, tags, and authors
- Categories: hierarchical content organization
- Tags: flexible content labeling
- Authors: multi-author support
- Uses Goose for migrations in `migrations/` directory

### Frontend Assets
- CSS in `assets/css/`: base styles, reset, and custom styles
- JavaScript in `assets/js/`: theme switcher, hamburger menu, marquee animations
- Images and fonts in respective `assets/` subdirectories

## Development Workflow

1. Use `make cms-build && ./cms` to manage content
2. Use `make site-build-draft` to generate site with drafts
3. Use `make serve` to preview locally
4. Use `make site-build-release` for production builds
5. Use `make site-deploy` to deploy to production server

## Key Dependencies
- `tview`: Terminal UI framework for CMS
- `sqlx`: SQL extensions for Go
- `goldmark`: Markdown parser
- `urfave/cli/v3`: CLI framework for CMS commands
- `goose`: Database migration tool