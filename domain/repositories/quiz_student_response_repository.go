package repositories

import "qlass-be/domain/entities"

type QuizStudentResponseRepository interface {
	Create(entity *entities.QuizStudentResponse) error
	GetByGameLogIDAndUserID(gameLogID uint, userID uint) ([]*entities.QuizStudentResponse, error)
}
