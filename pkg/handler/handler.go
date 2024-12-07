package handler

import (
	"net/http"

	"github.com/sathishs-dev/pismo-transactions/pkg/repository"
)

type handler struct {
	repo repository.PismoRepo
}

type Handler interface {
	CreateAccount() http.HandlerFunc
	GetAccount() http.HandlerFunc
	CreateTransaction() http.HandlerFunc
}

func NewHandler(repo repository.PismoRepo) Handler {
	return &handler{
		repo,
	}
}

func (h *handler) CreateAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (h *handler) GetAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (h *handler) CreateTransaction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
