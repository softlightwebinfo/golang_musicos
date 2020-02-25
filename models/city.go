package models

import (
	. "../libs"
)

type City struct {
	Id        int    `json:"id"`
	Country   string `json:"country"`
	Region    string `json:"region"`
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

//Busca la informacion de los categories(todos)
func GetCities() (cities []City, err error) {
	q := `SELECT id, country, region, name, latitude, longitude, slug FROM cities`
	db := GetConnection()
	defer db.Close()

	rows, err := db.Query(q)

	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := City{}
		err = rows.Scan(
			&c.Id,
			&c.Country,
			&c.Region,
			&c.Name,
			&c.Latitude,
			&c.Longitude,
			&c.Slug,
		)
		if err != nil {
			return
		}
		cities = append(cities, c)

	}
	return cities, nil
}
func GetRegionCountryCities(codeCountry string, codeRegion string) (cities []City, err error) {
	db := GetConnection()
	defer db.Close()
	q := `SELECT id, country, region, name, latitude, longitude, slug FROM cities WHERE country=$1 AND region=$2`

	rows, err := db.Query(q, codeCountry, codeRegion)

	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		city := City{}
		err = rows.Scan(
			&city.Id,
			&city.Country,
			&city.Region,
			&city.Name,
			&city.Latitude,
			&city.Longitude,
			&city.Slug,
		)
		if err != nil {
			return
		}
		cities = append(cities, city)

	}
	return cities, nil
}
