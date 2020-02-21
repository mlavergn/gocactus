package cactus

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestParseDNS(t *testing.T) {
	cactus := NewCactus()
	file, err := os.Open("./test/dns.html")
	if err != nil {
		t.Error("failed to open test content", err)
	}
	data, rerr := ioutil.ReadAll(file)
	if rerr != nil {
		t.Error("failed to read test content", rerr)
	}

	html := string(data)

	actual := cactus.parseDNS(html)
	expected := "87.117.205.136"

	if actual != expected {
		t.Error("actual did not match expected value", actual, expected)
	}
}
