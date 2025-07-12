# Blog Platform

A self-contained blog platform with a terminal-based CMS and static site generator. The platform uses SQLite for content storage and provides tools for managing blog content through a command-line interface.

## Project Structure

- `cms/` - Terminal-based content management system for interacting with the SQLite database
- `site/` - Static site generator that builds the blog from database content
- `internal/` - Core application logic and models
- `migrations/` - Database migration files using Goose
- `templates/` - HTML templates for rendering the site
- `assets/` - Static assets (CSS, JS, images, fonts)
- `dist/` - Generated static site output

## Features

- **Content Management**: Create, update, delete, and list blog entries, categories, tags, and authors
- **Markdown Support**: Content is written in Markdown and rendered to HTML
- **Static Site Generation**: Generates a complete static site from database content
- **Tagging System**: Organize content with tags
- **Author Management**: Support for multiple authors
- **Draft Mode**: Preview draft content before publishing
- **Database Migrations**: Structured database schema changes using Goose

## Usage

### Build the Project

```bash
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
# Build the site (production)
./site

# Build the site with drafts enabled
./site -d

# Clean the dist directory
make site-clean

# Watch for changes and rebuild
make site-watch

# Deploy to production
make site-deploy
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

- **Run tests**: `go test ./...`
- **Format code**: `make goimports`
- **Debug CMS**: `make cms-debug`
- **Debug Site Generator**: `make site-debug`

## License

This project is licensed under the MIT License.
