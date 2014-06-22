package ghastly

import (
	"fmt"
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

// Create a new domain for a particular version of a service. Possible
// parameters are "name" and "comment".
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
	url := v.baseURL("domain/check_all")
	resp, err := v.service.ghastly.Get(url)
	if err != nil {
		return nil, err
	}
	dData, err := ParseJsonArray(resp.Body)
	if err != nil {
		return nil, err
	}
	checkedDomains := make([]*DomainCheck, len(dData))
	for i, dc := range dData {
		d, err := v.populateDomainCheck(dc.([]interface{}))
		if err != nil {
			return nil, err
		}
		checkedDomains[i] = d
	}
	return checkedDomains, nil
}

// Check one domain associated with a service.
func (v *Version) CheckDomain(name string) (*DomainCheck, error) {
	task := fmt.Sprintf("domain/%s/check", name)
	url := v.baseURL(task)
	resp, err := v.service.ghastly.Get(url)
	if err != nil {
		return nil, err
	}
	dData, err := ParseJsonArray(resp.Body)
	if err != nil {
		return nil, err
	}
	return v.populateDomainCheck(dData)
}

// List all domains associated with a version of a service.
func (v *Version) ListDomains() ([]*Domain, error) {
	url := v.baseURL("domain")
	resp, err := v.service.ghastly.Get(url)
	if err != nil {
		return nil, err
	}
	dData, err := ParseJsonArray(resp.Body)
	if err != nil {
		return nil, err
	}
	domains := make([]*Domain, len(dData))
	for i, dc := range dData {
		d, err := v.populateDomain(dc.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		domains[i] = d
	}
	return domains, nil
}

// Get a domain associated with this version
func (v *Version) GetDomain(name string) (*Domain, error) {
	task := fmt.Sprintf("domain/%s", name)
	url := v.baseURL(task)
	resp, err := v.service.ghastly.Get(url)
	if err != nil {
		return nil, err
	}
	dData, err := ParseJson(resp.Body)
	if err != nil {
		return nil, err
	}

	return v.populateDomain(dData)
}

// Delete a domain, for the version the domain belongs to.
func (d *Domain) Delete() error {
	task := fmt.Sprintf("domain/%s", d.Name)
	url := d.version.baseURL(task)
	_, err := d.version.service.ghastly.Delete(url)
	if err != nil {
		return err
	}
	return nil
}

// Update a domain, for the version the domain belongs to. Possible parameters
// for the domain are "name" and "comment".
func (d *Domain) Update(params map[string]string) error {
	task := fmt.Sprintf("domain/%s", d.Name)
	url := d.version.baseURL(task)
	_, err := d.version.service.ghastly.PutParams(url, params)
	if err != nil {
		return err
	}
	d.Name = params["name"]
	if c, ok := params["comment"]; ok {
		d.Comment = c
	}
	return nil
}

func (v *Version) populateDomain(domainData map[string]interface{}) (*Domain, error) {
	locked, _ := domainData["locked"].(bool)
	comment, _ := domainData["comment"].(string)
	return &Domain{Name: domainData["name"].(string), Comment: comment, Locked: locked, ServiceId: domainData["service_id"].(string), Version: int64(domainData["version"].(float64)), version: v}, nil
}

func (v *Version) populateDomainCheck(domainData []interface{}) (*DomainCheck, error) {
	d, err := v.populateDomain(domainData[0].(map[string]interface{}))
	if err != nil {
		return nil, err
	}
	cname, _ := domainData[1].(string)
	proper, _ := domainData[2].(bool)
	return &DomainCheck{d, cname, proper}, nil
}
