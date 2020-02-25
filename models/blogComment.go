package models

import (
	. "../libs"
	"errors"
	"time"
)

type BlogComment struct {
	BlogCommentCreated
	Id        int        `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
type BlogCommentCreated struct {
	FkIdBlog  int64   `json:"fk_id_blog"`
	FkUserId  *int    `json:"fk_user_id"`
	UserName  *string `json:"user_name"`
	UserEmail *string `json:"user_email"`
	Comment   string  `json:"comment"`
}

//CreateCategory Insert database Category
func CreateBlogComment(c BlogCommentCreated) (blog BlogComment, err error) {
	q := `INSERT INTO blogs_comments(fk_id_blog, id, fk_user_id, user_name, user_email, comment) 
		  VALUES($1, 1, $2, $3, $4, $5) RETURNING fk_id_blog, id, fk_user_id, user_name, user_email, comment, created_at, updated_at`
	db := GetConnection()

	defer db.Close()
	err = db.QueryRow(q, c.FkIdBlog, c.FkUserId, c.UserName, c.UserEmail, c.Comment).Scan(
		&blog.FkIdBlog,
		&blog.Id,
		&blog.FkUserId,
		&blog.UserName,
		&blog.UserEmail,
		&blog.Comment,
		&blog.CreatedAt,
		&blog.UpdatedAt,
	)
	if err != nil {
		return
	}

	return
}

//Busca la informacion de los categories(todos)
func GetBlogComment() (comments []BlogComment, err error) {
	q := `SELECT fk_id_blog, id, fk_user_id, user_name, user_email, comment, created_at, updated_at FROM blogs_comments`
	db := GetConnection()
	defer db.Close()

	rows, err := db.Query(q)

	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := BlogComment{}
		err = rows.Scan(
			&c.FkIdBlog,
			&c.Id,
			&c.FkUserId,
			&c.UserName,
			&c.UserEmail,
			&c.Comment,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return
		}
		comments = append(comments, c)

	}
	return
}
func GetOneBlogComments(fkIdBlog int64, id int64) (comments BlogComment, err error) {
	db := GetConnection()
	defer db.Close()
	err = db.QueryRow(`SELECT fk_id_blog, id, fk_user_id, user_name, user_email, comment, created_at, updated_at  FROM blogs_comments WHERE fk_id_blog=$1 and id=$2`, fkIdBlog, id).
		Scan(
			&comments.FkIdBlog,
			&comments.Id,
			&comments.FkUserId,
			&comments.UserName,
			&comments.UserEmail,
			&comments.Comment,
			&comments.CreatedAt,
			&comments.UpdatedAt,
		)
	if err != nil {
		err = errors.New(`No se ha encontrado el comment`)
		return
	}
	return
}
func GetBlogComments(fkIdBlog int64) (comments []BlogComment, err error) {
	q := `SELECT 
				c.fk_id_blog, c.id, c.fk_user_id, c.user_name, c.user_email, c.comment, c.created_at, c.updated_at,
				u.name as user_name, u.email as user_email
			from blogs_comments c 
			INNER JOIN users u on c.fk_user_id = u.id or c.fk_user_id is null
			where fk_id_blog=$1 ORDER BY c.fk_id_blog DESC, c.id DESC`
	db := GetConnection()
	defer db.Close()

	rows, err := db.Query(q, fkIdBlog)

	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := BlogComment{}
		var (
			userName  string
			userEmail string
		)
		err = rows.Scan(
			&c.FkIdBlog,
			&c.Id,
			&c.FkUserId,
			&c.UserName,
			&c.UserEmail,
			&c.Comment,
			&c.CreatedAt,
			&c.UpdatedAt,
			&userName,
			&userEmail,
		)
		if c.UserName == nil {
			c.UserName = &userName
		}
		if c.UserEmail == nil {
			c.UserEmail = &userEmail
		}
		if err != nil {
			return
		}
		comments = append(comments, c)

	}
	return comments, nil
}

// UpdateUser permite actualizar un registro de la db
func UpdateBlogComments(u BlogCommentCreated, fkIdBlog int, id int64) error {
	q := `UPDATE blogs_comments 
			SET fk_id_blog=$1, fk_user_id=$2, user_name=$3, user_email=$4,comment=$5 WHERE fk_id_blog=$6 and id=$7`

	db := GetConnection()

	defer db.Close()
	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	r, err := stmt.Exec(u.FkIdBlog, u.FkUserId, u.UserName, u.UserEmail, u.Comment, fkIdBlog, id)
	if err != nil {
		return err
	}
	a, _ := r.RowsAffected()
	if a != 1 {
		return errors.New("Error: Se esperaba 1 fila afectada")
	}
	return nil
}
func DeleteBlogComments(fkIdBlog int, id int64) error {
	q := `DELETE FROM blogs_comments WHERE fk_id_blog=$1 AND id=$2`
	db := GetConnection()
	defer db.Close()

	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	r, err := stmt.Exec(fkIdBlog, id)
	if err != nil {
		return err
	}
	raf, _ := r.RowsAffected()
	if raf != 1 {
		return errors.New("Se espera 1 fila afectada")
	}
	return nil
}
func DeleteBlogCommentsAll(fkIdBlog int) error {
	q := `DELETE FROM blogs_comments WHERE fk_id_blog=$1`
	db := GetConnection()
	defer db.Close()

	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	r, err := stmt.Exec(fkIdBlog)
	if err != nil {
		return err
	}
	raf, _ := r.RowsAffected()
	if raf != 1 {
		return errors.New("Se espera 1 fila afectada")
	}
	return nil
}
