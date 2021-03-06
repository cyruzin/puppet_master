package http

import (
	"net/http"
	"time"

	"github.com/cyruzin/puppet_master/domain"
	"github.com/cyruzin/puppet_master/pkg/enc"
	"github.com/cyruzin/puppet_master/pkg/validation"
	"github.com/go-chi/chi/v5"
)

// PermissionHandler  represent the http handler for permission

type PermissionHandler struct {
	PermissionUseCase domain.PermissionUsecase
}

// NewArticleHandler will initialize the articles/ resources endpoint
func NewArticleHandler(c *chi.Mux, p domain.PermissionUsecase) {
	handler := &PermissionHandler{
		PermissionUseCase: p,
	}

	c.Route("/permission", func(r chi.Router) {
		r.Get("/", handler.Fetch)
		r.Get("/{id}", handler.GetByID)
		r.Post("/", handler.Store)
		r.Put("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
	})
}

func (p *PermissionHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	permissions, err := p.PermissionUseCase.Fetch(r.Context())
	if err != nil {
		enc.EncodeError(w, r, err, http.StatusInternalServerError)
		return
	}

	enc.EncodeJSON(w, http.StatusOK, permissions)
}

func (p *PermissionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := enc.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		enc.EncodeError(w, r, domain.ErrIDParam, http.StatusBadRequest)
		return
	}

	permission, err := p.PermissionUseCase.GetByID(r.Context(), id)
	if err != nil {
		enc.EncodeError(w, r, err, http.StatusInternalServerError)
		return
	}

	enc.EncodeJSON(w, http.StatusOK, permission)
}

func (p *PermissionHandler) Store(w http.ResponseWriter, r *http.Request) {
	payload := &domain.Permission{}

	if err := enc.DecodeJSON(r.Body, payload); err != nil {
		enc.EncodeError(w, r, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	if err := validation.IsAValidSchema(ctx, payload); err != nil {
		enc.EncodeError(w, r, err, http.StatusBadRequest)
		return
	}

	payload.CreatedAt = time.Now()
	payload.UpdatedAt = time.Now()

	permissions, err := p.PermissionUseCase.Store(ctx, payload)
	if err != nil {
		enc.EncodeError(w, r, err, http.StatusBadRequest)
		return
	}

	enc.EncodeJSON(w, http.StatusCreated, permissions)
}

func (p *PermissionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := enc.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		enc.EncodeError(w, r, domain.ErrIDParam, http.StatusBadRequest)
		return
	}

	payload := &domain.Permission{}

	if err := enc.DecodeJSON(r.Body, payload); err != nil {
		enc.EncodeError(w, r, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	if err := validation.IsAValidSchema(ctx, payload); err != nil {
		enc.EncodeError(w, r, err, http.StatusBadRequest)
		return
	}

	payload.ID = id
	payload.UpdatedAt = time.Now()

	permissions, err := p.PermissionUseCase.Update(ctx, payload)
	if err != nil {
		enc.EncodeError(w, r, err, http.StatusBadRequest)
		return
	}

	enc.EncodeJSON(w, http.StatusOK, permissions)
}

func (p *PermissionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := enc.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		enc.EncodeError(w, r, domain.ErrIDParam, http.StatusBadRequest)
		return
	}

	if err := p.PermissionUseCase.Delete(r.Context(), id); err != nil {
		enc.EncodeError(w, r, err, http.StatusBadRequest)
		return
	}

	enc.EncodeJSON(w, http.StatusOK, &enc.APIMessage{Message: "permission deleted", Status: http.StatusOK})
}
