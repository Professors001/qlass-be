package usecases

import (
	"qlass-be/domain/entities"
	"qlass-be/domain/repositories"
	"qlass-be/dtos"
	"strings"
)

type QuizUseCase interface {
	CreateQuiz(dto dtos.CreateQuizDto) (uint, error)
	GetQuizByID(id uint) (*dtos.GetQuizResponseDto, error)
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
func (u *QuizUsecase) CreateQuiz(dto dtos.CreateQuizDto) (uint, error) {
	// 1. Check if quiz exists for this ClassMaterialID
	quizID := uint(0)
	quiz, err := u.quizRepo.GetByClassMaterialID(dto.ClassMaterialID)
	if err != nil {
		if !strings.Contains(err.Error(), "record not found") {
			return 0, err
		}
		quiz = nil
	}

	if quiz != nil {
		quizID = quiz.ID
		quiz.Title = dto.Title
		quiz.Description = dto.Description
		quiz.DefaultTimePerQuestion = dto.DefaultTimePerQuestion

		if err := u.quizRepo.Update(quiz); err != nil {
			return 0, err
		}

		if err := u.quizQuestionRepo.DeleteByQuizID(quiz.ID); err != nil {
			return 0, err
		}
	} else {
		quiz = &entities.Quiz{
			ClassMaterialID:        dto.ClassMaterialID,
			Title:                  dto.Title,
			Description:            dto.Description,
			DefaultTimePerQuestion: dto.DefaultTimePerQuestion,
		}

		if quizID, err = u.quizRepo.Create(quiz); err != nil {
			return 0, err
		}
	}

	for _, qDto := range dto.Questions {
		question := &entities.QuizQuestion{
			QuizID:            quizID,
			QuestionText:      qDto.QuestionText,
			PointsMultiplier:  qDto.PointsMultiplier,
			TimeLimitSeconds:  qDto.TimeLimitSeconds,
			OrderIndex:        qDto.OrderIndex,
			MediaAttachmentID: qDto.MediaAttachmentID,
		}
		if _, err = u.quizQuestionRepo.Create(question); err != nil {
			return 0, err
		}

		for _, oDto := range qDto.Options {
			option := &entities.QuizOption{
				QuestionID: question.ID,
				OptionText: oDto.OptionText,
				IsCorrect:  oDto.IsCorrect,
				OrderIndex: oDto.OrderIndex,
			}
			if _, err := u.quizOptionRepo.Create(option); err != nil {
				return 0, err
			}
		}
	}

	return quizID, nil
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
		if question.MediaAttachmentID != nil {
			attachment, err := u.attachmentRepo.GetByID(*question.MediaAttachmentID)
			if err != nil {
				return nil, err
			}
			attachmentDto, err = u.attachmentUseCase.GetAttachmentByID(attachment.ID)
			if err != nil {
				return nil, err
			}
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
		ClassMaterialID:        quiz.ClassMaterialID,
		Title:                  quiz.Title,
		Description:            quiz.Description,
		DefaultTimePerQuestion: quiz.DefaultTimePerQuestion,
		Questions:              questionDtos,
	}

	return quizDto, nil
}
