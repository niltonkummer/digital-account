package handler

import (
	"digital-account/application/repository"

	"gorm.io/gorm"
)

type Handler interface {
	Repository() repository.Repository
}

type handler struct {
	repository repository.Repository
}

func (h handler) Repository() repository.Repository {
	return h.repository
}

func New(db *gorm.DB) Handler {
	return &handler{repository: repository.Config(db)}
}
