package ubl

import (
	"FacturacionSunat/internal/domain"
	"encoding/xml"
	"fmt"
	"strconv"
)

// BuildInvoice transforms a domain.Invoice into a UBL Invoice structure ready for XML marshalling.
func BuildInvoice(inv *domain.Invoice) (*Invoice, error) {
	// Default namespaces and attributes
	ublInvoice := &Invoice{
		Xmlns:           "urn:oasis:names:specification:ubl:schema:xsd:Invoice-2",
		XmlnsCAC:        CAC,
		XmlnsCBC:        CBC,
		XmlnsCCTS:       CCTS,
		XmlnsDS:         DS,
		XmlnsEXT:        EXT,
		XmlnsQDT:        QDT,
		XmlnsUDT:        UDT,
		XmlnsXSI:        XSI,
		UBLVersionID:    "2.1",
		CustomizationID: "2.0",
		ProfileID: &ProfileID{
			SchemeName:       "SUNAT:Identificador de Tipo de Operación",
			SchemeAgencyName: "PE:SUNAT",
			SchemeURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo17",
			Value:            "0101", // For Boleta, Venta Interna
		},
		ID:        fmt.Sprintf("%s-%d", inv.Series, inv.Number),
		IssueDate: inv.IssueDate.Format("2006-01-02"),
		IssueTime: inv.IssueDate.Format("15:04:05"), // HH:MM:SS
		InvoiceTypeCode: &InvoiceTypeCode{
			ListAgencyName: "PE:SUNAT",
			ListName:       "SUNAT:Identificador de Tipo de Documento",
			ListURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo01",
			Value:          inv.Type, // e.g., "03" for Boleta, "01" for Factura
		},
		DocumentCurrencyCode: &DocumentCurrencyCode{
			ListID:         "ISO 4217 Alpha",
			ListName:       "Currency",
			ListAgencyName: "United Nations Economic Commission for Europe",
			Value:          inv.Currency,
		},
		Signature: &Signature{
			ID: "IDSignSP",
			SignatoryParty: &SignatoryParty{
				PartyIdentification: &PartyIdentification{ID: inv.Issuer.RUC},
				PartyName:           &PartyName{Name: inv.Issuer.Name},
			},
			DigitalSignatureAttachment: &DigitalSignatureAttachment{
				ExternalReference: &ExternalReference{URI: "#IDSignSP"},
			},
		},
		AccountingSupplierParty: &Supplier{
			CustomerAssignedAccountID: inv.Issuer.RUC,
			AdditionalAccountID:       "6", // RUC
			Party: &Party{
				PartyName:        &PartyName{Name: inv.Issuer.Name},
				PartyLegalEntity: &PartyLegalEntity{RegistrationName: inv.Issuer.Name},
				PartyTaxScheme: &PartyTaxScheme{
					RegistrationName: inv.Issuer.Name,
					CompanyID: &CompanyID{
						SchemeID:         "6",
						SchemeName:       "SUNAT:Identificador de Documento de Identidad",
						SchemeAgencyName: "PE:SUNAT",
						SchemeURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo06",
						Value:            inv.Issuer.RUC,
					},
					RegistrationAddress: &RegistrationAddress{AddressTypeCode: "0000"}, // Default for fiscal address
					TaxScheme:           &TaxScheme{ID: "-"},
				},
			},
		},
		AccountingCustomerParty: &Customer{
			CustomerAssignedAccountID: inv.Recipient.DocNum,
			AdditionalAccountID:       getDocType(inv.Recipient.DocType),
			Party: &Party{
				PartyLegalEntity: &PartyLegalEntity{RegistrationName: inv.Recipient.Name},
				PartyTaxScheme: &PartyTaxScheme{
					RegistrationName: inv.Recipient.Name,
					CompanyID: &CompanyID{
						SchemeID:         getDocType(inv.Recipient.DocType),
						SchemeName:       "SUNAT:Identificador de Documento de Identidad",
						SchemeAgencyName: "PE:SUNAT",
						SchemeURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo06",
						Value:            inv.Recipient.DocNum,
					},
					TaxScheme: &TaxScheme{ID: "-"},
				},
			},
		},
		LegalMonetaryTotal: &MonetaryTotal{
			PayableAmount:       &Amount{CurrencyID: inv.Currency, Value: inv.Totals.Total},
			LineExtensionAmount: &Amount{CurrencyID: inv.Currency, Value: inv.Totals.Gross},
			TaxInclusiveAmount:  &Amount{CurrencyID: inv.Currency, Value: inv.Totals.Total + inv.Totals.IGV}, // Assuming Total is Net + IGV
		},
	}

	// Tax Totals
	igvTaxTotal := &TaxTotal{
		TaxAmount: &Amount{CurrencyID: inv.Currency, Value: inv.Totals.IGV},
		TaxSubtotal: []*TaxSubtotal{{
			TaxableAmount: &Amount{CurrencyID: inv.Currency, Value: inv.Totals.Gross},
			TaxAmount:     &Amount{CurrencyID: inv.Currency, Value: inv.Totals.IGV},
			TaxCategory: &TaxCategory{
				ID:               "S",
				SchemeID:         "UN/ECE 5305",
				SchemeName:       "Tax Category Identifier",
				SchemeAgencyName: "United Nations Economic Commission for Europe",
				Percent:          18.00, // Assuming 18% IGV
				TaxExemptionReasonCode: &TaxExemptionReasonCode{
					ListAgencyName: "PE:SUNAT",
					ListName:       "SUNAT:Codigo de Tipo de Afectación del IGV",
					ListURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo07",
					Value:          "10", // Gravado - Operación Onerosa
				},
				TaxScheme: &TaxScheme{
					ID:             "1000",
					SchemeID:       "UN/ECE 5153",
					SchemeAgencyID: "6",
					Name:           "IGV",
					TaxTypeCode:    "VAT",
				},
			},
		}},
	}
	ublInvoice.TaxTotals = append(ublInvoice.TaxTotals, igvTaxTotal)

	// Invoice Lines
	for i, line := range inv.Lines {
		ublLine := &InvoiceLine{
			ID:                  strconv.Itoa(i + 1),
			InvoicedQuantity:    &Quantity{UnitCode: "NIU", Value: line.Quantity}, // Assuming NIU, should be configurable
			LineExtensionAmount: &Amount{CurrencyID: inv.Currency, Value: line.TotalValue},
			PricingReference: &PricingReference{
				AlternativeConditionPrice: []*Price{{
					PriceAmount: &Amount{CurrencyID: inv.Currency, Value: line.UnitPrice},
					PriceTypeCode: &PriceTypeCode{
						ListAgencyName: "PE:SUNAT",
						ListName:       "SUNAT:Indicador de Tipo de Precio",
						ListURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo16",
						Value:          "01", // Precio unitario (incluye IGV)
					},
				}},
			},
			Item: &Item{
				Description: line.Description,
				// SellersItemIdentification: &SellersItemIdentification{ID: ""},
				// CommodityClassification: &CommodityClassification{ItemClassificationCode: &ItemClassificationCode{Value: ""}},
			},
			Price: &Price{PriceAmount: &Amount{CurrencyID: inv.Currency, Value: line.UnitPrice}},
		}

		// Line Tax Totals (simplified for now, assuming only IGV)
		lineTaxTotal := &TaxTotal{
			TaxAmount: &Amount{CurrencyID: inv.Currency, Value: line.IGV},
			TaxSubtotal: []*TaxSubtotal{{
				TaxableAmount: &Amount{CurrencyID: inv.Currency, Value: line.TotalValue},
				TaxAmount:     &Amount{CurrencyID: inv.Currency, Value: line.IGV},
				TaxCategory: &TaxCategory{
					ID:               "S",
					SchemeID:         "UN/ECE 5305",
					SchemeName:       "Tax Category Identifier",
					SchemeAgencyName: "United Nations Economic Commission for Europe",
					Percent:          18.00,
					TaxExemptionReasonCode: &TaxExemptionReasonCode{
						ListAgencyName: "PE:SUNAT",
						ListName:       "SUNAT:Codigo de Tipo de Afectación del IGV",
						ListURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo07",
						Value:          "10", // Gravado - Operación Onerosa
					},
					TaxScheme: &TaxScheme{
						ID:             "1000",
						SchemeID:       "UN/ECE 5153",
						SchemeAgencyID: "6",
						Name:           "IGV",
						TaxTypeCode:    "VAT",
					},
				},
			}},
		}
		ublLine.TaxTotals = append(ublLine.TaxTotals, lineTaxTotal)

		ublInvoice.InvoiceLines = append(ublInvoice.InvoiceLines, ublLine)
	}

	// Set UBLExtensions for signature
	ublInvoice.UBLExtensions = &UBLExtensions{
		UBLExtension: &UBLExtension{
			ExtensionContent: &ExtensionContent{XMLName: xml.Name{Local: "ds:Signature"}}, // Placeholder for goxades to fill
		},
	}

	return ublInvoice, nil
}

