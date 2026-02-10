package repositories

import "qlass-be/domain/entities"

type SubmissionRepository interface {
	Create(submission *entities.Submission) error
	Update(submission *entities.Submission) error
	GetByID(id uint) (*entities.Submission, error)
	GetByClassMaterialIDAndStudentID(classMaterialID uint, studentID uint) (*entities.Submission, error)
	GetByClassMaterialID(classMaterialID uint) ([]*entities.Submission, error)
	GetByStudentID(studentID uint) ([]*entities.Submission, error)
}
