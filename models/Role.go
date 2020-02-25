package models

import (
	. "../libs"
	"errors"
)

type Role struct {
	Id int `json:"id"`
	RoleCreated
}
type RoleCreated struct {
	Name     string `json:"name"`
	ParentId *int   `json:"parent_id"`
}

//CreateCategory Insert database role
func CreateRole(c RoleCreated) (int, error) {
	q := `INSERT INTO roles(name, parent_id) 
		  VALUES($1, $2) RETURNING id`
	db := GetConnection()

	defer db.Close()
	var id int = 0
	err := db.QueryRow(q, c.Name, c.ParentId).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

//Busca la informacion de los roles(todos)
func GetRoles() (categories []Role, err error) {
	q := `SELECT id, name, parent_id FROM roles ORDER BY parent_id asc, id asc, name ASC`
	db := GetConnection()
	defer db.Close()

	rows, err := db.Query(q)

	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := Role{}
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
func GetRole(id int) (role Role, err error) {
	db := GetConnection()
	defer db.Close()
	err = db.QueryRow(`SELECT id, name, parent_id FROM roles WHERE id=$1`, id).
		Scan(
			&role.Id,
			&role.Name,
			&role.ParentId,
		)
	if err != nil {
		err = errors.New(`No se ha encontrado el role`)
		return
	}
	return
}

// UpdateUser permite actualizar un registro de la db
func UpdateRole(u RoleCreated, id int) (role Role, err error) {
	q := `UPDATE roles 
			SET name=$1, parent_id=$2 WHERE id=$3 RETURNING id, name, parent_id`

	db := GetConnection()

	defer db.Close()

	err = db.QueryRow(
		q,
		u.Name,
		u.ParentId,
		id,
	).Scan(
		&role.Id,
		&role.Name,
		&role.ParentId,
	)
	if err != nil {
		return
	}

	return
}
func DeleteRole(id int) error {
	q := `DELETE FROM roles WHERE id=$1`
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
