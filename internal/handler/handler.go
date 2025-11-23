package handler

import (
	"github.com/go-chi/chi/v5"
	"pr-reviewer-2025/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services,
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Route("/team", func(r chi.Router) {
		r.Post("/add", h.createTeam)
		r.Get("/get", h.getTeam)
	})

	router.Route("/users", func(r chi.Router) {
		r.Post("/setIsActive", h.setIsActive)
		r.Get("/getReview", h.getUserReview)
	})

	router.Route("/pullRequest", func(r chi.Router) {
		r.Post("/create", h.createPR)
		r.Post("/merge", h.mergePR)
		r.Post("/reassign", h.reassign)
	})

	return router
}
