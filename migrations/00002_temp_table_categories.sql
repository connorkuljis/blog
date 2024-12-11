-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS categories_copy (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    title       TEXT UNIQUE NOT NULL,
    description TEXT DEFAULT ''
);

INSERT INTO categories_copy (title, description)
    SELECT title, description FROM categories;

DROP TABLE categories;

ALTER TABLE categories_copy RENAME TO categories;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE categories_copy
-- +goose StatementEnd
