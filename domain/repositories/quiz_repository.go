package repositories

import "qlass-be/domain/entities"

type QuizRepository interface {
	Create(quiz *entities.Quiz) error
	Update(quiz *entities.Quiz) error
	GetByID(id uint) (*entities.Quiz, error)
	GetByClassMaterialID(classMaterialID uint) (*entities.Quiz, error)
}