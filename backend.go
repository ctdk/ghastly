package ghastly

import (
	"fmt"
)

type Backend struct {
	Name                string
	Address             string
	Port                uint16
	UseSSL              bool
	ConnectTimeout      int
	FirstByteTimeout    int
	BetweenBytesTimeout int
	ErrorThreshold      int
	MaxConn             int
	Weight              int
	AutoLoadbalance     int
	RequestCondition    string
	Healthcheck         string
	SSLClientCert string
	SSLClientKey string
	SSLHostname string
	SSLCACert string
	ClientCert string
	Comment string
	Hostname string
	Ipv4 string
	Ipv6 string
	version             *Version
}

// Create a new backend for a service and version.
func (v *Version) NewBackend(params map[string]string) (*Backend, error) {
	url := v.baseURL("backend")
	resp, err := v.service.ghastly.PostFormParams(url, params)
	if err != nil {
		return nil, err
	}
	bData, err := ParseJson(resp.Body)
	if err != nil {
		return nil, err
	}

	return v.populateBackend(bData)
}

// Get a backend associated with this version.
func (v *Version) GetBackend(name string) (*Backend, error) {
	task := fmt.Sprintf("backend/%s", name)
	url := v.baseURL(task)
	resp, err := v.service.ghastly.Get(url)
	if err != nil {
		return nil, err
	}
	bData, err := ParseJson(resp.Body)
	if err != nil {
		return nil, err
	}

	return v.populateBackend(bData)
}

// Delete a backend.
func (b *Backend) Delete() error {
	task := fmt.Sprintf("backend/%s", b.Name)
	url := b.version.baseURL(task)
	_, err := b.version.service.ghastly.Delete(url)
	if err != nil {
		return err
	}
	return nil
}

func (v *Version) populateBackend(backendData map[string]interface{}) (*Backend, error) {
	fmt.Printf("backend data: %v\n", backendData)
	name, ok := backendData["name"].(string)
	if !ok {
		err := fmt.Errorf("backend name invalid")
		return nil, err
	}
	address, ok := backendData["address"].(string)
	if !ok {
		err := fmt.Errorf("backend address invalid")
		return nil, err
	}
	var port uint16
	if p, ok := backendData["port"].(float64); !ok {
		err := fmt.Errorf("backend port invalid")
		return nil, err
	} else {
		port = uint16(p)
	}

	comment, _ := backendData["comment"].(string)
	ipv4, _ := backendData["ipv4"].(string)
	ipv6, _ := backendData["ipv6"].(string)
	hostname, _ := backendData["hostname"].(string)

	fmt.Printf("backend stuff: %v\n", backendData)
	return &Backend{ Name: name, Address: address, Port: port, Comment: comment, Hostname: hostname, Ipv4: ipv4, Ipv6: ipv6 }, nil
}
