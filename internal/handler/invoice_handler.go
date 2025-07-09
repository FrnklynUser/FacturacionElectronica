package handler

import (
	"FacturacionSunat/internal/domain"
	"FacturacionSunat/internal/platform/sunat"
	"encoding/json"
	"fmt"
	"net/http"
)

// IInvoiceService defines the interface for invoice services.
// We define it here to avoid circular dependencies.
type IInvoiceService interface {
	Create(invoice *domain.Invoice) (*domain.Invoice, error)
	CreateCreditNote(cn *domain.CreditNote) (*domain.CreditNote, error)
	CreateDebitNote(dn *domain.DebitNote) (*domain.DebitNote, error)
	GetDocumentStatus(id string) (string, error)
	GetDocumentStatusCdr(ruc, docType, series, number string) (*sunat.StatusCdr, error)
}

// InvoiceHandler handles the HTTP requests for invoices.
type InvoiceHandler struct {
	service IInvoiceService
}

// NewInvoiceHandler creates a new InvoiceHandler.
func NewInvoiceHandler(s IInvoiceService) *InvoiceHandler {
	return &InvoiceHandler{service: s}
}

// CreateInvoice handles the creation of a new invoice.
func (h *InvoiceHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	var invoice domain.Invoice

	// 1. Decode the incoming JSON request.
	if err := json.NewDecoder(r.Body).Decode(&invoice); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 2. Call the service layer to process the invoice.
	createdInvoice, err := h.service.Create(&invoice)
	if err != nil {
		// In a real app, you'd check the error type and return the appropriate status code.
		http.Error(w, "Failed to create invoice", http.StatusInternalServerError)
		return
	}

	// 3. Encode the response and send it back.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdInvoice); err != nil {
		// This is less likely to happen, but good to handle.
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

// CreateCreditNote handles the creation of a new credit note.
func (h *InvoiceHandler) CreateCreditNote(w http.ResponseWriter, r *http.Request) {
	var cn domain.CreditNote
	if err := json.NewDecoder(r.Body).Decode(&cn); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdCn, err := h.service.CreateCreditNote(&cn)
	if err != nil {
		http.Error(w, "Failed to create credit note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdCn)
}

// CreateDebitNote handles the creation of a new debit note.
func (h *InvoiceHandler) CreateDebitNote(w http.ResponseWriter, r *http.Request) {
	var dn domain.DebitNote
	if err := json.NewDecoder(r.Body).Decode(&dn); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdDn, err := h.service.CreateDebitNote(&dn)
	if err != nil {
		http.Error(w, "Failed to create debit note", http.StatusInternalServerError)
		return	
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdDn)
}

// GetDocumentStatus handles the request to get the status of a document.
func (h *InvoiceHandler) GetDocumentStatus(w http.ResponseWriter, r *http.Request) {
	// Extract the document ID from the URL path.
	id := r.URL.Path[len("/api/v1/documents/"):len(r.URL.Path)-len("/status")]
	if id == "" {
		http.Error(w, "Document ID is required", http.StatusBadRequest)
		return
	}

	status, err := h.service.GetDocumentStatus(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get document status: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}

// GetDocumentStatusCdr handles the request to get the status and CDR of a document.
func (h *InvoiceHandler) GetDocumentStatusCdr(w http.ResponseWriter, r *http.Request) {
	// Extract parameters from query string.
	ruc := r.URL.Query().Get("ruc")
	docType := r.URL.Query().Get("docType")
	series := r.URL.Query().Get("series")
	number := r.URL.Query().Get("number")

	if ruc == "" || docType == "" || series == "" || number == "" {
		http.Error(w, "ruc, docType, series, and number are required query parameters", http.StatusBadRequest)
		return
	}

	statusCdr, err := h.service.GetDocumentStatusCdr(ruc, docType, series, number)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get document CDR status: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statusCdr)
}