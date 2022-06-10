package models

type User struct {
	Model
	Name     string    `json:"name"`
	CPF      string    `json:"cpf"`
	Secret   string    `json:"-"`
	Accounts []Account `json:"accounts,omitempty"`
}
