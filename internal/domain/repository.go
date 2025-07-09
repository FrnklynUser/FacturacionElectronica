package domain

import "context"

// InvoiceRepository defines the persistence interface for Invoices.
type InvoiceRepository interface {
	// Save saves a given invoice to the repository.
	Save(ctx context.Context, invoice *Invoice) error

	// FindByID retrieves an invoice by its ID.
	FindByID(ctx context.Context, id string) (*Invoice, error)

	// UpdateStatus updates the status of a given invoice.
	UpdateStatus(ctx context.Context, id string, status string) error
}
