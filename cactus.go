package cactus

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

// Version export
const Version = "0.1.0"

// Cactus export
type Cactus struct {
	Client  *http.Client
	URLRoot string
}

// NewCactus export
func NewCactus() *Cactus {
	id := &Cactus{
		URLRoot: "https://billing.cactusvpn.com",
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

	id.Client = &http.Client{
		Transport: httpTransport,
	}

	return id
}

// Token export
func (id *Cactus) Token() string {
	log.Println("Cactus.Token")

	req, rerr := http.NewRequest("GET", id.URLRoot+"/clientarea.php", nil)
	if rerr != nil {
		log.Println("Cactus.Token failed request", rerr)
		return ""
	}

	req.Header.Add("Connection", "close")

	resp, derr := id.Client.Do(req)
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

	// TODO
	// .*csrfToken = '([^']*)'.*
	return string(data)
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

	req, rerr := http.NewRequest("POST", id.URLRoot+"/dologin.php", body)
	if rerr != nil {
		log.Println("Cactus.Login failed request", rerr)
		return false
	}

	req.Header.Add("Connection", "close")

	resp, derr := id.Client.Do(req)
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

	req, rerr := http.NewRequest("GET", id.URLRoot+"/clientarea.php", nil)
	if rerr != nil {
		log.Println("Cactus.ValidateIP failed request", rerr)
		return false
	}

	req.Header.Add("Referrer", id.URLRoot+"/settings.php")
	req.Header.Add("Connection", "close")

	resp, derr := id.Client.Do(req)
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
	if len(data) == 0 {
		return false
	}

	return true
}

// GetDNS export
func (id *Cactus) GetDNS() string {
	log.Println("Cactus.GetDNS")

	req, rerr := http.NewRequest("GET", id.URLRoot+"/dns-servers.php", nil)
	if rerr != nil {
		log.Println("Cactus.GetDNS failed request", rerr)
		return ""
	}

	req.Header.Add("Connection", "close")

	resp, derr := id.Client.Do(req)
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

	// TODO
	// .*id="dns[1]" value="([^"]+)".*
	return string(data)
}
