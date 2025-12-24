package http

import (
	"errors"
	"net/http"
	"strconv"

	"portal-data-backend/internal/analytics/usecase"
	"portal-data-backend/infrastructure/http/response"
	pkgErrors "portal-data-backend/pkg/errors"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	analyticsUsecase usecase.Usecase
}

func NewHandler(analyticsUsecase usecase.Usecase) *Handler {
	return &Handler{
		analyticsUsecase: analyticsUsecase,
	}
}

func (h *Handler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	dashboard, err := h.analyticsUsecase.GetDashboard(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Dashboard data retrieved successfully", dashboard)
}

func (h *Handler) GetDatasetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.analyticsUsecase.GetDatasetStats(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Dataset stats retrieved successfully", stats)
}

func (h *Handler) GetOrganizationStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.analyticsUsecase.GetOrganizationStats(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Organization stats retrieved successfully", stats)
}

func (h *Handler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.analyticsUsecase.GetUserStats(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "User stats retrieved successfully", stats)
}

func (h *Handler) GetPopularDatasets(w http.ResponseWriter, r *http.Request) {
	limit := parseIntQuery(r, "limit", 10)

	datasets, err := h.analyticsUsecase.GetPopularDatasets(r.Context(), limit)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Popular datasets retrieved successfully", datasets)
}

func (h *Handler) GetPopularTags(w http.ResponseWriter, r *http.Request) {
	limit := parseIntQuery(r, "limit", 10)

	tags, err := h.analyticsUsecase.GetPopularTags(r.Context(), limit)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Popular tags retrieved successfully", tags)
}

func (h *Handler) GetDatasetTrend(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	limit := parseIntQuery(r, "limit", 30)

	trend, err := h.analyticsUsecase.GetDatasetTrend(r.Context(), period, limit)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.OK(w, response.CodeSuccess, "Dataset trend retrieved successfully", trend)
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, pkgErrors.ErrNotFound):
		response.NotFound(w, response.CodeNotFound, "Resource not found", nil)
	default:
		response.InternalError(w, response.CodeInternalServerError, "Internal server error", nil)
	}
}

func parseIntQuery(r *http.Request, key string, defaultValue int) int {
	if value := r.URL.Query().Get(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Route("/analytics", func(r chi.Router) {
		r.Get("/dashboard", handler.GetDashboard)
		r.Get("/stats/datasets", handler.GetDatasetStats)
		r.Get("/stats/organizations", handler.GetOrganizationStats)
		r.Get("/stats/users", handler.GetUserStats)
		r.Get("/popular/datasets", handler.GetPopularDatasets)
		r.Get("/popular/tags", handler.GetPopularTags)
		r.Get("/trend/datasets", handler.GetDatasetTrend)
	})
}
