package dto

import "github.com/connorkuljis/blog/internal/model"

func CategoriesToDTO(categories []*model.Category) []*Category {
	var dtoCategories []*Category
	for _, c := range categories {
		dtoCategories = append(dtoCategories, NewCategoryDTO(*c))
	}
	return dtoCategories
}

func EntriesToDTO(entries []*model.Entry) []*Entry {
	var dtoEntries []*Entry
	for _, e := range entries {
		dtoEntries = append(dtoEntries, NewEntryDTO(*e))
	}
	return dtoEntries
}
