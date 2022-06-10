package handler

import "digital-account/application/repository"

type Handler interface {
	Repository() repository.Repository
}

type handler struct {
	repository repository.Repository
}
