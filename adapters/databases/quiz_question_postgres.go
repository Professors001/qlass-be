package databases

import (
	"qlass-be/domain/entities"

	"gorm.io/gorm"
)

type postgresQuizQuestionRepository struct {
	db *gorm.DB
}

func NewPostgresQuizQuestionRepository(db *gorm.DB) *postgresQuizQuestionRepository {
	return &postgresQuizQuestionRepository{db: db}
}

func (r *postgresQuizQuestionRepository) Create(question *entities.QuizQuestion) error {
	if err := r.db.Create(question).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresQuizQuestionRepository) Update(question *entities.QuizQuestion) error {
	if err := r.db.Save(question).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresQuizQuestionRepository) GetByID(id uint) (*entities.QuizQuestion, error) {
	var question entities.QuizQuestion

	if err := r.db.Preload("Options").Where("id = ?", id).First(&question).Error; err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *postgresQuizQuestionRepository) GetByQuizID(quizID uint) ([]*entities.QuizQuestion, error) {
	var questions []*entities.QuizQuestion

	if err := r.db.Preload("Options").Where("quiz_id = ?", quizID).Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}
