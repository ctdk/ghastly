package ghastly

import (
	"time"
	"fmt"
)

// A Service is a particular website, app, etc. served through Fastly. They are
// configured with different versions, which have backends, domains, etc.
type Service struct {
	Id string
	Name string
	CustomerId string
	PublishKey string
	Comment string
	ActiveVersion int64
	UpdatedAt time.Time
	CreatedAt time.Time
	versions map[int]*Version
	ghastly *Ghastly
}

func (g *Ghastly)GetService(id string) (*Service, error) {
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

func (g *Ghastly)ListServices() ([]*Service, error) {
	return nil, nil
}

func (g *Ghastly)NewService(name string) (*Service, error) {
	params := map[string]string{ "name": name }
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

func (g *Ghastly)populateService(serviceData map[string]interface{}) (*Service, error) {
	createdAt, err := time.Parse(time.RFC3339, serviceData["created_at"].(string))
	if err != nil {
		return nil, err
	}
	updatedAt, err := time.Parse(time.RFC3339, serviceData["updated_at"].(string))
	if err != nil {
		return nil, err
	}
	
	return &Service{ Id: serviceData["id"].(string), Name: serviceData["name"].(string), CustomerId: serviceData["customer_id"].(string), PublishKey: serviceData["publish_key"].(string), Comment: serviceData["comment"].(string), CreatedAt: createdAt, UpdatedAt: updatedAt, ghastly: g }, nil
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

func (s *Service)TaskURL(taskPath string) string {
	serviceURL := makeServiceURL(s.Id)
	url := fmt.Sprintf("%s/%s", serviceURL, taskPath)
	return url
}

func makeServiceURL(id string) string {
	return fmt.Sprintf("/service/%s", id)
}
