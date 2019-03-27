package payment

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-http-utils/headers"
	"github.com/tenzi-muraru/paymentApp/payment/model"
)

// Handler - handles HTTP requests for payments
type Handler struct {
	Repository
}

// NewHandler - creates a new Handler with a DB repository
func NewHandler(repository Repository) *Handler {
	return &Handler{
		repository,
	}
}

// GetPaymentByID - retrieves the payment having the provided ID
func (h Handler) GetPaymentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	paymentID := vars["paymentID"]

	payment, err := h.GetByID(paymentID)
	if err != nil {
		respondError("get payment by ID", w, err)
		return
	}

	respondJSON(w, http.StatusOK, payment)
}

// AddPayment - creates a new payment
func (h Handler) AddPayment(w http.ResponseWriter, r *http.Request) {
	var payment model.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		log.Printf("Error encountered while processing add payment: %v\n", err)
		respondJSON(w, http.StatusBadRequest, model.NewBadRequestError(err.Error()))
		return
	}

	newPayment, err := h.Add(payment)
	if err != nil {
		respondError("add payment", w, err)
		return
	}

	respondJSON(w, http.StatusCreated, newPayment)
}

// GetAllPayments - retrieves all payments
func (h Handler) GetAllPayments(w http.ResponseWriter, r *http.Request) {
	payments, err := h.GetAll()
	if err != nil {
		respondError("get all payments", w, err)
		return
	}
	respondJSON(w, http.StatusOK, payments)
}

// DeletePayment - deletes the payment with the provided ID
func (h Handler) DeletePayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	paymentID := vars["paymentID"]

	err := h.Delete(paymentID)
	if err != nil {
		respondError("delete payment", w, err)
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}

// UpdatePayment - updates the payment with the provided ID
func (h Handler) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	var payment model.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		log.Printf("Error encountered while processing update payment: %v\n", err)
		respondJSON(w, http.StatusBadRequest, model.NewBadRequestError(err.Error()))
		return
	}

	vars := mux.Vars(r)
	payment.ID = vars["paymentID"]
	if err := h.Update(payment); err != nil {
		respondError("update payment", w, err)
		return
	}

	respondJSON(w, http.StatusOK, nil)
}

// respondError - creates the error response with JSON payload
func respondError(callName string, w http.ResponseWriter, err error) {
	log.Printf("Error encountered while processing %s: %v\n", callName, err)

	if err == ErrPaymentNotFound {
		respondJSON(w, http.StatusNotFound, model.PaymentNotFound)
		return
	}
	respondJSON(w, http.StatusInternalServerError, model.GenericError)
}

// respondJSON - creates the response with the provided JSON payload
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		respond(w, http.StatusInternalServerError, []byte(err.Error()))
		return
	}

	respond(w, status, response)
}

func respond(w http.ResponseWriter, status int, payload []byte) {
	w.Header().Set(headers.ContentType, "application/json")
	w.WriteHeader(status)
	w.Write(payload)
}
