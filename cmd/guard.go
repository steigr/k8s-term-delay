package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/steigr/k8s-term-delay/pkg/health"
	"github.com/steigr/k8s-term-delay/pkg/runner"
	"strings"
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
			r.Guard()
			h.Run()
			gracePeriodOver := make(chan bool)
			r.Run(gracePeriodOver)
			<-gracePeriodOver
		},
	}
	livenessUrl, readinessUrl, healthBind string
	guardInterval                         int
)

func init() {
	// livenessUrl = viper.GetString("liveness")
	// log.Println("Got", livenessUrl, "as livenessUrl")
	// readinessUrl = viper.GetString("readiness-url")
	// healthBind = viper.GetString("health-bind")
	guardCmd.Flags().StringVar(&healthBind, "health-bind", "[::]:8081", "healthcheck bind address")
	guardCmd.Flags().StringVar(&livenessUrl, "liveness-url", "", "actual Livecheck URL")
	guardCmd.Flags().StringVar(&readinessUrl, "readiness-url", "", "actual Readiness URL")
	guardCmd.Flags().IntVar(&guardInterval, "guard-interval", DEFAULT_KTD_GUARD_INTERVAL, "delay until SIGTERM or SIGINT are forwarded to service, the guard claims to be not ready within this interval")
	rootCmd.AddCommand(guardCmd)
}
