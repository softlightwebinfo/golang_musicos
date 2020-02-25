package models

import (
	. "../libs"
	"../settings"
	"errors"
	"fmt"
	pg "github.com/lib/pq"
	"strings"
	"time"
)

type User struct {
	UserSimple
	UserCreated
}

type UserSimple struct {
	Id        int        `json:"id" xml:"id"`
	CreatedAt time.Time  `json:"createdAt" xml:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt" xml:"createdAt"`
}

type UserDescription struct {
	Phone    *string   `form:"phone" json:"phone"`
	Name     string    `form:"name" json:"name"`
	Email    string    `form:"email" json:"email"`
	FkRoleId int       `json:"fk_role_id"`
	Active   bool      `form:"active" json:"active"`
	Age      time.Time `form:"age" json:"age" time_format:"2006-01-02T15:04:05Z" time_utc:"1" example:"2006-01-02T00:00:00Z"`
	Slug     *string   `json:"slug"`
	Avatar   *string   `json:"avatar"`
}

type UserCreated struct {
	UserDescription
	Password string `form:"password" json:"password"`
}
type UserRegister struct {
	UserCreated
	CIF         string `json:"cif"`
	Direction   string `json:"direction"`
	RazonSocial string `json:"razon_social"`
}
type UserDetail struct {
	UserSimple
	UserDescription
}

type UserRecovery struct {
	Email string `json:"email"`
}
type UserRecoveryEmail struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}
type UserProfile struct {
	FkIdUser     int64   `json:"fk_id_user"`
	Slug         string  `json:"slug"`
	Description  string  `json:"description"`
	UrlFacebook  *string `json:"url_facebook"`
	UrlTwitter   *string `json:"url_twitter"`
	UrlYoutube   *string `json:"url_youtube"`
	UrlInstagram *string `json:"url_instagram"`
	UrlWebsite   *string `json:"url_website"`
	Phone        *string `json:"phone"`
	Address      *string `json:"address"`
	Email        *string `json:"email"`
	Occupation   *string `json:"occupation"`
	Banner       *string `json:"banner"`
}
type UserProfileGallery struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserProfileGalleryCreate
}
type UserProfileGalleryCreate struct {
	Image       string `json:"image"`
	FkIdProfile int64  `json:"fk_id_profile"`
	FkIdGroup   int64  `json:"fk_id_group"`
	Title       string `json:"title"`
}

func AuthLogin(userCredential AutCredentials) (user User, err error) {
	q := `SELECT u.id, u.name, u.email, u.password, u.created_at, u.active, u.age, u.updated_at, u.fk_role_id, u.phone,up.slug, u.avatar FROM users u LEFT JOIN users_profile up ON up.fk_id_user=u.id WHERE u.email=$1`
	db := GetConnection()
	defer db.Close()
	timeAgeNull := pg.NullTime{}
	err = db.QueryRow(q, userCredential.Email).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.Active,
		&timeAgeNull,
		&user.UpdatedAt,
		&user.FkRoleId,
		&user.Phone,
		&user.Slug,
		&user.Avatar,
	)
	user.Age = timeAgeNull.Time

	if !ComparePasswords(user.Password, GetPwd(userCredential.Password)) {
		err = errors.New("La contraseña no es valida")
	}

	return
}

