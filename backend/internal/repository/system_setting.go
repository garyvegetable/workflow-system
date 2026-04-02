package repository

import (
	"workflow-system/internal/domain/system"

	"gorm.io/gorm"
)

type SystemSettingRepository struct {
	db *gorm.DB
}

func NewSystemSettingRepository(db *gorm.DB) *SystemSettingRepository {
	return &SystemSettingRepository{db: db}
}

func (r *SystemSettingRepository) Get(key string) (string, error) {
	var setting system.SystemSetting
	err := r.db.First(&setting, "key = ?", key).Error
	if err != nil {
		return "", err
	}
	return setting.Value, nil
}

func (r *SystemSettingRepository) Set(key, value string) error {
	// 使用 upsert：如果不存在则创建，存在则更新
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existing system.SystemSetting
		err := tx.First(&existing, "key = ?", key).Error
		if err == gorm.ErrRecordNotFound {
			return tx.Create(&system.SystemSetting{Key: key, Value: value}).Error
		} else if err != nil {
			return err
		}
		return tx.Model(&existing).Update("value", value).Error
	})
}

func (r *SystemSettingRepository) Delete(key string) error {
	return r.db.Delete(&system.SystemSetting{}, "key = ?", key).Error
}

func (r *SystemSettingRepository) ListAll() ([]system.SystemSetting, error) {
	var settings []system.SystemSetting
	err := r.db.Find(&settings).Error
	return settings, err
}
