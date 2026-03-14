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
	// Preload User and their Profile Image
	if err := r.db.Preload("User.ProfileImgAttachment").First(&submission, id).Error; err != nil {
		return nil, err
	}
	return &submission, nil
}

func (r *postgresSubmissionRepository) GetByClassMaterialIDAndStudentID(classMaterialID uint, studentID uint) (*entities.Submission, error) {
	var submission entities.Submission

	// Preload User and their Profile Image
	err := r.db.Preload("User.ProfileImgAttachment").
		Where("class_material_id = ? AND user_id = ?", classMaterialID, studentID).
		First(&submission).Error

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
	// Preload for every submission in the list
	if err := r.db.Preload("User.ProfileImgAttachment").
		Where("class_material_id = ?", classMaterialID).
		Find(&submissions).Error; err != nil {
		return nil, err
	}
	return submissions, nil
}

func (r *postgresSubmissionRepository) GetByStudentID(studentID uint) ([]*entities.Submission, error) {
	var submissions []*entities.Submission
	// Preload for the student's own view
	if err := r.db.Preload("User.ProfileImgAttachment").
		Where("user_id = ?", studentID).
		Find(&submissions).Error; err != nil {
		return nil, err
	}
	return submissions, nil
}

func (r *postgresSubmissionRepository) FirstOrCreate(submission *entities.Submission, classMaterialID uint, studentID uint) (*entities.Submission, error) {
	var result entities.Submission
	err := r.db.Preload("User.ProfileImgAttachment").
		Where("class_material_id = ? AND user_id = ?", classMaterialID, studentID).
		Attrs(submission).
		FirstOrCreate(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}
