package storage

import (
	"FacturacionSunat/internal/domain"
	"context"
	"fmt"
	"sync"
)

// InvoiceMemoryRepo is an in-memory implementation of the InvoiceRepository.
type InvoiceMemoryRepo struct {
	mu       sync.RWMutex
	invoices map[string]*domain.Invoice
}

// NewInvoiceMemoryRepo creates a new InvoiceMemoryRepo.
func NewInvoiceMemoryRepo() *InvoiceMemoryRepo {
	return &InvoiceMemoryRepo{
		invoices: make(map[string]*domain.Invoice),
	}
}

// Save implements the domain.InvoiceRepository interface.
func (r *InvoiceMemoryRepo) Save(ctx context.Context, invoice *domain.Invoice) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.invoices[invoice.ID]; ok {
		return fmt.Errorf("factura con ID %s ya existe", invoice.ID)
	}
	r.invoices[invoice.ID] = invoice
	fmt.Printf("GUARDANDO factura %s en memoria...\n", invoice.ID)
	return nil
}

// FindByID implements the domain.InvoiceRepository interface.
func (r *InvoiceMemoryRepo) FindByID(ctx context.Context, id string) (*domain.Invoice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	invoice, ok := r.invoices[id]
	if !ok {
		return nil, fmt.Errorf("factura con ID %s no encontrada", id)
	}
	fmt.Printf("BUSCANDO factura %s en memoria...\n", id)
	return invoice, nil
}

// UpdateStatus implements the domain.InvoiceRepository interface.
func (r *InvoiceMemoryRepo) UpdateStatus(ctx context.Context, id string, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	invoice, ok := r.invoices[id]
	if !ok {
		return fmt.Errorf("factura con ID %s no encontrada para actualizar estado", id)
	}
	invoice.Status = status
	fmt.Printf("ACTUALIZANDO estado de factura %s a %s en memoria...\n", id, status)
	return nil
}
