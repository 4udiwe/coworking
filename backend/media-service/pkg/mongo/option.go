package mongodb

import "time"

type Option func(*MongoDB)

func ConnAttempts(a int) Option {
	return func(m *MongoDB) {
		m.connAttempts = a
	}
}

func TimeOut(t time.Duration) Option {
	return func(m *MongoDB) {
		m.connTimeout = t
	}
}
