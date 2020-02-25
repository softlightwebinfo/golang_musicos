package models

import (
	. "../libs"
	"errors"
)

type Category struct {
	Id int `json:"id"`
	CategoryCreated
}
type CategoryCreated struct {
	Name string `json:"name"`
}

//CreateCategory Insert database Category
func CreateCategory(c CategoryCreated) (int64, error) {
	q := `INSERT INTO categories(name) 
		  VALUES($1) RETURNING id`
	db := GetConnection()

	defer db.Close()
	var id int64 = 0
	err := db.QueryRow(q, c.Name).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

//Busca la informacion de los categories(todos)
func GetCategories() (categories []Category, err error) {
	q := `SELECT id, name FROM categories ORDER BY name ASC`
	db := GetConnection()
	defer db.Close()

	rows, err := db.Query(q)

	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := Category{}
		err = rows.Scan(
			&c.Id,
			&c.Name,
		)
		if err != nil {
			return
		}
		categories = append(categories, c)

	}
	return categories, nil
}
func GetCategory(id int64) (category Category, err error) {
	db := GetConnection()
	defer db.Close()
	err = db.QueryRow(`SELECT id, name FROM categories WHERE id=$1`, id).
		Scan(
			&category.Id,
			&category.Name,
		)
	if err != nil {
		err = errors.New(`No se ha encontrado la categoria`)
		return
	}
	return
}

// UpdateUser permite actualizar un registro de la db
func UpdateCategory(u Category) error {
	q := `UPDATE categories 
			SET name=$1 WHERE id=$4`

	db := GetConnection()

	defer db.Close()
	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	r, err := stmt.Exec(u.Name, u.Id)
	if err != nil {
		return err
	}
	a, _ := r.RowsAffected()
	if a != 1 {
		return errors.New("Error: Se esperaba 1 fila afectada")
	}
	return nil
}
func DeleteCategory(id int64) error {
	q := `DELETE FROM categories WHERE id=$1`
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
