
---

## Tecnologías utilizadas

- **Go (Golang)** - Backend principal
- **UBL 2.1** - Estandarización XML para comprobantes electrónicos
- **Certificados X.509** - Firma digital de documentos
- **HTTP API** - Exposición de servicios REST

---

## Requisitos previos

- Go 1.20 o superior
- Certificado digital en formato `.pem` o `.pfx`
- [Postman](https://www.postman.com/) o similar para pruebas
- Conexión a SUNAT (si se desea integrar en producción)

---

##  Cómo ejecutar

1. Clonar el repositorio:

```bash
git clone https://github.com/FrnklynUser/FacturacionElectronica.git
cd FacturacionElectronica
