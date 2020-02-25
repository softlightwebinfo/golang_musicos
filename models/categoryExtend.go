package models

import (
	. "../libs"
)

type CategoryExtend struct {
	IdCategory    int    `json:"id_category"`
	IdSubcategory int    `json:"id_subcategory"`
	Category      string `json:"category"`
	Subcategory   string `json:"subcategory"`
	Count         int    `json:"count"`
	Slug          string `json:"slug"`
}
type CategoryExtendGroup struct {
	Category   CategoryExtendCount `json:"category"`
	Categories []CategoryExtend    `json:"categories"`
}
type CategoryExtendCount struct {
	Category
	Count int `json:"count"`
}

//Busca la informacion de los categories(todos)
func GetCategoriesExtend() (categories []CategoryExtend, err error) {
	q := `SELECT id_category, id_subcategory, category, subcategory, count_items, slug FROM view_mat_categories_subcategories ORDER BY id_category ASC, id_subcategory ASC`
	db := GetConnection()
	defer db.Close()

	rows, err := db.Query(q)

	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := CategoryExtend{}
		err = rows.Scan(
			&c.IdCategory,
			&c.IdSubcategory,
			&c.Category,
			&c.Subcategory,
			&c.Count,
			&c.Slug,
		)
		if err != nil {
			return
		}
		categories = append(categories, c)
	}
	return categories, nil
}
func GetCategoryExtendSlug(slug string) (c CategoryExtend, err error) {
	q := `SELECT id_category, id_subcategory, category, subcategory, count_items, slug FROM view_mat_categories_subcategories WHERE slug=$1`
	db := GetConnection()
	defer db.Close()

	err = db.QueryRow(q, slug).Scan(
		&c.IdCategory,
		&c.IdSubcategory,
		&c.Category,
		&c.Subcategory,
		&c.Count,
		&c.Slug,
	)
	return
}

func GetCategoriesExtendGroup() (categories []CategoryExtendGroup, err error) {
	categoriesArray, err := GetCategoriesExtend()
	var unique []CategoryExtendGroup
	if err != nil {
		return
	}
	for _, categoryRow := range categoriesArray {
		var exist bool = false
		subcategory := CategoryExtend{
			IdCategory:    categoryRow.IdCategory,
			IdSubcategory: categoryRow.IdSubcategory,
			Category:      categoryRow.Category,
			Subcategory:   categoryRow.Subcategory,
			Count:         categoryRow.Count,
			Slug:          categoryRow.Slug,
		}
		for i, categoryUnique := range unique {
			if categoryUnique.Category.Id == categoryRow.IdCategory {
				exist = true
				unique[i].Categories = append(unique[i].Categories, subcategory)
				unique[i].Category.Count += subcategory.Count
				break
			}
		}
		if !exist {
			unique = append(unique, CategoryExtendGroup{
				Category: CategoryExtendCount{
					Count: categoryRow.Count,
					Category: Category{
						Id: categoryRow.IdCategory,
						CategoryCreated: CategoryCreated{
							Name: categoryRow.Category,
						},
					},
				},
				Categories: []CategoryExtend{subcategory},
			})
		}
	}
	categories = unique
	err = nil
	return
}
