# Kubernetes Termination Delay

This tool helps during rolling updates to catch every service request. When doing a rolling update on some deployment in Kubernetes there are (short) intervals when the old Pod is not ready and the service routing has not taken this Pod out of scheduling. Even if the rolling update enforces `minUnavilable==0` the old Pod has to be removed as instance as well as routable endpoint for a particular service.

This tool is a wrapper around your service running inside a container. When there are readiness or liveness probes, `k8s-term-delay` can check them too.

If the pod gets descheduled, the wrapper catches the signal (e.g. `SIGTERM` or `SIGINT`) and immediatly marks the service as not ready. Then it will wait a configurable amount of time and send the signal to the real service. During this time Kubernetes will remove the endpoint as the pod is not ready. All requests in between will get handled properly.

Make sure the readinessProbe interval is at most half the guard time plus some additional time for endpoint reconfiguration. The default guardTime is 10s. The default grace period for a Pod is 30s. Good timings depend on each individual application and its environment.

## Usage

`k8s-term-delay guard -- $SERVICE_EXECUTABLE $SERVICE_ARGS`

## Parameters / Environment Variables

* `--liveness-url` | `KTD_LIVENESS_URL`: check this URL for livenessProbes.
* `--readiness-url` | `KTD_READINESS_URL`: check this URL for readinessProbes.
* `--health-bind` | `KTD_HEALTH_BIND`: bind to this address for healthchecks.
* `--guard-interval` | `KTD_GUARD_INTERVAL`: delay termination signal for n seconds.