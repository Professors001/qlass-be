package databases

import (
	"qlass-be/domain/entities"
	"qlass-be/domain/repositories"

	"gorm.io/gorm"
)

type postgresClassMaterialRepository struct {
	db *gorm.DB
}

func NewPostgresClassMaterialRepository(db *gorm.DB) repositories.ClassMaterialRepository {
	return &postgresClassMaterialRepository{db: db}
}

func (r *postgresClassMaterialRepository) Create(classMaterial *entities.ClassMaterial) error {
	if err := r.db.Create(classMaterial).Error; err != nil {
		return err
	}
	return nil
}
func (r *postgresClassMaterialRepository) GetByID(id uint) (*entities.ClassMaterial, error) {
	var classMaterial entities.ClassMaterial
	err := r.db.Preload("Class.Owner.ProfileImgAttachment").
		First(&classMaterial, id).Error

	if err != nil {
		return nil, err
	}
	return &classMaterial, nil
}

func (r *postgresClassMaterialRepository) GetByClassID(classID uint) ([]*entities.ClassMaterial, error) {
	var classMaterials []*entities.ClassMaterial
	err := r.db.Preload("Class.Owner.ProfileImgAttachment").
		Where("class_id = ?", classID).
		Order("created_at desc").
		Find(&classMaterials).Error

	if err != nil {
		return nil, err
	}
	return classMaterials, nil
}
func (r *postgresClassMaterialRepository) Update(classMaterial *entities.ClassMaterial) error {
	if err := r.db.Save(classMaterial).Error; err != nil {
		return err
	}
	return nil
}
func (r *postgresClassMaterialRepository) Delete(id uint) error {
	var classMaterial entities.ClassMaterial
	classMaterial.ID = uint(id)
	if err := r.db.Delete(&classMaterial).Error; err != nil {
		return err
	}
	return nil
}
