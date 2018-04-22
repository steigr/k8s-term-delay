package health

import (
	"log"
	"net/http"
	"time"

	"github.com/heptiolabs/healthcheck"
)

const (
	HTTP_DEFAULT_TIMEOUT = time.Second * 5
)

type Check struct {
	Handler healthcheck.Handler
	Bind    string
	Name    string
}

func New(service, bind string, guardCheck *bool) (check *Check) {
	check = &Check{}
	check.Bind = bind
	check.Handler = healthcheck.NewHandler()
    check.Handler.AddReadinessCheck("guard",check.GuardingCheck(guardCheck))
	return check
}

func (s *Check) SetReadinessUrl(url string) (err error) {
	log.Println("Pass readinessProbe to", url, "unless terminating")
	s.Handler.AddLivenessCheck(s.Name, healthcheck.HTTPGetCheck(url, HTTP_DEFAULT_TIMEOUT))
	return nil
}

func (s *Check) SetLivenessUrl(url string) (err error) {
	log.Println("Pass livenessProbe to", url)
	s.Handler.AddReadinessCheck(s.Name, healthcheck.HTTPGetCheck(url, HTTP_DEFAULT_TIMEOUT))
	return nil
}

func (s *Check) Run() {
	log.Println("Answering healthchecks at", s.Bind)
	go http.ListenAndServe(s.Bind, s.Handler)
}
