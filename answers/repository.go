package answers

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

func (r *Repository) GetAll() ([]models.Answer, error) {
	var answers []models.Answer
	query := `SELECT answer_id, topic_id, Username, answer, created_date, updated_date FROM answers`
	answers, err := r.fetch(query, "")
	return answers, err
}

func (r *Repository) GetAllByTopic(topicId string) ([]models.Answer, error) {
	var answers []models.Answer
	query := `SELECT answer_id, topic_id, Username, answer, created_date, updated_date FROM answers WHERE topic_id = $1`
	answers, err := r.fetch(query, topicId)
	return answers, err
}

func (r *Repository) GetById(id string) (models.Answer, error) {
	var answer models.Answer

	query := `SELECT answer_id, topic_id, Username, answer, created_date, updated_date FROM answers WHERE answer_id = $1`
	answer, err := r.getOne(query, id)
	return answer, err
}

func (r *Repository) Post(answer *models.Answer) (*models.Answer, error) {
	statement := `INSERT INTO answers (answer_id, topic_id, answer, created_date, UserId, Username) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.dbClient.Exec(statement, answer.ID, answer.TopicID, answer.Text, time.Now(), answer.UserID, answer.Username)
	return answer, err
}

func (r *Repository) Update(answer *models.Answer) (models.Answer, error) {
	query := `UPDATE answers SET answer = $1, updated_date = $2 WHERE answer_id = $3`
	_, err := r.dbClient.Exec(query, answer.Text, time.Now(), answer.ID)

	return *answer, err
}

func (r *Repository) Delete(answer models.Answer) (models.Answer, error) {
	query := `DELETE FROM answers WHERE answer_id = $1`
	_, err := r.dbClient.Exec(query, answer.ID)
	return answer, err
}

func (r *Repository) fetch(query string, topicID string) ([]models.Answer, error) {
	var rows *sql.Rows
	var err error
	if len(topicID) > 0 {
		rows, err = r.dbClient.Query(query, topicID)
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
	result := make([]models.Answer, 0)
	for rows.Next() {
		answerDB := models.AnswerDB{}
		err := rows.Scan(&answerDB.ID, &answerDB.TopicID, &answerDB.Username, &answerDB.Text, &answerDB.CreatedDate, &answerDB.UpdatedDate)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			log.Infof("Fehler beim Lesen der Daten: %v", err)
			return result, err
		}
		result = append(result, answerDB.GetAnswer())
	}
	return result, nil
}

func (r *Repository) getOne(query string, id string) (models.Answer, error) {
	answerDB := models.AnswerDB{}
	err := r.dbClient.QueryRow(query, id).Scan(&answerDB.ID, &answerDB.TopicID, &answerDB.Username, &answerDB.Text, &answerDB.CreatedDate, &answerDB.UpdatedDate)
	if err != nil && err != sql.ErrNoRows {
		log.Infof("Fehler beim Lesen der Daten: %v", err)
	}
	return answerDB.GetAnswer(), err
}
