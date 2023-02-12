package domainverifier

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode"
)

const (
	domainNamePattern = `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`
	httpPrefix        = "http://"
	httpsPrefix       = "https://"
)

func getMetaTagInstruction(name, content string) string {
	var sb strings.Builder
	sb.WriteString("Copy and paste the <meta> tag into your site's home page.\n")
	sb.WriteString("It should go in the <head> section, before the first <body> section.\n")
	sb.WriteString(getMetaTagContent(name, content))
	sb.WriteString("\n* To stay verified, don't remove the meta tag even after verification succeeds.")
	return sb.String()
}

func getMetaTagContent(name, content string) string {
	return fmt.Sprintf(`<meta name="%s" content="%s" />`, name, content)
}

func getJsonInstruction(filename, key, value string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`Create a JSON file named %s with the content`, filename))
	sb.WriteString("\n")
	sb.WriteString(getJsonContent(key, value))
	sb.WriteString("\nand upload it to the root of your site.")
	return sb.String()
}

func getJsonContent(key, value string) string {
	return fmt.Sprintf(`{"%s": "%s"}`, key, value)
}

func getXmlInstruction(xmlFileName, rootName, code string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`Create an XML file named %s with the content:`, xmlFileName))
	sb.WriteString("\n")
	sb.WriteString(getXmlContent(rootName, code))
	sb.WriteString("\nand upload it to the root of your site.")
	return sb.String()
}

func getXmlContent(rootName, code string) string {
	var sb strings.Builder
	sb.WriteString(xml.Header)
	sb.WriteString(fmt.Sprintf(`<%s><code>%s</code></%s>`, rootName, code, rootName))
	return sb.String()
}

// IsValidDomainName checks if a string is a valid domain name.
func IsValidDomainName(domain string) bool {
	if len(strings.TrimSpace(domain)) < 1 {
		return false
	}
	regex := regexp.MustCompile(domainNamePattern)
	return regex.MatchString(domain)
}

// IsSecure returns a boolean value indicating whether the specified domain supports a secure connection over HTTPS or not.
// If the domain is not reachable, an error is returned.
func IsSecure(domain string, timeout time.Duration) (bool, error) {
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(fmt.Sprintf("%s%s", httpsPrefix, domain))
	if err != nil {
		resp, err = client.Get(fmt.Sprintf("%s%s", httpPrefix, domain))
		if err != nil {
			return false, err
		}
	}
	defer resp.Body.Close()

	return resp.TLS != nil, nil
}

// sanitizeString removes all non-alphanumeric characters from a string,
// removes all spaces and converts the string to lowercase.
//
// Parameters:
//   - str: the string to sanitize
//
// Example:
//
//	input := "My Super App"
//	fmt.Println(sanitizeString(input))
//	fmt.Println(input)
//	// Output:
//	// mysuperapp
func sanitizeString(str string) string {
	// Trim the string
	str = strings.TrimSpace(str)

	// Create a mapping function that removes non-alphabetic characters and spaces, and converts characters to lowercase
	mapping := func(r rune) rune {
		if (!unicode.IsLetter(r) && !unicode.IsNumber(r)) || r == ' ' {
			return -1
		}
		return unicode.ToLower(r)
	}

	// Apply the mapping function to each character in the string
	return strings.Map(mapping, str)
}

// ensureFileExtension ensures that a file name has a specific extension.
// If the file name already has the extension, it is returned as is.
// Otherwise, the extension is appended to the file name.
func ensureFileExtension(filename, extension string) string {
	if !strings.HasSuffix(filename, extension) {
		return filename + extension
	}
	return filename
}

// makeHttpCall makes an HTTP call to the specified URL.
func makeHttpCall(url string) (*http.Response, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s", httpsPrefix, url))
	if err != nil {
		resp, err = http.Get(fmt.Sprintf("%s%s", httpPrefix, url))
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}
