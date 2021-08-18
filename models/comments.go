package models

import (
	"database/sql"
	"libria/utility"
)

type Comment struct {
	ID          string `json:"id"`
	AnswerID    string `json:"answer_id"`
	UserID      string `json:"user_id"`
	Text        string `json:"body"`
	Username    string `json:"username"`
	CreatedDate string `json:"created_date"`
	UpdatedDate string `json:"updated_date"`
}

type CommentDB struct {
	ID          string
	AnswerID    string
	UserID      sql.NullString
	Text        sql.NullString
	Username    sql.NullString
	CreatedDate sql.NullString
	UpdatedDate sql.NullString
}

func (dbV *CommentDB) GetComment() (c Comment) {
	c.ID = dbV.ID
	c.UserID = utility.GetStringValue(dbV.UserID)
	c.AnswerID = dbV.AnswerID
	c.Text = utility.GetStringValue(dbV.Text)
	c.Username = utility.GetStringValue(dbV.Username)
	c.CreatedDate = utility.GetStringValue(dbV.CreatedDate)
	c.UpdatedDate = utility.GetStringValue(dbV.UpdatedDate)
	return c
}
