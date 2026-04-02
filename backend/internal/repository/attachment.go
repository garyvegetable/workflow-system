package repository

import (
	"workflow-system/internal/domain/attachment"

	"gorm.io/gorm"
)

type AttachmentRepository struct {
	db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) *AttachmentRepository {
	return &AttachmentRepository{db: db}
}

func (r *AttachmentRepository) Create(att *attachment.Attachment) error {
	return r.db.Create(att).Error
}

func (r *AttachmentRepository) GetByID(id int64) (*attachment.Attachment, error) {
	var att attachment.Attachment
	err := r.db.First(&att, id).Error
	if err != nil {
		return nil, err
	}
	return &att, nil
}

func (r *AttachmentRepository) ListByInstanceID(instanceID int64) ([]attachment.Attachment, error) {
	var atts []attachment.Attachment
	err := r.db.Where("instance_id = ?", instanceID).Find(&atts).Error
	return atts, err
}

func (r *AttachmentRepository) Delete(id int64) error {
	return r.db.Delete(&attachment.Attachment{}, id).Error
}
