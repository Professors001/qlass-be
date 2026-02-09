package databases

import (
	"qlass-be/domain/entities"
	"qlass-be/domain/repositories"

	"gorm.io/gorm"
)

type postgresAttachmentRepository struct {
	db *gorm.DB
}

func NewPostgresAttachmentRepository(db *gorm.DB) repositories.AttachmentRepository {
	return &postgresAttachmentRepository{db: db}
}

func (r *postgresAttachmentRepository) Create(attachment *entities.Attachment) error {
	if err := r.db.Create(attachment).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresAttachmentRepository) GetByID(id uint) (*entities.Attachment, error) {
	var attachment entities.Attachment
	if err := r.db.First(&attachment, id).Error; err != nil {
		return nil, err
	}
	return &attachment, nil
}

func (r *postgresAttachmentRepository) GetByOwnerTypeAndOwnerID(ownerType string, ownerID uint) ([]*entities.Attachment, error) {
	var attachments []*entities.Attachment
	if err := r.db.Where("owner_type = ? AND owner_id = ?", ownerType, ownerID).Find(&attachments).Error; err != nil {
		return nil, err
	}
	return attachments, nil
}

func (r *postgresAttachmentRepository) Update(attachment *entities.Attachment) error {
	if err := r.db.Save(attachment).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresAttachmentRepository) Delete(id uint) error {
	// Soft delete
	var attachment entities.Attachment
	attachment.ID = uint(id)

	if err := r.db.Delete(&attachment).Error; err != nil {
		return err
	}
	return nil
}
