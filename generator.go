package domainverifier

import (
	"errors"
	"fmt"
	"github.com/egbakou/domainverifier/config"
	"github.com/segmentio/ksuid"
	"strings"
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

// GenerateHtmlMetaFromConfig generates the HTML meta tag verification method instructions.
// It uses the provided config.HmlMetaTagGenerator to generate the instructions.
// If useInternalCode is true, internal K-Sortable Globally Unique ID will be generated.
// Otherwise, the code in the config.HmlMetaTagGenerator will be used.
func GenerateHtmlMetaFromConfig(config *config.HmlMetaTagGenerator, useInternalCode bool) (*HtmlMetaInstruction, error) {
	if config != nil && useInternalCode {
		config.Code = ksuid.New().String()
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &HtmlMetaInstruction{
		Code:   getMetaTagContent(config.TagName, config.Code),
		Action: getMetaTagInstruction(config.TagName, config.Code),
	}, nil
}

// GenerateHtmlMeta generates the HTML meta tag verification method instructions.
// appName is the name of the app that is requesting the verification (e.g. msvalidate.01, mysuperapp, etc.).
// It will be used as the name of the meta tag.
// Note that the appName will be sanitized to non-alphanumeric characters.
func GenerateHtmlMeta(appName string, sanitizeAppName bool) (*HtmlMetaInstruction, error) {
	if strings.TrimSpace(appName) == "" {
		return nil, InvalidAppNameError
	}
	if sanitizeAppName {
		appName = sanitizeString(appName)
	}
	htmTagConfig := &config.HmlMetaTagGenerator{
		TagName: appName,
		Code:    ksuid.New().String(),
	}
	return GenerateHtmlMetaFromConfig(htmTagConfig, false)
}

// GenerateJsonFromConfig generates the JSON verification method instructions.
// It uses the provided config.JsonGenerator to generate the instructions.
// If useInternalCode is true, internal K-Sortable Globally Unique ID will be generated.
// Otherwise, the code in the config.JsonGenerator will be used.
func GenerateJsonFromConfig(config *config.JsonGenerator, useInternalCode bool) (*FileInstruction, error) {
	if config != nil && useInternalCode {
		config.Code = ksuid.New().String()
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	fileName := ensureFileExtension(config.FileName, ".json")
	return &FileInstruction{
		FileName:    config.FileName,
		FileContent: getJsonContent(config.Attribute, config.Code),
		Action:      getJsonInstruction(fileName, config.Attribute, config.Code),
	}, nil
}

// GenerateJson generates the JSON verification method instructions.
// appName is the name of the app that is requesting the verification (e.g. google, bing, etc.).
// It will be used as prefix of the file name and the attribute name.
// Note that the appName will be sanitized to non-alphanumeric characters.
func GenerateJson(appName string) (*FileInstruction, error) {
	if strings.TrimSpace(appName) == "" {
		return nil, InvalidAppNameError
	}

	appName = sanitizeString(appName)

	jsonConfig := &config.JsonGenerator{
		FileName:  fmt.Sprintf("%s%s", appName, jsonFileNameSuffix),
		Attribute: fmt.Sprintf("%s%s", appName, jsonKeySuffix),
		Code:      ksuid.New().String(),
	}
	return GenerateJsonFromConfig(jsonConfig, false)
}
