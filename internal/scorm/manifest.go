package scorm

import (
	"encoding/xml"
	"os"
	"path/filepath"

	"github.com/jeffory/scorm-generator-poc/pkg/models"
)

func (g *Generator) generateManifest(dir string, quiz *models.Quiz) error {
	manifest := Manifest{
		Identifier: "SCORM_QUIZ_" + sanitizeFilename(quiz.Subject),
		Version:    "1.0",
		Metadata: Metadata{
			Schema:         "ADL SCORM",
			SchemaVersion:  "1.2",
		},
		Organizations: Organizations{
			Default: "ORG-001",
			Organization: Organization{
				Identifier: "ORG-001",
				Title:      quiz.Subject + " Quiz",
				Item: Item{
					Identifier:     "ITEM-001",
					IdentifierRef:  "RESOURCE-001",
					Title:          quiz.Subject + " Quiz",
				},
			},
		},
		Resources: Resources{
			Resource: Resource{
				Identifier: "RESOURCE-001",
				Type:       "webcontent",
				Href:       "index.html",
				AdlcpScormType: "sco",
				Files: []File{
					{Href: "index.html"},
					{Href: "scormapi.js"},
					{Href: "styles.css"},
				},
			},
		},
	}

	xmlData, err := xml.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}

	xmlString := xml.Header + string(xmlData)
	return os.WriteFile(filepath.Join(dir, "imsmanifest.xml"), []byte(xmlString), 0644)
}

// Manifest XML structure for SCORM 1.2
type Manifest struct {
	XMLName       xml.Name      `xml:"manifest"`
	Identifier    string        `xml:"identifier,attr"`
	Version       string        `xml:"version,attr"`
	Xmlns         string        `xml:"xmlns,attr"`
	XmlnsAdlcp    string        `xml:"xmlns:adlcp,attr"`
	XmlnsXsi      string        `xml:"xmlns:xsi,attr"`
	XsiSchemaLocation string    `xml:"xsi:schemaLocation,attr"`
	Metadata      Metadata      `xml:"metadata"`
	Organizations Organizations `xml:"organizations"`
	Resources     Resources     `xml:"resources"`
}

type Metadata struct {
	Schema        string `xml:"schema"`
	SchemaVersion string `xml:"schemaversion"`
}

type Organizations struct {
	Default      string       `xml:"default,attr"`
	Organization Organization `xml:"organization"`
}

type Organization struct {
	Identifier string `xml:"identifier,attr"`
	Title      string `xml:"title"`
	Item       Item   `xml:"item"`
}

type Item struct {
	Identifier    string `xml:"identifier,attr"`
	IdentifierRef string `xml:"identifierref,attr"`
	Title         string `xml:"title"`
}

type Resources struct {
	Resource Resource `xml:"resource"`
}

type Resource struct {
	Identifier     string `xml:"identifier,attr"`
	Type           string `xml:"type,attr"`
	Href           string `xml:"href,attr"`
	AdlcpScormType string `xml:"adlcp:scormtype,attr"`
	Files          []File `xml:"file"`
}

type File struct {
	Href string `xml:"href,attr"`
}

func init() {
	// Set default namespace values
}

func (m *Manifest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type Alias Manifest
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	// Set namespaces
	start.Attr = []xml.Attr{
		{Name: xml.Name{Local: "identifier"}, Value: m.Identifier},
		{Name: xml.Name{Local: "version"}, Value: m.Version},
		{Name: xml.Name{Local: "xmlns"}, Value: "http://www.imsproject.org/xsd/imscp_rootv1p1p2"},
		{Name: xml.Name{Local: "xmlns:adlcp"}, Value: "http://www.adlnet.org/xsd/adlcp_rootv1p2"},
		{Name: xml.Name{Local: "xmlns:xsi"}, Value: "http://www.w3.org/2001/XMLSchema-instance"},
		{Name: xml.Name{Local: "xsi:schemaLocation"}, Value: "http://www.imsproject.org/xsd/imscp_rootv1p1p2 imscp_rootv1p1p2.xsd http://www.imsglobal.org/xsd/imsmd_rootv1p2p1 imsmd_rootv1p2p1.xsd http://www.adlnet.org/xsd/adlcp_rootv1p2 adlcp_rootv1p2.xsd"},
	}

	return e.EncodeElement(aux, start)
}
