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

func (r *postgresQuizStudentResponseRepository) Create(entity *entities.QuizStudentResponse) error {
	return r.db.Create(entity).Error
}

func (r *postgresQuizStudentResponseRepository) GetByGameLogIDAndUserID(gameLogID uint, userID uint) ([]*entities.QuizStudentResponse, error) {
	var responses []*entities.QuizStudentResponse
	err := r.db.Where("quiz_game_log_id = ? AND user_id = ?", gameLogID, userID).Find(&responses).Error
	return responses, err
}
