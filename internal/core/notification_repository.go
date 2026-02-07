package database

import (
	"context"
	"fmt"
)

// NotificationRepository handles database operations for notifications.
// This abstraction keeps SQL queries in the Core layer (proper 3-layer architecture).
type NotificationRepository struct{}

// NewNotificationRepository creates a new NotificationRepository.
func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{}
}

// CountUnreadByUserID returns the count of unread notifications for a user.
func (r *NotificationRepository) CountUnreadByUserID(ctx context.Context, userID int) (int, error) {
	db := GetPool()
	if db == nil {
		return 0, fmt.Errorf("database connection not available")
	}

	var count int
	err := db.QueryRow(ctx, `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND read = false`, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count unread notifications: %w", err)
	}

	return count, nil
}
