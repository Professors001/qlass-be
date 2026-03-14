package repositories

import "qlass-be/domain/entities"

type QuizGameLogRepository interface {
	Create(log *entities.QuizGameLog) error
	Update(log *entities.QuizGameLog) error
	GetByID(id uint) (*entities.QuizGameLog, error)
	GetByClassMaterialID(classMaterialID uint) ([]*entities.QuizGameLog, error)
}
