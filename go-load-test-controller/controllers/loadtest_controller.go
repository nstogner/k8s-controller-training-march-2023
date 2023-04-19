/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	platformv1 "github.com/nstogner/k8s-controller-training-march-2023/go-load-test-controller/api/v1"
	"github.com/nstogner/k8s-controller-training-march-2023/go-load-test-controller/internal/loadtest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LoadTestReconciler reconciles a LoadTest object
type LoadTestReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	Runner *loadtest.Runner
}

//+kubebuilder:rbac:groups=platform.mycompany.com,resources=loadtests,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=platform.mycompany.com,resources=loadtests/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=platform.mycompany.com,resources=loadtests/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LoadTest object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *LoadTestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	defer log.Info("Done Reconciling LoadTest")

	var lt platformv1.LoadTest
	if err := r.Client.Get(ctx, req.NamespacedName, &lt); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Reconciling LoadTest", "address", lt.Spec.Address)

	if lt.Status.Phase != "" {
		log.Info("LoadTest already run, ignoring")
		return ctrl.Result{}, nil
	}

	fmt.Println("1. -----pre-status-update-1-resource-version: ", lt.ObjectMeta.ResourceVersion)

	log.Info("Updating LoadTest status to Running")
	lt.Status.Phase = platformv1.PhaseRunning
	lt.Status.StartTime = metav1.Time{Time: time.Now()}
	if err := r.Client.Status().Update(ctx, &lt); err != nil {
		// NOTE: This pre-run update will fail if the .Get() above
		// read from a stale cache. This protects against re-running
		// a load test because of a stale cache read.
		return ctrl.Result{}, fmt.Errorf("updating status: %w", err)
	}

	fmt.Println("2. -----pre-run-resource-version: ", lt.ObjectMeta.ResourceVersion)

	log.Info("Running LoadTest")
	out := r.Runner.Run(ctx, loadtest.Input{
		ID:        string(lt.UID),
		URL:       lt.Spec.Address,
		Method:    lt.Spec.Method,
		Duration:  lt.Spec.Duration.Duration,
		ReqPerSec: 10,
	})

	lt.Status.Phase = platformv1.PhaseCompleted
	lt.Status.RequestCount = out.RequestCount
	lt.Status.SuccessCount = out.SuccessCount
	lt.Status.EndTime = metav1.Time{Time: time.Now()}

	if err := r.Client.Status().Update(ctx, &lt); err != nil {
		return ctrl.Result{}, fmt.Errorf("updating status: %w", err)
	}

	fmt.Println("3.-----post-run--resource-version: ", lt.ObjectMeta.ResourceVersion)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LoadTestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&platformv1.LoadTest{}).
		Complete(r)
}
