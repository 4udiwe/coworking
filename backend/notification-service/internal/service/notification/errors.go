package notification_service

import "errors"

var (
	ErrCannotCreateNotification = errors.New("cannot create notification")
	ErrCannotRegisterDevice     = errors.New("cannot register device")
	ErrCannotMarkRead           = errors.New("cannot mark notification as read")
	ErrNotificationNotFound     = errors.New("notification not found")
	ErrCannotFetchNotification  = errors.New("cannot fetch notification")
)
