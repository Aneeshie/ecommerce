package token

type Manager struct {
	secret []byte
}

func NewManager(secret string) *Manager {
	return &Manager{
		secret: []byte(secret),
	}
}
