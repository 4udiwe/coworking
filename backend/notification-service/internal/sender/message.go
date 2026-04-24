package sender

type PushMessage struct {
	Token string

	Title string
	Body  string

	Data map[string]string
}
