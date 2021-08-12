package topics

import (
	"fmt"
	"libria/answers"
	"libria/models"
	"math/rand"

	uuid "github.com/nu7hatch/gouuid"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	topicRepo     Repository
	answerService answers.Service
}

func NewService(topicRepository Repository, answerService answers.Service) Service {
	return Service{
		topicRepo:     topicRepository,
		answerService: answerService,
	}
}

func (s *Service) GetAll() ([]models.Topic, error) {
	topics, err := s.topicRepo.GetAll()
	if err != nil {
		log.Warnf("topicsService GetAll(), could not load topics: %s", err)
	}
	for index, topic := range topics {
		answers, err := s.answerService.GetAllByTopic(topic.ID, "")
		if err != nil {
			log.Warnf("topicsService GetAll(), could not load answers: %s", err)
		}
		topics[index].AmountOfAnswers = len(answers)
	}
	return topics, err
}

func (s *Service) GetById(id string) (models.Topic, error) {
	return s.topicRepo.GetById(id)
}

func (s *Service) GetRandom() (models.Topic, error) {
	topics, err := s.GetAll()
	if err != nil {
		log.Warnf("topicsService GetRandom(), could not load topics: %s", err)
	}
	randomIndex := rand.Intn(len(topics))
	randomTopic := topics[randomIndex]
	return randomTopic, err
}

func (s *Service) Post(topic *models.Topic) (*models.Topic, error) {
	id, err := uuid.NewV4()
	if err != nil {
		fmt.Print(err)
		return topic, err
	}
	topic.ID = id.String()
	newTopic, err := s.topicRepo.Post(topic)

	var answer models.Answer
	answer.Text = topic.Body
	answer.TopicID = id.String()
	answer.UserID = topic.UserID
	answer.Username = topic.Username
	if topic.Body != "" {
		newAnswer, err := s.answerService.Post(&answer)
		if err != nil {
			log.Warnf("TopicService.Post() Could not post bestAnswer: %s", err)
			return newTopic, err
		}
		newTopic.Body = newAnswer.Text
		newTopicAfterUpdate, err := s.Update(newTopic.ID, newTopic)
		if err != nil {
			log.Warnf("TopicService.Post() Could not Update Topic: %s", err)
			return newTopic, err
		}
		return &newTopicAfterUpdate, err
	}
	return newTopic, err
}

func (s *Service) Update(id string, topic *models.Topic) (models.Topic, error) {
	topic.ID = id
	return s.topicRepo.Update(topic)
}

func (s *Service) Delete(id string) (models.Topic, error) {
	var topic models.Topic
	topic.ID = id
	return s.topicRepo.Delete(topic)
}
