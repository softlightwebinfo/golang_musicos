package models

import (
	. "../libs"
	"errors"
)

type Country struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}
//Busca la informacion de los categories(todos)
func GetCountries() (country []Country, err error) {
	q := `SELECT code, name, latitude, longitude, slug FROM countries ORDER BY code ASC`
	db := GetConnection()
	defer db.Close()

	rows, err := db.Query(q)

	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := Country{}
		err = rows.Scan(
			&c.Code,
			&c.Name,
			&c.Latitude,
			&c.Longitude,
			&c.Slug,
		)
		if err != nil {
			return
		}
		country = append(country, c)

	}
	return country, nil
}
func GetCountry(code string) (country Country, err error) {
	db := GetConnection()
	defer db.Close()
	err = db.QueryRow(`SELECT code, name, latitude, longitude, slug FROM countries WHERE code=$1`, code).
		Scan(
			&country.Code,
			&country.Name,
			&country.Latitude,
			&country.Longitude,
			&country.Slug,
		)
	if err != nil {
		err = errors.New(`No se ha encontrado el country`)
		return
	}
	return
}
