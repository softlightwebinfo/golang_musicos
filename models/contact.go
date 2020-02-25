package models

import (
	. "../libs"
	"errors"
	"time"
)

type Contact struct {
	Id        int        `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	ContactCreated
}
type ContactCreated struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

//CreateContact Insert database Category
func CreateContact(c ContactCreated) (contact Contact, err error) {
	q := `INSERT INTO contacts(name, email, subject, message) 
		  VALUES($1,$2,$3,$4) RETURNING id, name, email, subject, message, created_at`
	db := GetConnection()
	defer db.Close()
	err = db.QueryRow(
		q,
		c.Name,
		c.Email,
		c.Subject,
		c.Message,
	).Scan(
		&contact.Id,
		&contact.Name,
		&contact.Email,
		&contact.Subject,
		&contact.Message,
		&contact.CreatedAt,
	)
	return
}

//Busca la informacion de los categories(todos)
func GetContacts() (contacts []Contact, err error) {
	q := `SELECT id, name, email, subject, message, created_at  FROM contacts ORDER BY id DESC`
	db := GetConnection()
	defer db.Close()

	rows, err := db.Query(q)

	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := Contact{}
		err = rows.Scan(
			&c.Id,
			&c.Name,
			&c.Email,
			&c.Subject,
			&c.Message,
			&c.CreatedAt,
		)
		if err != nil {
			return
		}
		contacts = append(contacts, c)
	}
	return
}
func GetContact(id int64) (c Contact, err error) {
	db := GetConnection()
	defer db.Close()
	err = db.QueryRow(`SELECT id, name, email, subject, message, created_at FROM contacts WHERE id=$1`, id).
		Scan(
			&c.Id,
			&c.Name,
			&c.Email,
			&c.Subject,
			&c.Message,
			&c.CreatedAt,
		)
	if err != nil {
		err = errors.New(`No se ha encontrado el mensaje`)
		return
	}
	return
}

// UpdateUser permite actualizar un registro de la db
func UpdateContact(u Contact) error {
	q := `UPDATE contacts 
			SET name=$1, email=$2, subject=$3, message=$4 WHERE id=$5`

	db := GetConnection()

	defer db.Close()
	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	r, err := stmt.Exec(
		u.Name,
		u.Email,
		u.Subject,
		u.Message,
		u.Id,
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
func DeleteContact(id int64) error {
	q := `DELETE FROM contacts WHERE id=$1`
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
