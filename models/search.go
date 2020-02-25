package models

import (
	. "../libs"
	"time"
)

type SearchWeb struct {
	Primary     string    `json:"primary"`
	Id          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
	Type        string    `json:"type"`
	Slug        string    `json:"slug"`
	Data        *string   `json:"data"`
}
type InitSearchWeb struct {
	Count int64       `json:"count"`
	Items []SearchWeb `json:"items"`
}

type InitFilterWeb struct {
	Limit  int   `json:"limit"`
	Offset int64 `json:"offset"`
}
type SearchImage struct {
	Id          int64     `json:"id"`
	UpdatedAt   time.Time `json:"updated_at"`
	Image       string    `json:"image"`
	FkIdProfile int64     `json:"fk_id_profile"`
	FkIdGroup   int64     `json:"fk_id_group"`
	Title       string    `json:"title"`
	Route       string    `json:"route"`
	Slug        string    `json:"slug"`
}

func GetSearchWeb() (items []SearchWeb, err error) {
	q := `SELECT "primary", id, title, description, updated_at, type, slug FROM view_mat_web`
	db := GetConnection()
	defer db.Close()

	rows, err := db.Query(q)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := SearchWeb{}
		err = rows.Scan(
			&c.Primary,
			&c.Id,
			&c.Title,
			&c.Description,
			&c.UpdatedAt,
			&c.Type,
			&c.Slug,
		)
		if err != nil {
			return
		}

		items = append(items, c)

	}
	return items, nil
}
func GetSearchImages() (items []SearchImage, err error) {
	q := `SELECT g.id,g.title,g.updated_at, g.fk_id_profile, g.image, g.fk_id_group, 'profile-public' as route, up.slug from users_profile_gallery g INNER JOIN users_profile up ON up.fk_id_user=g.fk_id_profile ORDER BY updated_at DESC`
	db := GetConnection()
	defer db.Close()

	rows, err := db.Query(q)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := SearchImage{}
		err = rows.Scan(
			&c.Id,
			&c.Title,
			&c.UpdatedAt,
			&c.FkIdProfile,
			&c.Image,
			&c.FkIdGroup,
			&c.Route,
			&c.Slug,
		)
		if err != nil {
			return
		}

		items = append(items, c)

	}
	return items, nil
}
func GetSearchWebLimit(filter InitFilterWeb) (items []SearchWeb, err error) {
	q := `SELECT "primary", id, title, description, updated_at, type, slug, data FROM view_mat_web LIMIT $1 OFFSET $2`
	db := GetConnection()
	defer db.Close()
	rows, err := db.Query(q, filter.Limit, filter.Offset)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := SearchWeb{}
		err = rows.Scan(
			&c.Primary,
			&c.Id,
			&c.Title,
			&c.Description,
			&c.UpdatedAt,
			&c.Type,
			&c.Slug,
			&c.Data,
		)
		if err != nil {
			return
		}

		items = append(items, c)

	}
	return items, nil
}
func GetSearchWebCount() (count int64, err error) {
	q := `SELECT count("primary") as total FROM view_mat_web`
	db := GetConnection()
	defer db.Close()
	err = db.QueryRow(q).
		Scan(&count)
	return
}
