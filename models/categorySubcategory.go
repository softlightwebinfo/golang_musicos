package models

import (
	. "../libs"
	"errors"
)

type CategorySubcategory struct {
	FkIdCategory    int    `json:"fk_id_category"`
	FkIdSubcategory int    `json:"fk_id_subcategory"`
	Slug            string `json:"slug"`
}

//CreateCategory Insert database Category
func CreateCategorySubcategory(c CategorySubcategory) error {
	q := `INSERT INTO categories_subcategories(fk_id_category, fk_id_subcategory) 
		  VALUES($1, $2)`
	db := GetConnection()

	defer db.Close()
	_, err := db.Exec(q, c.FkIdCategory, c.FkIdSubcategory)
	if err != nil {
		return err
	}

	return nil
}

//Busca la informacion de los categories(todos)
func GetCategoriesSubcategories() (categories []CategorySubcategory, err error) {
	q := `SELECT * FROM categories_subcategories ORDER BY fk_id_category ASC, fk_id_subcategory ASC`
	db := GetConnection()
	defer db.Close()

	rows, err := db.Query(q)

	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := CategorySubcategory{}
		err = rows.Scan(
			&c.FkIdCategory,
			&c.FkIdSubcategory,
			&c.Slug,
		)
		if err != nil {
			return
		}
		categories = append(categories, c)

	}
	return categories, nil
}
func GetCategorySubcategory(categoryId int64, subcategoryId int64) (category CategorySubcategory, err error) {
	db := GetConnection()
	defer db.Close()
	err = db.QueryRow(`SELECT * FROM categories_subcategories WHERE fk_id_category=$1 AND fk_id_subcategory=$2`, categoryId, subcategoryId).
		Scan(
			&category.FkIdCategory,
			&category.FkIdSubcategory,
			&category.Slug,
		)
	if err != nil {
		err = errors.New(`No se ha encontrado la categoria-subcategoria`)
		return
	}
	return
}
func GetCategorySubcategories(categoryId int64) (categories []CategorySubcategory, err error) {
	db := GetConnection()
	defer db.Close()
	q := `SELECT * FROM categories_subcategories WHERE fk_id_category=$1`

	rows, err := db.Query(q, categoryId)

	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := CategorySubcategory{}
		err = rows.Scan(
			&c.FkIdCategory,
			&c.FkIdSubcategory,
			&c.Slug,
		)
		if err != nil {
			return
		}
		categories = append(categories, c)

	}
	return categories, nil
}

// UpdateUser permite actualizar un registro de la db
func UpdateCategorySubcategory(fk_id_category int64, fk_id_subcategory int64, u CategorySubcategory) error {
	q := `UPDATE categories_subcategories 
			SET fk_id_category=$1,fk_id_subcategory=$2 WHERE fk_id_category=$3 AND fk_id_subcategory=$4`

	db := GetConnection()

	defer db.Close()
	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	r, err := stmt.Exec(fk_id_category, fk_id_subcategory, u.FkIdCategory, u.FkIdSubcategory)
	if err != nil {
		return err
	}
	a, _ := r.RowsAffected()
	if a != 1 {
		return errors.New("Error: Se esperaba 1 fila afectada")
	}
	return nil
}
func DeleteCategorySubcategory(fk_id_category int64, fk_id_subcategory int64) error {
	q := `DELETE FROM categories_subcategories WHERE fk_id_category=$1 AND fk_id_subcategory=$2`
	db := GetConnection()
	defer db.Close()

	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	r, err := stmt.Exec(fk_id_category, fk_id_subcategory)
	if err != nil {
		return err
	}
	raf, _ := r.RowsAffected()
	if raf != 1 {
		return errors.New("Se espera 1 fila afectada")
	}
	return nil
}
func DeleteCategorySubcategories(fk_id_category int64) error {
	q := `DELETE FROM categories_subcategories WHERE fk_id_category=$1`
	db := GetConnection()
	defer db.Close()

	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	r, err := stmt.Exec(fk_id_category)
	if err != nil {
		return err
	}
	raf, _ := r.RowsAffected()
	if raf != 1 {
		return errors.New("Se espera 1 fila afectada")
	}
	return nil
}
