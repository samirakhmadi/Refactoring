package service

import "net/http"

type CreateUserRequest struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

// TODO добавить валидацию на поля
func (c *CreateUserRequest) Bind(r *http.Request) error {
	return nil
}

type UpdateUserRequest struct {
	DisplayName string `json:"display_name"`
}

func (c *UpdateUserRequest) Bind(r *http.Request) error {
	return nil
}

type Status struct {
	id      string
	success string
}
