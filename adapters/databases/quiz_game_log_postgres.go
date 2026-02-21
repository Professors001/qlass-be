package databases

import (
	"qlass-be/domain/entities"
	"qlass-be/domain/repositories"

	"gorm.io/gorm"
)

type postgresQuizGameLogRepository struct {
	db *gorm.DB
}

func NewPostgresQuizGameLogRepository(db *gorm.DB) repositories.QuizGameLogRepository {
	return &postgresQuizGameLogRepository{db: db}
}

func (r *postgresQuizGameLogRepository) Create(log *entities.QuizGameLog) error {
	if err := r.db.Create(log).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresQuizGameLogRepository) Update(log *entities.QuizGameLog) error {
	if err := r.db.Save(log).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresQuizGameLogRepository) GetByID(id uint) (*entities.QuizGameLog, error) {
	var log entities.QuizGameLog
	if err := r.db.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *postgresQuizGameLogRepository) GetByClassMaterialID(classMaterialID uint) ([]*entities.QuizGameLog, error) {
	var logs []*entities.QuizGameLog
	if err := r.db.Where("class_material_id = ?", classMaterialID).Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
