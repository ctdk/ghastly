package ghastly

import (
	"fmt"
)

// Purge a URL from the CDN.
func (g *Ghastly) PurgeURL(url string) (string, error) {
	resp, err := g.Purge(url)
	if err != nil {
		return "", err
	}
	pData, err := ParseJson(resp.Body)
	if err != nil {
		return "", err
	}
	if pData["status"].(string) != "ok" {
		err = fmt.Errorf("Status was not ok with purging '%s'. The content of the reply was %v.", url, pData)
		return "", err
	}
	return pData["id"].(string), nil
}

// Purge everything from a service.
func (s *Service) PurgeAll() error {
	purl := s.TaskURL("purge_all")
	resp, err := s.ghastly.Post(purl, "application/json", nil)
	if err != nil {
		return err
	}
	pData, err := ParseJson(resp.Body)
	if err != nil {
		return err
	}
	if pData["status"].(string) != "ok" {
		err = fmt.Errorf("Status was not ok with purging all items from service %s. The content of the reply was %v.", s.Name, pData)
		return err
	}
	return nil
}

// Purge a service of items tagged with a particular key.
func (s *Service) PurgeKey(key string) error {
	pkey := fmt.Sprintf("purge/%s", key)
	purl := s.TaskURL(pkey)
	resp, err := s.ghastly.Post(purl, "application/json", nil)
	if err != nil {
		return err
	}
	pData, err := ParseJson(resp.Body)
	if err != nil {
		return err
	}
	if pData["status"].(string) != "ok" {
		err = fmt.Errorf("Status was not ok with purging items keyed with %s from service %s. The content of the reply was %v.", key, s.Name, pData)
		return err
	}
	return nil
}
