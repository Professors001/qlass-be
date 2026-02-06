package databases

import (
	"qlass-be/domain/entities"

	"gorm.io/gorm"
)

type postgresQuizOptionRepository struct {
	db *gorm.DB
}

func NewPostgresQuizOptionRepository(db *gorm.DB) *postgresQuizOptionRepository {
	return &postgresQuizOptionRepository{db: db}
}

func (r *postgresQuizOptionRepository) Create(option *entities.QuizOption) error {
	if err := r.db.Create(option).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresQuizOptionRepository) Update(option *entities.QuizOption) error {
	if err := r.db.Save(option).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresQuizOptionRepository) GetByID(id uint) (*entities.QuizOption, error) {
	var option entities.QuizOption

	if err := r.db.Where("id = ?", id).First(&option).Error; err != nil {
		return nil, err
	}
	return &option, nil
}

func (r *postgresQuizOptionRepository) GetByQuestionID(questionID uint) ([]*entities.QuizOption, error) {
	var options []*entities.QuizOption

	if err := r.db.Where("quiz_question_id = ?", questionID).Find(&options).Error; err != nil {
		return nil, err
	}
	return options, nil
}
