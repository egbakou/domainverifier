//go:build exclude
// +build exclude

package main

import (
	"fmt"
	"github.com/egbakou/domainverifier"
)

type OwnershipVerification struct {
	XMLName struct{} `xml:"verification"` // fix: https://stackoverflow.com/questions/12398925/go-xml-marshalling-and-the-root-element
	Code    string   `xml:"code" json:"myapp_site_verification"`
}

func main() {
	isVerified, err := domainverifier.CheckHtmlMetaTag("the-domain-to-verify.com",
		"myapp-site-verification",
		"verification-code")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("Html Meta Tag is verified:", isVerified)

	data := OwnershipVerification{
		Code: "dcf56hgvghy674fc",
	}
	isVerified, err = domainverifier.CheckJsonFile("domainverify.lioncoding.workers.dev",
		"myapp-site-verification.json",
		data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Json File is verified:", isVerified)

	isVerified, err = domainverifier.CheckXmlFile("domainverify.lioncoding.workers.dev",
		"myappSiteAuth.xml",
		data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Xml File is verified:", isVerified)

	isVerified, err = domainverifier.CheckTxtRecord("lioncoding.com",
		"@",
		"ownership-demo-app=random000454")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Dns Txt Record is verified:", isVerified)

	isVerified, err = domainverifier.CheckCnameRecord("lioncoding.com",
		"random000454",
		"ownership-demo-app.com")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("Dns Cname Record is verified:", isVerified)
}
