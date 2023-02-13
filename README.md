# domainverifier

[![Go Reference](https://pkg.go.dev/badge/github.com/egbakou/domainverifier.svg)](https://pkg.go.dev/github.com/egbakou/domainverifier) [![CI](https://github.com/egbakou/domainverifier/actions/workflows/ci.yml/badge.svg)](https://github.com/egbakou/domainverifier/actions/workflows/ci.yml)

domainverifier is a Go package that provides a simple and easy way to verify domain name ownership. It also includes a generator module, which makes it easier for developers who are new to DNS verification to quickly set up and integrate the verification process into their applications.

## Supported verification methods

| Method             | Verification supported ? | Generator supported  ? |
| :----------------- | :----------------------: | :--------------------: |
| `HTML Meta Tag`    |            ✅             |           ✅            |
| `JSON File Upload` |            ✅             |           ✅            |
| `XML File Upload`  |            ✅             |           ✅            |
| `DNS TXT record`   |            ✅             |           ✅            |
| `DNS CNAME record` |            ✅             |           ✅            |
| `HTML File Upload` |            ❌             |           ❌            |

## Installation

To get the package use the standard:

```bash
go get -u github.com/egbakou/domainverifier
```

Using Go modules is recommended.
