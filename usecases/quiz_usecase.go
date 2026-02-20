package usecases

import (
	"qlass-be/domain/repositories"
	"qlass-be/dtos"
	"qlass-be/transforms"
	"qlass-be/utils"
)

type QuizUseCase interface {
	CreateQuiz(dto dtos.SaveQuizDto, userID uint) (uint, error)
	UpdateQuiz(dto dtos.SaveQuizDto, quizID uint) error
	SaveQuizQuestion(dto dtos.SaveQuizQuestionDtoRequest, quizID uint) error
	GetQuizByID(id uint) (*dtos.GetQuizResponseDto, error)
	GetQuizzesByUserID(userID uint) ([]*dtos.GetQuizResponseDto, error)
}

type QuizUsecase struct {
	quizRepo          repositories.QuizRepository
	quizQuestionRepo  repositories.QuizQuestionRepository
	quizOptionRepo    repositories.QuizOptionRepository
	attachmentRepo    repositories.AttachmentRepository
	attachmentUseCase AttachmentUseCase
}

func NewQuizUseCase(quizRepo repositories.QuizRepository, quizQuestionRepo repositories.QuizQuestionRepository, quizOptionRepo repositories.QuizOptionRepository, attachmentRepo repositories.AttachmentRepository, attachmentUseCase AttachmentUseCase) QuizUseCase {
	return &QuizUsecase{
		quizRepo:          quizRepo,
		quizQuestionRepo:  quizQuestionRepo,
		quizOptionRepo:    quizOptionRepo,
		attachmentRepo:    attachmentRepo,
		attachmentUseCase: attachmentUseCase,
	}
}

// Implementations of the QuizUseCase methods would go here
func (u *QuizUsecase) CreateQuiz(dto dtos.SaveQuizDto, userID uint) (uint, error) {
	// Create the Quiz entity
	quiz := transforms.SaveQuizDtoToQuizEntity(dto, userID)

	quizID, err := u.quizRepo.Create(quiz)
	if err != nil {
		return 0, err
	}

	return quizID, nil
}

func (u *QuizUsecase) UpdateQuiz(dto dtos.SaveQuizDto, quizID uint) error {
	// Update the Quiz entity
	quiz := transforms.SaveQuizDtoToQuizEntity(dto, quizID)

	err := u.quizRepo.Update(quiz)
	if err != nil {
		return err
	}

	return nil
}

func (u *QuizUsecase) SaveQuizQuestion(dto dtos.SaveQuizQuestionDtoRequest, quizID uint) error {
	// Check is these exits?
	questions, err := u.quizQuestionRepo.GetByQuizID(quizID)
	if err != nil {
		return err
	}

	if len(questions) > 0 {
		for _, q := range questions {
			attachments, err := u.attachmentRepo.GetByOwnerTypeAndOwnerID("quiz_question", q.ID)
			if err == nil {
				for _, att := range attachments {
					att.OwnerID = nil
					att.OwnerType = nil
					if err := u.attachmentRepo.Update(att); err != nil {
						return err
					}
				}
			}
		}

		err = u.quizQuestionRepo.DeleteByQuizID(quizID)
		if err != nil {
			return err
		}
	}

	// Creata a question
	for _, question := range dto.Questions {
		questionEntity := transforms.SaveQuizQuestionDtoToQuizQuestionEntity(question, quizID)

		questionID, err := u.quizQuestionRepo.Create(questionEntity)
		if err != nil {
			return err
		}

		attachmentID := question.MediaAttachmentID
		if attachmentID != nil {
			attachment, err := u.attachmentRepo.GetByID(*attachmentID)
			if err != nil {
				return err
			}

			attachment.OwnerType = utils.Ptr("quiz_question")
			attachment.OwnerID = &questionID

			err = u.attachmentRepo.Update(attachment)
			if err != nil {
				return err
			}
		}

		// Create Options
		for _, option := range question.Options {
			optionEntity := transforms.SaveQuizOptionDtoToQuizOptionEntity(option, questionID)

			_, err := u.quizOptionRepo.Create(optionEntity)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (u *QuizUsecase) GetQuizByID(id uint) (*dtos.GetQuizResponseDto, error) {
	quiz, err := u.quizRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	questions, err := u.quizQuestionRepo.GetWithOptionsByQuizID(quiz.ID)
	if err != nil {
		return nil, err
	}

	var questionDtos []dtos.GetQuizQuestionResponse
	for _, question := range questions {
		var attachmentDto *dtos.GetAttachmentResponseDto

		attachments, err := u.attachmentUseCase.GetAttachmentsByOwner("quiz_question", question.ID)
		if err != nil {
			return nil, err
		}

		if len(attachments) > 0 {
			attachmentDto = attachments[0]
		}

		var optionDtos []dtos.GetQuizOptionResponse
		for _, option := range question.Options {
			optionDtos = append(optionDtos, dtos.GetQuizOptionResponse{
				ID:         option.ID,
				OptionText: option.OptionText,
				IsCorrect:  option.IsCorrect,
				OrderIndex: option.OrderIndex,
			})
		}

		questionDtos = append(questionDtos, dtos.GetQuizQuestionResponse{
			ID:               question.ID,
			QuestionText:     question.QuestionText,
			MediaAttachment:  attachmentDto,
			PointsMultiplier: question.PointsMultiplier,
			TimeLimitSeconds: question.TimeLimitSeconds,
			OrderIndex:       question.OrderIndex,
			Options:          optionDtos,
		})
	}

	quizDto := &dtos.GetQuizResponseDto{
		ID:                     quiz.ID,
		UserID:                 quiz.UserID,
		Title:                  quiz.Title,
		Description:            quiz.Description,
		DefaultTimePerQuestion: quiz.DefaultTimePerQuestion,
		Questions:              questionDtos,
	}

	return quizDto, nil
}

func (u *QuizUsecase) GetQuizzesByUserID(userID uint) ([]*dtos.GetQuizResponseDto, error) {

	quizzes, err := u.quizRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var quizDtos []*dtos.GetQuizResponseDto
	for _, quiz := range quizzes {
		quizDto, err := u.GetQuizByID(quiz.ID)
		if err != nil {
			return nil, err
		}
		quizDtos = append(quizDtos, quizDto)
	}

	return quizDtos, nil

}
