# Blog Platform

A self-contained blog platform with a terminal-based CMS and static site generator. Built with Go, it uses SQLite for content storage and provides tools for managing blog content through a command-line interface.

## Project Structure

- `cmd/` - Command-line applications
  - `cms/` - Terminal-based content management system
  - `site/` - Static site generator
- `internal/` - Core application logic
  - `markdown/` - Markdown processing
  - `model/` - Data models (entry, category, tag, stats)
  - `site/` - Site generation logic
  - `store/` - Database access layer
  - `templates/` - Template rendering
  - `util/` - Utility functions
- `migrations/` - Database migration files using Goose
- `templates/` - HTML templates for rendering the site
- `assets/` - Static assets
  - `css/` - Stylesheets (base, reset, styles)
  - `js/` - JavaScript (hamburger menu, marquee animation, theme switcher)
  - `fonts/` - Web fonts
  - `images/` - Image assets
- `database/` - SQLite database file
- `dist/` - Generated static site output (created on build)
- `logs/` - Build logs

## Features

- **Content Management**: Create, update, delete, and list blog entries, categories, tags, and authors
- **Markdown Support**: Content is written in Markdown and rendered to HTML
- **Static Site Generation**: Generates a complete static site from database content
- **Tagging System**: Organize content with tags
- **Author Management**: Support for multiple authors
- **Draft Mode**: Preview draft content before publishing
- **Database Migrations**: Structured database schema changes using Goose
- **Theme Switching**: Built-in theme switcher with CSS variables
- **Responsive Design**: Mobile-friendly with hamburger menu
- **Hot Reload**: Watch mode for development with automatic rebuilds
- **Local Development Server**: Built-in Python server for testing

## Prerequisites

- Go 1.19 or higher
- SQLite3
- Python 3 (for local development server)
- [Goose](https://pressly.github.io/goose/installation/) (for database migrations) - included in go.mod as github.com/pressly/goose/v3
- entr (optional, for watch mode)

## Configuration

The blog is configured via `config.toml`:

```toml
title = "Your Blog Title"
domain = "yourdomain.com"
author = "Your Name"
dir_build = "dist"
dir_assets = "assets"
```

## Usage

### Build the Project

```bash
# Build the CMS binary
make cms-build

# Or build everything
make
```

### Content Management System (CMS)

```bash
# Show CMS help
./cms --help

# Manage entries
./cms entries [create|update|list|delete]

# Manage categories
./cms categories [create|update|delete]

# Manage tags
./cms tags [create|update|delete]

# Manage authors
./cms authors [create|update|delete]
```

### Site Generation

```bash
# Build the site (production - excludes drafts)
make site-build-release

# Build the site with drafts enabled
make site-build-draft

# Clean the dist directory
make site-clean

# Watch for changes and rebuild automatically
make site-watch

# Deploy to production (requires rsync and SSH access)
make site-deploy

# Start local development server (port 3000)
make serve
```

### Database Migrations

1. [Install Goose](https://pressly.github.io/goose/installation/)
2. Set environment variables:
   ```bash
   source migrations/goose.env
   ```
3. Create a new migration:
   ```bash
   goose create <migration_name> sql
   ```
4. Apply migrations:
   ```bash
   goose up
   ```

## Development

### Commands

- **Run tests**: `go test ./...`
- **Format code**: `make goimports`
- **Debug CMS**: `make cms-debug`
- **Debug Site Generator**: `make site-debug`
- **List all commands**: `make list-commands`

### Development Workflow

1. Start the watch mode: `make site-watch`
2. In another terminal, start the dev server: `make serve`
3. Make changes to templates, assets, or Go code
4. The site will automatically rebuild on file changes
5. Refresh your browser to see updates

### Go Tools

The project includes several Go tools in its go.mod file:
- `goose` - Database migrations (github.com/pressly/goose/v3)
- `goimports` - Code formatting (golang.org/x/tools/cmd/goimports)
- `dlv` - Debugger (github.com/go-delve/delve/cmd/dlv)

### Project Dependencies

The project uses standard Go modules. Key dependencies include:
- SQLite driver for database operations
- Markdown parser for content rendering
- Template engine for HTML generation

Run `go mod download` to fetch all dependencies.

## Deployment

The Makefile includes a deployment target that uses rsync to deploy the built site:

```bash
# Deploy to production (configure DEPLOY_HOST in Makefile)
make site-deploy
```

By default, it deploys to `prod@kuljis.xyz:/home/prod/www/kuljis.xyz/dist/`. Update the `DEPLOY_HOST` variable in the Makefile to match your deployment target.

## License

This project is licensed under the MIT License.
