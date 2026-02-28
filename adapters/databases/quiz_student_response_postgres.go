package databases

import (
	"qlass-be/domain/entities"
	"qlass-be/domain/repositories"

	"gorm.io/gorm"
)

type postgresQuizStudentResponseRepository struct {
	db *gorm.DB
}

func NewPostgresQuizStudentResponseRepository(db *gorm.DB) repositories.QuizStudentResponseRepository {
	return &postgresQuizStudentResponseRepository{db: db}
}

func (r *postgresQuizStudentResponseRepository) Create(response *entities.QuizStudentResponse) error {
	if err := r.db.Create(response).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresQuizStudentResponseRepository) Update(response *entities.QuizStudentResponse) error {
	if err := r.db.Save(response).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresQuizStudentResponseRepository) GetByID(id uint) (*entities.QuizStudentResponse, error) {
	var response entities.QuizStudentResponse
	if err := r.db.First(&response, id).Error; err != nil {
		return nil, err
	}
	return &response, nil
}

func (r *postgresQuizStudentResponseRepository) GetByStudentID(studentID uint) ([]*entities.QuizStudentResponse, error) {
	var responses []*entities.QuizStudentResponse

	if err := r.db.Where("student_id = ?", studentID).Find(&responses).Error; err != nil {
		return nil, err
	}
	return responses, nil
}

func (r *postgresQuizStudentResponseRepository) GetByQuestionID(questionID uint) ([]*entities.QuizStudentResponse, error) {
	var responses []*entities.QuizStudentResponse

	if err := r.db.Where("question_id = ?", questionID).Find(&responses).Error; err != nil {
		return nil, err
	}
	return responses, nil
}

func (r *postgresQuizStudentResponseRepository) GetByQuizGameLogID(quizGameLogID uint) ([]*entities.QuizStudentResponse, error) {
	var responses []*entities.QuizStudentResponse
	if err := r.db.Where("quiz_game_log_id = ?", quizGameLogID).Find(&responses).Error; err != nil {
		return nil, err
	}
	return responses, nil
}
