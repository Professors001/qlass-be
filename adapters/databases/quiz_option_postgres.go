package databases

import (
	"qlass-be/domain/entities"
	"qlass-be/domain/repositories"

	"gorm.io/gorm"
)

type postgresQuizOptionRepository struct {
	db *gorm.DB
}

func NewPostgresQuizOptionRepository(db *gorm.DB) repositories.QuizOptionRepository {
	return &postgresQuizOptionRepository{db: db}
}

func (r *postgresQuizOptionRepository) Create(option *entities.QuizOption) (uint, error) {
	if err := r.db.Create(option).Error; err != nil {
		return 0, err
	}
	return option.ID, nil
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

	if err := r.db.Where("question_id = ?", questionID).Find(&options).Error; err != nil {
		return nil, err
	}
	return options, nil
}

func (r *postgresQuizOptionRepository) DeleteByQuestionID(questionID uint) error {
	if err := r.db.Where("question_id = ?", questionID).Delete(&entities.QuizOption{}).Error; err != nil {
		return err
	}
	return nil
}
