package models

type UserRequest struct {
	Message string `json:"message"`
}

type Prompt struct {
	Content string
	Message string
}

type Answer struct {
	AIAnswer string `json:"answer"`
}
