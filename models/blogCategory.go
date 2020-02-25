package models

import (
	. "../libs"
	"errors"
)

type BlogCategory struct {
	Id int `json:"id"`
	BlogCategoryCreated
}
type BlogCategoryCreated struct {
	Name     string  `json:"name"`
	ParentId *string `json:"parent_id"`
}

//CreateCategory Insert database Category
func CreateBlogCategory(c BlogCategoryCreated) (int64, error) {
	q := `INSERT INTO blogs_categories(name, parent_id) 
		  VALUES($1) RETURNING id`
	db := GetConnection()

	defer db.Close()
	var id int64 = 0
	err := db.QueryRow(q, c.Name, c.ParentId).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

//Busca la informacion de los categories(todos)
func GetBlogCategories() (categories []BlogCategory, err error) {
	q := `SELECT id, name, parent_id FROM blogs_categories ORDER BY parent_id ASC, name ASC`
	db := GetConnection()
	defer db.Close()

	rows, err := db.Query(q)

	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := BlogCategory{}
		err = rows.Scan(
			&c.Id,
			&c.Name,
			&c.ParentId,
		)
		if err != nil {
			return
		}
		categories = append(categories, c)

	}
	return categories, nil
}
func GetBlogCategory(id int64) (category BlogCategory, err error) {
	db := GetConnection()
	defer db.Close()
	err = db.QueryRow(`SELECT id, name, parent_id FROM blogs_categories WHERE id=$1`, id).
		Scan(
			&category.Id,
			&category.Name,
			&category.ParentId,
		)
	if err != nil {
		err = errors.New(`No se ha encontrado la categoria`)
		return
	}
	return
}

// UpdateUser permite actualizar un registro de la db
func UpdateBlogCategory(u BlogCategory) error {
	q := `UPDATE blogs_categories 
			SET name=$1, parent_id=$2 WHERE id=$3`

	db := GetConnection()

	defer db.Close()
	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	r, err := stmt.Exec(u.Name, u.ParentId, u.Id)
	if err != nil {
		return err
	}
	a, _ := r.RowsAffected()
	if a != 1 {
		return errors.New("Error: Se esperaba 1 fila afectada")
	}
	return nil
}
func DeleteBlogCategory(id int64) error {
	q := `DELETE FROM blogs_categories WHERE id=$1`
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
