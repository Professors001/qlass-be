package databases

import (
	"qlass-be/domain/entities"

	"gorm.io/gorm"
)

type postgresEnrollRepository struct {
	db *gorm.DB
}

// NewPostgresEnrollRepository creates a new instance of the repository
func NewPostgresEnrollRepository(db *gorm.DB) *postgresEnrollRepository {
	return &postgresEnrollRepository{db: db}
}

func (r *postgresEnrollRepository) EnrollStudent(classID uint, studentID uint) error {
	enrollment := &entities.ClassEnrollment{
		ClassID: classID,
		UserID:  studentID,
	}

	if err := r.db.Create(enrollment).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresEnrollRepository) EnrollWithRole(classID uint, userID uint, role string) error {
	enrollment := &entities.ClassEnrollment{
		ClassID: classID,
		UserID:  userID,
		Role:    role,
	}

	if err := r.db.Create(enrollment).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresEnrollRepository) GetEnrolledStudentsByClassID(classID uint) ([]entities.ClassEnrollment, error) {
	var enrollments []entities.ClassEnrollment

	if err := r.db.Preload("User").Where("class_id = ?", classID).Find(&enrollments).Error; err != nil {
		return nil, err
	}
	return enrollments, nil
}

func (r *postgresEnrollRepository) RemoveStudent(classID uint, studentID uint) error {
	if err := r.db.Where("class_id = ? AND user_id = ?", classID, studentID).Delete(&entities.ClassEnrollment{}).Error; err != nil {
		return err
	}
	return nil
}
