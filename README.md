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

This section assumes prior knowledge of generating instructions for domain name ownership verification. If not, please review the examples code for each verification method, which include the process of instruction generation.

### 🚀 **Html Meta Tag method**

`HTML Meta Tag` is an element that provide metadata about a web page. It is placed in the head section of an HTML document and provides information about the page to search engines and web browsers and others.

This method requires the ability for the user to edit the HTML source code of his site's homepage.

#### 💻 Generation

The generator module contains two functions for generating HTML Meta tags to verify ownership of a specific domain name:

`func GenerateHtmlMetaFromConfig(config *config.HmlMetaTagGenerator, useInternalCode bool) (*HtmlMetaInstruction, error)`

```go
cf := &config.HmlMetaTagGenerator{
		TagName: "example-tag",
		Code:    "external-code", // unique random string, optionnal if useInternalCode is true
	}

// If useInternalCode is set to true, cf.code will be automatically filled with an internal K-Sortable Globally Unique ID
instruction, err := domainverifier.GenerateHtmlMetaFromConfig(cf, false)

if err == nil {
	fmt.Println("Html Code", instruction.Code)
	// Output: <meta name="example-tag" content="external-code" />
	fmt.Println("Indication to provide to the user", instruction.Ation)
	// Output:
	// Copy and paste the <meta> tag into your site's home page.
	// It should go in the <head> section, before the first <body> section.
	// <meta name="example-tag" content="external-code" />
	// * To stay verified, don't remove the meta tag even after verification succeeds.
}
```

<aside> 💡 Ensure that you have stored the Tag Name and code in the database as they will be necessary parameters for the subsequent ownership verification process.

</aside>

`func GenerateHtmlMeta(appName string, sanitizeAppName bool) (*HtmlMetaInstruction, error) {`

This function offers a straightforward approach to generating instructions for the HTML meta tag method.

The `appName` serves as `TagName` **appended by** `-site-verification`**.** If `sanitizeAppName` is set to true, non-alphanumeric characters will be removed from the `appName`.

This function is the simple way to generate instruction for the HTML meta tag method.

```go
instruction, err := domainverifier.GenerateHtmlMeta("your app name", true)

if err == nil {
	fmt.Println("Html Code", instruction.Code)
	// Output: <meta name="yourappname-site-verification" content="random K-Sortable unique code" />
	fmt.Println("Indication to provide to the user", instruction.Ation)
	// Output:
	// Copy and paste the <meta> tag into your site's home page.
	// It should go in the <head> section, before the first <body> section.
	// <meta name="yourappname-site-verification" content="random K-Sortable unique code" />
	// * To stay verified, don't remove the meta tag even after verification succeeds.
}
```

#### 🔎 Verification

The verification process is fast and sample. It requires:

- The domain name for which you’ve generated the the verification instruction
- The Html Meta `Tag Name` and it `value`(Code) you have stored somewhere

```go
isVerified, err := domainverifier.CheckHtmlMetaTag("the-domain-to-verify.com",
		"tag name",
		"verification-code")

if err != nil {
	fmt.Print(err)
}

fmt.Println("Is onwershsip verified:", isVerified)
```
