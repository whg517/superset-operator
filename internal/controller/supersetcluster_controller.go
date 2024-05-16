/*
Copyright 2024 zncdata-labs.

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

package controller

import (
	"context"

	"github.com/zncdata-labs/superset-operator/internal/controller/cluster"
	"github.com/zncdata-labs/superset-operator/pkg/reconciler"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	supersetv1alpha1 "github.com/zncdata-labs/superset-operator/api/v1alpha1"
)

// SupersetClusterReconciler reconciles a SupersetCluster object
type SupersetClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var (
	logger = ctrl.Log.WithName("common").WithName("reconciler")
)

// +kubebuilder:rbac:groups=superset.zncdata.dev,resources=supersetclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=superset.zncdata.dev,resources=supersetclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=superset.zncdata.dev,resources=supersetclusters/finalizers,verbs=update

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *SupersetClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	logger.V(0).Info("Reconciling SupersetCluster")

	instance := &supersetv1alpha1.SupersetCluster{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			logger.V(1).Info("SupersetCluster resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	resourceClient := reconciler.ResourceClient{
		Client: r.Client,
	}

	clusterRreconciler := cluster.NewReconciler(resourceClient, instance)

	if err := clusterRreconciler.RegisterResources(ctx); err != nil {
		return ctrl.Result{}, err
	}

	if result := clusterRreconciler.Reconcile(); result.RequeueOrNot() {
		return result.Result()
	}

	if result := clusterRreconciler.Ready(); result.RequeueOrNot() {
		return result.Result()
	}

	logger.V(0).Info("Reconcile completed")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SupersetClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&supersetv1alpha1.SupersetCluster{}).
		Complete(r)
}
