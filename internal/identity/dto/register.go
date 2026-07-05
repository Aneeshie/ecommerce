package dto

type RegisterRequest struct {
	Name     string
	Email    string
	Password string
}

type RegisterResponse struct {
	Message string
}
