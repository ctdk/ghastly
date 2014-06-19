package ghastly

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

var G *Ghastly
var S *Service

func makeServiceName() string {
	return fmt.Sprintf("ghastly-test-%d-%d", os.Getpid(), rand.Int31())
}

func TestBadGhastlyLogin(t *testing.T) {
	// no username/pass
	login_opts := make(map[string]string)
	_, err := New(login_opts)
	if err == nil {
		t.Errorf("Logging into fastly unexpectedly succeeded with no username or password.")
	}
	// invalid username/pass
	login_opts["user"] = "invalid@example.com"
	login_opts["password"] = "12345"
	_, err = New(login_opts)
	if err == nil {
		t.Errorf("Logging into fastly unexpectedly succeeded with bad username and password.")
	}
}

func TestGhastlyLogin(t *testing.T) {
	login_opts := make(map[string]string)
	login_opts["user"] = os.Getenv("FASTLY_TEST_USER")
	login_opts["password"] = os.Getenv("FASTLY_TEST_PASSWORD")
	_, err := New(login_opts)
	if err != nil {
		t.Errorf("Error logging into fastly: %s", err.Error())
	}
}

func TestService(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	login_opts := make(map[string]string)
	login_opts["user"] = os.Getenv("FASTLY_TEST_USER")
	login_opts["password"] = os.Getenv("FASTLY_TEST_PASSWORD")
	g, err := New(login_opts)
	if err != nil {
		t.Errorf("Error logging into fastly: %s", err.Error())
	}
	serviceName := makeServiceName()
	s, err := g.NewService(serviceName)
	if err != nil {
		t.Errorf(err.Error())
	}
	s2, err := g.GetService(s.Id)
	if err != nil {
		t.Errorf(err.Error())
	}
	if s.Id != s2.Id {
		t.Errorf("Somehow the test service %s was re-fetched, and had a mismatched id: %s vs. %s", s.Name, s.Id, s2.Id)
	}
	err = s.Delete()
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestServiceDetail(t *testing.T) {
	login_opts := make(map[string]string)
	login_opts["user"] = os.Getenv("FASTLY_TEST_USER")
	login_opts["password"] = os.Getenv("FASTLY_TEST_PASSWORD")
	g, err := New(login_opts)
	if err != nil {
		t.Errorf("Error logging into fastly: %s", err.Error())
	}
	serviceName := makeServiceName()
	s, err := g.NewService(serviceName)
	defer s.Delete()
	if err != nil {
		t.Errorf(err.Error())
	}
	s2, err := g.GetService(s.Id)
	if err != nil {
		t.Errorf(err.Error())
	}
	s3, err := s2.Details()
	if err != nil {
		t.Errorf(err.Error())
	}
	if s2.Id != s3.Id {
		t.Errorf("The normal and detailed services were not the same service.")
	}
	if s2.CreatedAt == s3.CreatedAt {
		t.Errorf("The detailed service does not seem to actually have different details.")
	}
}

func TestListServices(t *testing.T) {
	login_opts := make(map[string]string)
	login_opts["user"] = os.Getenv("FASTLY_TEST_USER")
	login_opts["password"] = os.Getenv("FASTLY_TEST_PASSWORD")
	g, err := New(login_opts)
	if err != nil {
		t.Errorf("Error logging into fastly: %s", err.Error())
	}
	serviceName := makeServiceName()
	serviceName2 := makeServiceName()
	g1, err := g.NewService(serviceName)
	g2, err := g.NewService(serviceName2)
	defer g1.Delete()
	defer g2.Delete()
	services, err := g.ListServices()
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(services) == 0 {
		t.Errorf("Expected services to be listed, but got zero back.")
	}
}

func TestSearchServices(t *testing.T) {
	login_opts := make(map[string]string)
	login_opts["user"] = os.Getenv("FASTLY_TEST_USER")
	login_opts["password"] = os.Getenv("FASTLY_TEST_PASSWORD")
	g, err := New(login_opts)
	if err != nil {
		t.Errorf("Error logging into fastly: %s", err.Error())
	}
	serviceName := makeServiceName()
	//serviceName2 := makeServiceName()
	s1, err := g.NewService(serviceName)
	defer s1.Delete()
	sought, err := g.SearchServices(s1.Name)
	if err != nil {
		t.Errorf(err.Error())
	}
	if sought.Id != s1.Id {
		t.Errorf("Searching for a service named %s returned id %s, but expected %s", s1.Name, sought.Id, s1.Id)
	}
	searched, err := g.SearchServices("omg-totally-fake")
	if err == nil {
		t.Errorf("Searching for 'omg-totally-fake' should have failed, but unexpectedly succeeded, returning id %s.", searched.Id)
	}
}

func TestUpdateService(t *testing.T) {
	login_opts := make(map[string]string)
	login_opts["user"] = os.Getenv("FASTLY_TEST_USER")
	login_opts["password"] = os.Getenv("FASTLY_TEST_PASSWORD")
	g, err := New(login_opts)
	if err != nil {
		t.Errorf("Error logging into fastly: %s", err.Error())
	}
	serviceName := makeServiceName()
	serviceName2 := makeServiceName()
	s1, err := g.NewService(serviceName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer s1.Delete()
	s2, _ := g.GetService(s1.Id)
	params := map[string]string{"name": serviceName2}
	err = s2.Update(params)
	if err != nil {
		t.Errorf(err.Error())
	}
	if s2.Name != serviceName2 {
		t.Errorf("Service name should have been %s, but was %s instead", serviceName2, s2.Name)
	}
	s3, _ := g.GetService(s1.Id)
	if s3.Name == s1.Name {
		t.Errorf("Service name did not update at source, expected %s, got %s", serviceName2, s3.Name)
	}
}

// one service from here on out
func TestSetupService(t *testing.T) {
	login_opts := make(map[string]string)
	login_opts["user"] = os.Getenv("FASTLY_TEST_USER")
	login_opts["password"] = os.Getenv("FASTLY_TEST_PASSWORD")
	var err error
	G, err = New(login_opts)
	if err != nil {
		t.Errorf("Error logging into fastly: %s", err.Error())
	}
	serviceName := makeServiceName()
	S, err = G.NewService(serviceName)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestVersion(t *testing.T) {
	v, err := S.NewVersion()
	if err != nil {
		t.Errorf(err.Error())
	}
	_, err = S.GetVersion(v.Number)
	if err != nil {
		t.Errorf(err.Error())
	}
	_, err = v.Clone()
	if err != nil {
		t.Errorf(err.Error())
	}

}

func TestDomain(t *testing.T) {
	v, err := S.NewVersion()
	if err != nil {
		t.Errorf(err.Error())
	}
	domainName := "oiweruaklsjfas.com"
	domainParams := map[string]string{"name": domainName}
	d, err := v.NewDomain(domainParams)
	if err != nil {
		t.Errorf(err.Error())
	}
	if d.Name != domainName {
		t.Errorf("Created domain name did not match: expected %s, got %s", domainName, d.Name)
	}
	if d.Version != v.Number {
		t.Errorf("Created domain version did not match, expected %d, got %d", v.Number, d.Version)
	}
}

func TestPurge(t *testing.T) {
	pid, err := G.PurgeURL("http://localhost/img.png")
	if err != nil {
		t.Errorf(err.Error())
	}
	if pid == "" {
		t.Errorf("Purge id after purging was unexpectedly nil")
	}

	// Waiting for being able to activate versions for this
	/*
		err = S.PurgeAll()
		if err != nil {
			t.Errorf(err.Error())
		}
	*/
	err = S.PurgeKey("uhok")
	if err != nil {
		t.Errorf(err.Error())
	}
}

// post-test cleanup
func TestCleanup(t *testing.T) {
	S.Delete()
}
