package repositories

import "qlass-be/domain/entities"

type ClassMaterialRepository interface {
	Create(classMaterial *entities.ClassMaterial) error
	GetByID(id uint) (*entities.ClassMaterial, error)
	GetByClassID(classID uint) ([]*entities.ClassMaterial, error)
	Update(classMaterial *entities.ClassMaterial) error
	Delete(id uint) error
}
