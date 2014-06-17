package ghastly

import (
	"fmt"
	"time"
)

// A version of the configuration for a particular service. Backends, domains,
// etc. belong to a particular version.
type Version struct {
	Number    int64
	ServiceId string
	Active    bool
	Locked    bool
	Comment   string
	Testing   bool
	Staging   bool
	Deployed  bool
	Network   VersionNetwork
	deletedAt time.Time
	service   *Service
}

type VersionNetwork struct {
	Name                string
	Description         string
	AvailableAll        bool
	AvailableRestricted bool
	AvailablePrivate    bool
	CustomerId          string
}

func (s *Service) populateVersion(versionData map[string]interface{}) (*Version, error) {
	active, _ := versionData["active"].(bool)
	locked, _ := versionData["locked"].(bool)
	testing, _ := versionData["testing"].(bool)
	staging, _ := versionData["staging"].(bool)
	deployed, _ := versionData["deployed"].(bool)
	comment, _ := versionData["comment"].(string)
	return &Version{Number: int64(versionData["number"].(float64)), ServiceId: versionData["service_id"].(string), Active: active, Locked: locked, Comment: comment, Testing: testing, Staging: staging, Deployed: deployed, service: s}, nil
}

// Create a brand new, pristine version of a service, with nothing in it.
func (s *Service) NewVersion() (*Version, error) {
	params := map[string]string{"service": s.Id}
	url := s.TaskURL("/version")
	resp, err := s.ghastly.PostFormParams(url, params)
	if err != nil {
		return nil, err
	}
	vData, err := ParseJson(resp.Body)
	if err != nil {
		return nil, err
	}

	return s.populateVersion(vData)
}

// Get a particular version of this service identified by the version number.
func (s *Service) GetVersion(number int64) (*Version, error) {
	u := fmt.Sprintf("/version/%d", number)
	url := s.TaskURL(u)
	resp, err := s.ghastly.Get(url)
	if err != nil {
		return nil, err
	}
	vData, err := ParseJson(resp.Body)
	if err != nil {
		return nil, err
	}

	return s.populateVersion(vData)
}

/*
// List all versions belonging to a service.
func (s *Service) ListVersions() ([]*Version, error) {

}
*/

// Clone this version of the service, returning the new version.
func (v *Version) Clone() (*Version, error) {
	u := fmt.Sprintf("/version/%d/clone", v.Number)
	url := v.service.TaskURL(u)
	resp, err := v.service.ghastly.Put(url, nil)
	if err != nil {
		return nil, err
	}
	vData, err := ParseJson(resp.Body)
	if err != nil {
		return nil, err
	}

	return v.service.populateVersion(vData)
}

func (v *Version) Activate() error {

}
