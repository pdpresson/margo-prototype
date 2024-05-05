/*
Copyright 2024.

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
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	symphonyv1 "ghcr.io/pdpresson/symphony/orchestration-operator/api/v1"
)

// SolutionReconciler reconciles a Solution object
type SolutionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=solution.symphony,resources=solutions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=solution.symphony,resources=solutions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=solution.symphony,resources=solutions/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Solution object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *SolutionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	finalizer := "solution.solution.symphony/finalizer"

	log := log.FromContext(ctx)
	log.Info("Reconciling solution...")

	solution := &symphonyv1.Solution{}
	err := r.Get(ctx, req.NamespacedName, solution)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("The solution resource does not exist")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get solution")
		return ctrl.Result{}, err
	}

	if solution.ObjectMeta.DeletionTimestamp.IsZero() { // update
		log.Info("Adding finalizer")
		if !controllerutil.ContainsFinalizer(solution, finalizer) {
			controllerutil.AddFinalizer(solution, finalizer)
			if err := r.Client.Update(ctx, solution); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	isToBeDeleted := solution.GetDeletionTimestamp() != nil
	if isToBeDeleted {
		log.Info(fmt.Sprintf("solution %s to be deleted", solution.Name))
		if controllerutil.ContainsFinalizer(solution, finalizer) {
			log.Info("Performing Finalizer Operations for solution before deleting resource")
		}

		log.Info("TODO: Remove the helm chart that was installed")

		log.Info("Removing finalizer for solution after successfully perform the operations")
		controllerutil.RemoveFinalizer(solution, finalizer)
		if err := r.Client.Update(ctx, solution); err != nil {
			log.Error(err, "Removing finalizer failed")
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	log.Info(fmt.Sprintf("solution details: %v", solution))

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SolutionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&symphonyv1.Solution{}).
		Complete(r)
}
