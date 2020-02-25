package models

import (
	"../settings"
)

func CronRenovated() {
	db := settings.InstanceDb
	_, e := db.Exec(
		"SELECT renovated_users_month()",
	)
	if e != nil {
		println(e.Error())
	}
}
