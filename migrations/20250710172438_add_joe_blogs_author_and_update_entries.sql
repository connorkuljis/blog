-- +goose Up
-- +goose StatementBegin
-- Insert Joe Blogs author
INSERT INTO authors (name, email, bio) VALUES (
    'Joe Blogs',
    'joe@example.com',
    'Default author for blog posts'
);

-- Add author_id column to entries table
ALTER TABLE entries ADD COLUMN author_id INTEGER NOT NULL DEFAULT 1;

-- Add foreign key constraint
-- Note: SQLite doesn't support adding foreign key constraints to existing tables
-- So we'll create a new table and migrate data
CREATE TABLE entries_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT,
    description TEXT,
    featured_image_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_draft INTEGER NOT NULL DEFAULT 1,
    author_id INTEGER NOT NULL DEFAULT 1,
    FOREIGN KEY (category_id) REFERENCES categories(id),
    FOREIGN KEY (author_id) REFERENCES authors(id)
);

-- Copy data from old table to new table
INSERT INTO entries_new (id, category_id, title, content, description, featured_image_url, created_at, updated_at, is_draft, author_id)
SELECT id, category_id, title, content, description, featured_image_url, created_at, updated_at, is_draft, 1
FROM entries;

-- Drop old table and rename new table
DROP TABLE entries;
ALTER TABLE entries_new RENAME TO entries;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Recreate original entries table structure
CREATE TABLE entries_old (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT,
    description TEXT,
    featured_image_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_draft INTEGER NOT NULL DEFAULT 1,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

-- Copy data back (excluding author_id)
INSERT INTO entries_old (id, category_id, title, content, description, featured_image_url, created_at, updated_at, is_draft)
SELECT id, category_id, title, content, description, featured_image_url, created_at, updated_at, is_draft
FROM entries;

-- Drop new table and rename old table back
DROP TABLE entries;
ALTER TABLE entries_old RENAME TO entries;

-- Remove Joe Blogs author
DELETE FROM authors WHERE name = 'Joe Blogs' AND email = 'joe@example.com';
-- +goose StatementEnd
