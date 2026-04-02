package service

import (
	"workflow-system/internal/domain/notification"
	"workflow-system/internal/repository"
	"workflow-system/internal/pkg/email"
)

type NotificationService struct {
	repo        *repository.NotificationRepository
	emailSvc    *email.EmailService
}

func NewNotificationService(repo *repository.NotificationRepository, emailSvc *email.EmailService) *NotificationService {
	return &NotificationService{
		repo:     repo,
		emailSvc: emailSvc,
	}
}

func (s *NotificationService) Create(userID int64, title, content, notifType string) error {
	notif := &notification.Notification{
		UserID:  userID,
		Title:   title,
		Content: content,
		Type:    notifType,
		IsRead:  false,
	}
	if err := s.repo.Create(notif); err != nil {
		return err
	}

	// 发送邮件通知（Mock）
	go s.emailSvc.Send("user@example.com", title, content)

	return nil
}

func (s *NotificationService) SendToUser(userID int64, message string) error {
	return s.Create(userID, "工作流通知", message, "workflow")
}

func (s *NotificationService) ListByUser(userID int64) ([]notification.Notification, error) {
	return s.repo.ListByUser(userID)
}

func (s *NotificationService) MarkAsRead(id int64) error {
	return s.repo.MarkAsRead(id)
}

func (s *NotificationService) MarkAllAsRead(userID int64) error {
	return s.repo.MarkAllAsRead(userID)
}

func (s *NotificationService) GetUnreadCount(userID int64) (int64, error) {
	return s.repo.GetUnreadCount(userID)
}
