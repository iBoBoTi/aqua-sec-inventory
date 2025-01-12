package usecase

import (
    "errors"

    "github.com/iBoBoTi/aqua-sec-inventory/internal/notification-service/domain"
    "github.com/iBoBoTi/aqua-sec-inventory/internal/notification-service/repository"
)

type NotificationUsecase interface {
    CreateNotification(userID int64, message string) (*domain.Notification, error)
    GetAllNotifications(userID int64) ([]domain.Notification, error)
    ClearNotification(notificationID int64) error
    ClearAllNotifications(userID int64) error
}

type notificationUC struct {
    notificationRepo repository.NotificationRepository
}

func NewNotificationUsecase(notificationRepo repository.NotificationRepository) NotificationUsecase {
    return &notificationUC{
        notificationRepo: notificationRepo,
    }
}

func (uc *notificationUC) CreateNotification(userID int64, message string) (*domain.Notification, error) {
    if userID <= 0 {
        return nil, errors.New("invalid user id")
    }
    if message == "" {
        return nil, errors.New("empty notification message")
    }

    n := &domain.Notification{
        UserID:  userID,
        Message: message,
    }
    if err := uc.notificationRepo.Create(n); err != nil {
        return nil, err
    }
    return n, nil
}

func (uc *notificationUC) GetAllNotifications(userID int64) ([]domain.Notification, error) {
    if userID <= 0 {
        return nil, errors.New("invalid user id")
    }
    return uc.notificationRepo.GetAllByUserID(userID)
}

func (uc *notificationUC) ClearNotification(notificationID int64) error {
    if notificationID <= 0 {
        return errors.New("invalid notification id")
    }
    return uc.notificationRepo.DeleteByID(notificationID)
}

func (uc *notificationUC) ClearAllNotifications(userID int64) error {
    if userID <= 0 {
        return errors.New("invalid user id")
    }
    return uc.notificationRepo.DeleteAllByUserID(userID)
}
