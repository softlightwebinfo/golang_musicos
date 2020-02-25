package models

import (
	. "../libs"
	"errors"
)

type Subcategory struct {
	Id int `json:"id"`
	SubCategoryCreated
}
type SubCategoryCreated struct {
	Name string `json:"name"`
}

//CreateCategory Insert database Category
func CreateSubCategory(s SubCategoryCreated) (int64, error) {
	q := `INSERT INTO subcategories(name) 
		  VALUES($1) RETURNING id`
	db := GetConnection()

	defer db.Close()
	var id int64 = 0
	err := db.QueryRow(q, s.Name).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

//Busca la informacion de los categories(todos)
func GetSubCategories() (subcategories []Subcategory, err error) {
	q := `SELECT id, name FROM subcategories ORDER BY name ASC`
	db := GetConnection()
	defer db.Close()

	rows, err := db.Query(q)

	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		s := Subcategory{}
		err = rows.Scan(
			&s.Id,
			&s.Name,
		)
		if err != nil {
			return
		}
		subcategories = append(subcategories, s)

	}
	return subcategories, nil
}
func GetSubCategory(id int64) (category Category, err error) {
	db := GetConnection()
	defer db.Close()
	err = db.QueryRow(`SELECT id, name FROM subcategories WHERE id=$1`, id).
		Scan(
			&category.Id,
			&category.Name,
		)
	if err != nil {
		err = errors.New(`No se ha encontrado la subcategoria`)
		return
	}
	return
}

// UpdateUser permite actualizar un registro de la db
func UpdateSubCategory(s Subcategory) error {
	q := `UPDATE subcategories 
			SET name=$1 WHERE id=$4`

	db := GetConnection()

	defer db.Close()
	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	r, err := stmt.Exec(s.Name, s.Id)
	if err != nil {
		return err
	}
	a, _ := r.RowsAffected()
	if a != 1 {
		return errors.New("Error: Se esperaba 1 fila afectada")
	}
	return nil
}
func DeleteSubCategory(id int64) error {
	q := `DELETE FROM subcategories WHERE id=$1`
	db := GetConnection()
	defer db.Close()

	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	r, err := stmt.Exec(id)
	if err != nil {
		return err
	}
	raf, _ := r.RowsAffected()
	if raf != 1 {
		return errors.New("Se espera 1 fila afectada")
	}
	return nil
}
