package topics

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

func (r *Repository) GetAll() ([]models.Topic, error) {
	var topics []models.Topic
	query := `SELECT * FROM topics`
	topics, err := r.fetch(query)
	return topics, err
}

func (r *Repository) GetById(id string) (models.Topic, error) {
	var topic models.Topic

	query := `SELECT * FROM topics WHERE topic_id = $1`
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

func (r *Repository) Delete(topic models.Topic) (models.Topic, error) {
	query := `DELETE FROM topics WHERE topic_id = $1`
	_, err := r.dbClient.Exec(query, topic.ID)
	return topic, err
}

func (r *Repository) fetch(query string) ([]models.Topic, error) {
	rows, err := r.dbClient.Query(query)
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
		err := rows.Scan(&topicDB.ID, &topicDB.UserID, &topicDB.Username, &topicDB.Title, &topicDB.Body, &topicDB.CreatedDate, &topicDB.UpdatedDate)
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
	err := r.dbClient.QueryRow(query, id).Scan(&topicDB.ID, &topicDB.UserID, &topicDB.Username, &topicDB.Title, &topicDB.Body, &topicDB.CreatedDate, &topicDB.UpdatedDate)
	if err != nil && err != sql.ErrNoRows {
		log.Infof("Fehler beim Lesen der Daten: %v", err)
	}
	return topicDB.GetTopic(), err
}
