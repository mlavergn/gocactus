package cactus

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

// Version export
const Version = "0.1.0"

// Cactus export
type Cactus struct {
	httpClient *http.Client
	urlRoot    string
}

// NewCactus export
func NewCactus() *Cactus {
	id := &Cactus{
		urlRoot: "https://billing.cactusvpn.com",
	}

	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// currently based on Linux CA location
	caCert, err := ioutil.ReadFile("/etc/ssl/ca-bundle.crt")
	if err == nil {
		rootCAs.AppendCertsFromPEM(caCert)
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		RootCAs:            rootCAs,
	}

	httpTransport := &http.Transport{
		TLSClientConfig: tlsConfig,
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
	}

	id.httpClient = &http.Client{
		Transport: httpTransport,
	}

	return id
}

// Token export
func (id *Cactus) Token() string {
	log.Println("Cactus.Token")

	req, rerr := http.NewRequest("GET", id.urlRoot+"/clientarea.php", nil)
	if rerr != nil {
		log.Println("Cactus.Token failed request", rerr)
		return ""
	}

	req.Header.Add("Connection", "close")

	resp, derr := id.httpClient.Do(req)
	if derr != nil {
		log.Println("Cactus.Token failed send", derr)
		return ""
	}
	defer resp.Body.Close()

	data, berr := ioutil.ReadAll(resp.Body)
	if berr != nil {
		log.Println("Cactus.Token failed read", derr)
		return ""
	}
	html := string(data)

	// find the CSRF token
	re, err := regexp.Compile(".*csrfToken = '([^']*)'.*")
	if err != nil {
		fmt.Println("failed to compile", err)
		return ""
	}
	match := re.FindStringSubmatch(html)
	if len(match) < 2 {
		fmt.Println("failed to match")
		return ""
	}

	return match[1]
}

// Login export
func (id *Cactus) Login(email string, password string) bool {
	log.Println("Cactus.Login")

	token := id.Token()

	params := url.Values{}
	params.Set("username", email)
	params.Set("password", password)
	params.Set("token", token)
	body := bytes.NewBufferString(params.Encode())

	req, rerr := http.NewRequest("POST", id.urlRoot+"/dologin.php", body)
	if rerr != nil {
		log.Println("Cactus.Login failed request", rerr)
		return false
	}

	req.Header.Add("Connection", "close")

	resp, derr := id.httpClient.Do(req)
	if derr != nil {
		log.Println("Cactus.Login failed send", derr)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 {
		return true
	}

	return true
}

// ValidateIP export
func (id *Cactus) ValidateIP() bool {
	log.Println("Cactus.ValidateIP")

	req, rerr := http.NewRequest("GET", id.urlRoot+"/clientarea.php", nil)
	if rerr != nil {
		log.Println("Cactus.ValidateIP failed request", rerr)
		return false
	}

	req.Header.Add("Referrer", id.urlRoot+"/settings.php")
	req.Header.Add("Connection", "close")

	resp, derr := id.httpClient.Do(req)
	if derr != nil {
		log.Println("Cactus.ValidateIP failed send", derr)
		return false
	}
	defer resp.Body.Close()

	data, berr := ioutil.ReadAll(resp.Body)
	if berr != nil {
		log.Println("Cactus.ValidateIP failed read", derr)
		return false
	}

	// TODO
	// can use use status code or does the page have to be parsed
	if len(data) == 0 {
		return false
	}

	return true
}

// GetDNS export
func (id *Cactus) GetDNS() string {
	log.Println("Cactus.GetDNS")

	req, rerr := http.NewRequest("GET", id.urlRoot+"/dns-servers.php", nil)
	if rerr != nil {
		log.Println("Cactus.GetDNS failed request", rerr)
		return ""
	}

	req.Header.Add("Connection", "close")

	resp, derr := id.httpClient.Do(req)
	if derr != nil {
		log.Println("Cactus.GetDNS failed send", derr)
		return ""
	}
	defer resp.Body.Close()

	data, berr := ioutil.ReadAll(resp.Body)
	if berr != nil {
		log.Println("Cactus.GetDNS failed read", derr)
		return ""
	}
	html := string(data)
	return id.parseDNS(html)
}

func (id *Cactus) parseDNS(html string) string {
	// find the DNS entry
	re, err := regexp.Compile(".*id=\"dns[1]\" value=\"([^\"]+)\".*")
	if err != nil {
		fmt.Println("failed to compile", err)
		return ""
	}
	match := re.FindStringSubmatch(html)
	if len(match) < 2 {
		fmt.Println("failed to match")
		return ""
	}

	return match[1]
}
