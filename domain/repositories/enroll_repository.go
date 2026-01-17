package repositories

import "qlass-be/domain/entities"

type EnrollRepository interface {
	EnrollStudent(classID uint, studentID uint) error
	GetStudentsByClassID(classID uint) (*entities.User, error)
	RemoveStudent(classID uint, studentID uint) error
}
