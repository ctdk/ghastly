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

func (v *Version) populateDomain(domainData map[string]interface{}) (*Domain, error) {
	locked, _ := domainData["locked"].(bool)
	comment, _ := domainData["comment"].(string)
	return &Domain{Name: domainData["name"].(string), Comment: comment, Locked: locked, ServiceId: domainData["service_id"].(string), Version: int64(domainData["version"].(float64)), version: v}, nil
}
