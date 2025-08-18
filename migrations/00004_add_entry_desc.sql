-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS entries_with_desc (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id INTEGER NOT NULL,
    title       TEXT NOT NULL,
    content     TEXT,
	description TEXT, -- add description
	featured_image_url TEXT, -- add featured image url
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_draft    INTEGER NOT NULL DEFAULT 1, -- change publish -> is_draft
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

INSERT INTO entries_with_desc (
	id, 
	category_id, 
	title, 
	content,
	created_at, 
	updated_at, 
	is_draft
)
SELECT
	id, 
	category_id, 
	title, 
	content,
	created_at, 
	updated_at, 
	1 - publish
FROM 
	entries
;

ALTER TABLE entries RENAME TO entries_backup;
ALTER TABLE entries_with_desc RENAME TO entries;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE entries RENAME TO entries_with_desc;
ALTER TABLE entries_backup RENAME TO entries;
-- +goose StatementEnd
