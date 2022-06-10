package repository

type Repository interface {
	Account() Account
	Login() Login
}

type repo struct {
	login   Login
	account Account
}

func (r repo) Account() Account {
	return r.account
}

func (r repo) Login() Login {
	return r.login
}
