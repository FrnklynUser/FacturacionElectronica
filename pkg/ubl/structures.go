package ubl

import "encoding/xml"

// Namespaces used in UBL XML
const (
	DS   = "http://www.w3.org/2000/09/xmldsig#"
	CAC  = "urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2"
	CBC  = "urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2"
	EXT  = "urn:oasis:names:specification:ubl:schema:xsd:CommonExtensionComponents-2"
	SAC  = "urn:sunat:names:specification:ubl:peru:schema:xsd:SunatAggregateComponents-1"
	CCTS = "urn:un:unece:uncefact:documentation:2"
	QDT  = "urn:oasis:names:specification:ubl:schema:xsd:QualifiedDatatypes-2"
	UDT  = "urn:un:unece:uncefact:data:specification:UnqualifiedDataTypesSchemaModule:2"
	XSI  = "http://www.w3.org/2001/XMLSchema-instance"
)

// Invoice is the top-level UBL Invoice structure
type Invoice struct {
	XMLName                     xml.Name                       `xml:"Invoice"`
	Xmlns                       string                         `xml:"xmlns,attr"`
	XmlnsCAC                    string                         `xml:"xmlns:cac,attr"`
	XmlnsCBC                    string                         `xml:"xmlns:cbc,attr"`
	XmlnsCCTS                   string                         `xml:"xmlns:ccts,attr"`
	XmlnsDS                     string                         `xml:"xmlns:ds,attr"`
	XmlnsEXT                    string                         `xml:"xmlns:ext,attr"`
	XmlnsQDT                    string                         `xml:"xmlns:qdt,attr"`
	XmlnsUDT                    string                         `xml:"xmlns:udt,attr"`
	XmlnsXSI                    string                         `xml:"xmlns:xsi,attr"`
	UBLExtensions               *UBLExtensions                 `xml:"ext:UBLExtensions"`
	UBLVersionID                string                         `xml:"cbc:UBLVersionID"`
	CustomizationID             string                         `xml:"cbc:CustomizationID"`
	ProfileID                   *ProfileID                     `xml:"cbc:ProfileID"`
	ID                          string                         `xml:"cbc:ID"` // Serie-Numero
	IssueDate                   string                         `xml:"cbc:IssueDate"`
	IssueTime                   string                         `xml:"cbc:IssueTime"`
	InvoiceTypeCode             *InvoiceTypeCode               `xml:"cbc:InvoiceTypeCode"`
	DocumentCurrencyCode        *DocumentCurrencyCode          `xml:"cbc:DocumentCurrencyCode"`
	LineCountNumeric            string                         `xml:"cbc:LineCountNumeric"`            // New for Invoice
	OrderReference              *OrderReference                `xml:"cac:OrderReference"`              // New for Invoice
	DespatchDocumentReference   []*DespatchDocumentReference   `xml:"cac:DespatchDocumentReference"`   // New for Invoice
	AdditionalDocumentReference []*AdditionalDocumentReference `xml:"cac:AdditionalDocumentReference"` // New for Invoice
	Signature                   *Signature                     `xml:"cac:Signature"`
	AccountingSupplierParty     *Supplier                      `xml:"cac:AccountingSupplierParty"`
	AccountingCustomerParty     *Customer                      `xml:"cac:AccountingCustomerParty"`
	TaxTotals                   []*TaxTotal                    `xml:"cac:TaxTotal"`
	LegalMonetaryTotal          *MonetaryTotal                 `xml:"cac:LegalMonetaryTotal"`
	InvoiceLines                []*InvoiceLine                 `xml:"cac:InvoiceLine"`
}

// UBLExtensions is the container for UBL extensions.
type UBLExtensions struct {
	UBLExtension *UBLExtension `xml:"ext:UBLExtension"`
}

// UBLExtension contains the extension content.
type UBLExtension struct {
	ExtensionContent *ExtensionContent `xml:"ext:ExtensionContent"`
}

// ExtensionContent holds the actual content of the extension, e.g., the digital signature.
type ExtensionContent struct {
	XMLName xml.Name `xml:",any"` // This will capture the ds:Signature element
}

