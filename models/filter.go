package models

type FilterPage struct {
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	All     bool   `json:"all"`
	Slug    string `json:"slug"`
	IdUser  int    `json:"id_user"`
	Country string `json:"country"`
	Region  string `json:"region"`
	City    string `json:"city"`
}
