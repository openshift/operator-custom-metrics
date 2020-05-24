module github.com/openshift/operator-custom-metrics

go 1.13

require (
	github.com/coreos/prometheus-operator v0.38.0
	github.com/openshift/api v3.9.1-0.20190424152011-77b8897ec79a+incompatible
	github.com/operator-framework/operator-sdk v0.17.0
	github.com/prometheus/client_golang v1.5.1
	github.com/sirupsen/logrus v1.5.0
	k8s.io/api v0.17.4
	k8s.io/apimachinery v0.17.4
	sigs.k8s.io/controller-runtime v0.5.2
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM
	k8s.io/client-go => k8s.io/client-go v0.17.4 // Required by prometheus-operator
)
