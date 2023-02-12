package config

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
)

var InvalidConfigError = errors.New("config cannot be nil")

// HmlMetaTagGenerator is the required config to generate HTML Meta verification method instructions.
type HmlMetaTagGenerator struct {
	TagName string
	Code    string
}

func (h *HmlMetaTagGenerator) Validate() error {
	if h == nil {
		return InvalidConfigError
	}

	if strings.TrimSpace(h.TagName) == "" {
		return fmt.Errorf("tag name cannot be empty")
	}

	if strings.TrimSpace(h.Code) == "" {
		return fmt.Errorf("code cannot be empty")
	}
	return nil
}

// JsonGenerator is the required config to generate JSON verification method instructions.
type JsonGenerator struct {
	FileName  string
	Attribute string
	Code      string
}

func (j *JsonGenerator) Validate() error {
	if j == nil {
		return InvalidConfigError
	}

	if strings.TrimSpace(j.FileName) == "" {
		return fmt.Errorf("file name cannot be empty")
	}

	if strings.TrimSpace(j.Attribute) == "" {
		return fmt.Errorf("attribute cannot be empty")
	}

	if strings.TrimSpace(j.Code) == "" {
		return fmt.Errorf("code cannot be empty")
	}
	return nil
}

// XmlGenerator is the required config to generate XML verification method instructions.
type XmlGenerator struct {
	FileName string
	RootName string
	Code     string
}

func (x *XmlGenerator) Validate() error {
	if x == nil {
		return InvalidConfigError
	}

	if strings.TrimSpace(x.FileName) == "" {
		return fmt.Errorf("file name cannot be empty")
	}

	if strings.TrimSpace(x.RootName) == "" {
		return fmt.Errorf("root name cannot be empty")
	}

	if strings.TrimSpace(x.Code) == "" {
		return fmt.Errorf("code cannot be empty")
	}
	return nil
}

func (x *XmlGenerator) ToXml() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`%s`, xml.Header))
	sb.WriteString(fmt.Sprintf(`<%s><code>%s</code></%s>`, x.RootName, x.Code, x.RootName))
	return sb.String()
}

// TxtRecordGenerator is the required config to generate TXT record verification method instructions.
type TxtRecordGenerator struct {
	HostName             string // @ or the domain name to verify or unique generated code.
	RecordAttribute      string
	RecordAttributeValue string
}

func (t *TxtRecordGenerator) Validate() error {
	if t == nil {
		return InvalidConfigError
	}

	if strings.TrimSpace(t.HostName) == "" {
		return fmt.Errorf("host name cannot be empty")
	}

	if strings.TrimSpace(t.RecordAttribute) == "" {
		return fmt.Errorf("record attribute cannot be empty")
	}

	if strings.TrimSpace(t.RecordAttributeValue) == "" {
		return fmt.Errorf("record attribute value cannot be empty")
	}
	return nil
}

// CnameRecordGenerator is the required config to generate CNAME record verification method instructions.
type CnameRecordGenerator struct {
	RecordName   string
	RecordTarget string
}

func (c *CnameRecordGenerator) Validate() error {
	if c == nil {
		return InvalidConfigError
	}

	if strings.TrimSpace(c.RecordName) == "" {
		return fmt.Errorf("record name cannot be empty")
	}

	if strings.TrimSpace(c.RecordTarget) == "" {
		return fmt.Errorf("record target cannot be empty")
	}
	return nil
}
