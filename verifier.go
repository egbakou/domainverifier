package domainverifier

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
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
