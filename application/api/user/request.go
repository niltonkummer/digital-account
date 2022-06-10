package user

type LoginRequest struct {
	CPF    string `json:"cpf" binding:"required"`
	Secret string `json:"secret" binding:"required"`
}
