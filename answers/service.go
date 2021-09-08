package answers

import (
	"libria/models"
	"libria/votes"

	uuid "github.com/nu7hatch/gouuid"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	answerRepo  Repository
	voteService votes.Service
}

func NewService(answerRepository Repository, voteService votes.Service) Service {
	return Service{
		answerRepo:  answerRepository,
		voteService: voteService,
	}
}

func (s *Service) GetAll() ([]models.Answer, error) {
	answers, err := s.answerRepo.GetAll()
	if err != nil {
		log.Warnf("AnswerService.GetAll() Could not Load Answers: %s", err)
		return answers, err
	}
	for index, answer := range answers {
		votesByAnswer, err := s.voteService.GetAllByAnswer(answer.ID, false)
		if err != nil {
			log.Warnf("AnswerService.GetAllByTopic() Could no tLoad Votes by Answer: %s", err)
			return answers, err
		}
		answers[index].Votes = votesByAnswer
	}
	return answers, err
}

func (s *Service) GetReported() ([]models.Answer, error) {
	answers, err := s.answerRepo.GetReported()
	if err != nil {
		log.Warnf("AnswerService.GetAll() Could not Load Answers: %s", err)
		return answers, err
	}
	for index, answer := range answers {
		votesByAnswer, err := s.voteService.GetAllByAnswer(answer.ID, false)
		if err != nil {
			log.Warnf("AnswerService.GetAllByTopic() Could no tLoad Votes by Answer: %s", err)
			return answers, err
		}
		answers[index].Votes = votesByAnswer
	}
	return answers, err
}

func (s *Service) GetAllByTopic(topicId string, userId string) ([]models.Answer, error) {
	answers, err := s.answerRepo.GetAllByTopic(topicId)
	if err != nil {
		log.Warnf("AnswerService.GetAllByTopic() Could not Load Answers by topic: %s", err)
		return answers, err
	}
	if userId != "" {
		for index, answer := range answers {
			votesByAnswer, err := s.voteService.GetAllByAnswer(answer.ID, true)
			if err != nil {
				log.Warnf("AnswerService.GetAllByTopic() Could not Load Votes by Answer: %s", err)
				return answers, err
			}
			for _, vote := range votesByAnswer {
				if vote.UserID == userId && vote.Upvote == "true" {
					answers[index].UpvotedByMe = true
					answers[index].DownvotedByMe = false
				} else if vote.UserID == userId && vote.Upvote == "false" {
					answers[index].DownvotedByMe = true
					answers[index].UpvotedByMe = false
				}
				vote.UserID = ""
			}
			answers[index].Votes = votesByAnswer
		}
	} else {
		for index, answer := range answers {
			votesByAnswer, err := s.voteService.GetAllByAnswer(answer.ID, false)
			if err != nil {
				log.Warnf("AnswerService.GetAllByTopic() Could notLoad Votes by Answer: %s", err)
				return answers, err
			}
			answers[index].Votes = votesByAnswer
		}
	}
	return answers, err
}

func (s *Service) GetById(id string) (models.Answer, error) {
	answer, err := s.answerRepo.GetById(id)
	if err != nil {
		log.Warnf("AnswerService.GetById() Could not Load answer by id: %s", err)
		return answer, err
	}
	votesByAnswer, err := s.voteService.GetAllByAnswer(answer.ID, false)
	if err != nil {
		log.Warnf("AnswerService.GetAllByTopic() Could no tLoad Votes by Answer: %s", err)
		return answer, err
	}
	answer.Votes = votesByAnswer
	return answer, err
}

func (s *Service) Post(answer *models.Answer) (*models.Answer, error) {
	id, err := uuid.NewV4()
	if err != nil {
		log.Warnf("AnswerService.Post() Could not create new uuid: %s", err)
		return answer, err
	}
	answer.ID = id.String()
	answer, err = s.answerRepo.Post(answer)
	if err != nil {
		log.Warnf("AnswerService.Post() Could not Post Answer: %s", err)
		return answer, err
	}
	var vote models.Vote
	vote.AnswerID = answer.ID
	vote.UserID = answer.UserID
	vote.Upvote = "true"
	_, err = s.voteService.Post(&vote)
	return answer, err
}

func (s *Service) Update(id string, answer *models.Answer) (models.Answer, error) {
	answer.ID = id
	_, err := s.answerRepo.Update(answer)
	if err != nil {
		log.Warnf("answerService.Update() Could not update answer")
	}
	*answer, err = s.GetById(answer.ID)
	if err != nil {
		log.Warnf("answerService.Update() Could not get answerById")
	}
	return *answer, nil
}

func (s *Service) Delete(id string) (models.Answer, error) {
	var answer models.Answer
	answer.ID = id
	answer, err := s.answerRepo.Delete(answer)
	if err != nil {
		log.Warnf("AnswerService.DeleteReported Could not delete answer")
		return answer, err
	}
	return answer, nil
}

func (s *Service) DeleteReported(id string) (models.Answer, error) {
	answer, err := s.GetById(id)
	if err != nil {
		log.Warnf("answerService.DeleteReported() Could not get by id: %s", err)
		return answer, err
	}
	if answer.Reported == true {
		return s.answerRepo.DeleteReported(answer)
	}
	answer, err = s.answerRepo.DeleteReported(answer)
	if err != nil {
		log.Warnf("AnswerService.DeleteReported Could not delete reported")
		return answer, err
	}
	return answer, nil
}
