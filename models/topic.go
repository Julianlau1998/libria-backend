package models

import (
	"database/sql"
	"libria/utility"
)

type Topic struct {
	ID              string `json:"id"`
	UserID          string `json:"user_id"`
	Title           string `json:"title"`
	Body            string `json:"body"`
	Username        string `json:"username"`
	CreatedDate     string `json:"created_date"`
	UpdatedDate     string `json:"updated_date"`
	AmountOfAnswers int    `json:"amount_of_answers"`
	Reported        bool   `json:"reported"`
}

type TopicDB struct {
	ID          string
	UserID      sql.NullString
	Title       sql.NullString
	Body        sql.NullString
	Username    sql.NullString
	CreatedDate sql.NullString
	UpdatedDate sql.NullString
	Reported    sql.NullBool
}

func (dbV *TopicDB) GetTopic() (t Topic) {
	t.ID = dbV.ID
	t.UserID = utility.GetStringValue(dbV.UserID)
	t.Title = utility.GetStringValue(dbV.Title)
	t.Body = utility.GetStringValue(dbV.Body)
	t.Username = utility.GetStringValue(dbV.Username)
	t.CreatedDate = utility.GetStringValue(dbV.CreatedDate)
	t.UpdatedDate = utility.GetStringValue(dbV.UpdatedDate)
	t.Reported = dbV.Reported.Bool
	return t
}
