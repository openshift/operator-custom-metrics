# Operator Custom Metrics
This library assists in collecting operator custom-metrics from operators and exposing them to a prometheus instance. The user can define their own set of metrics for a specific operator, and the library would help in integrating these custom metrics with a prometheus instance for monitoring purposes.

## Description
This library is designed to provide the user with an easy procedure to register custom metrics with a running prometheus instance. The user can provide the **metrics configuration** based on which the metrics are collected. The library currently supports openshift-operators and hence, to access the metrics routes are created. 
The library simplifies the following processes:
1. Registering: Register custom metrics provided by the user to prometheus.
2. Publishing: Publish metrics at a specified endpoint and path.
3. Creating resources:
    1. Create a service resource. 
    2. Create an external [route](https://docs.openshift.com/container-platform/3.9/architecture/networking/routes.html) so that the metrics can be accessed. 

The generated service is created with a label selector which targets pods with the labels `name: serviceName` where
`serviceName` is the name of the service passed to the configuration builder during initialization.

The following are the parameters of the metrics configuration (`metricsConfig`) for the custom metrics:

- Metrics Path:
This is the endpoint at which the metrics would be exposed. The default metrics endpoint is `/customMetrics`.

- Metrics Port:
This is the port at which the metrics would be published. The default metrics port is `:8089`

- List of collectors:
The user can provide the collectors which are to be registered with prometheus. If no collectors are passed, the default prometheus server metrics are exposed.

- Record metrics function:
This is a user-defined function which describes the process in which the custom metrics are to be collected. This can be passed to the library or can be executed within the operator code based on the desired implementation.

## Prerequisites
The library can be integrated by downloading the same, using the following command:

```bash
go get -d github.com/openshift/operator-custom-metrics
```

## Using operator-custom metrics library

The following functions of the library can be called by the user to create a metrics configuration which is to be passed to the library:
1. `NewBuilder(namespace, serviceName)` - Sets the parameters of the metricsConfig Object to default values.
2. `WithPort(port string)` - Updates the default value of port in the metricsConfig object.
3. `WithPath(path string)` - Updates the default value of path in the metricsConfig object.
4. `WithCollector(collector prometheus.Collector)` - Creates a list of prometheus collectors which are to be registered.
5. `WithCollectors([]prometheus.Collector)` - Updates the list of collectors in the metricsConfig object.
6. `GetConfig()` - Returns the reference to metricsConfig which is built with the configuration set by the user.
