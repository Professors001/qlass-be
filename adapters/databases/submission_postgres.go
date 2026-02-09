package databases

import (
	"errors"
	"qlass-be/domain/entities"
	"qlass-be/domain/repositories"

	"gorm.io/gorm"
)

type postgresSubmissionRepository struct {
	db *gorm.DB
}

func NewPostgresSubmissionRepository(db *gorm.DB) repositories.SubmissionRepository {
	return &postgresSubmissionRepository{db: db}
}

func (r *postgresSubmissionRepository) Create(submission *entities.Submission) error {
	if err := r.db.Create(submission).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresSubmissionRepository) Update(submission *entities.Submission) error {
	if err := r.db.Save(submission).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresSubmissionRepository) GetByID(id uint) (*entities.Submission, error) {
	var submission entities.Submission
	if err := r.db.First(&submission, id).Error; err != nil {
		return nil, err
	}
	return &submission, nil
}

func (r *postgresSubmissionRepository) GetByClassMaterialIDAndStudentID(classMaterialID uint, studentID uint) (*entities.Submission, error) {
	var submission entities.Submission

	// Query the database
	err := r.db.Debug().Where("class_material_id = ? AND user_id = ?", classMaterialID, studentID).First(&submission).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}
	return &submission, nil
}

func (r *postgresSubmissionRepository) GetByClassMaterialID(classMaterialID uint) ([]*entities.Submission, error) {
	var submissions []*entities.Submission
	if err := r.db.Where("class_material_id = ?", classMaterialID).Find(&submissions).Error; err != nil {
		return nil, err
	}
	return submissions, nil
}

func (r *postgresSubmissionRepository) GetByStudentID(studentID uint) ([]*entities.Submission, error) {
	var submissions []*entities.Submission
	if err := r.db.Where("user_id = ?", studentID).Find(&submissions).Error; err != nil {
		return nil, err
	}
	return submissions, nil
}
