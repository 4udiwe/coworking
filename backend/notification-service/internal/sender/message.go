package sender

type PushMessage struct {
	Token          string
	Title          string
	Body           string
	NotificationID string
	ActionURL      string
	Data           map[string]string
}
