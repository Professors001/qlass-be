package repositories

import "qlass-be/domain/entities"

type QuizOptionRepository interface {
	Create(option *entities.QuizOption) (uint, error)
	Update(option *entities.QuizOption) error
	GetByID(id uint) (*entities.QuizOption, error)
	GetByQuestionID(questionID uint) ([]*entities.QuizOption, error)
	DeleteByQuestionID(questionID uint) error
}
