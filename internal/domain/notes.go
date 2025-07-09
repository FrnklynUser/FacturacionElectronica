package domain

import "time"

// DiscrepancyResponse describes the reason for a credit/debit note.
type DiscrepancyResponse struct {
	ReferenceID string `json:"nro_comprobante_afectado"` // e.g., F001-123
	TypeCode    string `json:"codigo_motivo"`            // Catalog 09 for Credit, 10 for Debit
	Description string `json:"descripcion_motivo"`
}

// CreditNote represents the main electronic credit note document.
type CreditNote struct {
	ID                  string               `json:"id"`
	Type                string               `json:"tipo_comprobante"` // 07: Credit Note
	Series              string               `json:"serie"`
	Number              int                  `json:"numero"`
	IssueDate           time.Time            `json:"fecha_emision"`
	Currency            string               `json:"moneda"` // PEN, USD
	Issuer              Issuer               `json:"emisor"`
	Recipient           Recipient            `json:"receptor"`
	DiscrepancyResponse DiscrepancyResponse  `json:"motivo_o_sustento"`
	Lines               []InvoiceLine        `json:"items"`
	Totals              Totals               `json:"totales"`
	Status              string               `json:"estado"` // (aceptado, rechazado, etc.)
	TicketID            string               `json:"ticket_id,omitempty"` // SUNAT ticket ID for tracking
}

// DebitNote represents the main electronic debit note document.
// Structure is very similar to CreditNote, just different TypeCode.
type DebitNote struct {
	ID                  string               `json:"id"`
	Type                string               `json:"tipo_comprobante"` // 08: Debit Note
	Series              string               `json:"serie"`
	Number              int                  `json:"numero"`
	IssueDate           time.Time            `json:"fecha_emision"`
	Currency            string               `json:"moneda"` // PEN, USD
	Issuer              Issuer               `json:"emisor"`
	Recipient           Recipient            `json:"receptor"`
	DiscrepancyResponse DiscrepancyResponse  `json:"motivo_o_sustento"`
	Lines               []InvoiceLine        `json:"items"`
	Totals              Totals               `json:"totales"`
	Status              string               `json:"estado"` // (aceptado, rechazado, etc.)
	TicketID            string               `json:"ticket_id,omitempty"` // SUNAT ticket ID for tracking
}