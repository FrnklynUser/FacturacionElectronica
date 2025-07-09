package service

import (
	"FacturacionSunat/internal/domain"
	"FacturacionSunat/internal/platform/signer"
	"FacturacionSunat/internal/platform/sunat"
	"FacturacionSunat/pkg/ubl"
	"context"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// InvoiceService is the service for handling invoice business logic.
type InvoiceService struct {
	invoiceRepo domain.InvoiceRepository
	signer      *signer.XMLSigner
	sunatClient *sunat.Client
}

// NewInvoiceService creates a new InvoiceService.
func NewInvoiceService(repo domain.InvoiceRepository, signer *signer.XMLSigner, sunatClient *sunat.Client) *InvoiceService {
	return &InvoiceService{
	invoiceRepo: repo,
	signer:      signer,
	sunatClient: sunatClient,
	}
}

// Create processes a new invoice.
func (s *InvoiceService) Create(invoice *domain.Invoice) (*domain.Invoice, error) {
	// 1. Basic validation.
	if invoice.Series == "" || invoice.Number == 0 {
		return nil, fmt.Errorf("serie y número son requeridos")
	}

	// 2. Set server-side fields.
	invoice.ID = uuid.New().String()
	invoice.Status = "RECIBIDO"
	invoice.IssueDate = time.Now()
	if err := s.invoiceRepo.Save(context.Background(), invoice); err != nil {
		return nil, fmt.Errorf("error al guardar la factura: %w", err)
	}

	// 4. Build the UBL structure.
	ublInvoice, err := ubl.BuildInvoice(invoice)
	if err != nil {
		return nil, fmt.Errorf("error al construir el UBL: %w", err)
	}

	// 5. Generate the unsigned XML.
	unsignedXML, err := xml.MarshalIndent(ublInvoice, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error al generar el XML: %w", err)
	}

	// 6. Sign the XML.
	signedXML, err := s.signer.Sign(unsignedXML)
	if err != nil {
		return nil, fmt.Errorf("error al firmar el XML: %w", err)
	}
	invoice.Status = "FIRMADO"

	// 7. Send the bill to SUNAT.
	fileName := fmt.Sprintf("%s-%s-%s-%d.xml", invoice.Issuer.RUC, invoice.Type, invoice.Series, invoice.Number)
	ticket, err := s.sunatClient.SendBill(fileName, signedXML)
	if err != nil {
		// In a real app, you would handle specific SUNAT errors here.
		invoice.Status = "RECHAZADO"
		return nil, fmt.Errorf("error al enviar a SUNAT: %w", err)
	}

	// 8. Store the ticket ID for status tracking.
	invoice.TicketID = ticket // Assuming domain.Invoice has a TicketID field
	fmt.Printf("Ticket recibido de SUNAT: %s\n", ticket)

	invoice.Status = "ENVIADO"
	// You would update the status in the database here.
	// s.invoiceRepo.UpdateStatus(context.Background(), invoice.ID, invoice.Status)

	return invoice, nil
}

// CreateCreditNote processes a new credit note.
func (s *InvoiceService) CreateCreditNote(cn *domain.CreditNote) (*domain.CreditNote, error) {
	cn.ID = uuid.New().String()
	cn.Status = "RECIBIDO"
	cn.IssueDate = time.Now()

	ublCreditNote, err := ubl.BuildCreditNote(cn)
	if err != nil {
		return nil, fmt.Errorf("error al construir UBL de nota de crédito: %w", err)
	}

	unsignedXML, err := xml.MarshalIndent(ublCreditNote, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error al generar XML de nota de crédito: %w", err)
	}

	signedXML, err := s.signer.Sign(unsignedXML)
	if err != nil {
		return nil, fmt.Errorf("error al firmar XML de nota de crédito: %w", err)
	}
	cn.Status = "FIRMADO"

	fileName := fmt.Sprintf("%s-%s-%s-%d.xml", cn.Issuer.RUC, cn.Type, cn.Series, cn.Number)
	ticket, err := s.sunatClient.SendBill(fileName, signedXML)
	if err != nil {
		cn.Status = "RECHAZADO"
		return nil, fmt.Errorf("error al enviar nota de crédito a SUNAT: %w", err)
	}

	cn.TicketID = ticket // Assuming domain.CreditNote has a TicketID field
	fmt.Printf("Ticket de Nota de Crédito recibido de SUNAT: %s\n", ticket)

	cn.Status = "ENVIADO"
	return cn, nil
}

// CreateDebitNote processes a new debit note.
func (s *InvoiceService) CreateDebitNote(dn *domain.DebitNote) (*domain.DebitNote, error) {
	dn.ID = uuid.New().String()
	dn.Status = "RECIBIDO"
	dn.IssueDate = time.Now()

	ublDebitNote, err := ubl.BuildDebitNote(dn)
	if err != nil {
		return nil, fmt.Errorf("error al construir UBL de nota de débito: %w", err)
	}

	unsignedXML, err := xml.MarshalIndent(ublDebitNote, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error al generar XML de nota de débito: %w", err)
	}

	signedXML, err := s.signer.Sign(unsignedXML)
	if err != nil {
		return nil, fmt.Errorf("error al firmar XML de nota de débito: %w", err)
	}
	dn.Status = "FIRMADO"

	fileName := fmt.Sprintf("%s-%s-%s-%d.xml", dn.Issuer.RUC, dn.Type, dn.Series, dn.Number)
	ticket, err := s.sunatClient.SendBill(fileName, signedXML)
	if err != nil {
		dn.Status = "RECHAZADO"
		return nil, fmt.Errorf("error al enviar nota de débito a SUNAT: %w", err)
	}

	dn.TicketID = ticket // Assuming domain.DebitNote has a TicketID field
	fmt.Printf("Ticket de Nota de Débito recibido de SUNAT: %s\n", ticket)

	dn.Status = "ENVIADO"
	return dn, nil
}

// GetDocumentStatus retrieves the status of a document from SUNAT.
func (s *InvoiceService) GetDocumentStatus(id string) (string, error) {
	statusResp, err := s.sunatClient.GetStatus(id)
	if err != nil {
		return "", fmt.Errorf("error al consultar estado en SUNAT: %w", err)
	}
	// Here you would parse statusResp.Content (Base64 encoded CDR zip) if needed.
	return statusResp.StatusCode, nil
}

// GetDocumentStatusCdr retrieves the status and CDR of a document using its full details from SUNAT.
func (s *InvoiceService) GetDocumentStatusCdr(ruc, docType, series, number string) (*sunat.StatusCdr, error) {
	statusCdrResp, err := s.sunatClient.GetStatusCdr(ruc, docType, series, number)
	if err != nil {
		return nil, fmt.Errorf("error al consultar CDR en SUNAT: %w", err)
	}
	// Here you would parse statusCdrResp.Content (Base64 encoded CDR zip) if needed.
	return statusCdrResp, nil
}
