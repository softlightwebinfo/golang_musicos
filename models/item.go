package models

import (
	. "../libs"
	"database/sql"
	"errors"
	"time"
)

type Item struct {
	ItemSimple
	ItemCreated
}

type ItemCreated struct {
	ItemDescription
	FkUserId        int64   `json:"fk_user_id"`
	FkIdCategory    int64   `json:"fk_id_category"`
	FkIdSubCategory int64   `json:"fk_id_subcategory"`
	FkCity          int     `json:"fk_city"`
}
type ItemSimple struct {
	Id        int64      `json:"id"`
	Image     string     `json:"image"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
type ItemDescription struct {
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	ContactName  string  `json:"contact_name"`
	ContactPhone string  `json:"contact_phone"`
}
type ItemDetail struct {
	ItemSimple
	ItemDescription
	ItemCity    string `json:"item_city"`
	ItemCountry string `json:"item_country"`
	ItemRegion  string `json:"item_region"`
	CityName    string `json:"city_name"`
	RegionName  string `json:"region_name"`
	CountryName string `json:"country_name"`
}

//CreateCategory Insert database Category
func CreateItem(c ItemCreated) (Item, error) {
	q := `INSERT INTO items(title, description, price, contact_name, contact_phone, fk_user_id, fk_id_category, fk_id_subcategory, fk_city) 
		  VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, title, description, price, contact_name, contact_phone, image, fk_user_id, fk_id_category, fk_id_subcategory, created_at, updated_at, fk_city`
	db := GetConnection()

	defer db.Close()
	var item Item
	isNullImage := sql.NullString{}
	err := db.QueryRow(q,
		&c.Title,
		&c.Description,
		&c.Price,
		&c.ContactName,
		&c.ContactPhone,
		&c.FkUserId,
		&c.FkIdCategory,
		&c.FkIdSubCategory,
		&c.FkCity,
	).Scan(
		&item.Id,
		&item.Title,
		&item.Description,
		&item.Price,
		&item.ContactName,
		&item.ContactPhone,
		&isNullImage,
		&item.FkUserId,
		&item.FkIdCategory,
		&item.FkIdSubCategory,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.FkCity,
	)
	item.Image = isNullImage.String

	if err != nil {
		return item, err
	}

	return item, nil
}
func UploadItemImage(c Item, image *string) error {
	q := `INSERT INTO items_images(image, fk_item_id) 
		  VALUES($1, $2)`
	db := GetConnection()
	defer db.Close()
	_, err := db.Query(q,
		image,
		&c.Id,
	)
	return err
}

//Busca la informacion de los categories(todos)
func GetItems() (items []Item, err error) {
	q := `SELECT id, title, description, price, contact_name, contact_phone, image, fk_user_id, fk_id_category, fk_id_subcategory, created_at, updated_at, fk_city FROM items ORDER BY updated_at ASC`
	db := GetConnection()
	defer db.Close()

	rows, err := db.Query(q)
	isNullImage := sql.NullString{}
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := Item{}
		err = rows.Scan(
			&c.Id,
			&c.Title,
			&c.Description,
			&c.Price,
			&c.ContactName,
			&c.ContactPhone,
			&isNullImage,
			&c.FkUserId,
			&c.FkIdCategory,
			&c.FkIdSubCategory,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.FkCity,
		)
		if err != nil {
			return
		}
		c.Image = isNullImage.String

		items = append(items, c)

	}
	return items, nil
}
func GetItem(id int64) (item Item, err error) {
	db := GetConnection()
	defer db.Close()
	isNullImage := sql.NullString{}
	err = db.QueryRow(`SELECT id, title, description, price, contact_name, contact_phone, image, fk_user_id, fk_id_category, fk_id_subcategory, created_at, updated_at, fk_city FROM items WHERE id=$1`, id).
		Scan(
			&item.Id,
			&item.Title,
			&item.Description,
			&item.Price,
			&item.ContactName,
			&item.ContactPhone,
			&isNullImage,
			&item.FkUserId,
			&item.FkIdCategory,
			&item.FkIdSubCategory,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.FkCity,
		)
	item.Image = isNullImage.String

	if err != nil {
		err = errors.New(`No se ha encontrado el item`)
		return
	}
	return
}

// UpdateUser permite actualizar un registro de la db
func UpdateItem(u ItemCreated, id int64) (Item, error) {
	q := `UPDATE items 
			SET title=$1, description=$2, price=$3, contact_name=$4, contact_phone=$5, fk_user_id=$6, fk_id_category=$7, fk_id_subcategory=$8, fk_city=$9 WHERE id=$10
			RETURNING id, title, description, price, contact_name, contact_phone, image, fk_user_id, fk_id_category, fk_id_subcategory, created_at, updated_at, fk_city
			`
	db := GetConnection()
	defer db.Close()
	var item Item
	isNullImage := sql.NullString{}
	err := db.QueryRow(q,
		&u.Title,
		&u.Description,
		&u.Price,
		&u.ContactName,
		&u.ContactPhone,
		&u.FkUserId,
		&u.FkIdCategory,
		&u.FkIdSubCategory,
		&u.FkCity,
		id,
	).Scan(
		&item.Id,
		&item.Title,
		&item.Description,
		&item.Price,
		&item.ContactName,
		&item.ContactPhone,
		&isNullImage,
		&item.FkUserId,
		&item.FkIdCategory,
		&item.FkIdSubCategory,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.FkCity,
	)
	item.Image = isNullImage.String
	if err != nil {
		return item, err
	}
	return item, nil
}
func DeleteItem(id int64) error {
	q := `DELETE FROM items WHERE id=$1`
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
