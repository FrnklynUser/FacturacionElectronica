package main

import (
	"FacturacionSunat/internal/handler"
	"FacturacionSunat/internal/platform/signer"
	"FacturacionSunat/internal/platform/storage"
	"FacturacionSunat/internal/platform/sunat"
	"FacturacionSunat/internal/service"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno desde .env (si existe)
	_ = godotenv.Load()

	// =============================================================================
	// CONFIGURACIÓN - Lee desde variables de entorno
	// =============================================================================
	certPath := "./internal/certificados/certificate.p12" // Ruta actualizada
	certPass := os.Getenv("CERT_PASS")
	sunatUsername := os.Getenv("SUNAT_USER")
	sunatPassword := os.Getenv("SUNAT_PASS")

	if certPass == "" || sunatUsername == "" || sunatPassword == "" {
		log.Fatal("Las variables de entorno CERT_PASS, SUNAT_USER y SUNAT_PASS son requeridas. Puedes definirlas en un archivo .env")
	}

	// URL del servicio de SUNAT (beta por defecto)
	sunatBaseURL := "https://e-beta.sunat.gob.pe/ol-ti-itcpfegem-beta/billService?wsdl" // WSDL for sendBill and getStatus

	// 1. Initialize dependencies (the "platform" layer).
	invoiceRepo := storage.NewInvoiceMemoryRepo() // Usando el repositorio en memoria
	signer, err := signer.NewXMLSigner(certPath, certPass)
	if err != nil {
		log.Fatalf("Error al inicializar el firmador digital: %v", err)
	}
	sunatClient, err := sunat.NewClient(sunatBaseURL, sunatUsername, sunatPassword)
	if err != nil {
		log.Fatalf("Error al inicializar el cliente de SUNAT: %v", err)
	}

	// 2. Initialize the core logic (the "service" layer).
	invoiceService := service.NewInvoiceService(invoiceRepo, signer, sunatClient)

	// 3. Initialize the entrypoint (the "handler" layer).
	invoiceHandler := handler.NewInvoiceHandler(invoiceService)

	// 4. Register API routes.
	apiV1 := http.NewServeMux()
	apiV1.HandleFunc("/api/v1/invoices", invoiceHandler.CreateInvoice)
	apiV1.HandleFunc("/api/v1/credit-notes", invoiceHandler.CreateCreditNote)
	apiV1.HandleFunc("/api/v1/debit-notes", invoiceHandler.CreateDebitNote)
	apiV1.HandleFunc("/api/v1/documents/", invoiceHandler.GetDocumentStatus)       // Handles /api/v1/documents/{id}/status
	apiV1.HandleFunc("/api/v1/documents/cdr", invoiceHandler.GetDocumentStatusCdr) // Handles /api/v1/documents/cdr?ruc=...&docType=...&series=...&number=...
	apiV1.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	// 5. Start the server.
	server := &http.Server{
		Addr:    ":8080",
		Handler: apiV1,
	}

	fmt.Println("Servidor escuchando en http://localhost:8080")
	fmt.Println("Asegúrate de haber creado un archivo .env con CERT_PASS, SUNAT_USER y SUNAT_PASS, o de haberlas exportado.")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error al iniciar el servidor: %s\n", err)
	}
}
