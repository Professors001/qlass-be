package repositories

import "qlass-be/domain/entities"

type ClassRepository interface {
	Create(class *entities.Class) error
	GetByID(id uint) (*entities.Class, error)
	GetByInviteCode(code string) (*entities.Class, error)
	GetByUserID(userID uint) ([]entities.ClassEnrollment, error)
}
