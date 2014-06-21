package ghastly

import (
	"fmt"
	"time"
)

// A Service is a particular website, app, etc. served through Fastly. They are
// configured with different versions, which have backends, domains, etc.
type Service struct {
	Id            string
	Name          string
	CustomerId    string
	PublishKey    string
	Comment       string
	ActiveVersion int64
	UpdatedAt     time.Time
	CreatedAt     time.Time
	versions      map[int]*Version
	ghastly       *Ghastly
}

// Get a service with the ID string.
func (g *Ghastly) GetService(id string) (*Service, error) {
	url := makeServiceURL(id)
	resp, err := g.Get(url)
	if err != nil {
		return nil, err
	}
	sData, err := ParseJson(resp.Body)
	if err != nil {
		return nil, err
	}
	return g.populateService(sData)
}

// List the current services.
func (g *Ghastly) ListServices() ([]*Service, error) {
	resp, err := g.Get("/service")
	if err != nil {
		return nil, err
	}

	servicesData, err := ParseJsonArray(resp.Body)
	if err != nil {
		return nil, err
	}

	services := make([]*Service, len(servicesData))
	for i, v := range servicesData {
		ss, err := g.populateService(v.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		services[i] = ss
	}
	return services, nil

	return nil, nil
}

// Search for a service by name. The API does not appear to permit wildcards at
// this time.
func (g *Ghastly) SearchServices(searchStr string) (*Service, error) {
	params := map[string]string{"name": searchStr}
	searchURL := makeServiceURL("search")
	resp, err := g.GetParams(searchURL, params)
	if err != nil {
		return nil, err
	}

	s, err := ParseJson(resp.Body)
	return g.populateService(s)
}

// Create a new service.
func (g *Ghastly) NewService(name string) (*Service, error) {
	params := map[string]string{"name": name}
	resp, err := g.PostFormParams("/service", params)
	if err != nil {
		return nil, err
	}
	sData, err := ParseJson(resp.Body)
	if err != nil {
		return nil, err
	}

	return g.populateService(sData)
}

func (g *Ghastly) populateService(serviceData map[string]interface{}) (*Service, error) {
	s := new(Service)
	s.Id = serviceData["id"].(string)
	s.Name = serviceData["name"].(string)
	s.CustomerId = serviceData["customer_id"].(string)
	s.PublishKey, _ = serviceData["publish_key"].(string)
	s.Comment, _ = serviceData["comment"].(string)
	s.ghastly = g

	if cc, ok := serviceData["created_at"].(string); ok {
		createdAt, err := time.Parse(time.RFC3339, cc)
		if err != nil {
			return nil, err
		}
		s.CreatedAt = createdAt
	}
	if uc, ok := serviceData["created_at"].(string); ok {
		updatedAt, err := time.Parse(time.RFC3339, uc)
		if err != nil {
			return nil, err
		}
		s.UpdatedAt = updatedAt
	}

	return s, nil
}

// Delete a service and everything attached to it.
func (s *Service) Delete() error {
	url := makeServiceURL(s.Id)
	_, err := s.ghastly.Delete(url)
	if err != nil {
		return err
	}
	return nil
}

// Make the base URL for this service for performing tasks.
func (s *Service) TaskURL(taskPath string) string {
	serviceURL := makeServiceURL(s.Id)
	url := fmt.Sprintf("%s/%s", serviceURL, taskPath)
	return url
}

// Get detailed information about a service.
func (s *Service) Details() (*Service, error) {
	detailURL := makeServiceURL(fmt.Sprintf("%s/details", s.Id))
	resp, err := s.ghastly.Get(detailURL)
	if err != nil {
		return nil, err
	}
	sData, err := ParseJson(resp.Body)
	if err != nil {
		return nil, err
	}
	return s.ghastly.populateService(sData)
}

// Update a service's attributes with a map of strings for params. Currently the
// only meaningful attribute that can be updated here is the service's name.
func (s *Service) Update(params map[string]string) error {
	url := makeServiceURL(s.Id)
	_, err := s.ghastly.PutParams(url, params)
	if err != nil {
		return nil
	}
	s.Name = params["name"]
	return nil
}

func makeServiceURL(id string) string {
	return fmt.Sprintf("/service/%s", id)
}
