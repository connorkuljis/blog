-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS categories (
    title       TEXT PRIMARY KEY UNIQUE NOT NULL,
    description TEXT DEFAULT ''
);

CREATE TABLE IF NOT EXISTS entries (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    category    TEXT NOT NULL,
    title       TEXT NOT NULL,
    content     TEXT DEFAULT '',
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    publish     INTEGER DEFAULT 0,
    FOREIGN KEY (category) REFERENCES categories(title)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
