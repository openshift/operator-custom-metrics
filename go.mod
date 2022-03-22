module github.com/openshift/operator-custom-metrics

go 1.16

require (
	github.com/openshift/api v0.0.0-20220124143425-d74727069f6f
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.55.0
	github.com/prometheus/client_golang v1.12.1
	github.com/sirupsen/logrus v1.8.1
	k8s.io/api v0.23.5
	k8s.io/apimachinery v0.23.5
	sigs.k8s.io/controller-runtime v0.11.1
)
