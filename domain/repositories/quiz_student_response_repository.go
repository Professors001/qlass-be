package repositories

import "qlass-be/domain/entities"

type QuizStudentResponseRepository interface {
	Create(response *entities.QuizStudentResponse) error
	Update(response *entities.QuizStudentResponse) error
	GetByID(id uint) (*entities.QuizStudentResponse, error)
	GetByStudentID(studentID uint) ([]*entities.QuizStudentResponse, error)
	GetByQuestionID(questionID uint) ([]*entities.QuizStudentResponse, error)
	GetByQuizGameLogID(quizGameLogID uint) ([]*entities.QuizStudentResponse, error)
}
