package cactus

import (
	"log"
	"os/exec"
	"strings"
)

// OpenWRT type
type OpenWRT struct {
}

// NewOpenWRT export
func NewOpenWRT() *OpenWRT {
	return &OpenWRT{}
}

// Service export
func (id *OpenWRT) Service(service string, command string) (string, error) {
	cmd := exec.Command("service", service, command)
	err := cmd.Run()
	if err != nil {
		log.Println("service failed to exec", service, command, err)
		return "", err
	}
	return "", nil
}

// UCI export
func (id *OpenWRT) UCI(command string, arg string) (string, error) {
	cmd := exec.Command("uci", command, arg)
	err := cmd.Run()
	if err != nil {
		log.Println("uci failed to exec", command, arg, err)
		return "", err
	}
	return "", nil
}

// setDns export
func (id *OpenWRT) setDNS(servers []string) bool {
	serverString := strings.Join(servers, " ")
	_, err := id.UCI("set", "network.wan.dns=\""+serverString+"\"")
	if err != nil {
		log.Println("setDNS failed to uci set")
		return false
	}
	_, err = id.UCI("commit", "network")
	if err != nil {
		log.Println("setDNS failed to uci commit")
		return false
	}
	_, err = id.Service("network", "reload")
	if err != nil {
		log.Println("setDNS failed to network reload")
		return false
	}
	return true
}
