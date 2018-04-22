package runner

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
    "time"
)

type Service struct {
	Command         *exec.Cmd
	Name            string
	GuardInterval   time.Duration
    GracePeriodOver chan bool
    GuardActive     bool
}

func New(arguments []string) (service *Service) {
	service = &Service{}
	service.Command = &exec.Cmd{}
	service.Name = filepath.Base(arguments[0])
    service.GuardActive = false

	if service.Name == arguments[0] {
		if lp, err := exec.LookPath(arguments[0]); err == nil {
			service.Command.Path = lp
		}
	}

	service.Command.Args = arguments

	return service
}

func (s *Service) Run(gracePeriodOver chan bool) {
    s.GracePeriodOver = gracePeriodOver
	s.Command.Env = os.Environ()
	s.Command.Stdout = os.Stdout
	s.Command.Stderr = os.Stderr
    s.Command.Stdin = os.Stdin

    waitCh := make(chan error, 1)

	if err := s.Command.Start(); err != nil {
		log.Fatal(err)
	}

	go func() {
		waitCh <- s.Command.Wait()
		close(waitCh)
	}()
	for {
		select {
		case err := <-waitCh:
			var waitStatus syscall.WaitStatus
			if exitError, ok := err.(*exec.ExitError); ok {
				waitStatus = exitError.Sys().(syscall.WaitStatus)
				os.Exit(waitStatus.ExitStatus())
			}
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}
}
