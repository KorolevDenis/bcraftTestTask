package models

import "database/sql"

type DTORecipe struct {
	Id          int64     `json:"id"`
	Title       string    `json:"title"`
	Info        string    `json:"info"`
	Ingredients string    `json:"ingredients"`
	Steps       []DTOStep `json:"steps"`
}
type DTOStep struct {
	Id   sql.NullInt64  `json:"id"`
	Info sql.NullString `json:"info"`
	Time sql.NullInt64  `json:"time"`
}
