package libs

import (
	"../settings"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func GetConnection() *sql.DB {
	var database settings.DatabaseConfig = settings.Database()
	dns := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		database.User,
		database.Password,
		database.Host,
		database.Port,
		database.Database,
	)

	db, err := sql.Open("postgres", dns)

	if err != nil {
		log.Fatal(err)
	}
	return db
}
