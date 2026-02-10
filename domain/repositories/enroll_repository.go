package repositories

import "qlass-be/domain/entities"

type EnrollRepository interface {
	EnrollStudent(classID uint, studentID uint) error
	EnrollWithRole(classID uint, userID uint, role string) error
	GetEnrolledStudentsByClassID(classID uint) ([]entities.ClassEnrollment, error)
	RemoveStudent(classID uint, studentID uint) error
	IsEnrolled(classID uint, userID uint) (bool, error)
}
