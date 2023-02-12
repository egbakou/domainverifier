package domainverifier

import (
	"errors"
)

const (
	jsonKeySuffix            = "_site_verification"
	jsonFileNameSuffix       = "-site-verification.json"
	xmlRootName              = "verification"
	xmlFileNameSuffix        = "SiteAuth.xml"
	txtRecordAttributeSuffix = "-site-verification"
)

// InvalidAppNameError indicates that the app name is invalid.
var InvalidAppNameError = errors.New("app name cannot be empty")

// HtmlMetaInstruction is the verification method that uses HTML meta tag.
type HtmlMetaInstruction struct {
	Code   string
	Action string
}

// FileInstruction is the verification method that uses JSON or XML file.
type FileInstruction struct {
	FileName    string
	FileContent string
	Action      string
}

// DnsRecordInstruction is the verification method that uses TXT or CNAME record.
type DnsRecordInstruction struct {
	HostName string
	Record   string
	Action   string
}
