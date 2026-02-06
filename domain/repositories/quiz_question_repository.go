package repositories

import "qlass-be/domain/entities"

type QuizQuestionRepository interface {
	Create(quiz *entities.QuizQuestion) error
	Update(quiz *entities.QuizQuestion) error
	GetByID(id uint) (*entities.QuizQuestion, error)
	GetByQuizID(quizID uint) ([]*entities.QuizQuestion, error)
	GetWithOptionsByQuizID(quizID uint) ([]*entities.QuizQuestion, error)
}
