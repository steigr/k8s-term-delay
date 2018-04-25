package runner

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func (s *Service) Guard() (err error) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)
	go func() {
		sig := <-sigs
		if sig == os.Interrupt {
			s.DelayTermination(s.OverrideSignal(sig))
		} else if (sig == syscall.SIGTERM) || (sig == syscall.SIGINT) || (sig == syscall.SIGHUP) || (sig == syscall.SIGQUIT) {
			go s.DelayTermination(s.OverrideSignal(sig))
		} else {
			s.Command.Process.Signal(sig)
		}
	}()

	return nil
}

func (s *Service) OverrideSignal(sig os.Signal) os.Signal {
	// log.Println(spew.Sdump(s.OverriddenSignal))
	if s.OverriddenSignal != nil {
		log.Println("Override signal with", s.OverriddenSignal)
		return s.OverriddenSignal
	}
	return sig
}

func (s *Service) SetOverrideSignal(sig string) {
	if strings.Compare("INT", sig) == 0 {
		s.OverriddenSignal = syscall.SIGINT
	}
	if strings.Compare("QUIT", sig) == 0 {
		s.OverriddenSignal = syscall.SIGQUIT
	}
	if strings.Compare("TERM", sig) == 0 {
		s.OverriddenSignal = syscall.SIGTERM
	}
	if strings.Compare("HUP", sig) == 0 {
		s.OverriddenSignal = syscall.SIGHUP
	}
	if strings.Compare("KILL", sig) == 0 {
		s.OverriddenSignal = syscall.SIGKILL
	}
}

func (s *Service) DelayTermination(sig os.Signal) (err error) {
	s.GuardActive = true
	log.Println("delaying signal", sig, "for", s.GuardInterval, "seconds")
	time.Sleep(s.GuardInterval)
	log.Println("Send", sig, "to child")
	s.Command.Process.Signal(sig)
	s.GracePeriodOver <- true
	return nil
}

func (s *Service) SetGuardInterval(interval int) {
	s.GuardInterval = time.Second * time.Duration(interval)
}
