package models

import (
	. "../libs"
	"errors"
	"time"
)

type ItemContact struct {
	Id        int        `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	ItemContactCreated
}
type ItemContactCreated struct {
	FkIdItem    int64  `json:"fk_id_item"`
	UserName    string `json:"user_name"`
	UserEmail   string `json:"user_email"`
	UserPhone   string `json:"user_phone"`
	UserMessage string `json:"user_message"`
}

//CreateContact Insert database Category
func CreateItemContact(contact ItemContactCreated) (c ItemContact, err error) {
	q := `INSERT INTO items_contacts(fk_id_item, user_name, user_email, user_phone, user_message) 
		  VALUES($1,$2,$3,$4,$5) RETURNING fk_id_item, id, user_name, user_email, user_phone, user_message, created_at`
	db := GetConnection()
	defer db.Close()
	err = db.QueryRow(
		q,
		contact.FkIdItem,
		contact.UserName,
		contact.UserEmail,
		contact.UserPhone,
		contact.UserMessage,
	).Scan(
		&c.FkIdItem,
		&c.Id,
		&c.UserName,
		&c.UserEmail,
		&c.UserPhone,
		&c.UserMessage,
		&c.CreatedAt,
	)
	return
}

//Busca la informacion de los categories(todos)
func GetItemContacts() (contacts []ItemContact, err error) {
	q := `SELECT fk_id_item, id, user_name, user_email, user_phone, user_message, created_at  FROM items_contacts ORDER BY created_at DESC`
	db := GetConnection()
	defer db.Close()

	rows, err := db.Query(q)

	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := ItemContact{}
		err = rows.Scan(
			&c.FkIdItem,
			&c.Id,
			&c.UserName,
			&c.UserEmail,
			&c.UserPhone,
			&c.UserMessage,
			&c.CreatedAt,
		)
		if err != nil {
			return
		}
		contacts = append(contacts, c)
	}
	return
}

//GetItemContact
func GetItemContact(fkIdItem int64, id int64) (c ItemContact, err error) {
	db := GetConnection()
	defer db.Close()
	err = db.QueryRow(`SELECT fk_id_item, id, user_name, user_email, user_phone, user_message, created_at FROM items_contacts WHERE fk_id_item=$1 and id=$2`, fkIdItem, id).
		Scan(
			&c.FkIdItem,
			&c.Id,
			&c.UserName,
			&c.UserEmail,
			&c.UserPhone,
			&c.UserMessage,
			&c.CreatedAt,
		)
	if err != nil {
		err = errors.New(`No se ha encontrado el contacto`)
		return
	}
	return
}

// UpdateItemContact permite actualizar un registro de la db
func UpdateItemContact(c ItemContact) error {
	q := `UPDATE items_contacts 
			SET fk_id_item=$1,user_email=$2,user_name=$3,user_phone=$4,user_message=$5, WHERE fk_id_item=$6 and id=$7`

	db := GetConnection()

	defer db.Close()
	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	r, err := stmt.Exec(
		c.FkIdItem,
		c.UserEmail,
		c.UserName,
		c.UserPhone,
		c.UserMessage,
		c.FkIdItem,
		c.Id,
	)
	if err != nil {
		return err
	}
	a, _ := r.RowsAffected()
	if a != 1 {
		return errors.New("Error: Se esperaba 1 fila afectada")
	}
	return nil
}
func DeleteItemContact(fkIdItem int64, id int64) error {
	q := `DELETE FROM items_contacts WHERE fk_id_item=$1 and id=$2`
	db := GetConnection()
	defer db.Close()

	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	r, err := stmt.Exec(fkIdItem, id)
	if err != nil {
		return err
	}
	raf, _ := r.RowsAffected()
	if raf != 1 {
		return errors.New("Se espera 1 fila afectada")
	}
	return nil
}
