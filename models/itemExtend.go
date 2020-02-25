package models

import (
	. "../libs"
	"../settings"
	"database/sql"
	"fmt"
	"time"
)

type ItemExtend struct {
	Item     ItemDetail     `json:"item"`
	User     UserDetail     `json:"user"`
	Category CategoryExtend `json:"category"`
}

type ItemExtendArray struct {
	Count int          `json:"count"`
	Items []ItemExtend `json:"items"`
}

type ItemExtendCreated struct {
	Item ItemCreated `json:"item"`
	User UserCreated `json:"user"`
}
type ItemAction struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Views     int16     `json:"views"`
	Click     int16     `json:"click"`
	FkItemId  int64     `json:"fk_item_id"`
}

func GetItemsExtend(filter FilterPage) (items ItemExtendArray, err error) {
	sqlSelect := ` 
			item_id, item_title, item_description, item_image, item_price, item_category, item_subcategory,
			item_contact_name, item_contact_phone, item_created_at, item_updated_at,
			user_id, user_name, user_email, user_age, user_created_at, user_updated_at, user_active,
			item_id_category, item_id_subcategory, user_phone, item_category_slug, item_city, item_country, item_region, city_name, region_name, country_name		
	`
	sqlFrom := `view_mat_items`
	sqlLimit := `LIMIT $1 OFFSET $2`
	sqlWhere := ``
	sqlWhereCount := ``
	if filter.All {
		sqlLimit = ""
	}
	if len(filter.Slug) > 0 {
		sqlWhere = `WHERE item_category_slug=$3`
		sqlWhereCount = `WHERE item_category_slug=$1`
	}
	if len(filter.Country) > 0 {
		sqlWhere += fmt.Sprintf(" AND country_slug='%s'", filter.Country)
		sqlWhereCount += fmt.Sprintf(" AND country_slug='%s'", filter.Country)
	}
	if len(filter.Region) > 0 {
		sqlWhere += fmt.Sprintf(" AND region_slug='%s'", filter.Region)
		sqlWhereCount += fmt.Sprintf(" AND region_slug='%s'", filter.Region)
	}
	if len(filter.City) > 0 {
		sqlWhere += fmt.Sprintf(" AND city_slug='%s'", filter.City)
		sqlWhereCount += fmt.Sprintf(" AND city_slug='%s'", filter.City)
	}
	if len(sqlWhereCount) > 0 {
		if filter.IdUser > 0 {
			sqlWhereCount += fmt.Sprintf(" AND user_id=%d", filter.IdUser)
		}
	} else {
		if filter.IdUser > 0 {
			sqlWhereCount += fmt.Sprintf("WHERE user_id=%d", filter.IdUser)
		}
	}
	if len(sqlWhere) > 0 {
		if filter.IdUser > 0 {
			sqlWhere += fmt.Sprintf(" AND user_id=%d", filter.IdUser)
		}
	} else {
		if filter.IdUser > 0 {
			sqlWhere += fmt.Sprintf("WHERE user_id=%d", filter.IdUser)
		}
	}
	q := fmt.Sprintf(`SELECT %s FROM %s %s %s`, sqlSelect, sqlFrom, sqlWhere, sqlLimit)
	qCount := fmt.Sprintf(`SELECT count(item_id) as total FROM %s %s`, sqlFrom, sqlWhereCount)
	db := GetConnection()
	defer db.Close()
	var rows *sql.Rows
	if !filter.All {
		var filterPage int = filter.Page
		var filterOffset int = filter.Limit * filterPage
		if filterPage == 0 {
			filterOffset = 0
		}
		if len(filter.Slug) > 0 {
			rows, err = db.Query(q, filter.Limit, filterOffset, filter.Slug)
		} else {
			rows, err = db.Query(q, filter.Limit, filterOffset)
		}
	} else {
		rows, err = db.Query(q)
	}
	if err != nil {
		return
	}
	defer rows.Close()
	var err2 error
	if len(filter.Slug) > 0 {
		err2 = db.QueryRow(qCount, filter.Slug).Scan(
			&items.Count,
		)
	} else {
		err2 = db.QueryRow(qCount).Scan(
			&items.Count,
		)
	}
	if err2 != nil {
		err = err2
		return
	}
	isNullImage := sql.NullString{}
	for rows.Next() {
		c := ItemExtend{}
		err = rows.Scan(
			&c.Item.Id,
			&c.Item.Title,
			&c.Item.Description,
			&isNullImage,
			&c.Item.Price,
			&c.Category.Category,
			&c.Category.Subcategory,
			&c.Item.ContactName,
			&c.Item.ContactPhone,
			&c.Item.CreatedAt,
			&c.Item.UpdatedAt,
			&c.User.Id,
			&c.User.Name,
			&c.User.Email,
			&c.User.Age,
			&c.User.CreatedAt,
			&c.User.UpdatedAt,
			&c.User.Active,
			&c.Category.IdCategory,
			&c.Category.IdSubcategory,
			&c.User.Phone,
			&c.Category.Slug,
			&c.Item.ItemCity,
			&c.Item.ItemCountry,
			&c.Item.ItemRegion,
			&c.Item.CityName,
			&c.Item.RegionName,
			&c.Item.CountryName,
		)
		c.Item.Image = isNullImage.String
		var image string = Base64Image(isNullImage.String)
		if isNullImage.String == "" {
			image = settings.GetImage("no-foto.png")
		} else if len(isNullImage.String) < 80 {
			image = fmt.Sprintf("%s/%s", "/images", isNullImage.String)
		}
		c.Item.Image = image
		if err != nil {
			return
		}
		items.Items = append(items.Items, c)

	}
	return items, nil
}
func GetItemExtend(idItem int64) (item ItemExtend, err error) {
	q := `
		SELECT 
			item_id, item_title, item_description, item_image, item_price, item_category, item_subcategory,
			item_contact_name, item_contact_phone, item_created_at, item_updated_at,
			user_id, user_name, user_email, user_age, user_created_at, user_updated_at, user_active,
			item_id_category, item_id_subcategory, user_phone, item_city, item_country, item_region, city_name, region_name, country_name	
		FROM view_mat_items WHERE item_id=$1
	`
	db := GetConnection()
	defer db.Close()
	isNullImage := sql.NullString{}

	err = db.QueryRow(q, idItem).Scan(
		&item.Item.Id,
		&item.Item.Title,
		&item.Item.Description,
		&isNullImage,
		&item.Item.Price,
		&item.Category.Category,
		&item.Category.Subcategory,
		&item.Item.ContactName,
		&item.Item.ContactPhone,
		&item.Item.CreatedAt,
		&item.Item.UpdatedAt,
		&item.User.Id,
		&item.User.Name,
		&item.User.Email,
		&item.User.Age,
		&item.User.CreatedAt,
		&item.User.UpdatedAt,
		&item.User.Active,
		&item.Category.IdCategory,
		&item.Category.IdSubcategory,
		&item.User.Phone,
		&item.Item.ItemCity,
		&item.Item.ItemCountry,
		&item.Item.ItemRegion,
		&item.Item.CityName,
		&item.Item.RegionName,
		&item.Item.CountryName,
	)
	var image string = Base64Image(isNullImage.String)
	if isNullImage.String == "" {
		image = settings.GetImage("no-foto.png")
	} else if len(isNullImage.String) < 80 {
		image = fmt.Sprintf("%s/%s", "/images", isNullImage.String)
	}
	item.Item.Image = image
	if err != nil {
		return
	}
	isNullImage = sql.NullString{}
	return
}
func PublishRowActionItem(id int64, typeValue string) (err error) {
	const q = `INSERT into items_actions(fk_id_item, type_col) VALUES($1,$2)`
	db := GetConnection()
	_, err = db.Exec(q, &id, &typeValue)
	defer db.Close()
	return
}
func PublishClick(id int64) (err error) {
	err = PublishRowActionItem(id, "click")
	return
}
func PublishView(id int64) (err error) {
	err = PublishRowActionItem(id, "view")
	return
}

func ItemSaveFile(idUser int64, idItem int64, name string) (err error) {
	const q = `UPDATE items SET image=$3 WHERE fk_user_id=$1 AND id=$2`
	db := settings.InstanceDb
	_, err = db.Exec(q, &idUser, &idItem, name)
	return
}
