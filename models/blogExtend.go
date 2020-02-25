package models

import (
	. "../libs"
	"../settings"
	"database/sql"
	"errors"
)

type BlogExtend struct {
	Blog         Blog          `json:"blog"`
	Category     BlogCategory  `json:"category"`
	User         UserDetail    `json:"user"`
	BlogTags     []BlogTag     `json:"tags"`
	BlogComments []BlogComment `json:"comments"`
}

//Busca la informacion de los categories(todos)
func GetBlogsExtend() (categories []BlogExtend, err error) {
	q := `SELECT id,fk_user_id,cat_id,cat_name,cat_parent_id,created_at,description,title,	
		image,updated_at,user_created_at,user_email,user_name,user_phone,user_updated_at, user_active, user_age
 		FROM view_mat_blogs`
	db := GetConnection()
	defer db.Close()

	rows, err := db.Query(q)

	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := BlogExtend{}
		imageNul := sql.NullString{}
		err = rows.Scan(
			&c.Blog.Id,
			&c.User.Id,
			&c.Category.Id,
			&c.Category.Name,
			&c.Category.ParentId,
			&c.Blog.CreatedAt,
			&c.Blog.Description,
			&c.Blog.Title,
			&imageNul,
			&c.Blog.UpdatedAt,
			&c.User.CreatedAt,
			&c.User.Email,
			&c.User.Name,
			&c.User.Phone,
			&c.User.UpdatedAt,
			&c.User.Active,
			&c.User.Age,
		)
		c.Blog.FkUserId = int64(c.User.Id)
		c.Blog.FkIdCategory = int64(c.Category.Id)
		var image string = settings.GetImage(imageNul.String)
		if imageNul.String == "" {
			image = settings.GetImage("no-foto.png")

		}
		c.Blog.Image = &image;
		c.BlogTags, err = GetBlogTags(c.Blog.Id)
		c.BlogComments, err = GetBlogComments(c.Blog.Id)

		if err != nil {
			return
		}
		categories = append(categories, c)
	}
	return categories, nil
}
func GetBlogsExtendUser(userId int64) (categories []BlogExtend, err error) {
	q := `SELECT id,fk_user_id,cat_id,cat_name,cat_parent_id,created_at,description,title,	
		image,updated_at,user_created_at,user_email,user_name,user_phone,user_updated_at, user_active, user_age
 		FROM view_mat_blogs WHERE fk_user_id=$1`
	db := GetConnection()
	defer db.Close()

	rows, err := db.Query(q, userId)

	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := BlogExtend{}
		imageNul := sql.NullString{}
		err = rows.Scan(
			&c.Blog.Id,
			&c.User.Id,
			&c.Category.Id,
			&c.Category.Name,
			&c.Category.ParentId,
			&c.Blog.CreatedAt,
			&c.Blog.Description,
			&c.Blog.Title,
			&imageNul,
			&c.Blog.UpdatedAt,
			&c.User.CreatedAt,
			&c.User.Email,
			&c.User.Name,
			&c.User.Phone,
			&c.User.UpdatedAt,
			&c.User.Active,
			&c.User.Age,
		)
		c.Blog.FkUserId = int64(c.User.Id)
		c.Blog.FkIdCategory = int64(c.Category.Id)
		var image string = settings.GetImage(imageNul.String)
		if imageNul.String == "" {
			image = settings.GetImage("no-foto.png")

		}
		c.Blog.Image = &image;
		c.BlogTags, err = GetBlogTags(c.Blog.Id)
		c.BlogComments, err = GetBlogComments(c.Blog.Id)

		if err != nil {
			return
		}
		categories = append(categories, c)
	}
	return categories, nil
}
func GetBlogsGetExtend(id int64) (c BlogExtend, err error) {
	db := GetConnection()
	defer db.Close()
	imageNul := sql.NullString{}

	err = db.QueryRow(`SELECT id,fk_user_id,cat_id,cat_name,cat_parent_id,created_at,description,title,	
		image,updated_at,user_created_at,user_email,user_name,user_phone,user_updated_at, user_active, user_age
 		FROM view_mat_blogs WHERE id=$1`, id).
		Scan(
			&c.Blog.Id,
			&c.User.Id,
			&c.Category.Id,
			&c.Category.Name,
			&c.Category.ParentId,
			&c.Blog.CreatedAt,
			&c.Blog.Description,
			&c.Blog.Title,
			&imageNul,
			&c.Blog.UpdatedAt,
			&c.User.CreatedAt,
			&c.User.Email,
			&c.User.Name,
			&c.User.Phone,
			&c.User.UpdatedAt,
			&c.User.Active,
			&c.User.Age,
		)
	c.BlogTags, err = GetBlogTags(c.Blog.Id)
	c.BlogComments, err = GetBlogComments(c.Blog.Id)
	if err != nil {
		err = errors.New(`No se ha encontrado el blog`)
		return
	}
	c.Blog.FkUserId = int64(c.User.Id)
	c.Blog.FkIdCategory = int64(c.Category.Id)
	var image string = settings.GetImage(imageNul.String)
	if imageNul.String == "" {
		image = settings.GetImage("no-foto.png")

	}
	c.Blog.Image = &image;
	return
}
