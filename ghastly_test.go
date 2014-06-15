package ghastly

import (
	"testing"
	"os"
	"fmt"
	"math/rand"
	"time"
)

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

func TestVersion(t *testing.T) {
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
	v, err := s.NewVersion()
	if err != nil {
		t.Errorf(err.Error())
	}
	_, err = s.GetVersion(v.Number)
	if err != nil {
		t.Errorf(err.Error())
	}
	_, err = v.Clone()
	if err != nil {
		t.Errorf(err.Error())
	}

}
