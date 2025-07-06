package store

import (
	"github.com/jmoiron/sqlx"
)

type Tag struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type TagRepo struct {
	db *sqlx.DB
}

func NewTagRepo(db *sqlx.DB) *TagRepo {
	return &TagRepo{db: db}
}

func (r *TagRepo) CreateTag(tag *Tag) error {
	res, err := r.db.Exec("INSERT INTO tags (name) VALUES (?)", tag.Name)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	tag.ID = id
	return nil
}

func (r *TagRepo) GetTag(id int64) (*Tag, error) {
	var tag Tag
	err := r.db.Get(&tag, "SELECT * FROM tags WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepo) GetTagByName(name string) (*Tag, error) {
	var tag Tag
	err := r.db.Get(&tag, "SELECT * FROM tags WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepo) GetTags() ([]*Tag, error) {
	var tags []*Tag
	err := r.db.Select(&tags, "SELECT * FROM tags")
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *TagRepo) GetTagsForEntry(entryID int64) ([]*Tag, error) {
	var tags []*Tag
	err := r.db.Select(&tags, `
		SELECT t.*
		FROM tags t
		JOIN entry_tags et ON t.id = et.tag_id
		WHERE et.entry_id = ?
	`, entryID)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *TagRepo) UpdateTag(tag *Tag) error {
	_, err := r.db.Exec("UPDATE tags SET name = ? WHERE id = ?", tag.Name, tag.ID)
	return err
}

func (r *TagRepo) DeleteTag(id int64) error {
	_, err := r.db.Exec("DELETE FROM tags WHERE id = ?", id)
	return err
}
