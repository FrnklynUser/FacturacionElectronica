package signer

import (
	"fmt"
	"os"

	"github.com/artemkunich/goxades"
	"github.com/beevik/etree"
)

// XMLSigner handles the digital signature of XML documents.
type XMLSigner struct {
	signer *goxades.Signer
}

// NewXMLSigner creates a new XMLSigner.
// It requires the path to a PKCS#12 certificate file and its password.
func NewXMLSigner(p12FilePath, password string) (*XMLSigner, error) {
	// In a real application, these values should come from a secure configuration.
	if p12FilePath == "" {
		return nil, fmt.Errorf("ruta del certificado no puede estar vacía")
	}

	cert, err := os.ReadFile(p12FilePath)
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer el archivo del certificado: %w", err)
	}

	signer, err := goxades.NewSigner(cert, password, goxades.SignaturePolicy{
		Identifier: goxades.Identifier{Value: "https://www.sunat.gob.pe/policies/factura-electronica"},
	})
	if err != nil {
		return nil, fmt.Errorf("error al crear el firmador: %w", err)
	}

	return &XMLSigner{signer: signer}, nil
}

// Sign applies a digital signature to an XML document.
func (s *XMLSigner) Sign(xmlContent []byte) ([]byte, error) {
	fmt.Println("FIRMANDO XML con la librería goxades...")

	// 1. Parse the XML content into an etree document.
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlContent); err != nil {
		return nil, fmt.Errorf("error al parsear el XML para firmar: %w", err)
	}

	// 2. Sign the document.
	if err := s.signer.Sign(doc); err != nil {
		return nil, fmt.Errorf("error al firmar el documento XML: %w", err)
	}

	// 3. Serialize the signed document back to bytes.
	signedXML, err := doc.WriteToBytes()
	if err != nil {
		return nil, fmt.Errorf("error al serializar el XML firmado: %w", err)
	}

	return signedXML, nil
}
