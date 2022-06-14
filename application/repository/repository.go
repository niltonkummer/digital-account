package repository

type Repository interface {
	Account() Account
	User() User
	Transfer() Transfer
}

type repo struct {
	user     User
	account  Account
	transfer Transfer
}

func (r repo) Account() Account {
	return r.account
}

func (r repo) User() User {
	return r.user
}

func (r repo) Transfer() Transfer {
	return r.transfer
}
