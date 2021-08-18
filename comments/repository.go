package comments

import (
	"database/sql"
	"libria/models"
	"time"

	log "github.com/sirupsen/logrus"
)

type Repository struct {
	dbClient *sql.DB
}

func NewRepository(dbClient *sql.DB) Repository {
	return Repository{dbClient: dbClient}
}

func (r *Repository) GetAll() ([]models.Comment, error) {
	var comments []models.Comment
	query := `SELECT * FROM comments`
	comments, err := r.fetch(query, "")
	return comments, err
}

func (r *Repository) GetAllByAnswer(answerId string) ([]models.Comment, error) {
	var comments []models.Comment
	query := `SELECT * FROM comments WHERE answer_id = $1`
	comments, err := r.fetch(query, answerId)
	return comments, err
}

func (r *Repository) GetById(id string) (models.Comment, error) {
	var comment models.Comment

	query := `SELECT * FROM comments WHERE comment_id = $1`
	comment, err := r.getOne(query, id)
	return comment, err
}

func (r *Repository) Post(comment *models.Comment) (*models.Comment, error) {
	statement := `INSERT INTO comments (id, answer_id, comment_text, created_date, UserId, Username) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.dbClient.Exec(statement, comment.ID, comment.AnswerID, comment.Text, time.Now(), comment.UserID, comment.Username)
	return comment, err
}

func (r *Repository) Update(comment *models.Comment) (models.Comment, error) {
	query := `UPDATE comments SET comment_text = $1, updated_date = $2 WHERE comment_id = $3`
	_, err := r.dbClient.Exec(query, comment.Text, time.Now(), comment.ID)

	return *comment, err
}

func (r *Repository) Delete(comment models.Comment) (models.Comment, error) {
	query := `DELETE FROM comments WHERE id = $1`
	_, err := r.dbClient.Exec(query, comment.ID)
	return comment, err
}

func (r *Repository) fetch(query string, answerID string) ([]models.Comment, error) {
	var rows *sql.Rows
	var err error
	if len(answerID) > 0 {
		rows, err = r.dbClient.Query(query, answerID)
	} else {
		rows, err = r.dbClient.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Errorf("Datenbankverbindung konnte nicht korrekt geschlossen werden: %v", err)
		}
	}()
	result := make([]models.Comment, 0)
	for rows.Next() {
		commentDB := models.CommentDB{}
		err := rows.Scan(&commentDB.ID, &commentDB.AnswerID, &commentDB.UserID, &commentDB.Username, &commentDB.Text, &commentDB.CreatedDate, &commentDB.UpdatedDate)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			log.Infof("Fehler beim Lesen der Daten: %v", err)
			return result, err
		}
		result = append(result, commentDB.GetComment())
	}
	return result, nil
}

func (r *Repository) getOne(query string, id string) (models.Comment, error) {
	commentDB := models.CommentDB{}
	err := r.dbClient.QueryRow(query, id).Scan(&commentDB.ID, &commentDB.AnswerID, &commentDB.UserID, &commentDB.Username, &commentDB.Text, &commentDB.CreatedDate, &commentDB.UpdatedDate)
	if err != nil && err != sql.ErrNoRows {
		log.Infof("Fehler beim Lesen der Daten: %v", err)
	}
	return commentDB.GetComment(), err
}
