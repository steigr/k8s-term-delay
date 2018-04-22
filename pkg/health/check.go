package health

import (
    "fmt"
    
    "github.com/heptiolabs/healthcheck"
)

func (s *Check) GuardingCheck(guard *bool) healthcheck.Check {
    return func() error {
        if *guard {
            return fmt.Errorf("readiness guard is active.")
        }
        return nil
    }
}