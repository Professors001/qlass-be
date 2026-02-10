package repositories

import "qlass-be/domain/entities"

type QuizRepository interface {
	Create(quiz *entities.Quiz) (uint, error)
	Update(quiz *entities.Quiz) error
	GetByID(id uint) (*entities.Quiz, error)
	GetByUserID(userID uint) ([]entities.Quiz, error)
}
