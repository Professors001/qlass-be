package repositories

type EnrollRepository interface {
	EnrollStudent(classID uint, studentID uint) error
	GetStudentsByClassID(classID uint) ([]uint, error)
	RemoveStudent(classID uint, studentID uint) error
}