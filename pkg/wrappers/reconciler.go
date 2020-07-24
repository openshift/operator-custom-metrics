package wrappers

/*

# Time your Reconcile.

## Usage

### Create the metric
reconcileTimer := NewReconcileTimer("wizbang-operator")

### Include the metric in your builder
metricsServer := metrics.NewBuilder().WithCollector(reconcileTimer)

### Wrap your Reconciler implementation when bootstrapping your controller
reconciler := &TimedReconciler{
	WrappedReconciler: &ReconcileWiz{ ... },
	Metric:            reconcileTimer,
	ControllerName:    "wiz-controller",
	Logger:            logf.Log...,
}

### Profit
\o/

*/

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// NewReconcileTimer produces a HistogramVec Collector that records the duration of a Reconcile,
// labeled by controller.
// The Collector's name is {prefix}_reconcile_duration_seconds, where the {prefix} is derived from
// the operatorName parameter.
func NewReconcileTimer(operatorName string) *prometheus.HistogramVec {
	name := fmt.Sprintf("%s_reconcile_duration_seconds", strings.ReplaceAll(operatorName, "-", "_"))
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        name,
		Help:        "Distribution of the number of seconds a Reconcile takes, broken down by controller",
		ConstLabels: prometheus.Labels{"name": operatorName},
		Buckets:     []float64{0.001, 0.01, 0.1, 1, 5, 10, 20},
	}, []string{"controller"})
}

// TimedReconciler is a wrapper that times how long Reconcile takes and reports it through the
// provided `metric`. Use the result of `NewReconcileTimer` in the `Metric` field.
type TimedReconciler struct {
	WrappedReconciler reconcile.Reconciler
	Metric            *prometheus.HistogramVec
	ControllerName    string
	Logger            logr.Logger
}

// Reconcile implements the Reconciler interface, timing the wrapped Reconciler's Reconcile call
// and reporting that metric. Also logs bracketing messages, where the closing message includes
// the duration of the wrapped Reconcile.
func (r *TimedReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	logger := r.Logger.
		WithValues(
			"Controller", r.ControllerName,
			"Request.Namespace", req.Namespace,
			"Request.Name", req.Name,
		)
	logger.Info("Reconciling")

	start := time.Now()
	result, err := r.WrappedReconciler.Reconcile(req)
	dur := time.Since(start)
	r.Metric.WithLabelValues(r.ControllerName).Observe(dur.Seconds())

	logger.WithValues("Duration", dur).Info("Reconcile complete")
	return result, err
}