// ProfileID defines the operation type.
type ProfileID struct {
	XMLName          xml.Name `xml:"cbc:ProfileID"`
	SchemeName       string   `xml:"schemeName,attr"`
	SchemeAgencyName string   `xml:"schemeAgencyName,attr"`
	SchemeURI        string   `xml:"schemeURI,attr"`
	Value            string   `xml:",chardata"`
}

// InvoiceTypeCode defines the type of document.
type InvoiceTypeCode struct {
	XMLName        xml.Name `xml:"cbc:InvoiceTypeCode"`
	ListAgencyName string   `xml:"listAgencyName,attr"`
	ListName       string   `xml:"listName,attr"`
	ListURI        string   `xml:"listURI,attr"`
	Value          string   `xml:",chardata"`
}

// DocumentCurrencyCode defines the currency of the document.
type DocumentCurrencyCode struct {
	XMLName        xml.Name `xml:"cbc:DocumentCurrencyCode"`
	ListID         string   `xml:"listID,attr"`
	ListName       string   `xml:"listName,attr"`
	ListAgencyName string   `xml:"listAgencyName,attr"`
	Value          string   `xml:",chardata"`
}

// OrderReference holds the purchase order number.
type OrderReference struct {
	XMLName xml.Name `xml:"cac:OrderReference"`
	ID      string   `xml:"cbc:ID"`
}

// DespatchDocumentReference holds reference to a despatch advice.
type DespatchDocumentReference struct {
	XMLName          xml.Name         `xml:"cac:DespatchDocumentReference"`
	ID               string           `xml:"cbc:ID"`
	DocumentTypeCode *InvoiceTypeCode `xml:"cbc:DocumentTypeCode"` // Reusing InvoiceTypeCode for simplicity
}

// AdditionalDocumentReference holds reference to other related documents.
type AdditionalDocumentReference struct {
	XMLName          xml.Name         `xml:"cac:AdditionalDocumentReference"`
	ID               string           `xml:"cbc:ID"`
	DocumentTypeCode *InvoiceTypeCode `xml:"cbc:DocumentTypeCode"` // Reusing InvoiceTypeCode for simplicity
}

// Signature holds the digital signature information
type Signature struct {
	ID                         string                      `xml:"cbc:ID"`
	SignatoryParty             *SignatoryParty             `xml:"cac:SignatoryParty"`
	DigitalSignatureAttachment *DigitalSignatureAttachment `xml:"cac:DigitalSignatureAttachment"`
}

// SignatoryParty represents a party in the signature
type SignatoryParty struct {
	PartyIdentification *PartyIdentification `xml:"cac:PartyIdentification"`
	PartyName           *PartyName           `xml:"cac:PartyName"`
}

// PartyIdentification identifies a party
type PartyIdentification struct {
	ID string `xml:"cbc:ID"`
}

// PartyName represents the name of a party
type PartyName struct {
	Name string `xml:"cbc:Name"`
}

// DigitalSignatureAttachment holds the reference to the signature
type DigitalSignatureAttachment struct {
	ExternalReference *ExternalReference `xml:"cac:ExternalReference"`
}

// ExternalReference holds the URI of the external reference
type ExternalReference struct {
	URI string `xml:"cbc:URI"`
}

// Supplier represents the party issuing the invoice
type Supplier struct {
	CustomerAssignedAccountID string `xml:"cbc:CustomerAssignedAccountID"` // RUC
	AdditionalAccountID       string `xml:"cbc:AdditionalAccountID"`       // 6 for RUC
	Party                     *Party `xml:"cac:Party"`
}

// Customer represents the party receiving the invoice
type Customer struct {
	CustomerAssignedAccountID string `xml:"cbc:CustomerAssignedAccountID"` // RUC/DNI
	AdditionalAccountID       string `xml:"cbc:AdditionalAccountID"`       // 6 for RUC, 1 for DNI
	Party                     *Party `xml:"cac:Party"`
}

