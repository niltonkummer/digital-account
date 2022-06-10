package accounts

type AccountRequest struct {
	Secret string `json:"secret" binding:"required"`
}