// BuildCreditNote transforms a domain.CreditNote into a UBL CreditNote structure.
func BuildCreditNote(cn *domain.CreditNote) (*CreditNote, error) {
	ublCreditNote := &CreditNote{
		Xmlns:     "urn:oasis:names:specification:ubl:schema:xsd:CreditNote-2",
		XmlnsCAC:  CAC,
		XmlnsCBC:  CBC,
		XmlnsCCTS: CCTS,
		XmlnsDS:   DS,
		XmlnsEXT:  EXT,
		XmlnsQDT:  QDT,
		XmlnsUDT:  UDT,
		XmlnsXSI:  XSI,
		UBLExtensions: &UBLExtensions{
			UBLExtension: &UBLExtension{
				ExtensionContent: &ExtensionContent{XMLName: xml.Name{Local: "ds:Signature"}},
			},
		},
		UBLVersionID:    "2.1",
		CustomizationID: "2.0",
		ID:              fmt.Sprintf("%s-%d", cn.Series, cn.Number),
		IssueDate:       cn.IssueDate.Format("2006-01-02"),
		IssueTime:       cn.IssueDate.Format("15:04:05"),
		DocumentCurrencyCode: &DocumentCurrencyCode{
			ListID:         "ISO 4217 Alpha",
			ListName:       "Currency",
			ListAgencyName: "United Nations Economic Commission for Europe",
			Value:          cn.Currency,
		},
		DiscrepancyResponse: &DiscrepancyResponse{
			ReferenceID:  cn.DiscrepancyResponse.ReferenceID,
			ResponseCode: cn.DiscrepancyResponse.TypeCode,
			Description:  cn.DiscrepancyResponse.Description,
		},
		BillingReference: &BillingReference{
			InvoiceDocumentReference: &InvoiceDocumentReference{ID: cn.DiscrepancyResponse.ReferenceID, DocumentTypeCode: &InvoiceTypeCode{Value: "01"}}, // Assuming 01 for Invoice
		},
		Signature: &Signature{
			ID: "IDSignSP",
			SignatoryParty: &SignatoryParty{
				PartyIdentification: &PartyIdentification{ID: cn.Issuer.RUC},
				PartyName:           &PartyName{Name: cn.Issuer.Name},
			},
			DigitalSignatureAttachment: &DigitalSignatureAttachment{
				ExternalReference: &ExternalReference{URI: "#IDSignSP"},
			},
		},
		AccountingSupplierParty: &Supplier{
			CustomerAssignedAccountID: cn.Issuer.RUC,
			AdditionalAccountID:       "6", // RUC
			Party: &Party{
				PartyName:        &PartyName{Name: cn.Issuer.Name},
				PartyLegalEntity: &PartyLegalEntity{RegistrationName: cn.Issuer.Name},
				PartyTaxScheme: &PartyTaxScheme{
					RegistrationName: cn.Issuer.Name,
					CompanyID: &CompanyID{
						SchemeID:         "6",
						SchemeName:       "SUNAT:Identificador de Documento de Identidad",
						SchemeAgencyName: "PE:SUNAT",
						SchemeURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo06",
						Value:            cn.Issuer.RUC,
					},
					RegistrationAddress: &RegistrationAddress{AddressTypeCode: "0000"},
					TaxScheme:           &TaxScheme{ID: "-"},
				},
			},
		},
		AccountingCustomerParty: &Customer{
			CustomerAssignedAccountID: cn.Recipient.DocNum,
			AdditionalAccountID:       getDocType(cn.Recipient.DocType),
			Party: &Party{
				PartyLegalEntity: &PartyLegalEntity{RegistrationName: cn.Recipient.Name},
				PartyTaxScheme: &PartyTaxScheme{
					RegistrationName: cn.Recipient.Name,
					CompanyID: &CompanyID{
						SchemeID:         getDocType(cn.Recipient.DocType),
						SchemeName:       "SUNAT:Identificador de Documento de Identidad",
						SchemeAgencyName: "PE:SUNAT",
						SchemeURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo06",
						Value:            cn.Recipient.DocNum,
					},
					TaxScheme: &TaxScheme{ID: "-"},
				},
			},
		},
		LegalMonetaryTotal: &MonetaryTotal{
			PayableAmount:       &Amount{CurrencyID: cn.Currency, Value: cn.Totals.Total},
			LineExtensionAmount: &Amount{CurrencyID: cn.Currency, Value: cn.Totals.Gross},
			TaxInclusiveAmount:  &Amount{CurrencyID: cn.Currency, Value: cn.Totals.Total + cn.Totals.IGV},
		},
	}

	// Tax Totals (simplified for now, assuming only IGV)
	igvTaxTotal := &TaxTotal{
		TaxAmount: &Amount{CurrencyID: cn.Currency, Value: cn.Totals.IGV},
		TaxSubtotal: []*TaxSubtotal{{
			TaxableAmount: &Amount{CurrencyID: cn.Currency, Value: cn.Totals.Gross},
			TaxAmount:     &Amount{CurrencyID: cn.Currency, Value: cn.Totals.IGV},
			TaxCategory: &TaxCategory{
				ID:               "S",
				SchemeID:         "UN/ECE 5305",
				SchemeName:       "Tax Category Identifier",
				SchemeAgencyName: "United Nations Economic Commission for Europe",
				Percent:          18.00, // Assuming 18% IGV
				TaxExemptionReasonCode: &TaxExemptionReasonCode{
					ListAgencyName: "PE:SUNAT",
					ListName:       "SUNAT:Codigo de Tipo de Afectación del IGV",
					ListURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo07",
					Value:          "10", // Gravado - Operación Onerosa
				},
				TaxScheme: &TaxScheme{
					ID:             "1000",
					SchemeID:       "UN/ECE 5153",
					SchemeAgencyID: "6",
					Name:           "IGV",
					TaxTypeCode:    "VAT",
				},
			},
		}},
	}
	ublCreditNote.TaxTotals = append(ublCreditNote.TaxTotals, igvTaxTotal)

	// Credit Note Lines
	for i, line := range cn.Lines {
		ublLine := &InvoiceLine{
			ID:                  strconv.Itoa(i + 1),
			InvoicedQuantity:    &Quantity{UnitCode: "NIU", Value: line.Quantity}, // Assuming NIU, should be configurable
			LineExtensionAmount: &Amount{CurrencyID: cn.Currency, Value: line.TotalValue},
			PricingReference: &PricingReference{
				AlternativeConditionPrice: []*Price{{
					PriceAmount: &Amount{CurrencyID: cn.Currency, Value: line.UnitPrice},
					PriceTypeCode: &PriceTypeCode{
						ListAgencyName: "PE:SUNAT",
						ListName:       "SUNAT:Indicador de Tipo de Precio",
						ListURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo16",
						Value:          "01", // Precio unitario (incluye IGV)
					},
				}},
			},
			Item: &Item{
				Description: line.Description,
			},
			Price: &Price{PriceAmount: &Amount{CurrencyID: cn.Currency, Value: line.UnitPrice}},
		}

		// Line Tax Totals (simplified for now, assuming only IGV)
		lineTaxTotal := &TaxTotal{
			TaxAmount: &Amount{CurrencyID: cn.Currency, Value: line.IGV},
			TaxSubtotal: []*TaxSubtotal{{
				TaxableAmount: &Amount{CurrencyID: cn.Currency, Value: line.TotalValue},
				TaxAmount:     &Amount{CurrencyID: cn.Currency, Value: line.IGV},
				TaxCategory: &TaxCategory{
					ID:               "S",
					SchemeID:         "UN/ECE 5305",
					SchemeName:       "Tax Category Identifier",
					SchemeAgencyName: "United Nations Economic Commission for Europe",
					Percent:          18.00,
					TaxExemptionReasonCode: &TaxExemptionReasonCode{
						ListAgencyName: "PE:SUNAT",
						ListName:       "SUNAT:Codigo de Tipo de Afectación del IGV",
						ListURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo07",
						Value:          "10", // Gravado - Operación Onerosa
					},
					TaxScheme: &TaxScheme{
						ID:             "1000",
						SchemeID:       "UN/ECE 5153",
						SchemeAgencyID: "6",
						Name:           "IGV",
						TaxTypeCode:    "VAT",
					},
				},
			}},
		}
		ublLine.TaxTotals = append(ublLine.TaxTotals, lineTaxTotal)

		ublCreditNote.CreditNoteLines = append(ublCreditNote.CreditNoteLines, ublLine)
	}

	return ublCreditNote, nil
}

