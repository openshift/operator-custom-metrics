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
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// StartMetrics starts the server based on the metricsConfig provided by the user.
func StartMetrics(config metricsConfig) error {
	if config.metricsRegisterer == nil {
		config.metricsRegisterer = prometheus.DefaultRegisterer
		config.metricsGatherer = prometheus.DefaultGatherer
	}
	err := RegisterMetrics(config.metricsRegisterer, config.collectorList)
	if err != nil {
		return err
	}
	metricsHandler := promhttp.InstrumentMetricHandler(
		config.metricsRegisterer, promhttp.HandlerFor(config.metricsGatherer, promhttp.HandlerOpts{}),
	)
	http.Handle(config.metricsPath, metricsHandler)
	log.Info(fmt.Sprintf("Port: %s", config.metricsPort))
	metricsPort := ":" + (config.metricsPort)
	go http.ListenAndServe(metricsPort, nil)
	return nil
}

// RegisterMetrics takes the list of metrics to be registered from the user and
// registers to prometheus.
func RegisterMetrics(metricsRegisterer prometheus.Registerer, list []prometheus.Collector) error {
	for _, metric := range list {
		err := metricsRegisterer.Register(metric)
		if err != nil {
			return err
		}
	}
	return nil
}