// Party is a common struct for both Supplier and Customer
type Party struct {
	PartyName        *PartyName        `xml:"cac:PartyName"`
	PartyLegalEntity *PartyLegalEntity `xml:"cac:PartyLegalEntity"`
	PartyTaxScheme   *PartyTaxScheme   `xml:"cac:PartyTaxScheme"`
}

// PartyLegalEntity holds the legal name of the party
type PartyLegalEntity struct {
	RegistrationName string `xml:"cbc:RegistrationName"`
}

// PartyTaxScheme holds tax scheme information for a party
type PartyTaxScheme struct {
	RegistrationName    string               `xml:"cbc:RegistrationName"`
	CompanyID           *CompanyID           `xml:"cbc:CompanyID"`
	RegistrationAddress *RegistrationAddress `xml:"cac:RegistrationAddress"`
	TaxScheme           *TaxScheme           `xml:"cac:TaxScheme"`
}

// CompanyID identifies the company
type CompanyID struct {
	XMLName          xml.Name `xml:"cbc:CompanyID"`
	SchemeID         string   `xml:"schemeID,attr"`
	SchemeName       string   `xml:"schemeName,attr"`
	SchemeAgencyName string   `xml:"schemeAgencyName,attr"`
	SchemeURI        string   `xml:"schemeURI,attr"`
	Value            string   `xml:",chardata"`
}

// RegistrationAddress holds the registration address details
type RegistrationAddress struct {
	AddressTypeCode string `xml:"cbc:AddressTypeCode"`
}

// TaxTotal aggregates the total tax amounts
type TaxTotal struct {
	TaxAmount   *Amount        `xml:"cbc:TaxAmount"`
	TaxSubtotal []*TaxSubtotal `xml:"cac:TaxSubtotal"`
}

// TaxSubtotal is a sub-total for a specific tax category
type TaxSubtotal struct {
	TaxableAmount *Amount      `xml:"cbc:TaxableAmount"`
	TaxAmount     *Amount      `xml:"cbc:TaxAmount"`
	TaxCategory   *TaxCategory `xml:"cac:TaxCategory"`
}

// TaxCategory defines the type of tax
type TaxCategory struct {
	XMLName                xml.Name                `xml:"cac:TaxCategory"`
	ID                     string                  `xml:"cbc:ID"`
	SchemeID               string                  `xml:"schemeID,attr"`
	SchemeName             string                  `xml:"schemeName,attr"`
	SchemeAgencyName       string                  `xml:"schemeAgencyName,attr"`
	Percent                float64                 `xml:"cbc:Percent"`
	TaxExemptionReasonCode *TaxExemptionReasonCode `xml:"cbc:TaxExemptionReasonCode"`
	TaxScheme              *TaxScheme              `xml:"cac:TaxScheme"`
}

// TaxExemptionReasonCode defines the reason for tax exemption
type TaxExemptionReasonCode struct {
	XMLName        xml.Name `xml:"cbc:TaxExemptionReasonCode"`
	ListAgencyName string   `xml:"listAgencyName,attr"`
	ListName       string   `xml:"listName,attr"`
	ListURI        string   `xml:"listURI,attr"`
	Value          string   `xml:",chardata"`
}

// TaxScheme provides details about the tax
type TaxScheme struct {
	XMLName        xml.Name `xml:"cac:TaxScheme"`
	ID             string   `xml:"cbc:ID"` // e.g., 1000
	SchemeID       string   `xml:"schemeID,attr"`
	SchemeAgencyID string   `xml:"schemeAgencyID,attr"`
	Name           string   `xml:"cbc:Name"`        // e.g., IGV
	TaxTypeCode    string   `xml:"cbc:TaxTypeCode"` // e.g., VAT
}

// MonetaryTotal represents the final total amount of the invoice
type MonetaryTotal struct {
	LineExtensionAmount  *Amount `xml:"cbc:LineExtensionAmount"`
	TaxInclusiveAmount   *Amount `xml:"cbc:TaxInclusiveAmount"`
	AllowanceTotalAmount *Amount `xml:"cbc:AllowanceTotalAmount"`
	ChargeTotalAmount    *Amount `xml:"cbc:ChargeTotalAmount"`
	PrepaidAmount        *Amount `xml:"cbc:PrepaidAmount"`
	PayableAmount        *Amount `xml:"cbc:PayableAmount"`
}

