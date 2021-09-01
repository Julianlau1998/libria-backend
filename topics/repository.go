package topics

import (
	"database/sql"
	"fmt"
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

func (r *Repository) GetAll(limit int, offset int) ([]models.Topic, error) {
	var topics []models.Topic
	var query string
	if limit == 0 {
		query = `SELECT topic_id, Username, title, body, reported, created_date, updated_date FROM topics ORDER BY created_date DESC`
	} else {
		query = `SELECT topic_id, Username, title, body, reported, created_date, updated_date FROM topics ORDER BY created_date DESC LIMIT $1 Offset $2`
	}

	topics, err := r.fetch(query, limit, offset)
	return topics, err
}

func (r *Repository) CountAll() (int, error) {
	var amount int
	row := r.dbClient.QueryRow("SELECT COUNT(*) FROM topics")
	err := row.Scan(&amount)
	return amount, err
}

func (r *Repository) GetReported() ([]models.Topic, error) {
	var topics []models.Topic
	query := `SELECT topic_id, Username, title, body, reported, created_date, updated_date FROM topics WHERE reported = true`
	topics, err := r.fetch(query, 0, 0)
	return topics, err
}

func (r *Repository) GetById(id string) (models.Topic, error) {
	var topic models.Topic
	query := `SELECT topic_id, Username, title, body, reported, created_date, updated_date FROM topics WHERE topic_id = $1`
	topic, err := r.getOne(query, id)
	return topic, err
}

func (r *Repository) Post(topic *models.Topic) (*models.Topic, error) {
	statement := `INSERT INTO topics (topic_id, title, body, created_date, UserId, Username) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.dbClient.Exec(statement, topic.ID, topic.Title, topic.Body, time.Now(), topic.UserID, topic.Username)
	return topic, err
}

func (r *Repository) Update(topic *models.Topic) (models.Topic, error) {
	query := `UPDATE topics SET title = $1, body = $2, updated_date = $3 WHERE topic_id = $4`
	_, err := r.dbClient.Exec(query, topic.Title, topic.Body, time.Now(), topic.ID)

	return *topic, err
}

func (r *Repository) UpdateBestAnswer(topic *models.Topic) (models.Topic, error) {
	query := `UPDATE topics SET body = $1, updated_date = $2 WHERE topic_id = $3`
	_, err := r.dbClient.Exec(query, topic.Body, time.Now(), topic.ID)

	return *topic, err
}

func (r *Repository) Delete(topic models.Topic) (models.Topic, error) {
	query := `DELETE FROM topics WHERE topic_id = $1 AND UserID = $2`
	_, err := r.dbClient.Exec(query, topic.ID, topic.UserID)
	return topic, err
}

func (r *Repository) DeleteAsAdmin(topic models.Topic) (models.Topic, error) {
	query := `DELETE FROM topics WHERE topic_id = $1`
	_, err := r.dbClient.Exec(query, topic.ID)
	return topic, err
}

func (r *Repository) fetch(query string, limit int, offset int) ([]models.Topic, error) {
	var rows *sql.Rows
	var err error
	if limit != 0 {
		fmt.Print(limit)
		rows, err = r.dbClient.Query(query, limit, offset)
	} else {
		fmt.Print(limit)
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
	result := make([]models.Topic, 0)
	for rows.Next() {
		topicDB := models.TopicDB{}
		err := rows.Scan(&topicDB.ID, &topicDB.Username, &topicDB.Title, &topicDB.Body, &topicDB.Reported, &topicDB.CreatedDate, &topicDB.UpdatedDate)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			log.Infof("Fehler beim Lesen der Daten: %v", err)
			return result, err
		}
		result = append(result, topicDB.GetTopic())
	}
	return result, nil
}

func (r *Repository) getOne(query string, id string) (models.Topic, error) {
	topicDB := models.TopicDB{}
	err := r.dbClient.QueryRow(query, id).Scan(&topicDB.ID, &topicDB.Username, &topicDB.Title, &topicDB.Body, &topicDB.Reported, &topicDB.CreatedDate, &topicDB.UpdatedDate)
	if err != nil && err != sql.ErrNoRows {
		log.Infof("Fehler beim Lesen der Daten: %v", err)
	}
	return topicDB.GetTopic(), err
}
