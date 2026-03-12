package payments

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/retrosnack-clothing/retrosnack/pkg/httputil"
	"github.com/retrosnack-clothing/retrosnack/pkg/middleware"
)

type Handler struct {
	svc           Service
	applicationID string
	locationID    string
	environment   string
}

func NewHandler(svc Service, applicationID, locationID, environment string) *Handler {
	return &Handler{svc: svc, applicationID: applicationID, locationID: locationID, environment: environment}
}

func (h *Handler) Register(r chi.Router) {
	r.Get("/payments/config", h.paymentConfig)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RateLimit(10, 1*time.Minute))
		r.Post("/checkout", h.createCheckout)
		r.Post("/payments/process", h.processPayment)
	})
	r.Post("/webhooks/square", h.squareWebhook)
}

func (h *Handler) paymentConfig(w http.ResponseWriter, r *http.Request) {
	httputil.JSON(w, http.StatusOK, map[string]string{
		"application_id": h.applicationID,
		"location_id":    h.locationID,
		"environment":    h.environment,
	})
}

func (h *Handler) createCheckout(w http.ResponseWriter, r *http.Request) {
	var req CreateCheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.ErrorMsg(w, http.StatusBadRequest, "invalid request body")
		return
	}

	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = "https://retrosnack.shop"
	}

	redirectURL := origin + "/orders/" + req.OrderID.String() + "/confirmation"

	sess, err := h.svc.CreateCheckout(r.Context(), req, redirectURL)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err)
		return
	}
	httputil.JSON(w, http.StatusOK, sess)
}

func (h *Handler) processPayment(w http.ResponseWriter, r *http.Request) {
	var req ProcessPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.ErrorMsg(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.SourceID == "" {
		httputil.ErrorMsg(w, http.StatusBadRequest, "source_id is required")
		return
	}
	if req.OrderID == (uuid.UUID{}) {
		httputil.ErrorMsg(w, http.StatusBadRequest, "order_id is required")
		return
	}

	result, err := h.svc.ProcessPayment(r.Context(), req)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err)
		return
	}
	httputil.JSON(w, http.StatusOK, result)
}

func (h *Handler) squareWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(io.LimitReader(r.Body, 65536))
	if err != nil {
		httputil.ErrorMsg(w, http.StatusBadRequest, "failed to read body")
		return
	}

	signature := r.Header.Get("x-square-hmacsha256-signature")
	if err := h.svc.HandleWebhook(r.Context(), payload, signature); err != nil {
		// log the error but return 200 so square doesn't retry permanent failures
		slog.Error("webhook processing failed", "error", err)
	}

	httputil.NoContent(w)
}
