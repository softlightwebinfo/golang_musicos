package models

type UserConfirmEmailToken struct {
	Success bool    `json:"success"`
	Token   string  `json:"token"`
	Err     *string `json:"err"`
}
