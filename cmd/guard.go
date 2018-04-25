package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/steigr/k8s-term-delay/pkg/health"
	"github.com/steigr/k8s-term-delay/pkg/runner"
)

const (
	DEFAULT_KTD_GUARD_INTERVAL = 10
)

// guardCmd represents the guard command
var (
	guardCmd = &cobra.Command{
		Use: "guard",
		Run: func(cmd *cobra.Command, args []string) {
			r := runner.New(args)
			if guardInterval == 30 {
				gi := viper.GetInt("guard-interval")
				if gi != 0 {
					guardInterval = gi
				}
			}
			r.SetGuardInterval(guardInterval)
			if strings.Compare(healthBind, "[::]:8081") == 0 {
				hb := viper.GetString("health-bind")
				if hb != "" {
					healthBind = hb
				}
			}
			h := health.New(r.Name, healthBind, &r.GuardActive)
			if strings.Compare(livenessUrl, "") == 0 {
				lU := viper.GetString("liveness-url")
				if lU != "" {
					livenessUrl = lU
				}
			}
			if len(livenessUrl) > 0 {
				h.SetLivenessUrl(livenessUrl)
			}
			if strings.Compare(readinessUrl, "") == 0 {
				rU := viper.GetString("readiness-url")
				if rU != "" {
					readinessUrl = rU
				}
			}
			if len(readinessUrl) > 0 {
				h.SetReadinessUrl(readinessUrl)
			}
			if strings.Compare(readinessUrl, "") == 0 {
				rU := viper.GetString("readiness-url")
				if rU != "" {
					readinessUrl = rU
				}
			}
			if len(readinessUrl) > 0 {
				h.SetReadinessUrl(readinessUrl)
			}
			if len(overrideSignal) > 0 {
				r.SetOverrideSignal(overrideSignal)
			}
			if len(viper.GetString("override-signal")) > 0 {
				r.SetOverrideSignal(viper.GetString("override-signal"))
			}
			r.Guard()
			h.Run()
			gracePeriodOver := make(chan bool)
			r.Run(gracePeriodOver)
			<-gracePeriodOver
		},
	}
	livenessUrl, readinessUrl, healthBind, overrideSignal string
	guardInterval                                         int
)

func init() {
	guardCmd.Flags().StringVar(&healthBind, "health-bind", "[::]:8081", "healthcheck bind address")
	guardCmd.Flags().StringVar(&livenessUrl, "liveness-url", "", "actual Livecheck URL")
	guardCmd.Flags().StringVar(&readinessUrl, "readiness-url", "", "actual Readiness URL")
	guardCmd.Flags().StringVar(&overrideSignal, "override-signal", "", "substitude INT/TERM/QUIT/HUP with this signal")
	guardCmd.Flags().IntVar(&guardInterval, "guard-interval", DEFAULT_KTD_GUARD_INTERVAL, "delay until SIGTERM or SIGINT are forwarded to service, the guard claims to be not ready within this interval")
	rootCmd.AddCommand(guardCmd)
}
