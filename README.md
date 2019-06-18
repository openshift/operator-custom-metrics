Current functions:
StartMetrics() - starts the metrics server over an endpoint and port. Should be configurable with a sane default

RegisterMetrics() - Registers all metrics that are given to metricsList with the metrics server. Need a way to pass a metricsList or to register individual metrics easily.

GenerateService() - Specfies port and name of servive. Should be configurable.

GenerateServiceMonitor() - Uses the service to create a servicemonitor. Should not be configurable, as it depends on Service

GenerateRoute() - Uses service to create a route. Possibly be configurable in the case of needing a specific hostname. Does NOT cover k8s objects which is needed

ConfigureMetrics() - Generates everything to get metrics server. Wrapper for GenerateService() and right now GenerateRouter(). Should take some config from somewhere to generate and create the needed objects.

createOrUpdateService() - Code that ensure that the service gets updated if the service already exists on the cluster.

createClient() - WIP function for wrapping the libaray in a client

Design Goals:
- Easily configurable. Set the settings (in-cluster or external promethes, metrics list, etc.) and call ConfigureMetrics() and metrics should come up.
- use interface/client to wrap around the functions so users don't have to call them directly
