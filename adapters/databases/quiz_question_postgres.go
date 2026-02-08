package databases

import (
	"qlass-be/domain/entities"
	"qlass-be/domain/repositories"

	"gorm.io/gorm"
)

type postgresQuizQuestionRepository struct {
	db *gorm.DB
}

func NewPostgresQuizQuestionRepository(db *gorm.DB) repositories.QuizQuestionRepository {
	return &postgresQuizQuestionRepository{db: db}
}

func (r *postgresQuizQuestionRepository) Create(question *entities.QuizQuestion) (uint, error) {
	if err := r.db.Create(question).Error; err != nil {
		return 0, err
	}
	return question.ID, nil
}

func (r *postgresQuizQuestionRepository) Update(question *entities.QuizQuestion) error {
	if err := r.db.Save(question).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresQuizQuestionRepository) GetByID(id uint) (*entities.QuizQuestion, error) {
	var question entities.QuizQuestion

	if err := r.db.Where("id = ?", id).First(&question).Error; err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *postgresQuizQuestionRepository) GetByQuizID(quizID uint) ([]*entities.QuizQuestion, error) {
	var questions []*entities.QuizQuestion

	if err := r.db.Where("quiz_id = ?", quizID).Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

func (r *postgresQuizQuestionRepository) GetWithOptionsByQuizID(quizID uint) ([]*entities.QuizQuestion, error) {
	var questions []*entities.QuizQuestion

	if err := r.db.Preload("Options").Where("quiz_id = ?", quizID).Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

func (r *postgresQuizQuestionRepository) DeleteByQuizID(quizID uint) error {
	if err := r.db.Where("quiz_id = ?", quizID).Delete(&entities.QuizQuestion{}).Error; err != nil {
		return err
	}
	return nil
}
