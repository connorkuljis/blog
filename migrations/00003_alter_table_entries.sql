-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS entries_copy (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id INTEGER NOT NULL,
    title       TEXT NOT NULL,
    content     TEXT DEFAULT '',
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    publish     INTEGER DEFAULT 0,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);
    
INSERT INTO entries_copy
    SELECT
        entries.id,
        categories.id AS category_id,
        entries.title,
        entries.content,
        entries.created_at,
        entries.updated_at,
        entries.publish
    FROM
        entries
    INNER JOIN
        categories ON categories.title = entries.category;

DROP TABLE entries;

ALTER TABLE entries_copy RENAME TO entries;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS entries_old (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    category    TEXT NOT NULL,
    title       TEXT NOT NULL,
    content     TEXT DEFAULT '',
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    publish     INTEGER DEFAULT 0,
    FOREIGN KEY (category) REFERENCES categories(title)
);

INSERT INTO entries_old
    SELECT
        entries.id,
        categories.title,
        entries.title,
        entries.content,
        entries.created_at,
        entries.updated_at,
        entries.publish
    FROM
        entries
    INNER JOIN
        categories ON categories.id = entries.category_id;

DROP TABLE entries;

ALTER TABLE entries_old RENAME TO entries;
-- +goose StatementEnd
