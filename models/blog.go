package models

import (
	. "../libs"
	"../settings"
	"errors"
	"time"
)

type Blog struct {
	BlogSimple
	BlogCreated
}

type BlogCreated struct {
	BlogDescription
	FkUserId     int64 `json:"fk_user_id"`
	FkIdCategory int64 `json:"fk_id_category"`
}
type BlogSimple struct {
	Id        int64      `json:"id"`
	Image     *string    `json:"image"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
type BlogDescription struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

//CreateCategory Insert database Category
func CreateBlog(c BlogCreated) (Blog, error) {
	q := `INSERT INTO blogs(title, description, fk_user_id, fk_category_id) 
		  VALUES($1, $2, $3, $4) RETURNING id, title, description, fk_user_id, fk_category_id, image, created_at, updated_at`
	db := GetConnection()

	defer db.Close()
	var item Blog

	err := db.QueryRow(q,
		&c.Title,
		&c.Description,
		&c.FkUserId,
		&c.FkIdCategory,
	).Scan(
		&item.Id,
		&item.Title,
		&item.Description,
		&item.FkUserId,
		&item.FkIdCategory,
		&item.Image,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return item, err
	}

	return item, nil
}

//Busca la informacion de los categories(todos)
func GetBlogs() (items []Blog, err error) {
	q := `SELECT id, title, description, image, fk_user_id, fk_category_id, created_at, updated_at FROM blogs ORDER BY updated_at ASC`
	db := GetConnection()
	defer db.Close()

	rows, err := db.Query(q)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := Blog{}
		err = rows.Scan(
			&c.Id,
			&c.Title,
			&c.Description,
			&c.Image,
			&c.FkUserId,
			&c.FkIdCategory,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return
		}

		items = append(items, c)

	}
	return items, nil
}
func GetBlog(id int64) (item Blog, err error) {
	db := GetConnection()
	defer db.Close()
	err = db.QueryRow(`SELECT id, title, description, image, fk_user_id, fk_category_id, created_at, updated_at FROM blogs WHERE id=$1`, id).
		Scan(
			&item.Id,
			&item.Title,
			&item.Description,
			&item.Image,
			&item.FkUserId,
			&item.FkIdCategory,
			&item.CreatedAt,
			&item.UpdatedAt,
		)

	if err != nil {
		err = errors.New(`No se ha encontrado el blog`)
		return
	}
	return
}

// UpdateUser permite actualizar un registro de la db
func UpdateBlog(u BlogCreated, id int64) (Blog, error) {
	q := `UPDATE blogs 
			SET title=$1, description=$2, fk_category_id=$3 WHERE id=$4
			RETURNING id, title, description, image, fk_user_id, fk_category_id, created_at, updated_at
			`
	db := GetConnection()
	defer db.Close()
	var item Blog
	err := db.QueryRow(q,
		&u.Title,
		&u.Description,
		&u.FkIdCategory,
		id,
	).Scan(
		&item.Id,
		&item.Title,
		&item.Description,
		&item.Image,
		&item.FkUserId,
		&item.FkIdCategory,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return item, err
	}
	return item, nil
}
func DeleteBlog(id int64) error {
	q := `DELETE FROM blogs WHERE id=$1`
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
func BlogSaveFile(userId int64, id int64, name string) error {
	q := `UPDATE blogs SET image=$1 WHERE fk_user_id=$2 and id=$3`
	db := settings.InstanceDb

	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	r, err := stmt.Exec(name, userId, id)
	if err != nil {
		return err
	}
	raf, _ := r.RowsAffected()
	if raf != 1 {
		return errors.New("Se espera 1 fila afectada")
	}
	return nil
}
