package domain

import "time"

// Issuer represents the company issuing the invoice.
type Issuer struct {
	RUC      string `json:"ruc"`
	Name     string `json:"razon_social"`
	Address  string `json:"direccion"`
}

// Recipient represents the customer receiving the invoice.
type Recipient struct {
	DocType string `json:"tipo_doc"` // DNI, RUC, CE
	DocNum  string `json:"num_doc"`
	Name    string `json:"nombre"`
}

// InvoiceLine represents a single item line in the invoice.
type InvoiceLine struct {
	ID          string  `json:"id"`
	Code        string  `json:"codigo"`
	Description string  `json:"descripcion"`
	Quantity    float64 `json:"cantidad"`
	UnitPrice   float64 `json:"valor_unitario"`
	TotalValue  float64 `json:"valor_total"`
	IGV         float64 `json:"igv"`
}

// Totals represents the monetary totals for the invoice.
type Totals struct {
	Gross      float64 `json:"gravado"`
	IGV        float64 `json:"igv"`
	Total      float64 `json:"total"`
}

// Invoice represents the main electronic invoice document.
type Invoice struct {
	ID          string        `json:"id"`
	Type        string        `json:"tipo_comprobante"` // 01: Factura, 03: Boleta
	Series      string        `json:"serie"`
	Number      int           `json:"numero"`
	IssueDate   time.Time     `json:"fecha_emision"`
	Currency    string        `json:"moneda"` // PEN, USD
	Issuer      Issuer        `json:"emisor"`
	Recipient   Recipient     `json:"receptor"`
	Lines       []InvoiceLine `json:"items"`
	Totals      Totals        `json:"totales"`
	Status      string        `json:"estado"` // (aceptado, rechazado, etc.)
	TicketID    string        `json:"ticket_id,omitempty"` // SUNAT ticket ID for tracking
}