package votes

import (
	"database/sql"
	"libria/models"

	log "github.com/sirupsen/logrus"
)

type Repository struct {
	dbClient *sql.DB
}

func NewRepository(dbClient *sql.DB) Repository {
	return Repository{dbClient: dbClient}
}

func (r *Repository) GetAllByAnswer(answerID string) ([]models.Vote, error) {
	var votes []models.Vote
	query := `SELECT * FROM votes WHERE answer_id = $1`
	votes, err := r.fetch(query, answerID)
	return votes, err
}

func (r *Repository) Post(vote *models.Vote) (*models.Vote, error) {
	statement := `INSERT INTO votes (vote_id, answer_id, userID, upvote) VALUES ($1, $2, $3, $4)`
	_, err := r.dbClient.Exec(statement, vote.ID, vote.AnswerID, vote.UserID, vote.Upvote)
	return vote, err
}

func (r *Repository) Update(vote *models.Vote) (models.Vote, error) {
	query := `UPDATE votes SET upvote = $1 WHERE vote_id = $2`
	_, err := r.dbClient.Exec(query, vote.Upvote, vote.ID)

	return *vote, err
}

func (r *Repository) fetch(query string, topicID string) ([]models.Vote, error) {
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
	result := make([]models.Vote, 0)
	for rows.Next() {
		voteDB := models.VoteDB{}
		err := rows.Scan(&voteDB.ID, &voteDB.AnswerID, &voteDB.UserID, &voteDB.Upvote)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			log.Infof("Fehler beim Lesen der Daten: %v", err)
			return result, err
		}
		result = append(result, voteDB.GetVote())
	}
	return result, nil
}

func (r *Repository) getOne(query string, id string) (models.Vote, error) {
	voteDB := models.VoteDB{}
	err := r.dbClient.QueryRow(query, id).Scan(&voteDB.ID, &voteDB.AnswerID, &voteDB.UserID, &voteDB.Upvote)
	if err != nil && err != sql.ErrNoRows {
		log.Infof("Fehler beim Lesen der Daten: %v", err)
	}
	return voteDB.GetVote(), err
}
