package repositories

import "qlass-be/domain/entities"

type AttachmentRepository interface {
	Create(attachment *entities.Attachment) error
	GetByID(id uint) (*entities.Attachment, error)
	GetByOwnerTypeAndOwnerID(ownerType string, ownerID uint) ([]*entities.Attachment, error)
	Update(attachment *entities.Attachment) error
	Delete(id uint) error
}
