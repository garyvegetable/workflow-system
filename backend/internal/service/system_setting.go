package service

import (
	"workflow-system/internal/domain/system"
	"workflow-system/internal/repository"
)

type SystemSettingService struct {
	repo *repository.SystemSettingRepository
}

func NewSystemSettingService(repo *repository.SystemSettingRepository) *SystemSettingService {
	return &SystemSettingService{repo: repo}
}

func (s *SystemSettingService) Get(key string) (string, error) {
	return s.repo.Get(key)
}

func (s *SystemSettingService) Set(key, value string) error {
	return s.repo.Set(key, value)
}

func (s *SystemSettingService) Delete(key string) error {
	return s.repo.Delete(key)
}

func (s *SystemSettingService) ListAll() ([]system.SystemSetting, error) {
	return s.repo.ListAll()
}

// GetSMTPConfig 获取 SMTP 配置
func (s *SystemSettingService) GetSMTPConfig() (map[string]string, error) {
	config := make(map[string]string)
	keys := []string{
		system.KeySMTPHost,
		system.KeySMTPPort,
		system.KeySMTPUser,
		system.KeySMTPPassword,
		system.KeySMTPFrom,
	}
	for _, key := range keys {
		val, err := s.repo.Get(key)
		if err != nil {
			config[key] = ""
		} else {
			config[key] = val
		}
	}
	return config, nil
}
