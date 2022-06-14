package models

type User struct {
	Model
	Name    string   `json:"name"`
	CPF     string   `gorm:"unique:cpf_index" json:"cpf"`
	Secret  string   `json:"-"`
	Account *Account `json:"account,omitempty"`
}
