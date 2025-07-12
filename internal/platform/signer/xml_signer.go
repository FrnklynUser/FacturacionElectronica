package signer

import (
	"crypto"
	"crypto/rsa"
	"crypto/tls"
	"fmt"
	"io/ioutil"

	"github.com/beevik/etree"
	dsig "github.com/russellhaering/goxmldsig"
	"golang.org/x/crypto/pkcs12"
	"encoding/pem"
)

// XMLSigner handles the digital signature of XML documents.
type XMLSigner struct {
	signer *dsig.SigningContext
}

// NewXMLSigner creates a new XMLSigner.
// It requires the path to a PKCS#12 certificate file and its password.
func NewXMLSigner(p12FilePath, password string) (*XMLSigner, error) {
	if p12FilePath == "" {
		return nil, fmt.Errorf("ruta del certificado no puede estar vacía")
	}

	p12, err := ioutil.ReadFile(p12FilePath)
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer el archivo del certificado: %w", err)
	}

	pemBlocks, err := pkcs12.ToPEM(p12, password)
	if err != nil {
		return nil, fmt.Errorf("no se pudo decodificar el archivo P12 a PEM: %w", err)
	}

	var certPEM, keyPEM []byte
	for _, b := range pemBlocks {
		if b.Type == "CERTIFICATE" {
			certPEM = pem.EncodeToMemory(b)
		}
		if b.Type == "PRIVATE KEY" {
			keyPEM = pem.EncodeToMemory(b)
		}
	}

	if keyPEM == nil || certPEM == nil {
		return nil, fmt.Errorf("no se encontró la llave privada o el certificado en el archivo P12")
	}

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, fmt.Errorf("error al cargar el par de llaves X509: %w", err)
	}
	
	rsaPrivateKey, ok := cert.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("la llave privada no es de tipo RSA")
	}

	ks := &customKeyStoreImpl{
		privateKey: rsaPrivateKey,
		cert: cert.Certificate[0],
	}

	signer := dsig.NewDefaultSigningContext(ks)
	signer.Hash = crypto.SHA256

	return &XMLSigner{signer: signer}, nil
}

type customKeyStoreImpl struct {
	privateKey *rsa.PrivateKey
	cert       []byte
}

func (ks *customKeyStoreImpl) GetKeyPair() (*rsa.PrivateKey, []byte, error) {
	return ks.privateKey, ks.cert, nil
}


// Sign applies a digital signature to an XML document according to UBL standards.
func (s *XMLSigner) Sign(xmlContent []byte) ([]byte, error) {
	fmt.Println("FIRMANDO XML con la librería goxmldsig...")

	// 1. Parse the XML content into an etree document.
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlContent); err != nil {
		return nil, fmt.Errorf("error al parsear el XML para firmar: %w", err)
	}

	root := doc.Root()
	if root == nil {
		return nil, fmt.Errorf("documento XML no tiene elemento raíz")
	}

	// 2. Find the placeholder for the signature and its parent.
	extensionContent := root.FindElement("./ext:UBLExtensions/ext:UBLExtension/ext:ExtensionContent")
	if extensionContent == nil {
		return nil, fmt.Errorf("no se encontró el elemento ext:UBLExtensions/ext:UBLExtension/ext:ExtensionContent")
	}
	
	// Remove any existing signature placeholder. The UBL builder creates one.
	if placeholder := extensionContent.SelectElement("ds:Signature"); placeholder != nil {
		extensionContent.RemoveChild(placeholder)
	}

	// 3. Sign the root element. This creates an enveloped signature.
	signed, err := s.signer.SignEnveloped(root)
	if err != nil {
		return nil, fmt.Errorf("error al firmar el documento XML: %w", err)
	}

	// 4. Extract the generated signature element.
	signature := signed.SelectElement("ds:Signature")
	if signature == nil {
		return nil, fmt.Errorf("la firma generada no se encontró en el documento firmado")
	}

	// 5. Add the real signature to the correct place in the original document.
	extensionContent.AddChild(signature.Copy())

	// 6. Serialize the original document, now with the signature correctly placed.
	signedXML, err := doc.WriteToBytes()
	if err != nil {
		return nil, fmt.Errorf("error al serializar el XML firmado: %w", err)
	}

	return signedXML, nil
}
