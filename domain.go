package ghastly

import (
//"fmt"
)

type Domain struct {
	Name      string
	Comment   string
	ServiceId string
	Version   int64
	Locked    bool
	version   *Version
}

type DomainCheck struct {
	*Domain
	Cname    string
	IsProper bool
}

// Create a new domain for a particular version of a service. Possible parameters
// are "name" and "comment".
func (v *Version) NewDomain(params map[string]string) (*Domain, error) {
	url := v.baseURL("domain")
	resp, err := v.service.ghastly.PostFormParams(url, params)
	if err != nil {
		return nil, err
	}
	dData, err := ParseJson(resp.Body)
	if err != nil {
		return nil, err
	}

	return v.populateDomain(dData)
}

// Check all domains associated with a version of a service.
func (v *Version) CheckAllDomains() ([]*DomainCheck, error) {

}

// Check one domain associated with a service.
func (v *Version) CheckDomain(name string) (*DomainCheck, error) {

}

// List all domains associated with a version of a service.
func (v *Version) ListDomains() ([]*Domain, error) {

}

// Get a domain associated with this version
func (v *Version) GetDomain(name string) (*Domain, error) {

}

// Delete a domain, for the version the domain belongs to.
func (d *Domain) Delete() error {

}

// Update a domain, for the version the domain belongs to. Possible parameters
// for the domain are "name" and "comment".
func (d *Domain) Update(params map[string]string) error {

}

func (v *Version) populateDomain(domainData map[string]interface{}) (*Domain, error) {
	locked, _ := domainData["locked"].(bool)
	comment, _ := domainData["comment"].(string)
	return &Domain{Name: domainData["name"].(string), Comment: comment, Locked: locked, ServiceId: domainData["service_id"].(string), Version: int64(domainData["version"].(float64)), version: v}, nil
}
