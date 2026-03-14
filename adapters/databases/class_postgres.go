package databases

import (
	"qlass-be/domain/entities"
	"qlass-be/domain/repositories"

	"gorm.io/gorm"
)

type postgresClassRepository struct {
	db *gorm.DB
}

// NewPostgresClassRepository creates a new instance of the repository
func NewPostgresClassRepository(db *gorm.DB) repositories.ClassRepository {
	return &postgresClassRepository{db: db}
}

func (r *postgresClassRepository) Create(class *entities.Class) error {
	// GORM handles the SQL Insert automatically
	if err := r.db.Create(class).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresClassRepository) GetByID(id uint) (*entities.Class, error) {
	var class entities.Class

	// Added nested preload for ProfileImgAttachment
	if err := r.db.Preload("Owner.ProfileImgAttachment").Where("id = ?", id).First(&class).Error; err != nil {
		return nil, err
	}
	return &class, nil
}

func (r *postgresClassRepository) GetByInviteCode(code string) (*entities.Class, error) {
	var class entities.Class

	// Added nested preload for ProfileImgAttachment
	if err := r.db.Preload("Owner.ProfileImgAttachment").Where("invite_code = ?", code).First(&class).Error; err != nil {
		return nil, err
	}
	return &class, nil
}

func (r *postgresClassRepository) GetByUserID(userID uint) ([]entities.ClassEnrollment, error) {
	var enrollments []entities.ClassEnrollment

	// Added nested preload: Class -> Owner -> ProfileImgAttachment
	if err := r.db.Preload("Class.Owner.ProfileImgAttachment").Where("user_id = ?", userID).Order("created_at desc").Find(&enrollments).Error; err != nil {
		return nil, err
	}
	return enrollments, nil
}

func (r *postgresClassRepository) Update(class *entities.Class) error {
	// GORM handles the SQL Update automatically
	if err := r.db.Save(class).Error; err != nil {
		return err
	}
	return nil
}
