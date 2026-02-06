package usecases

import (
	"qlass-be/domain/repositories"
	"qlass-be/dtos"
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

func (u *QuizUsecase) CreateQuiz(dto dtos.CreateQuizDto) (uint, error) {
	quizID, err := u.quizRepo.Create(dto)
	if err != nil {
		return 0, err
	}

	for _, questionDto := range dto.Questions {
		questionID, err := u.quizQuestionRepo.Create(quizID, questionDto)
		if err != nil {
			return 0, err
		}

		for _, optionDto := range questionDto.Options {
			_, err := u.quizOptionRepo.Create(questionID, optionDto)
			if err != nil {
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

	questions, err := u.quizQuestionRepo.GetByQuizID(id)
	if err != nil {
		return nil, err
	}

	for i, question := range questions {
		options, err := u.quizOptionRepo.GetByQuestionID(question.ID)
		if err != nil {
			return nil, err
		}
		questions[i].Options = options
	}

	return transform.EntityToGetQuizResponseDto(quiz, questions), nil
}
