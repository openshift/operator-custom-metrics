// Copyright 2019 RedHat
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

// useTLS is a boolean value, which indicates wheather https connection
// can be established based on the parameters provided by the user.
var useTLS bool

// StartMetrics starts the server based on the metricsConfig provided by the user.
// If the path to TLC certificate and key are provided by the user, establish a
// secure https connection, or else fall back to non-https
func StartMetrics(config metricsConfig) {
	// Register metrics only when the metric list is provided by the operator
	if config.collectorList != nil {
		RegisterMetrics(config.collectorList)
	}

	if config.tlsCertPath != "" && config.tlsKeyPath == "" || config.tlsCertPath == "" && config.tlsKeyPath != "" {
		log.Info("both --tls-key and --tls-crt must be provided for TLS to be enabled, falling back to non-https")
	} else if config.tlsCertPath == "" && config.tlsKeyPath == "" {
		log.Info("TLS keys not set, using non-https for metrics")
	} else {
		log.Info("TLS keys set, using https for metrics")
		useTLS = true
	}

	http.Handle(config.metricsPath, prometheus.Handler())
	log.Info(fmt.Sprintf("Port: %s", config.metricsPort))
	metricsPort := ":" + (config.metricsPort)

	if useTLS {
		go func() {
			err := http.ListenAndServeTLS(metricsPort, config.tlsCertPath, config.tlsKeyPath, nil)
			if err != nil {
				log.Info("Metrics (https) serving failed: %v", err)
			}
		}()
	} else {
		go func() {
			err := http.ListenAndServe(metricsPort, nil)
			if err != nil {
				log.Info("Metrics (http) serving failed: %v", err)
			}
		}()
	}
}

// RegisterMetrics takes the list of metrics to be registered from the user and
// registeres to prometheus.
func RegisterMetrics(list []prometheus.Collector) error {
	for _, metric := range list {
		err := prometheus.Register(metric)
		if err != nil {
			return err
		}
	}
	return nil
}
