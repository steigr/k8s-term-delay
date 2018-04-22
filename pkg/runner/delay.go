package runner

import (
	"log"
	"os"
    "time"
	"os/signal"
	"syscall"
)

func (s *Service) Guard() (err error) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)
	go func() {
		sig := <-sigs
		if ( sig == os.Interrupt) {
           s.DelayTermination(sig)
        } else if ( sig == syscall.SIGTERM ) || ( sig == syscall.SIGINT  ) || ( sig == syscall.SIGHUP  ) || ( sig == syscall.SIGQUIT ) { 
			go s.DelayTermination(sig)
        } else {
			s.Command.Process.Signal(sig)
		}
	}()

	return nil
}

func (s *Service) DelayTermination(sig os.Signal) (err error) {
    s.GuardActive = true
    log.Println("delaying for",s.GuardInterval,"seconds")
    time.Sleep(s.GuardInterval)
    s.Command.Process.Signal(sig)
    s.GracePeriodOver <- true
    return nil
}

func (s *Service) SetGuardInterval(interval int) {
	s.GuardInterval = time.Second * time.Duration(interval)
}