// InvoiceLine represents one line item on the invoice
type InvoiceLine struct {
	ID                  string            `xml:"cbc:ID"` // Line number
	InvoicedQuantity    *Quantity         `xml:"cbc:InvoicedQuantity"`
	LineExtensionAmount *Amount           `xml:"cbc:LineExtensionAmount"` // Total sin IGV
	PricingReference    *PricingReference `xml:"cac:PricingReference"`
	TaxTotals           []*TaxTotal       `xml:"cac:TaxTotal"`
	Item                *Item             `xml:"cac:Item"`
	Price               *Price            `xml:"cac:Price"`
}

// PricingReference holds pricing details
type PricingReference struct {
	AlternativeConditionPrice []*Price `xml:"cac:AlternativeConditionPrice"`
}

// Item is the product or service being sold
type Item struct {
	Description               string                     `xml:"cbc:Description"`
	SellersItemIdentification *SellersItemIdentification `xml:"cac:SellersItemIdentification"`
	CommodityClassification   *CommodityClassification   `xml:"cac:CommodityClassification"`
}

// SellersItemIdentification identifies the item by the seller
type SellersItemIdentification struct {
	ID string `xml:"cbc:ID"`
}

// CommodityClassification classifies the item
type CommodityClassification struct {
	ItemClassificationCode *ItemClassificationCode `xml:"cbc:ItemClassificationCode"`
}

// ItemClassificationCode holds the classification code
type ItemClassificationCode struct {
	XMLName        xml.Name `xml:"cbc:ItemClassificationCode"`
	ListID         string   `xml:"listID,attr"`
	ListAgencyName string   `xml:"listAgencyName,attr"`
	ListName       string   `xml:"listName,attr"`
	Value          string   `xml:",chardata"`
}

// Price contains the price of an item
type Price struct {
	PriceAmount   *Amount        `xml:"cbc:PriceAmount"`
	PriceTypeCode *PriceTypeCode `xml:"cbc:PriceTypeCode"`
}

// PriceTypeCode defines the type of price
type PriceTypeCode struct {
	XMLName        xml.Name `xml:"cbc:PriceTypeCode"`
	ListAgencyName string   `xml:"listAgencyName,attr"`
	ListName       string   `xml:"listName,attr"`
	ListURI        string   `xml:"listURI,attr"`
	Value          string   `xml:",chardata"`
}

// Amount is a numeric value with a currency attribute
type Amount struct {
	CurrencyID string  `xml:"currencyID,attr"`
	Value      float64 `xml:",chardata"`
}

// Quantity is a numeric value with a unit code attribute
type Quantity struct {
	UnitCode string  `xml:"unitCode,attr"`
	Value    float64 `xml:",chardata"`
}

// BillingReference references a previous document
type BillingReference struct {
	InvoiceDocumentReference *InvoiceDocumentReference `xml:"cac:InvoiceDocumentReference"`
}

// InvoiceDocumentReference holds the ID of the referenced invoice
type InvoiceDocumentReference struct {
	ID               string           `xml:"cbc:ID"`
	DocumentTypeCode *InvoiceTypeCode `xml:"cbc:DocumentTypeCode"` // Reusing InvoiceTypeCode for simplicity
}

// DiscrepancyResponse describes the reason for the note
type DiscrepancyResponse struct {
	ReferenceID  string `xml:"cbc:ReferenceID"`
	ResponseCode string `xml:"cbc:ResponseCode"`
	Description  string `xml:"cbc:Description"`
}