//CreateUser Insert database users
func CreateUser(u UserCreated) (int64, error) {
	q := `INSERT INTO users(name, email, password, active, age, phone, fk_role_id) 
		  VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	db := GetConnection()

	defer db.Close()
	var id int64 = 0
	pswd := GetPwd(u.Password)
	password := HashAndSalt(pswd)
	err := db.QueryRow(
		q,
		u.Name,
		strings.ToLower(u.Email),
		password,
		u.Active,
		u.Age,
		u.Phone,
		u.FkRoleId,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
} //CreateUser Insert database users
func CreateBusiness(id int64, u UserRegister) error {
	q := `INSERT INTO users_business(fk_id_user, razon_social, cif, direction) 
		  VALUES($1, $2, $3, $4)`
	db := GetConnection()
	defer db.Close()
	_, err := db.Exec(
		q,
		id,
		u.RazonSocial,
		u.CIF,
		u.Direction,
	)
	return err
}

//Busca la informacion de los estudiantes(todos)
func GetUsers() (users []User, err error) {
	q := `SELECT id, name, email, password, age, active, created_at, updated_at , phone
		  FROM users ORDER BY id ASC`
	timeAgeNull := pg.NullTime{}
	db := GetConnection()
	defer db.Close()

	rows, err := db.Query(q)

	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		u := User{}
		err = rows.Scan(
			&u.Id,
			&u.Name,
			&u.Email,
			&u.Password,
			&timeAgeNull,
			&u.Active,
			&u.CreatedAt,
			&u.UpdatedAt,
			&u.Phone,
		)
		if err != nil {
			return
		}
		u.Age = timeAgeNull.Time
		users = append(users, u)

	}
	return users, nil
}
func GetUser(id int64) (user User, err error) {
	db := GetConnection()
	defer db.Close()
	timeAgeNull := pg.NullTime{}
	q := `SELECT u.id, u.name, u.email, u.password, u.created_at, u.active, u.age, u.updated_at, u.fk_role_id, u.phone,up.slug,u.avatar FROM users u LEFT JOIN users_profile up ON up.fk_id_user=u.id WHERE u.id=$1`
	err = db.QueryRow(q, id).
		Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.Active,
			&timeAgeNull,
			&user.UpdatedAt,
			&user.FkRoleId,
			&user.Phone,
			&user.Slug,
			&user.Avatar,
		)

	user.Age = timeAgeNull.Time
	if err != nil {
		err = errors.New(`No se ha encontrado el usuario`)
		return
	}
	return
}
func GetUserProfile(id int64) (user UserProfile, err error) {
	db := GetConnection()
	defer db.Close()
	q := `SELECT fk_id_user, slug, description, url_facebook, url_twitter, url_youtube, url_instagram, url_website, phone, address, email, occupation, banner  FROM users_profile u WHERE u.fk_id_user=$1`
	err = db.QueryRow(q, id).
		Scan(
			&user.FkIdUser,
			&user.Slug,
			&user.Description,
			&user.UrlFacebook,
			&user.UrlTwitter,
			&user.UrlYoutube,
			&user.UrlInstagram,
			&user.UrlWebsite,
			&user.Phone,
			&user.Address,
			&user.Email,
			&user.Occupation,
			&user.Banner,
		)

	if err != nil {
		err = errors.New(`No se ha encontrado el usuario`)
		return
	}
	return
}

func IsExistUserEmail(email string) (user User, err error) {
	db := GetConnection()
	defer db.Close()
	timeAgeNull := pg.NullTime{}
	err = db.QueryRow(`SELECT id, name, email,password, age, active, created_at, updated_at, phone FROM users WHERE email=$1`, email).
		Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Password,
			&timeAgeNull,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Phone,
		)

	user.Age = timeAgeNull.Time
	if err == nil {
		err = errors.New(`El usuario ya existe con este email`)
		return
	} else {
		err = nil
	}
	return
}
func IsExistUserPhone(phone string) (user User, err error) {
	db := GetConnection()
	defer db.Close()
	timeAgeNull := pg.NullTime{}
	err = db.QueryRow(`SELECT id, name, email,password, age, active, created_at, updated_at, phone FROM users WHERE phone=$1`, phone).
		Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Password,
			&timeAgeNull,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Phone,
		)

	user.Age = timeAgeNull.Time
	if err == nil {
		err = errors.New(`El usuario ya existe con este telefono`)
		return
	} else {
		err = nil
	}
	return
}
func IsExistUserEmailPhone(email string, phone string) (user User, err error) {
	userEmail, errEmail := IsExistUserEmail(email)
	if errEmail != nil {
		err = errEmail
		user = userEmail
		return
	} else {
		err = nil
	}
	userPhone, errPhone := IsExistUserPhone(phone)
	if errPhone != nil {
		err = errPhone
		user = userPhone
		return
	} else {
		err = nil
	}
	return
}

// UpdateUser permite actualizar un registro de la db
func UpdateUser(u User) error {
	q := `UPDATE users 
			SET name=$1, email=$2, active=$3, fk_role_id=$4, updated_at = now(), phone=$5 WHERE id=$6`

	db := GetConnection()

	defer db.Close()
	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	r, err := stmt.Exec(&u.Name, &u.Email, &u.Active, &u.FkRoleId, &u.Phone, &u.Id)
	if err != nil {
		return err
	}
	a, _ := r.RowsAffected()
	if a != 1 {
		return errors.New("Error: Se esperaba 1 fila afectada")
	}
	return nil
}
func DeleteUser(id int64) error {
	q := `DELETE FROM users WHERE id=$1`
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
func UserLog(id int64) (err error) {
	q := `INSERT INTO users_login(fk_user_id) VALUES($1)`
	db := GetConnection()
	_, err = db.Exec(q, id)
	return
}
func GenerateTokenUserEmail(user UserCreated) string {
	return Hash(fmt.Sprintf("%s.%s.%s", user.Email, user.Name, user.Password), string(settings.JwtKey))
}
func GenerateTokenEmail(email string) string {
	return Hash(fmt.Sprintf("%s", email), string(settings.JwtKey))
}
func UserSaveTokenConfirmEmailCreate(userId int64, token string) (err error) {
	q := `INSERT INTO users_email_confirm(fk_user_id, token) VALUES($1, $2)`
	db := GetConnection()
	defer db.Close()
	_, err = db.Exec(
		q,
		userId,
		token,
	)
	return
}
func UserSaveTokenConfirmEmailConfirmed(token string) error {
	var id int64
	q := `SELECT fk_user_id FROM users_email_confirm WHERE token=$1 LIMIT 1`
	qu := `UPDATE users SET confirm_email=$1 where id=$2`
	db := GetConnection()
	defer db.Close()
	err := db.QueryRow(
		q,
		token,
	).Scan(&id)
	if err != nil {
		return errors.New("no existe el token")
	}
	_, err2 := db.Exec(
		qu,
		true,
		id,
	)
	if err2 != nil {
		return errors.New("No se ha podido confirmar el email")
	}

	return nil
}
func UserSaveTokenConfirmEmail(userId int64, token string, init bool) (errs error) {
	if init {
		errs = UserSaveTokenConfirmEmailCreate(userId, token)
	} else {
		errs = UserSaveTokenConfirmEmailConfirmed(token)
	}
	return
}
func SetUserRecovery(email, token string) (err error) {
	q := `INSERT INTO users_recovery(email, token) VALUES($1, $2)`
	db := GetConnection()
	defer db.Close()
	_, err = db.Exec(
		q,
		email,
		token,
	)
	return
}
func DeleteUserRecovery(token string) (email string, err error) {
	q := `DELETE FROM users_recovery WHERE token = $1 RETURNING email`
	db := GetConnection()
	defer db.Close()
	err = db.QueryRow(
		q,
		token,
	).Scan(&email)
	return
}

func ChangePassword(email, password string) (err error) {
	if email != "" {
		pswd := GetPwd(password)
		password := HashAndSalt(pswd)
		q := `UPDATE users SET password=$1 WHERE email=$2`
		db := GetConnection()
		defer db.Close()
		_, err = db.Exec(
			q,
			password,
			email,
		)
	}
	return
}
func UserSaveAvatar(idUser int64, avatar string) (err error) {
	q := `UPDATE users SET avatar=$1 WHERE id=$2`
	db := settings.InstanceDb
	_, err = db.Exec(
		q,
		avatar,
		idUser,
	)
	return
}
func UserSaveBanner(idUser int64, banner string) (err error) {
	q := `UPDATE users_profile SET banner=$1 WHERE fk_id_user=$2`
	db := settings.InstanceDb
	_, err = db.Exec(
		q,
		banner,
		idUser,
	)
	return
}
func UserSaveGallery(profile []UserProfileGalleryCreate) (err error) {
	q := `INSERT into users_profile_gallery(fk_id_profile,fk_id_group,title,image) values `
	db := settings.InstanceDb

	vals := []interface{}{}
	var count = 0;
	for _, row := range profile {
		q += fmt.Sprintf("($%d, $%d, $%d, $%d),", count+1, count+2, count+3, count+4)
		count += 4
		vals = append(vals, row.FkIdProfile, row.FkIdGroup, row.Title, row.Image)
	}
	//trim the last ,
	q = q[0 : len(q)-1]
	//prepare the statement
	_, err = db.Exec(q, vals...)
	return
}
