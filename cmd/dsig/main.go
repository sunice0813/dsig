package main

import (
	"crypto/x509"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/beevik/etree"
	x_dsig "github.com/russellhaering/goxmldsig"
	"github.com/ucarion/dsig"
)

type assertion struct {
	XMLName   xml.Name
	Attrs     []xml.Attr     `xml:",any,attr"`
	Signature dsig.Signature `xml:"Signature"`
	InnerXML  []byte         `xml:",innerxml"`
}

func main() {
	block, _ := pem.Decode([]byte(`-----BEGIN CERTIFICATE-----
MIIDnjCCAoagAwIBAgIGAV2VSLVQMA0GCSqGSIb3DQEBCwUAMIGPMQswCQYDVQQGEwJVUzETMBEG
A1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNjbzENMAsGA1UECgwET2t0YTEU
MBIGA1UECwwLU1NPUHJvdmlkZXIxEDAOBgNVBAMMB3NlZ21lbnQxHDAaBgkqhkiG9w0BCQEWDWlu
Zm9Ab2t0YS5jb20wHhcNMTcwNzMwMjA1NDU2WhcNMjcwNzMwMjA1NTU2WjCBjzELMAkGA1UEBhMC
VVMxEzARBgNVBAgMCkNhbGlmb3JuaWExFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28xDTALBgNVBAoM
BE9rdGExFDASBgNVBAsMC1NTT1Byb3ZpZGVyMRAwDgYDVQQDDAdzZWdtZW50MRwwGgYJKoZIhvcN
AQkBFg1pbmZvQG9rdGEuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0coLp1Rs
lJ0wA1EdBxrK1RhM7SGmOiygS0wJ9usYeeWpKjdqGKVbLJ9Yl2sZ4QcChYkSQN0VPGtgA4kQ0a2u
ErtU+HhFNeHO1sJtUPUkhCvBDhDdw1q+RO9h9NNn4LkA/VqEWlZKweppNS2qwwBn7as1ElwlwYlR
B70DhFSbiXcOL5tiT72ixvbkLhpeWu4uflKCAbVPT2vsCGV02UrUz+b3VxXkGK3T8dlFWwZxy5rs
q0Kx9FgmhSBryNAn9RSR+qj/XJkl/S72VwM896UfARb4bY1ThyN0LbowtdhN4/8vIr+G/BfSvfMM
uDDdLMlNUPj5mw1kOkGUXqu8aOmJZQIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQDNmsvXmeEyYCLR
kpGMoVAOSA8Y2gSX3D5oMz4bpq7loI635JOw8USom6w3vQoZ6nyEs/w81bUj0w0ZfioilXBiW1ei
e8OOZFTu+65yZztbkrQCnHaLyOiVB2QwJ+GetlRU22NlqsDoJdxNfi82yXgdhKhXa79v4e+HvxGo
LvSsVmQs5RA43jXppaZahAPThaEUjrwwZ0C1OUDnuotZy//gY2GtihpW2NeI6bWNQu4AcNGmXeuk
uErRpY4i6k/GBZiyDazqeZqWe9f5Z1Dqxa1s6YIezCE7pdUFn4guqfSbbkj9IinkX9AzI93XTOaM
ukq9hNcrvMJvBGeWim6vloPQ
-----END CERTIFICATE-----`))

	s := `<?xml version="1.0" encoding="UTF-8"?><saml2p:Response Destination="http://localhost:8080/acs" ID="id55456921196388931032386368" IssueInstant="2020-05-15T20:34:27.037Z" Version="2.0" xmlns:saml2p="urn:oasis:names:tc:SAML:2.0:protocol"><saml2:Issuer Format="urn:oasis:names:tc:SAML:2.0:nameid-format:entity" xmlns:saml2="urn:oasis:names:tc:SAML:2.0:assertion">http://www.okta.com/exkrbxxvs5Fw8wDHj0h7</saml2:Issuer><ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#"><ds:SignedInfo><ds:CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/><ds:SignatureMethod Algorithm="http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"/><ds:Reference URI="#id55456921196388931032386368"><ds:Transforms><ds:Transform Algorithm="http://www.w3.org/2000/09/xmldsig#enveloped-signature"/><ds:Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/></ds:Transforms><ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/><ds:DigestValue>uicsT0uZmeJf8VAG293+uzORVRRt8JeXMNirfZ9/9xM=</ds:DigestValue></ds:Reference></ds:SignedInfo><ds:SignatureValue>iCcr+on1UPGbDL3haAFF1M7WCb13+qYCWG7kkKpD3BLVjEhV8J2dQpoI/VHrORe+eJgBW+k9Gg9HJOp5iVcsdGFNiFAjHhQ5zsBpQ87QU6WvfPM28+owFptFtnHXo41pZx2d2AR5Z6xg8UjY61levfmkLM4+9GMmpExoID915+AonqoYGsZFGelEv5CIisn7TcZy2K5fUxJIdvnSnD3HkwS6Y8TsJA9DXUMKSMrnLicn+C0o253Ow5m7x1gHYn5NXPK9lCLN6DzePJuPfn0Q9NM5qbE0PwSt3cd6gm/qOkmD+58GoRNz53Fr1xkrltPoVDgaEpMCwhvODaTw9PBTZQ==</ds:SignatureValue><ds:KeyInfo><ds:X509Data><ds:X509Certificate>MIIDnjCCAoagAwIBAgIGAV2VSLVQMA0GCSqGSIb3DQEBCwUAMIGPMQswCQYDVQQGEwJVUzETMBEG
A1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNjbzENMAsGA1UECgwET2t0YTEU
MBIGA1UECwwLU1NPUHJvdmlkZXIxEDAOBgNVBAMMB3NlZ21lbnQxHDAaBgkqhkiG9w0BCQEWDWlu
Zm9Ab2t0YS5jb20wHhcNMTcwNzMwMjA1NDU2WhcNMjcwNzMwMjA1NTU2WjCBjzELMAkGA1UEBhMC
VVMxEzARBgNVBAgMCkNhbGlmb3JuaWExFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28xDTALBgNVBAoM
BE9rdGExFDASBgNVBAsMC1NTT1Byb3ZpZGVyMRAwDgYDVQQDDAdzZWdtZW50MRwwGgYJKoZIhvcN
AQkBFg1pbmZvQG9rdGEuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0coLp1Rs
lJ0wA1EdBxrK1RhM7SGmOiygS0wJ9usYeeWpKjdqGKVbLJ9Yl2sZ4QcChYkSQN0VPGtgA4kQ0a2u
ErtU+HhFNeHO1sJtUPUkhCvBDhDdw1q+RO9h9NNn4LkA/VqEWlZKweppNS2qwwBn7as1ElwlwYlR
B70DhFSbiXcOL5tiT72ixvbkLhpeWu4uflKCAbVPT2vsCGV02UrUz+b3VxXkGK3T8dlFWwZxy5rs
q0Kx9FgmhSBryNAn9RSR+qj/XJkl/S72VwM896UfARb4bY1ThyN0LbowtdhN4/8vIr+G/BfSvfMM
uDDdLMlNUPj5mw1kOkGUXqu8aOmJZQIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQDNmsvXmeEyYCLR
kpGMoVAOSA8Y2gSX3D5oMz4bpq7loI635JOw8USom6w3vQoZ6nyEs/w81bUj0w0ZfioilXBiW1ei
e8OOZFTu+65yZztbkrQCnHaLyOiVB2QwJ+GetlRU22NlqsDoJdxNfi82yXgdhKhXa79v4e+HvxGo
LvSsVmQs5RA43jXppaZahAPThaEUjrwwZ0C1OUDnuotZy//gY2GtihpW2NeI6bWNQu4AcNGmXeuk
uErRpY4i6k/GBZiyDazqeZqWe9f5Z1Dqxa1s6YIezCE7pdUFn4guqfSbbkj9IinkX9AzI93XTOaM
ukq9hNcrvMJvBGeWim6vloPQ</ds:X509Certificate></ds:X509Data></ds:KeyInfo></ds:Signature><saml2p:Status xmlns:saml2p="urn:oasis:names:tc:SAML:2.0:protocol"><saml2p:StatusCode Value="urn:oasis:names:tc:SAML:2.0:status:Success"/></saml2p:Status><saml2:Assertion ID="id55456921197340911775845224" IssueInstant="2020-05-15T20:34:27.037Z" Version="2.0" xmlns:saml2="urn:oasis:names:tc:SAML:2.0:assertion"><saml2:Issuer Format="urn:oasis:names:tc:SAML:2.0:nameid-format:entity" xmlns:saml2="urn:oasis:names:tc:SAML:2.0:assertion">http://www.okta.com/exkrbxxvs5Fw8wDHj0h7</saml2:Issuer><ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#"><ds:SignedInfo><ds:CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/><ds:SignatureMethod Algorithm="http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"/><ds:Reference URI="#id55456921197340911775845224"><ds:Transforms><ds:Transform Algorithm="http://www.w3.org/2000/09/xmldsig#enveloped-signature"/><ds:Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/></ds:Transforms><ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/><ds:DigestValue>5BSTzVaflJGYosfXYB1Q8uX/MH1PQBZi2sQ3EaVt2es=</ds:DigestValue></ds:Reference></ds:SignedInfo><ds:SignatureValue>KcB6wj6OGum5sjx78ElfsL4jc421a6h0JLBU4W7u9OTQS106e52IuwchAnZBmDQ0uvUP/qU8Ub/9lT13kfsGeRpgk7eCngMSFBivhlgQ7bwLTplo0i22worw1cLh7Wgdnj2vivmflCC1VABk5GF7RE3vNBQGpexAo/LFxAa6rA01/SvFq6L+XxbYxcYqSKCL/EwKampNqEnzdoNXHXLvciJW6n6G7Vl/g3C+GhuZayUiMX/suQZm83ueFgr7hr2juAXJ/l3EBf4aO3A8YlN7KbVp4V/fGsVgWHbr68H0jOu9aZU8Hr92b0a4HrEuhFtmB0FY5GF5a2GW3I6lmwkSyQ==</ds:SignatureValue><ds:KeyInfo><ds:X509Data><ds:X509Certificate>MIIDnjCCAoagAwIBAgIGAV2VSLVQMA0GCSqGSIb3DQEBCwUAMIGPMQswCQYDVQQGEwJVUzETMBEG
A1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNjbzENMAsGA1UECgwET2t0YTEU
MBIGA1UECwwLU1NPUHJvdmlkZXIxEDAOBgNVBAMMB3NlZ21lbnQxHDAaBgkqhkiG9w0BCQEWDWlu
Zm9Ab2t0YS5jb20wHhcNMTcwNzMwMjA1NDU2WhcNMjcwNzMwMjA1NTU2WjCBjzELMAkGA1UEBhMC
VVMxEzARBgNVBAgMCkNhbGlmb3JuaWExFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28xDTALBgNVBAoM
BE9rdGExFDASBgNVBAsMC1NTT1Byb3ZpZGVyMRAwDgYDVQQDDAdzZWdtZW50MRwwGgYJKoZIhvcN
AQkBFg1pbmZvQG9rdGEuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0coLp1Rs
lJ0wA1EdBxrK1RhM7SGmOiygS0wJ9usYeeWpKjdqGKVbLJ9Yl2sZ4QcChYkSQN0VPGtgA4kQ0a2u
ErtU+HhFNeHO1sJtUPUkhCvBDhDdw1q+RO9h9NNn4LkA/VqEWlZKweppNS2qwwBn7as1ElwlwYlR
B70DhFSbiXcOL5tiT72ixvbkLhpeWu4uflKCAbVPT2vsCGV02UrUz+b3VxXkGK3T8dlFWwZxy5rs
q0Kx9FgmhSBryNAn9RSR+qj/XJkl/S72VwM896UfARb4bY1ThyN0LbowtdhN4/8vIr+G/BfSvfMM
uDDdLMlNUPj5mw1kOkGUXqu8aOmJZQIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQDNmsvXmeEyYCLR
kpGMoVAOSA8Y2gSX3D5oMz4bpq7loI635JOw8USom6w3vQoZ6nyEs/w81bUj0w0ZfioilXBiW1ei
e8OOZFTu+65yZztbkrQCnHaLyOiVB2QwJ+GetlRU22NlqsDoJdxNfi82yXgdhKhXa79v4e+HvxGo
LvSsVmQs5RA43jXppaZahAPThaEUjrwwZ0C1OUDnuotZy//gY2GtihpW2NeI6bWNQu4AcNGmXeuk
uErRpY4i6k/GBZiyDazqeZqWe9f5Z1Dqxa1s6YIezCE7pdUFn4guqfSbbkj9IinkX9AzI93XTOaM
ukq9hNcrvMJvBGeWim6vloPQ</ds:X509Certificate></ds:X509Data></ds:KeyInfo></ds:Signature><saml2:Subject xmlns:saml2="urn:oasis:names:tc:SAML:2.0:assertion"><saml2:NameID Format="urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified">ulysse@segment.com</saml2:NameID><saml2:SubjectConfirmation Method="urn:oasis:names:tc:SAML:2.0:cm:bearer"><saml2:SubjectConfirmationData NotOnOrAfter="2020-05-15T20:39:27.037Z" Recipient="http://localhost:8080/acs"/></saml2:SubjectConfirmation></saml2:Subject><saml2:Conditions NotBefore="2020-05-15T20:29:27.037Z" NotOnOrAfter="2020-05-15T20:39:27.037Z" xmlns:saml2="urn:oasis:names:tc:SAML:2.0:assertion"><saml2:AudienceRestriction><saml2:Audience>ucarion-test</saml2:Audience></saml2:AudienceRestriction></saml2:Conditions><saml2:AuthnStatement AuthnInstant="2020-05-15T20:33:59.121Z" SessionIndex="id1589574867035.2142172084" xmlns:saml2="urn:oasis:names:tc:SAML:2.0:assertion"><saml2:AuthnContext><saml2:AuthnContextClassRef>urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport</saml2:AuthnContextClassRef></saml2:AuthnContext></saml2:AuthnStatement></saml2:Assertion></saml2p:Response>`

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic(err)
	}

	ctx := x_dsig.NewDefaultValidationContext(&x_dsig.MemoryX509CertificateStore{
		Roots: []*x509.Certificate{cert},
	})

	el := etree.NewDocument()
	el.ReadFromString(s)

	// It is important to only use the returned validated element.
	// See: https://www.w3.org/TR/xmldsig-bestpractices/#check-what-is-signed
	validated, err := ctx.Validate(el.Root())
	if err != nil {
		panic(err)
	}

	fmt.Println(validated)

	var a assertion
	if err := xml.Unmarshal([]byte(s), &a); err != nil {
		panic(err)
	}

	decoder := xml.NewDecoder(strings.NewReader(s))
	if err := a.Signature.Verify(cert, decoder); err != nil {
		panic(err)
	}
}
