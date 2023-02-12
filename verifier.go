package domainverifier

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net"
	"net/http"
	"reflect"
)

const rootDomain = "@"

// InvalidDomainError indicates that the domain name is invalid
var InvalidDomainError = errors.New("invalid domain name")

// InvalidResponseError indicates that the response is invalid
var InvalidResponseError = errors.New("invalid response status code returned by the server")

// CheckHtmlMetaTag checks if the html meta tag exists and has the expected value
//
// Parameters:
//
//	domain: the domain name to check
//	tagName: the name of the meta tag to check
//	tagContent: the expected value of the meta tag
//
// Returns:
//   - true if the ownership of the domain is verified
//   - error if any
func CheckHtmlMetaTag(domain, metaTagName, metaTagContent string) (bool, error) {
	if !IsValidDomainName(domain) {
		return false, InvalidDomainError
	}
	resp, err := makeHttpCall(domain)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, InvalidResponseError
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return false, err
	}

	// search for specified HTML meta tag and value
	metaTag := doc.Find(fmt.Sprintf("meta[name=%s]", metaTagName))
	if metaTag.Length() == 0 {
		return false, nil
	}
	content, exists := metaTag.Attr("content")
	if !exists {
		return false, nil
	}

	return content == metaTagContent, nil
}

// CheckJsonFile checks if the json file exists and has
// the expected content to verify ownership of the domain
//
// Parameters:
//   - domain: the domain name to check
//   - fileName: the name of the json file to check
//   - expectedValue: the expected content
//
// Returns:
//   - true if the ownership of the domain is verified
//   - error if any
//
// Example:
//
//	type OwnershipVerification struct {
//		Code string `json:"myapp_site_verification"`
//	}
//	data := OwnershipVerification{Code: "1234567890"}
//	domain := "website.com"
//	fileName := "myapp-site-verification.json" // excepted file content: {"myapp_site_verification": "1234567890"}
//	verified, err := domainverify.CheckJsonFile(domain, fileName, data)
func CheckJsonFile(domain, fileName string, expectedValue interface{}) (bool, error) {
	return checkXmlOrJsonFile(false, domain, fileName, expectedValue)
}

// CheckXmlFile checks if the xml file exists and has
// the expected content to verify ownership of the domain
//
// Parameters:
//   - domain: the domain name to check
//   - fileName: the name of the xml file to check
//   - expectedValue: the expected content
//
// Returns:
//   - true if the ownership of the domain is verified
//   - error if any
//
// Example:
//
//	type OwnershipVerification struct {
//	    XMLName struct{} `xml:"verification"`
//		Code string `xml:"code" json:"myapp_site_verification"`
//	}
//	data := OwnershipVerification{Code: "1234567890"}
//	domain := "website.com"
//	fileName := "myappSiteAuth.xml" // excepted file content: <verification><code>1234567890</code></verification>
//	verified, err := domainverify.CheckXmlFile(domain, fileName, data)
func CheckXmlFile(domain, fileName string, expectedValue interface{}) (bool, error) {
	return checkXmlOrJsonFile(true, domain, fileName, expectedValue)
}

// checkXmlOrJsonFile checks domain name ownership using Xml or Json method
func checkXmlOrJsonFile(useXmlMethod bool, domain, fileName string, expectedValue interface{}) (bool, error) {
	if !IsValidDomainName(domain) {
		return false, InvalidDomainError
	}

	// Only struct type is supported
	if reflect.TypeOf(expectedValue).Kind() != reflect.Struct {
		return false, errors.New("expectedValue must be a struct")
	}

	resp, err := makeHttpCall(fmt.Sprintf("%s/%s", domain, fileName))
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return false, InvalidResponseError
	}

	// Decode the XML response from the URL
	decodedValue := reflect.New(reflect.TypeOf(expectedValue)).Interface()
	if useXmlMethod {
		err = xml.NewDecoder(resp.Body).Decode(decodedValue)
	} else {
		err = json.NewDecoder(resp.Body).Decode(decodedValue)
	}

	actualValue := reflect.ValueOf(decodedValue).Elem()
	mustMatchValue := reflect.ValueOf(expectedValue)

	return actualValue.Interface() == mustMatchValue.Interface(), nil
}

// CheckTxtRecord checks if the domain has a DNS TXT record
// with the specified value to verify ownership of the domain
//
// Parameters:
//   - domain: the domain name to check
//   - recordContent: the name of the DNS TXT record to check
//
// Returns:
//   - true if the ownership of the domain is verified
//   - error if any
//
// Example:
//
//	domain := "website.com"
//	recordName := "myapp-site-verification=1234567890"
//	verified, err := domainverify.CheckDnsTxtRecord(domain, recordName)
func CheckTxtRecord(domain, hostName, recordContent string) (bool, error) {
	if !IsValidDomainName(domain) {
		return false, InvalidDomainError
	}

	if hostName != rootDomain && hostName != domain {
		domain = fmt.Sprintf("%s.%s", hostName, domain)
	}

	records, err := net.LookupTXT(domain)
	if err != nil {
		return false, err
	}

	// Check if the TXT record exists
	for _, record := range records {
		if record == recordContent {
			return true, nil
		}
	}

	return false, nil
}
