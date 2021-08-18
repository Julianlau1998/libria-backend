package comments

import (
	"libria/models"

	uuid "github.com/nu7hatch/gouuid"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	commentsRepo Repository
}

func NewService(commentsRepository Repository) Service {
	return Service{
		commentsRepo: commentsRepository,
	}
}

func (s *Service) GetAll() ([]models.Comment, error) {
	comments, err := s.commentsRepo.GetAll()
	if err != nil {
		log.Warnf("CommentService.GetAll() Could not Load Comments: %s", err)
		return comments, err
	}
	return comments, err
}

func (s *Service) GetAllByAnswer(answerId string, userId string) ([]models.Comment, error) {
	comments, err := s.commentsRepo.GetAllByAnswer(answerId)
	if err != nil {
		log.Warnf("CommentService.GetAllByAnswer() Could not Load Comments by answer: %s", err)
		return comments, err
	}
	return comments, err
}

func (s *Service) GetById(id string) (models.Comment, error) {
	comment, err := s.commentsRepo.GetById(id)
	if err != nil {
		log.Warnf("CommentService.GetById() Could not Load comment by id: %s", err)
		return comment, err
	}
	return comment, err
}

func (s *Service) Post(comment *models.Comment) (*models.Comment, error) {
	id, err := uuid.NewV4()
	if err != nil {
		log.Warnf("CommentService.Post() Could not create new uuid: %s", err)
		return comment, err
	}
	comment.ID = id.String()
	comment, err = s.commentsRepo.Post(comment)
	if err != nil {
		log.Warnf("CommentService.Post() Could not Post Comment: %s", err)
		return comment, err
	}
	return comment, err
}

func (s *Service) Update(id string, comment *models.Comment) (models.Comment, error) {
	comment.ID = id
	return s.commentsRepo.Update(comment)
}

func (s *Service) Delete(id string) (models.Comment, error) {
	var comment models.Comment
	comment.ID = id
	return s.commentsRepo.Delete(comment)
}
