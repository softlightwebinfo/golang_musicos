package models

import (
	"../settings"
	"time"
)

type Events struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Image     string    `json:"image"`
	Slug      string    `json:"slug"`
	EventsCreate
}
type EventsLocation struct {
	Country string `json:"country"`
	Region  string `json:"region"`
	City    string `json:"city"`
}
type EventsCreate struct {
	FkIdUser    int64     `json:"fk_id_user"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventStart  time.Time `json:"event_start"`
	EventEnd    time.Time `json:"event_end"`
	FkIdCity    int       `json:"fk_id_city"`
	IsPublic    bool      `json:"is_public"`
}
type EventsRequest struct {
	EventsCreate
	EventEnd   string `json:"event_end"`
	EventStart string `json:"event_start"`
}
type EventsList struct {
	Events   Events         `json:"events"`
	Location EventsLocation `json:"location"`
}
type EventsUser struct {
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}
type EventsListGet struct {
	EventsList
	User EventsUser `json:"user"`
}

func CreateEvent(event EventsRequest, image string) (err error) {
	q := `INSERT INTO events(fk_user_id, title, description, event_start, event_end, fk_id_city, image, is_public) 
		  VALUES($1, $2, $3, $4, $5, $6, $7, $8)`
	db := settings.InstanceDb
	_, err = db.Exec(q,
		&event.FkIdUser,
		&event.Title,
		&event.Description,
		&event.EventStart,
		&event.EventEnd,
		&event.FkIdCity,
		&image,
		&event.IsPublic,
	)
	return
}
func GetEventsUser(userId int64) (items []EventsList, err error) {
	q := `SELECT e.id, e.fk_user_id, e.title, e.description, e.created_at, e.updated_at, e.event_start, e.event_end, e.fk_id_city, e.image, e.is_public,
       c.name as city,r.name as region, co.name as country
FROM events e
INNER JOIN cities c on e.fk_id_city = c.id
INNER JOIN regions r ON c.region = r.code AND c.country=r.country
INNER JOIN countries co ON c.country=co.code
WHERE e.fk_user_id = $1 ORDER BY e.updated_at DESC`
	db := settings.InstanceDb

	rows, err := db.Query(q, userId)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := EventsList{}
		err = rows.Scan(
			&c.Events.Id,
			&c.Events.FkIdUser,
			&c.Events.Title,
			&c.Events.Description,
			&c.Events.CreatedAt,
			&c.Events.UpdatedAt,
			&c.Events.EventStart,
			&c.Events.EventEnd,
			&c.Events.FkIdCity,
			&c.Events.Image,
			&c.Events.IsPublic,
			&c.Location.City,
			&c.Location.Region,
			&c.Location.Country,
		)
		if err != nil {
			return
		}

		items = append(items, c)

	}
	return
}
func DeleteEvent(userId int64, id int64) (err error) {
	q := `DELETE FROM events WHERE fk_user_id=$1 AND id=$2`
	db := settings.InstanceDb
	_, err = db.Exec(q,
		&userId,
		&id,
	)
	return
}
func GetEventUser(fkUserId, id int64) (event Events, err error) {
	db := settings.InstanceDb
	err = db.QueryRow(`SELECT id, fk_user_id, title, description, created_at, updated_at, event_start, event_end, fk_id_city, image, is_public, slug  FROM events WHERE id=$1 AND fk_user_id=$2`, id, fkUserId).
		Scan(
			&event.Id,
			&event.FkIdUser,
			&event.Title,
			&event.Description,
			&event.CreatedAt,
			&event.UpdatedAt,
			&event.EventStart,
			&event.EventEnd,
			&event.FkIdCity,
			&event.Image,
			&event.IsPublic,
			&event.Slug,
		)

	return
}
func GetEvent(id int64) (event EventsListGet, err error) {
	db := settings.InstanceDb
	q := `SELECT e.id, e.fk_user_id, e.title, e.description, e.created_at, e.updated_at, e.event_start, e.event_end, e.fk_id_city, e.image, e.is_public,
       c.name as city,r.name as region, co.name as country,e.slug, u.name as user_name, u.phone,u.created_at as user_created_at
FROM events e
INNER JOIN cities c on e.fk_id_city = c.id
INNER JOIN regions r ON c.region = r.code AND c.country=r.country
INNER JOIN countries co ON c.country=co.code
INNER JOIN users u on e.fk_user_id = u.id
WHERE e.id = $1 ORDER BY e.updated_at DESC`
	err = db.QueryRow(q, id).
		Scan(
			&event.Events.Id,
			&event.Events.FkIdUser,
			&event.Events.Title,
			&event.Events.Description,
			&event.Events.CreatedAt,
			&event.Events.UpdatedAt,
			&event.Events.EventStart,
			&event.Events.EventEnd,
			&event.Events.FkIdCity,
			&event.Events.Image,
			&event.Events.IsPublic,
			&event.Location.City,
			&event.Location.Region,
			&event.Location.Country,
			&event.Events.Slug,
			&event.User.Name,
			&event.User.Phone,
			&event.User.CreatedAt,
		)

	return
}
func GetEventsGroup() (items []EventsList, err error) {
	q := `SELECT e.id, e.fk_user_id, e.title, concat(substring(e.description, 0, 200), '...') as description, e.created_at, e.updated_at, e.event_start, e.event_end, e.fk_id_city, e.image, e.is_public,
       c.name as city,r.name as region, co.name as country, e.slug
FROM events e
INNER JOIN cities c on e.fk_id_city = c.id
INNER JOIN regions r ON c.region = r.code AND c.country=r.country
INNER JOIN countries co ON c.country=co.code
ORDER BY EXTRACT(YEAR FROM e.event_start) DESC, EXTRACT(MONTH FROM e.event_start) DESC`
	db := settings.InstanceDb

	rows, err := db.Query(q)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := EventsList{}
		err = rows.Scan(
			&c.Events.Id,
			&c.Events.FkIdUser,
			&c.Events.Title,
			&c.Events.Description,
			&c.Events.CreatedAt,
			&c.Events.UpdatedAt,
			&c.Events.EventStart,
			&c.Events.EventEnd,
			&c.Events.FkIdCity,
			&c.Events.Image,
			&c.Events.IsPublic,
			&c.Location.City,
			&c.Location.Region,
			&c.Location.Country,
			&c.Events.Slug,
		)
		if err != nil {
			return
		}

		items = append(items, c)

	}
	return
}
func GetEventsGroupId(user int64) (items []EventsList, err error) {
	q := `SELECT e.id, e.fk_user_id, e.title, concat(substring(e.description, 0, 200), '...') as description, e.created_at, e.updated_at, e.event_start, e.event_end, e.fk_id_city, e.image, e.is_public,
       c.name as city,r.name as region, co.name as country, e.slug
FROM events e
INNER JOIN cities c on e.fk_id_city = c.id
INNER JOIN regions r ON c.region = r.code AND c.country=r.country
INNER JOIN countries co ON c.country=co.code
WHERE e.fk_user_id=$1
ORDER BY EXTRACT(YEAR FROM e.event_start) DESC, EXTRACT(MONTH FROM e.event_start) DESC`
	db := settings.InstanceDb

	rows, err := db.Query(q, user)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := EventsList{}
		err = rows.Scan(
			&c.Events.Id,
			&c.Events.FkIdUser,
			&c.Events.Title,
			&c.Events.Description,
			&c.Events.CreatedAt,
			&c.Events.UpdatedAt,
			&c.Events.EventStart,
			&c.Events.EventEnd,
			&c.Events.FkIdCity,
			&c.Events.Image,
			&c.Events.IsPublic,
			&c.Location.City,
			&c.Location.Region,
			&c.Location.Country,
			&c.Events.Slug,
		)
		if err != nil {
			return
		}

		items = append(items, c)

	}
	return
}
