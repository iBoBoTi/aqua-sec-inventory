package repository

import (
    "database/sql"

    "github.com/iBoBoTi/aqua-sec-inventory/internal/domain"
)

type NotificationRepository interface {
    Create(notification *domain.Notification) error
    GetAllByUserID(userID int64) ([]domain.Notification, error)
    DeleteByID(notificationID int64) error
    DeleteAllByUserID(userID int64) error
}

type notificationRepo struct {
    db *sql.DB
}

func NewNotificationRepository(db *sql.DB) NotificationRepository {
    return &notificationRepo{db: db}
}

func (r *notificationRepo) Create(n *domain.Notification) error {
    query := `INSERT INTO notifications (user_id, message, created_at) VALUES ($1, $2, NOW()) RETURNING id`
    return r.db.QueryRow(query, n.UserID, n.Message).Scan(&n.ID)
}

func (r *notificationRepo) GetAllByUserID(userID int64) ([]domain.Notification, error) {
    query := `SELECT id, user_id, message, created_at FROM notifications WHERE user_id = $1`
    rows, err := r.db.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var notifs []domain.Notification
    for rows.Next() {
        var n domain.Notification
        if err := rows.Scan(&n.ID, &n.UserID, &n.Message, &n.CreatedAt); err != nil {
            return nil, err
        }
        notifs = append(notifs, n)
    }
    return notifs, nil
}

func (r *notificationRepo) DeleteByID(notificationID int64) error {
    query := `DELETE FROM notifications WHERE id = $1`
    _, err := r.db.Exec(query, notificationID)
    return err
}

func (r *notificationRepo) DeleteAllByUserID(userID int64) error {
    query := `DELETE FROM notifications WHERE user_id = $1`
    _, err := r.db.Exec(query, userID)
    return err
}
