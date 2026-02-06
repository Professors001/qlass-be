package databases

import (
	"qlass-be/domain/entities"
	"qlass-be/domain/repositories"

	"gorm.io/gorm"
)

type postgresQuizRepository struct {
	db *gorm.DB
}

func NewPostgresQuizRepository(db *gorm.DB) repositories.QuizRepository {
	return &postgresQuizRepository{db: db}
}

func (r *postgresQuizRepository) Create(quiz *entities.Quiz) error {
	if err := r.db.Create(quiz).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresQuizRepository) Update(quiz *entities.Quiz) error {
	if err := r.db.Save(quiz).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresQuizRepository) GetByID(id uint) (*entities.Quiz, error) {
	var quiz entities.Quiz

	if err := r.db.Preload("Questions.Options").Where("id = ?", id).First(&quiz).Error; err != nil {
		return nil, err
	}
	return &quiz, nil
}

func (r *postgresQuizRepository) GetByClassMaterialID(classMaterialID uint) (*entities.Quiz, error) {
	var quiz entities.Quiz

	if err := r.db.Preload("Questions.Options").Where("class_material_id = ?", classMaterialID).First(&quiz).Error; err != nil {
		return nil, err
	}
	return &quiz, nil
}
