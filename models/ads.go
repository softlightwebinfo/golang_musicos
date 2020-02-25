package models

import (
	"../settings"
	"github.com/lib/pq"
	"time"
)

type ADSModel struct {
	Id          int64          `json:"id"`
	Title       *string        `json:"title"`
	Title2      *string        `json:"title_2"`
	Title3      *string        `json:"title_3"`
	Description *string        `json:"description"`
	Url         *string        `json:"url"`
	Address     *string        `json:"address"`
	Active      *bool          `json:"active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	FkIdMeta    *int           `json:"fk_id_meta"`
	Step        *int           `json:"step"`
	Category    *string        `json:"category"`
	Meta        *string        `json:"meta"`
	Name        *string        `json:"name"`
	ActiveData  ADSModelActive `json:"active_data"`
	Tags        []string       `json:"tags"`
}
type ADSModelActive struct {
	Impressions      *int `json:"impressions"`
	Clicks           *int `json:"clicks"`
	Price            *int `json:"price"`
	TotalImpressions *int `json:"total_impressions"`
	TotalClicks      *int `json:"total_clicks"`
}
type AdsResponseGetAll struct {
	ADSModel
}

type ADSRequestNewFirst struct {
	Meta int `json:"meta"`
}

type ADSResponseNewFirst struct {
	Id int64 `json:"id"`
}

type ADSResponseNewThree struct {
	Id int64 `json:"id"`
}

type ADSRequestNewSecond struct {
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	Id       int      `json:"id"`
}
type ADSResponseNewSecond struct {
	Id   int64    `json:"id"`
	Data ADSModel `json:"data"`
}
type ADSRequestNewThree struct {
	Id       int64 `json:"id"`
	Selected int   `json:"selected"`
}
type ADSRequestNewReview struct {
	Id   int64 `json:"id"`
	Name string
}

func GetAllADS(userId int64) (ads []AdsResponseGetAll, err error) {
	db := settings.InstanceDb

	q := `SELECT 
       c.id, c.title, c.title2, c.title3, c.description, 
       c.url, c.address, c.created_at,c.updated_at, c.active,
       c.fk_id_meta, c.step, c.category, a.impressions, a.clicks, 
       p.price, p.impressions as total_impressions, p.clicks as total_clicks,
       array(SELECT ct.tag from ads.campaigns_tags ct WHERE ct.fk_id_campaign=c.id),
       m.text as meta, c.name
	FROM ads.campaigns c 
	LEFT JOIN ads.campaigns_active a ON a.fk_id_campaign=c.id
	LEFT JOIN ads.prices p ON a.fk_id_prices=p.id
	LEFT JOIN ads.metas m ON c.fk_id_meta=m.id
	WHERE c.fk_user_id=$1  
	ORDER BY c.updated_at DESC`

	rows, err := db.Query(q, userId)
	if err != nil {
		return
	}
	for rows.Next() {
		c := AdsResponseGetAll{}
		err = rows.Scan(
			&c.Id,
			&c.Title,
			&c.Title2,
			&c.Title3,
			&c.Description,
			&c.Url,
			&c.Address,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.Active,
			&c.FkIdMeta,
			&c.Step,
			&c.Category,
			&c.ActiveData.Impressions,
			&c.ActiveData.Clicks,
			&c.ActiveData.Price,
			&c.ActiveData.TotalImpressions,
			&c.ActiveData.TotalClicks,
			pq.Array(&c.Tags),
			&c.Meta,
			&c.Name,
		)
		if err != nil {
			return
		}
		ads = append(ads, c)
	}
	return
}
func GetADS(userId, id int64) (c AdsResponseGetAll, err error) {
	db := settings.InstanceDb
	q := `SELECT 
       c.id, c.title, c.title2, c.title3, c.description, 
       c.url, c.address, c.created_at,c.updated_at, c.active,
       c.fk_id_meta, c.step, c.category, a.impressions, a.clicks, 
       p.price, p.impressions as total_impressions, p.clicks as total_clicks,
       array(SELECT ct.tag from ads.campaigns_tags ct WHERE ct.fk_id_campaign=c.id),
       m.text as meta, c.name
	FROM ads.campaigns c 
	LEFT JOIN ads.campaigns_active a ON a.fk_id_campaign=c.id
	LEFT JOIN ads.prices p ON a.fk_id_prices=p.id
	LEFT JOIN ads.metas m ON c.fk_id_meta=m.id
	WHERE c.fk_user_id=$1 and c.id=$2
	ORDER BY c.updated_at DESC`
	err = db.QueryRow(q, userId, id).Scan(
		&c.Id,
		&c.Title,
		&c.Title2,
		&c.Title3,
		&c.Description,
		&c.Url,
		&c.Address,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.Active,
		&c.FkIdMeta,
		&c.Step,
		&c.Category,
		&c.ActiveData.Impressions,
		&c.ActiveData.Clicks,
		&c.ActiveData.Price,
		&c.ActiveData.TotalImpressions,
		&c.ActiveData.TotalClicks,
		pq.Array(&c.Tags),
		&c.Meta,
		&c.Name,
	)
	return
}
func PostAdsNewFirst(userId int64, first ADSRequestNewFirst) (id int64, err error) {
	db := &settings.InstanceDb

	err = (*db).QueryRow(
		"INSERT INTO ads.campaigns(fk_user_id, fk_id_meta) VALUES($1,$2) RETURNING id",
		&userId,
		&first.Meta,
	).Scan(&id)
	return
}
func PostAdsNewSecond(userId int64, first ADSRequestNewSecond) (id int64, err error) {
	db := &settings.InstanceDbTX
	err = (*db).QueryRow(
		"UPDATE ads.campaigns SET category=$1, step=$2 WHERE fk_user_id=$3 and id=$4 RETURNING id",
		&first.Category,
		2,
		&userId,
		&first.Id,
	).Scan(&id)
	stmt, _ := (*db).Prepare(
		"INSERT INTO ads.campaigns_tags(fk_id_campaign, tag) VALUES($1, $2)",
	)
	for _, row := range first.Tags {
		_, _ = stmt.Exec(first.Id, row)
		if err != nil {
			_ = (*db).Rollback()
			return
		}
	}
	err = (*db).Commit()
	defer stmt.Close()
	return
}
func DeleteAds(id int64, fkUserId int64) (err error) {
	db := &settings.InstanceDb
	_, err = (*db).Exec(
		"DELETE FROM ads.campaigns WHERE id=$1 and fk_user_id=$2",
		id,
		fkUserId,
	)
	return
}
func AdsChangeStep(id int64, fkUserId int64, step int) {
	db := settings.InstanceDb
	_, _ = db.Exec(
		"UPDATE ads.campaigns SET step=$1 WHERE id=$2 and fk_user_id=$3",
		step,
		id,
		fkUserId,
	)
}
func AdsActiveStep(id int64, selected int) (err error) {
	db := settings.InstanceDb
	_, err = db.Exec(
		"INSERT INTO ads.campaigns_active(fk_id_campaign,fk_id_prices) VALUES($1, $2)",
		id,
		selected,
	)
	return
}
func PostAdsNewReview(userId int64, review ADSRequestNewReview) (err error) {
	db := settings.InstanceDb
	_, err = db.Exec(
		"UPDATE ads.campaigns SET name=$1, step=$4 WHERE id=$2 and fk_user_id=$3",
		&review.Name,
		&review.Id,
		&userId,
		nil,
	)
	return
}
func PostAdsNewReviewActive(id int64, active bool) (err error) {
	db := settings.InstanceDb
	_, err = db.Exec(
		"UPDATE ads.campaigns_active SET active=$1 WHERE fk_id_campaign=$2",
		&active,
		&id,
	)
	return
}
func PostAdsNewReviewActiveAds(id int64, active bool) (err error) {
	db := settings.InstanceDb
	_, err = db.Exec(
		"UPDATE ads.campaigns SET active=$1 WHERE id=$2",
		&active,
		&id,
	)
	return
}
