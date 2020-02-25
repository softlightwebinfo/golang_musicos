package models

import (
	. "../libs"
	"errors"
)

type Region struct {
	Id        int    `json:"id"`
	Country   string `json:"country"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

//Busca la informacion de los categories(todos)
func GetRegions() (regions []Region, err error) {
	q := `SELECT id, country, code, name, latitude, longitude, slug FROM regions`
	db := GetConnection()
	defer db.Close()

	rows, err := db.Query(q)

	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := Region{}
		err = rows.Scan(
			&c.Id,
			&c.Country,
			&c.Code,
			&c.Name,
			&c.Latitude,
			&c.Longitude,
			&c.Slug,
		)
		if err != nil {
			return
		}
		regions = append(regions, c)

	}
	return regions, nil
}
func GetRegionsCountry(codeCountry string) (regions []Region, err error) {
	q := `SELECT id, country, code, name, latitude, longitude, slug FROM regions WHERE country=$1`
	db := GetConnection()
	defer db.Close()

	rows, err := db.Query(q, codeCountry)

	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := Region{}
		err = rows.Scan(
			&c.Id,
			&c.Country,
			&c.Code,
			&c.Name,
			&c.Latitude,
			&c.Longitude,
			&c.Slug,
		)
		if err != nil {
			return
		}
		regions = append(regions, c)

	}
	return regions, nil
}
func GetRegionCountry(codeCountry string, codeRegion string) (region Region, err error) {
	db := GetConnection()
	defer db.Close()
	err = db.QueryRow(`SELECT id, country, code, name, latitude, longitude, slug FROM regions WHERE code=$1 AND country=$2`, codeRegion, codeCountry).
		Scan(
			&region.Id,
			&region.Country,
			&region.Code,
			&region.Name,
			&region.Latitude,
			&region.Longitude,
			&region.Slug,
		)
	if err != nil {
		err = errors.New(`No se ha encontrado el country`)
		return
	}
	return
}
