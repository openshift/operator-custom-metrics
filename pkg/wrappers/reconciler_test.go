package wrappers

import (
	"fmt"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	pcm "github.com/prometheus/client_model/go"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestTimedReconciler(t *testing.T) {
	const (
		controllerName = "the-controller"
	)

	rt := NewReconcileTimer("op-name")
	tr := &TimedReconciler{
		WrappedReconciler: &testReconciler{},
		ControllerName:    controllerName,
		Logger:            log.Log,
		Metric:            rt,
	}

	res, err := tr.Reconcile(reconcile.Request{NamespacedName: types.NamespacedName{Name: "foo", Namespace: "bar"}})

	if !res.Requeue {
		t.Fatal("Expected requeue!")
	}
	if err == nil {
		t.Fatal("Expected an error!")
	}
	expectedErr := "Failed to reconcile bar/foo"
	if err.Error() != expectedErr {
		t.Fatalf("Wrong error message\nExpected: %q\nActual:   %q", expectedErr, err.Error())
	}

	observer, err := rt.GetMetricWith(prometheus.Labels{"controller": controllerName})
	if err != nil {
		t.Fatalf("Failed to get metric: %s", err.Error())
	}
	hg := observer.(prometheus.Histogram)

	actualDesc := hg.Desc().String()
	expectedDesc := `Desc{fqName: "op_name_reconcile_duration_seconds", help: "Distribution of the number of seconds a Reconcile takes, broken down by controller", constLabels: {name="op-name"}, variableLabels: [controller]}`
	if actualDesc != expectedDesc {
		t.Fatalf("Wrong description\nExpected: %q\nActual:   %q", expectedDesc, actualDesc)
	}

	metric := &pcm.Metric{}
	hg.Write(metric)
	if *metric.Histogram.SampleCount != 1 {
		t.Fatalf("Expected one count but got %d", *metric.Histogram.SampleCount)
	}
	actualSeconds := metric.Histogram.GetSampleSum()
	if actualSeconds < (time.Millisecond * 5).Seconds() {
		t.Fatalf("Expected reconcile to take at least 5ms but it took %f", actualSeconds)
	}
}

type testReconciler struct{}

func (tr *testReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	// We can't rely on the Reconcile taking an exact amount of time; just make sure it's not zero
	time.Sleep(time.Millisecond * 5)
	// Deliberately contrived returns, just to prove the arguments make the full round trip
	return reconcile.Result{Requeue: true}, fmt.Errorf("Failed to reconcile %s/%s", req.Namespace, req.Name)
}
