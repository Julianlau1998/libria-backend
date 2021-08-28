package models

import (
	"database/sql"
	"libria/utility"
)

type Vote struct {
	ID       string `json:"id"`
	AnswerID string `json:"answer_id"`
	UserID   string `json:"user_id"`
	Upvote   string `json:"upvote"`
	Reported bool   `json:"reported"`
}

type VoteDB struct {
	ID       string
	AnswerID string
	UserID   sql.NullString
	Upvote   sql.NullString
	Reported sql.NullBool
}

func (dbV *VoteDB) GetVote() (a Vote) {
	a.ID = dbV.ID
	a.AnswerID = dbV.AnswerID
	a.UserID = utility.GetStringValue(dbV.UserID)
	a.Upvote = utility.GetStringValue(dbV.Upvote)
	a.Reported = dbV.Reported.Bool
	return a
}
