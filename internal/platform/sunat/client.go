package sunat

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"path/filepath"

	"github.com/tiaguinho/gosoap"
)

// Client is used to communicate with the SUNAT API via SOAP.
type Client struct {
	soapClient *gosoap.Client
	username   string
	password   string
}

// NewClient creates a new SUNAT API client.
// wsdlURL is the URL to the WSDL file (e.g., "https://e-beta.sunat.gob.pe/ol-ti-itcpfegem-beta/billService?wsdl").
func NewClient(wsdlURL, username, password string) *Client {
	soapClient := gosoap.NewClient(wsdlURL)
	return &Client{
		soapClient: soapClient,
		username:   username,
		password:   password,
	}
}

// SendBillRequest represents the SOAP request for sendBill operation.
type SendBillRequest struct {
	XMLName     xml.Name `xml:"ser:sendBill"`
	FileName    string   `xml:"fileName"`
	ContentFile string   `xml:"contentFile"` // Base64 encoded ZIP content
}

// SendBillResponse represents the SOAP response for sendBill operation.
type SendBillResponse struct {
	XMLName xml.Name `xml:"sendBillResponse"`
	Ticket  string   `xml:"ticket"`
}

// GetStatusRequest represents the SOAP request for getStatus operation.
type GetStatusRequest struct {
	XMLName xml.Name `xml:"ser:getStatus"`
	Ticket  string   `xml:"ticket"`
}

// GetStatusResponse represents the SOAP response for getStatus operation.
type GetStatusResponse struct {
	XMLName xml.Name `xml:"getStatusResponse"`
	Status  *Status  `xml:"status"`
}

// Status represents the status object returned by getStatus.
type Status struct {
	StatusCode    string `xml:"statusCode"`
	Content       string `xml:"content"` // Base64 encoded CDR zip
	StatusMessage string `xml:"statusMessage"`
}

// GetStatusCdrRequest represents the SOAP request for getStatusCdr operation.
type GetStatusCdrRequest struct {
	XMLName           xml.Name `xml:"ser:getStatusCdr"`
	RucComprobante    string   `xml:"rucComprobante"`
	TipoComprobante   string   `xml:"tipoComprobante"`
	SerieComprobante  string   `xml:"serieComprobante"`
	NumeroComprobante string   `xml:"numeroComprobante"`
}

// GetStatusCdrResponse represents the SOAP response for getStatusCdr operation.
type GetStatusCdrResponse struct {
	XMLName   xml.Name   `xml:"getStatusCdrResponse"`
	StatusCdr *StatusCdr `xml:"statusCdr"`
}

// StatusCdr represents the status object returned by getStatusCdr.
type StatusCdr struct {
	StatusCode    string `xml:"statusCode"`
	Content       string `xml:"content"` // Base64 encoded CDR zip
	StatusMessage string `xml:"statusMessage"`
}

// SendBill sends a signed XML, zipped, to the SUNAT bill service via SOAP.
// It returns the ticket ID from SUNAT.
func (c *Client) SendBill(fileName string, signedXML []byte) (string, error) {
	// 1. Create a ZIP archive in memory.
	zipBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuffer)
	xmlFile, err := zipWriter.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("error al crear el archivo XML en el zip: %w", err)
	}
	_, err = xmlFile.Write(signedXML)
	if err != nil {
		return "", fmt.Errorf("error al escribir el XML en el zip: %w", err)
	}
	if err := zipWriter.Close(); err != nil {
		return "", fmt.Errorf("error al cerrar el archivo zip: %w", err)
	}

	// 2. Base64 encode the ZIP content.
	encodedZip := base64.StdEncoding.EncodeToString(zipBuffer.Bytes())

	// 3. Prepare the SOAP request.
	req := &SendBillRequest{
		FileName:    fileName,
		ContentFile: encodedZip,
	}

	resp := &SendBillResponse{}

	fmt.Printf("Enviando %s a SUNAT via SOAP...\n", filepath.Base(fileName))
	// gosoap usa Map para los parámetros, pero también soporta structs.
	params := gosoap.Params{
		"fileName":    req.FileName,
		"contentFile": req.ContentFile,
	}
	soapResp, err := c.soapClient.Call("sendBill", params)
	if err != nil {
		return "", fmt.Errorf("error al llamar al servicio sendBill: %w", err)
	}
	// Decodifica la respuesta en tu struct
	if err := soapResp.Unmarshal(&resp); err != nil {
		return "", fmt.Errorf("error al decodificar respuesta sendBill: %w", err)
	}

	fmt.Printf("Ticket recibido de SUNAT: %s\n", resp.Ticket)
	return resp.Ticket, nil
}

// GetStatus checks the status of a previously sent document using its ticket ID via SOAP.
// It returns the status object from SUNAT.
func (c *Client) GetStatus(ticketID string) (*Status, error) {
	req := &GetStatusRequest{
		Ticket: ticketID,
	}
	resp := &GetStatusResponse{}

	fmt.Printf("Consultando estado del ticket %s en SUNAT via SOAP...\n", ticketID)
	params := gosoap.Params{
		"ticket": req.Ticket,
	}
	soapResp, err := c.soapClient.Call("getStatus", params)
	if err != nil {
		return nil, fmt.Errorf("error al llamar al servicio getStatus: %w", err)
	}
	if err := soapResp.Unmarshal(&resp); err != nil {
		return nil, fmt.Errorf("error al decodificar respuesta getStatus: %w", err)
	}

	fmt.Printf("Estado recibido para ticket %s: %s\n", ticketID, resp.Status.StatusCode)
	return resp.Status, nil
}

// GetStatusCdr retrieves the status and CDR of a document using its full details via SOAP.
func (c *Client) GetStatusCdr(ruc, docType, series, number string) (*StatusCdr, error) {
	req := &GetStatusCdrRequest{
		RucComprobante:    ruc,
		TipoComprobante:   docType,
		SerieComprobante:  series,
		NumeroComprobante: number,
	}
	resp := &GetStatusCdrResponse{}

	fmt.Printf("Consultando CDR para %s-%s-%s-%s en SUNAT via SOAP...\n", ruc, docType, series, number)
	params := gosoap.Params{
		"rucComprobante":    req.RucComprobante,
		"tipoComprobante":   req.TipoComprobante,
		"serieComprobante":  req.SerieComprobante,
		"numeroComprobante": req.NumeroComprobante,
	}
	soapResp, err := c.soapClient.Call("getStatusCdr", params)
	if err != nil {
		return nil, fmt.Errorf("error al llamar al servicio getStatusCdr: %w", err)
	}
	if err := soapResp.Unmarshal(&resp); err != nil {
		return nil, fmt.Errorf("error al decodificar respuesta getStatusCdr: %w", err)
	}

	fmt.Printf("Estado CDR recibido: %s\n", resp.StatusCdr.StatusCode)
	return resp.StatusCdr, nil
}
