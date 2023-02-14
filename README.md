# domainverifier

[![Go Reference](<https://pkg.go.dev/badge/github.com/egbakou/domainverifier.svg>)](<https://pkg.go.dev/github.com/egbakou/domainverifier>) [![CI](<https://github.com/egbakou/domainverifier/actions/workflows/ci.yml/badge.svg>)](<https://github.com/egbakou/domainverifier/actions/workflows/ci.yml>)

`domainverifier` is a Go package that provides a simple and easy way to verify domain name ownership. It also includes a generator module, which makes it easier for developers who are new to DNS verification to quickly set up and integrate the verification process into their applications.

The package offers support for 5 different verification methods: `HTML Meta Tag`, `JSON File Upload`, `XML File Upload`, `DNS TXT record` and `DNS CNAME record`.

## Installation

To get the package use the standard:

```erlang
go get -u github.com/egbakou/domainverifier
```

Using Go modules is recommended.

## Usage

This section assumes prior knowledge of generating instructions for domain name ownership verification. If not, please review the examples code for each verification method, which include the process of instructions generation.

### 🚀 Html Meta Tag method

`HTML Meta Tag` is an element that provide metadata about a web page. It is placed in the head section of an HTML document and provides information about the page to search engines and web browsers and others.

This method requires the ability for the user to edit the HTML source code of his site's homepage.

💻 Generation

The generator module contains two functions for generating HTML Meta tags to verify ownership of a specific domain name:

⤵️`func GenerateHtmlMetaFromConfig(config *config.HmlMetaTagGenerator, useInternalCode bool) (*HtmlMetaInstruction, error)`

```go
config := &config.HmlMetaTagGenerator{
	TagName: "example-tag",
	Code:    "external-code", // a unique random string, optional if useInternalCode is true
}

// If useInternalCode is set to true, config.Code will be automatically filled with an internal K-Sortable Globally Unique ID
instruction, err := domainverifier.GenerateHtmlMetaFromConfig(cf, false)

if err == nil {
	fmt.Println("Html Code", instruction.Code)
	// Output: <meta name="example-tag" content="external-code" />
	fmt.Println("Indication to provide to the user", instruction.Action)
	// Output:
	// Copy and paste the <meta> tag into your site's home page.
	// It should go in the <head> section, before the first <body> section.
	// <meta name="example-tag" content="external-code" />
	// * To stay verified, don't remove the meta tag even after verification succeeds.
}
```

> 💡 Ensure that you have stored the Tag Name and code in the database as they will be necessary parameters for the subsequent ownership verification process.
>

⤵️`func GenerateHtmlMeta(appName string, sanitizeAppName bool) (*HtmlMetaInstruction, error)`

This function offers a straightforward approach to generating instructions for the HTML meta tag method.

The `appName` serves as `TagName` appended by `-site-verification`. If `sanitizeAppName` is set to true, non-alphanumeric characters will be removed from the `appName`.

This function is the simple way to generate instructions for the HTML meta tag method.

```go
instruction, err := domainverifier.GenerateHtmlMeta("your app name", true)

if err == nil {
	fmt.Println("Html Code:", instruction.Code)
	// Output: 
	// <meta name="yourappname-site-verification" content="random K-Sortable unique code" />
	fmt.Println("Indication to provide to the user:", instruction.Action)
	// Output:
	// Copy and paste the <meta> tag into your site's home page.
	// It should go in the <head> section, before the first <body> section.
	// <meta name="yourappname-site-verification" content="random K-Sortable unique code" />
	// * To stay verified, don't remove the meta tag even after verification succeeds.
}
```

🔎 Verification

The verification process is fast and sample. It requires:

- The domain name for which you’ve generated the verification instructions
- The Html Meta `Tag Name` and it `value`(Code) you have stored somewhere

```go
isVerified, err := domainverifier.CheckHtmlMetaTag("the-domain-to-verify.com",
		"tag name",
		"verification-code")

fmt.Println("Is onwershsip verified:", isVerified)
```

### 🚀 JSON file upload method

In the JSON method,  you need to create a JSON file that contains a specific structure, including a key-value pair that proves ownership of the domain. User then upload the JSON file to his website's root directory.

Once the JSON file is uploaded, the ownership verification service can access the file and verify the contents to confirm that the user indeed own the domain.

💻 Generation

⤵️ `GenerateJsonFromConfig(config *config.JsonGenerator, useInternalCode bool) (*FileInstruction, error)`

```go
config := &config.JsonGenerator{
	FileName:  "example.json",
	Attribute: "code",
	Code:      "external-code", // optional if useInternalCode is true
},

// If useInternalCode is set to true, config.Code will be automatically filled with an internal K-Sortable Globally Unique ID
instruction, err := domainverifier.GenerateJsonFromConfig(config, false)

if err == nil {
	fmt.Println("FileName :", instruction.FileName)
	// Output: 
	// example.json
	fmt.Println("FileContent:", instruction.FileContent)
	// Output: 
	// {"code": "external-code"}
	fmt.Println("Indication to provide to the user", instruction.Action)
	// Output:
	// Create a JSON file named example.json with the content
	// {"code": "external-code"}
	// and upload it to the root of your site.
}
```

> 💡 It is important to store `FileName`, `Attribute`, and `Code` in the database, as this data will be essential for verifying ownership later.

⤵️ `func GenerateJson(appName string) (*FileInstruction, error)`

The `appName` serves as `Attribute` appended by `_site_verification`.

```go
instruction, err := domainverifier.GenerateJson("your app name")

if err == nil {
	fmt.Println("FileName :", instruction.FileName)
	// Output: 
	// yourappname-site_verification.json
	fmt.Println("FileContent:", instruction.FileContent)
	// Output: 
	// {"yourappname_site_verification": "random K-Sortable unique code"}
	fmt.Println("Indication to provide to the user", instruction.Action)
	// Output:
	// Create a JSON file named yourappname-site_verification.json with the content
	// {"yourappname_site_verification": "random K-Sortable unique code"}
	// and upload it to the root of your site.
}
```

🔎 Verification

Requirements:

- The domain name for which you have generated the verification instructions
- The JSON file name
- An object of type `struct` that matches the content of the JSON file.

```go
type ownershipVerification struct {
	Code string json:"code"
}

expectedValue := ownershipVerification{Code: "verification-code"}

isVerified, err := domainverifier.CheckJsonFile("the-domain-to-verify.com",
		"example.json",
		expectedValue)

fmt.Println("Is onwershsip verified:", isVerified)
```

### 🚀 XML  file upload method

This approach is similar to the JSON method. There are two functions you can use to provide verification instructions to users.

⤵️ `func GenerateXmlFromConfig(config *config.XmlGenerator, useInternalCode bool) (*FileInstruction, error)`

```go
config: &config.XmlGenerator{
	FileName: "example.xml",
	RootName: "example-root",
	Code:     "internal-code",
}

// If useInternalCode is set to true, config.Code will be automatically filled with an internal K-Sortable Globally Unique ID
instruction, err := GenerateXmlFromConfig(config, false)

if err == nil {
	fmt.Println("FileName :", instruction.FileName)
	// Output: 
	// example.xml
	fmt.Println("FileContent:", instruction.FileContent)
	// Output: 
	// <example-root><code>internal-code</code></example-root>
	fmt.Println("Indication to provide to the user", instruction.Action)
	// Output:
	// Create an XML file named example.xml with the content:
	// <example-root><code>internal-code</code></example-root>
	// and upload it to the root of your site.
}
```

> 💡 `FileName`, `RootName`, and `Code` must be stored in database.

⤵️ `func GenerateXml(appName string, sanitizeAppName bool) (*FileInstruction, error)`

The `appName` serves as `FileName` appended by `SiteAuth.xml`.

```go
// We advise setting the sanitizeAppName parameter to true
instruction, err := GenerateXml("your app name", true)

if err == nil {
	fmt.Println("FileName :", instruction.FileName)
	// Output: 
	// YourappnameSiteauth.xml
	fmt.Println("FileContent:", instruction.FileContent)
	// Output: 
	// <verification><code>random K-Sortable unique code</code></verification>
	fmt.Println("Indication to provide to the user", instruction.Action)
	// Output:
	// Create an XML file named example.xml with the content:
	// <verification><code>random K-Sortable unique code</code></verificationt>
	// and upload it to the root of your site.
}
```

🔎 Verification

Requirements:

- The domain name for which you have generated the verification instructions
- The XML file name
- An object of type `struct` that matches the content of the XML file.

```go
type ownershipVerification struct {
    XMLName struct{} `xml:"verification"`
	Code string `xml:"code"`
}

expectedValue := ownershipVerification{Code: "verification-code"}

isVerified, err := domainverifier.CheckXmlFile("the-domain-to-verify.com",
		"example.xml",
		expectedValue)

fmt.Println("Is onwershsip verified:", isVerified)
```

### 🚀 DNS TXT record method

With this method, user needs to add a specific TXT record to the DNS configuration of your domain. The TXT record contains a unique value that proves ownership of the domain.

💻 Generation

⤵️ `func GenerateTxtRecordFromConfig(config *config.TxtRecordGenerator, useInternalCode bool) (*DnsRecordInstruction, error)`

This function allows you to generate instructions for the DNS TXT record method.

```go
config: &config.TxtRecordGenerator{
	HostName:             "@", // or domain.com for example
	RecordAttribute:      "myapp",
	RecordAttributeValue: "random-code",
},

// If useInternalCode is set to true, cf.code will be automatically filled with an internal K-Sortable Globally Unique ID
instruction, err := domainverifier.GenerateTxtFromConfig(config, false)

if err == nil {
	fmt.Println("HostName:", instruction.HostName)
	// Output:
	// @ or domain.com
	fmt.Println("Record:", instruction.Record)
	// Output:
	// myapp=random-code
	fmt.Println("Indications to provide to the user", instruction.Action)
	// Output:
	// Create a TXT record with the name @ and the content myapp=random-code
}
```

> 💡 Remember to store the config object somewhere.

⤵️ `func GenerateTxtRecord(appName string) (*DnsRecordInstruction, error)`

The `appName` serves as `RecordAttribute` appended by  `-site-verification`

```go
instruction, err := domainverifier.GenerateTxtRecord("your app name")

if err == nil {
	fmt.Println("HostName:", instruction.HostName)
	// Output:
	// @
	fmt.Println("Record:", instruction.Record)
	// Output:
	// yourappname-site-verification=random K-Sortable unique code
	fmt.Println("Indications to provide to the user", instruction.Action)
	// Output:
	// Create a TXT record with the name @ and the content yourappname-site-verification=random K-Sortable unique code
}
```

🔎 Verification

Requirements:

- The DNS server to use, if empty `Cloudflare DNS` will be used.
- The domain name for which you have generated the verification instructions
- The Host name
- The record content

```go
isVerified, err := domainverifier.CheckTxtRecord(dnsresolver.GooglePublicDNS, "the-domain-to-verify.com", "@", "yapp=random-code")

fmt.Println("Is onwershsip verified:", isVerified)
```

### 🚀 DNS CNAME record method

💻 Generation

⤵️ `func GenerateCnameRecordFromConfig(config *config.CnameRecordGenerator) (*DnsRecordInstruction, error)`

```go
config: &config.CnameRecordGenerator{
	RecordName:   "random-code",
	RecordTarget:  "verify.example.com",
},

instruction, err := domainverifier.GenerateCnameFromConfig(config)

if err == nil {
	fmt.Println("HostName:", instruction.HostName)
	// Output:
	// random-code
	fmt.Println("Record:", instruction.Record)
	// Output:
	// verify.example.com
	fmt.Println("Indications to provide to the user", instruction.Action)
	// Output:
	// Add CNAME (alias) record with name random-code and value verify.example.com.
}
```

> 💡 Ensure to store the DNS CNAME record information generated, including the Record Name and Record Target.

🔍 Verification

Requirements:

- The DNS server to use, if empty `Cloudflare DNS` will be used.
- The domain name for which you have generated the verification instructions
- The record name
- The target value

```go
isVerified, err := domainverifier.CheckCnameRecord(dnsresolver.GooglePublicDNS, "the-domain-to-verify.com", "random-code", "verify.example.com")

fmt.Println("Is onwershsip verified:", isVerified)
```

## Utility functions

In addition to its main features, `domainverifier` provides some helper functions that can be used.

- `domainverifier.IsSecure(domain string, timeout time.Duration)` returns a Boolean value indicating whether the specified domain supports a secure connection over HTTPS or not.
- `domainverifier.IsValidDomainName(domain string)` checks if a string is a valid domain name.

## Contributions

We're always looking for contributions to make this project even better! If you're interested in helping out, please take a look at our open issues, or create a new one if you have an idea for a feature or bug fix. We appreciate any and all help, so don't hesitate to reach out if you want to get involved!
