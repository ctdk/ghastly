package ghastly

import (
	"time"
)
// A version of the configuration for a particular service. Backends, domains,
// etc. belong to a particular version.
type Version struct {
	Number int64
	ServiceId string
	Active bool
	Locked bool
	Comment string
	Testing bool
	Staging bool
	Deployed bool
	Network VersionNetwork
	deletedAt time.Time
	service *Service
}

type VersionNetwork struct {
	Name string
	Description string
	AvailableAll bool
	AvailableRestricted bool
	AvailablePrivate bool
	CustomerId string
}

func (s *Service) populateVersion(versionData map[string]interface{}) (*Version, error) {
	return &Version{}, nil
}

// Create a brand new, pristine version of a service, with nothing in it.
func (s *Service) NewVersion() (*Version, error) {
	params := map[string]string{ "service": s.Id }
	url := s.TaskURL("/version")
	resp, err := s.ghastly.PostFormParams(url, params)
	if err != nil {
		return nil, err
	}
	sData, err := ParseJson(resp.Body)
	if err != nil {
		return nil, err
	}

	return s.populateVersion(sData)
}
