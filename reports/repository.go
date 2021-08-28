package reports

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	dbClient *sql.DB
}

func NewRepository(dbClient *sql.DB) Repository {
	return Repository{dbClient: dbClient}
}

func (r *Repository) Report(id string, contentType string) error {
	var query string
	switch contentType {
	case "topics":
		query = `UPDATE topics SET reported = true WHERE topic_id = $1`
	case "answers":
		query = `UPDATE answers SET reported = true WHERE answer_id = $1`
	case "comments":
		query = `UPDATE comments SET reported = true WHERE id = $1`
	default:
		query = ""
	}
	_, err := r.dbClient.Exec(query, id)
	return err
}

func (r *Repository) Unreport(id string, contentType string) error {
	var query string
	fmt.Print(contentType)
	switch contentType {
	case "topics":
		query = `UPDATE topics SET reported = false WHERE topic_id = $1`
	case "answers":
		query = `UPDATE answers SET reported = false WHERE answer_id = $1`
	case "comments":
		query = `UPDATE comments SET reported = false WHERE id = $1`
	default:
		query = ""
	}
	_, err := r.dbClient.Exec(query, id)
	return err
}
