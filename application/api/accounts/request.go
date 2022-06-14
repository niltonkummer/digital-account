package accounts

type AccountRequest struct {
	Name   string `json:"name" binding:"required"`
	CPF    string `json:"cpf" binding:"required"`
	Secret string `json:"secret" binding:"required"`
}
