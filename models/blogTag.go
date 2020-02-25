package models

import (
	. "../libs"
	"errors"
)

type BlogTag struct {
	FkIdBlog int    `json:"fk_id_blog""`
	Name     string `json:"name"`
}

//CreateCategory Insert database Category
func CreateBlogTag(c BlogTag) (blog BlogTag, err error) {
	q := `INSERT INTO blogs_tags(fk_id_blog, name) 
		  VALUES($1, $2) RETURNING fk_id_blog, name`
	db := GetConnection()

	defer db.Close()
	err = db.QueryRow(q, c.FkIdBlog, c.Name).Scan(
		&blog.FkIdBlog,
		&blog.Name,
	)
	if err != nil {
		return
	}

	return
}

//Busca la informacion de los categories(todos)
func GetBlogTag() (categories []BlogTag, err error) {
	q := `SELECT fk_id_blog, name from blogs_tags`
	db := GetConnection()
	defer db.Close()

	rows, err := db.Query(q)

	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := BlogTag{}
		err = rows.Scan(
			&c.FkIdBlog,
			&c.Name,
		)
		if err != nil {
			return
		}
		categories = append(categories, c)

	}
	return categories, nil
}
func GetOneBlogTag(id int64, name string) (category BlogTag, err error) {
	db := GetConnection()
	defer db.Close()
	err = db.QueryRow(`SELECT fk_id_blog, name FROM blogs_tags WHERE fk_id_blog=$1 and name=$2`, id, name).
		Scan(
			&category.FkIdBlog,
			&category.Name,
		)
	if err != nil {
		err = errors.New(`No se ha encontrado el tag`)
		return
	}
	return
}
func GetBlogTags(fk_id_blog int64) (categories []BlogTag, err error) {
	q := `SELECT fk_id_blog, name from blogs_tags where fk_id_blog=$1`
	db := GetConnection()
	defer db.Close()

	rows, err := db.Query(q, fk_id_blog)

	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := BlogTag{}
		err = rows.Scan(
			&c.FkIdBlog,
			&c.Name,
		)
		if err != nil {
			return
		}
		categories = append(categories, c)

	}
	return categories, nil
}

// UpdateUser permite actualizar un registro de la db
func UpdateBlogTag(u BlogTag, fk_id_blog int, name string) error {
	q := `UPDATE blogs_tags 
			SET fk_id_blog=$1, name=$2 WHERE fk_id_blog=$3 and name=$4`

	db := GetConnection()

	defer db.Close()
	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	r, err := stmt.Exec(u.FkIdBlog, u.Name, fk_id_blog, name)
	if err != nil {
		return err
	}
	a, _ := r.RowsAffected()
	if a != 1 {
		return errors.New("Error: Se esperaba 1 fila afectada")
	}
	return nil
}
func DeleteBlogTag(fk_id_blog int, name string) error {
	q := `DELETE FROM blogs_tags WHERE fk_id_blog=$1 AND name=$2`
	db := GetConnection()
	defer db.Close()

	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	r, err := stmt.Exec(fk_id_blog, name)
	if err != nil {
		return err
	}
	raf, _ := r.RowsAffected()
	if raf != 1 {
		return errors.New("Se espera 1 fila afectada")
	}
	return nil
}