// CreditNote is the top-level UBL CreditNote structure
type CreditNote struct {
	XMLName                 xml.Name              `xml:"CreditNote"`
	Xmlns                   string                `xml:"xmlns,attr"`
	XmlnsCAC                string                `xml:"xmlns:cac,attr"`
	XmlnsCBC                string                `xml:"xmlns:cbc,attr"`
	XmlnsCCTS               string                `xml:"xmlns:ccts,attr"`
	XmlnsDS                 string                `xml:"xmlns:ds,attr"`
	XmlnsEXT                string                `xml:"xmlns:ext,attr"`
	XmlnsQDT                string                `xml:"xmlns:qdt,attr"`
	XmlnsUDT                string                `xml:"xmlns:udt,attr"`
	XmlnsXSI                string                `xml:"xmlns:xsi,attr"`
	UBLExtensions           *UBLExtensions        `xml:"ext:UBLExtensions"`
	UBLVersionID            string                `xml:"cbc:UBLVersionID"`
	CustomizationID         string                `xml:"cbc:CustomizationID"`
	ID                      string                `xml:"cbc:ID"` // Serie-Numero
	IssueDate               string                `xml:"cbc:IssueDate"`
	IssueTime               string                `xml:"cbc:IssueTime"`
	DocumentCurrencyCode    *DocumentCurrencyCode `xml:"cbc:DocumentCurrencyCode"`
	DiscrepancyResponse     *DiscrepancyResponse  `xml:"cac:DiscrepancyResponse"`
	BillingReference        *BillingReference     `xml:"cac:BillingReference"`
	Signature               *Signature            `xml:"cac:Signature"`
	AccountingSupplierParty *Supplier             `xml:"cac:AccountingSupplierParty"`
	AccountingCustomerParty *Customer             `xml:"cac:AccountingCustomerParty"`
	TaxTotals               []*TaxTotal           `xml:"cac:TaxTotal"`
	LegalMonetaryTotal      *MonetaryTotal        `xml:"cac:LegalMonetaryTotal"`
	CreditNoteLines         []*InvoiceLine        `xml:"cac:CreditNoteLine"` // Note: CreditNoteLine is same as InvoiceLine
}

// DebitNote is the top-level UBL DebitNote structure
type DebitNote struct {
	XMLName                 xml.Name              `xml:"DebitNote"`
	Xmlns                   string                `xml:"xmlns,attr"`
	XmlnsCAC                string                `xml:"xmlns:cac,attr"`
	XmlnsCBC                string                `xml:"xmlns:cbc,attr"`
	XmlnsCCTS               string                `xml:"xmlns:ccts,attr"`
	XmlnsDS                 string                `xml:"xmlns:ds,attr"`
	XmlnsEXT                string                `xml:"xmlns:ext,attr"`
	XmlnsQDT                string                `xml:"xmlns:qdt,attr"`
	XmlnsUDT                string                `xml:"xmlns:udt,attr"`
	XmlnsXSI                string                `xml:"xmlns:xsi,attr"`
	UBLExtensions           *UBLExtensions        `xml:"ext:UBLExtensions"`
	UBLVersionID            string                `xml:"cbc:UBLVersionID"`
	CustomizationID         string                `xml:"cbc:CustomizationID"`
	ID                      string                `xml:"cbc:ID"` // Serie-Numero
	IssueDate               string                `xml:"cbc:IssueDate"`
	IssueTime               string                `xml:"cbc:IssueTime"`
	DocumentCurrencyCode    *DocumentCurrencyCode `xml:"cbc:DocumentCurrencyCode"`
	DiscrepancyResponse     *DiscrepancyResponse  `xml:"cac:DiscrepancyResponse"`
	BillingReference        *BillingReference     `xml:"cac:BillingReference"`
	Signature               *Signature            `xml:"cac:Signature"`
	AccountingSupplierParty *Supplier             `xml:"cac:AccountingSupplierParty"`
	AccountingCustomerParty *Customer             `xml:"cac:AccountingCustomerParty"`
	TaxTotals               []*TaxTotal           `xml:"cac:TaxTotal"`
	LegalMonetaryTotal      *MonetaryTotal        `xml:"cac:LegalMonetaryTotal"`
	DebitNoteLines          []*InvoiceLine        `xml:"cac:DebitNoteLine"` // Note: DebitNoteLine is same as InvoiceLine
}