// BuildDebitNote transforms a domain.DebitNote into a UBL DebitNote structure.
func BuildDebitNote(dn *domain.DebitNote) (*DebitNote, error) {
	ublDebitNote := &DebitNote{
		Xmlns:     "urn:oasis:names:specification:ubl:schema:xsd:DebitNote-2",
		XmlnsCAC:  CAC,
		XmlnsCBC:  CBC,
		XmlnsCCTS: CCTS,
		XmlnsDS:   DS,
		XmlnsEXT:  EXT,
		XmlnsQDT:  QDT,
		XmlnsUDT:  UDT,
		XmlnsXSI:  XSI,
		UBLExtensions: &UBLExtensions{
			UBLExtension: &UBLExtension{
				ExtensionContent: &ExtensionContent{XMLName: xml.Name{Local: "ds:Signature"}},
			},
		},
		UBLVersionID:    "2.1",
		CustomizationID: "2.0",
		ID:              fmt.Sprintf("%s-%d", dn.Series, dn.Number),
		IssueDate:       dn.IssueDate.Format("2006-01-02"),
		IssueTime:       dn.IssueDate.Format("15:04:05"),
		DocumentCurrencyCode: &DocumentCurrencyCode{
			ListID:         "ISO 4217 Alpha",
			ListName:       "Currency",
			ListAgencyName: "United Nations Economic Commission for Europe",
			Value:          dn.Currency,
		},
		DiscrepancyResponse: &DiscrepancyResponse{
			ReferenceID:  dn.DiscrepancyResponse.ReferenceID,
			ResponseCode: dn.DiscrepancyResponse.TypeCode,
			Description:  dn.DiscrepancyResponse.Description,
		},
		BillingReference: &BillingReference{
			InvoiceDocumentReference: &InvoiceDocumentReference{ID: dn.DiscrepancyResponse.ReferenceID, DocumentTypeCode: &InvoiceTypeCode{Value: "01"}}, // Assuming 01 for Invoice
		},
		Signature: &Signature{
			ID: "IDSignSP",
			SignatoryParty: &SignatoryParty{
				PartyIdentification: &PartyIdentification{ID: dn.Issuer.RUC},
				PartyName:           &PartyName{Name: dn.Issuer.Name},
			},
			DigitalSignatureAttachment: &DigitalSignatureAttachment{
				ExternalReference: &ExternalReference{URI: "#IDSignSP"},
			},
		},
		AccountingSupplierParty: &Supplier{
			CustomerAssignedAccountID: dn.Issuer.RUC,
			AdditionalAccountID:       "6", // RUC
			Party: &Party{
				PartyName:        &PartyName{Name: dn.Issuer.Name},
				PartyLegalEntity: &PartyLegalEntity{RegistrationName: dn.Issuer.Name},
				PartyTaxScheme: &PartyTaxScheme{
					RegistrationName: dn.Issuer.Name,
					CompanyID: &CompanyID{
						SchemeID:         "6",
						SchemeName:       "SUNAT:Identificador de Documento de Identidad",
						SchemeAgencyName: "PE:SUNAT",
						SchemeURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo06",
						Value:            dn.Issuer.RUC,
					},
					RegistrationAddress: &RegistrationAddress{AddressTypeCode: "0000"},
					TaxScheme:           &TaxScheme{ID: "-"},
				},
			},
		},
		AccountingCustomerParty: &Customer{
			CustomerAssignedAccountID: dn.Recipient.DocNum,
			AdditionalAccountID:       getDocType(dn.Recipient.DocType),
			Party: &Party{
				PartyLegalEntity: &PartyLegalEntity{RegistrationName: dn.Recipient.Name},
				PartyTaxScheme: &PartyTaxScheme{
					RegistrationName: dn.Recipient.Name,
					CompanyID: &CompanyID{
						SchemeID:         getDocType(dn.Recipient.DocType),
						SchemeName:       "SUNAT:Identificador de Documento de Identidad",
						SchemeAgencyName: "PE:SUNAT",
						SchemeURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo06",
						Value:            dn.Recipient.DocNum,
					},
					TaxScheme: &TaxScheme{ID: "-"},
				},
			},
		},
		LegalMonetaryTotal: &MonetaryTotal{
			PayableAmount:       &Amount{CurrencyID: dn.Currency, Value: dn.Totals.Total},
			LineExtensionAmount: &Amount{CurrencyID: dn.Currency, Value: dn.Totals.Gross},
			TaxInclusiveAmount:  &Amount{CurrencyID: dn.Currency, Value: dn.Totals.Total + dn.Totals.IGV},
		},
	}

	// Tax Totals (simplified for now, assuming only IGV)
	igvTaxTotal := &TaxTotal{
		TaxAmount: &Amount{CurrencyID: dn.Currency, Value: dn.Totals.IGV},
		TaxSubtotal: []*TaxSubtotal{{
			TaxableAmount: &Amount{CurrencyID: dn.Currency, Value: dn.Totals.Gross},
			TaxAmount:     &Amount{CurrencyID: dn.Currency, Value: dn.Totals.IGV},
			TaxCategory: &TaxCategory{
				ID:               "S",
				SchemeID:         "UN/ECE 5305",
				SchemeName:       "Tax Category Identifier",
				SchemeAgencyName: "United Nations Economic Commission for Europe",
				Percent:          18.00, // Assuming 18% IGV
				TaxExemptionReasonCode: &TaxExemptionReasonCode{
					ListAgencyName: "PE:SUNAT",
					ListName:       "SUNAT:Codigo de Tipo de Afectación del IGV",
					ListURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo07",
					Value:          "10", // Gravado - Operación Onerosa
				},
				TaxScheme: &TaxScheme{
					ID:             "1000",
					SchemeID:       "UN/ECE 5153",
					SchemeAgencyID: "6",
					Name:           "IGV",
					TaxTypeCode:    "VAT",
				},
			},
		}},
	}
	ublDebitNote.TaxTotals = append(ublDebitNote.TaxTotals, igvTaxTotal)

	// Debit Note Lines
	for i, line := range dn.Lines {
		ublLine := &InvoiceLine{
			ID:                  strconv.Itoa(i + 1),
			InvoicedQuantity:    &Quantity{UnitCode: "NIU", Value: line.Quantity}, // Assuming NIU, should be configurable
			LineExtensionAmount: &Amount{CurrencyID: dn.Currency, Value: line.TotalValue},
			PricingReference: &PricingReference{
				AlternativeConditionPrice: []*Price{{
					PriceAmount: &Amount{CurrencyID: dn.Currency, Value: line.UnitPrice},
					PriceTypeCode: &PriceTypeCode{
						ListAgencyName: "PE:SUNAT",
						ListName:       "SUNAT:Indicador de Tipo de Precio",
						ListURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo16",
						Value:          "01", // Precio unitario (incluye IGV)
					},
				}},
			},
			Item: &Item{
				Description: line.Description,
			},
			Price: &Price{PriceAmount: &Amount{CurrencyID: dn.Currency, Value: line.UnitPrice}},
		}

		// Line Tax Totals (simplified for now, assuming only IGV)
		lineTaxTotal := &TaxTotal{
			TaxAmount: &Amount{CurrencyID: dn.Currency, Value: line.IGV},
			TaxSubtotal: []*TaxSubtotal{{
				TaxableAmount: &Amount{CurrencyID: dn.Currency, Value: line.TotalValue},
				TaxAmount:     &Amount{CurrencyID: dn.Currency, Value: line.IGV},
				TaxCategory: &TaxCategory{
					ID:               "S",
					SchemeID:         "UN/ECE 5305",
					SchemeName:       "Tax Category Identifier",
					SchemeAgencyName: "United Nations Economic Commission for Europe",
					Percent:          18.00,
					TaxExemptionReasonCode: &TaxExemptionReasonCode{
						ListAgencyName: "PE:SUNAT",
						ListName:       "SUNAT:Codigo de Tipo de Afectación del IGV",
						ListURI:        "urn:pe:gob:sunat:cpe:see:gem:catalogos:catalogo07",
						Value:          "10", // Gravado - Operación Onerosa
					},
					TaxScheme: &TaxScheme{
						ID:             "1000",
						SchemeID:       "UN/ECE 5153",
						SchemeAgencyID: "6",
						Name:           "IGV",
						TaxTypeCode:    "VAT",
					},
				},
			}},
		}
		ublDebitNote.TaxTotals = append(ublDebitNote.TaxTotals, lineTaxTotal)

		ublDebitNote.DebitNoteLines = append(ublDebitNote.DebitNoteLines, ublLine)
	}

	return ublDebitNote, nil
}

func getDocType(docType string) string {
	switch docType {
	case "DNI":
		return "1"
	case "RUC":
		return "6"
	case "CE":
		return "4"
	default:
		return "0"
	}
}