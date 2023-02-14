package domainverifier

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/egbakou/domainverifier/dnsresolver"
	"github.com/miekg/dns"
	"net/http"
	"reflect"
	"strings"
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
//	metaTagName: the name of the meta tag to check
//	metaTagContent: the expected value of the meta tag
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

	if err != nil {
		return false, err
	}

	actualValue := reflect.ValueOf(decodedValue).Elem()
	mustMatchValue := reflect.ValueOf(expectedValue)

	return actualValue.Interface() == mustMatchValue.Interface(), nil
}

// CheckTxtRecord checks if the domain has a DNS TXT record
// with the specified values to verify ownership of the domain
//
// Parameters:
//   - dnsResolver: the DNS server to use
//   - domain: the domain name to verify
//   - hostName: the TXT record name
//   - recordContent: the content of the TXT record
//
// Returns:
//   - true if the ownership of the domain is verified
//   - error if any
//
// Example:
//
//	dnsResolver = dnsresolver.CloudflareDNS
//	domain := "website.com"
//	hostName := "@"
//	recordContent := "myapp-site-verification=1234567890"
//	verified, err := domainverify.CheckDnsTxtRecord(dnsServer, domain, hostName, recordContent)
func CheckTxtRecord(dnsResolver, domain, hostName, recordContent string) (bool, error) {
	return checkDNSRecord(dnsResolver, domain, hostName, recordContent, dns.TypeTXT)
}

// CheckCnameRecord checks if the domain has a DNS CNAME record
// with the specified values to verify ownership of the domain
//
// Parameters:
//   - dnsResolver: the DNS server to use
//   - domain: the domain name to verify the ownership
//   - recordName: the name of the DNS CNAME record to check
//   - targetValue: the value of recordName
//
// Returns:
//   - true if the ownership of the domain is verified
//   - error if any
//
// Example:
//
//	dnsResolver = dnsresolver.CloudflareDNS
//	domain := "website.com"
//	recordName := "1234567890"
//	targetValue := "verify.myapp.com"
//	verified, err := domainverify.CheckDnsCnameRecord(dnsResolver, domain, recordName, targetValue)
func CheckCnameRecord(dnsResolver, domain, recordName, targetValue string) (bool, error) {
	return checkDNSRecord(dnsResolver, domain, recordName, targetValue, dns.TypeCNAME)
}

func checkDNSRecord(dnsResolver, domain, recordName, recordContent string, recordType uint16) (bool, error) {
	if strings.TrimSpace(dnsResolver) == "" {
		dnsResolver = dnsresolver.CloudflareDNS
	}

	if !IsValidDomainName(domain) {
		return false, InvalidDomainError
	}

	if recordName != rootDomain && recordName != domain {
		domain = fmt.Sprintf("%s.%s", recordName, domain)
	}

	c := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion(dns.Fqdn(domain), recordType)
	r, _, err := c.Exchange(&m, dnsResolver)
	if err != nil {
		return false, err
	}

	if r.Rcode != dns.RcodeSuccess {
		return false, nil
	}

	for _, a := range r.Answer {
		switch t := a.(type) {
		case *dns.TXT:
			for _, txt := range t.Txt {
				if txt == recordContent {
					return true, nil
				}
			}
		case *dns.CNAME:
			if !strings.HasSuffix(recordContent, ".") {
				recordContent = fmt.Sprintf("%s.", recordContent)
			}
			if t.Target == recordContent {
				return true, nil
			}
		}
	}

	return false, nil
}
