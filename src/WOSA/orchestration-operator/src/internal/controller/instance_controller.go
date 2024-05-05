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
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/registry"
	"helm.sh/helm/v3/pkg/release"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	symphonyv1 "ghcr.io/pdpresson/symphony/orchestration-operator/api/v1"
	"ghcr.io/pdpresson/symphony/orchestration-operator/api/v1/models"
)

// InstanceReconciler reconciles a Instance object
type InstanceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=solution.symphony,resources=instances,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=solution.symphony,resources=instances/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=solution.symphony,resources=instances/finalizers,verbs=update
//+kubebuilder:rbac:groups=*,resources=*,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Instance object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *InstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	finalizer := "instance.solution.symphony/finalizer"

	log := log.FromContext(ctx)
	log.Info("Reconciling instance...")

	instance := &symphonyv1.Instance{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("The instance resource does not exist")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get instance")
		return ctrl.Result{}, err
	}

	if instance.Status.Properties == nil {
		instance.Status.Properties = make(map[string]string)
	}

	if instance.ObjectMeta.DeletionTimestamp.IsZero() { // update
		log.Info("Adding finalizer")
		if !controllerutil.ContainsFinalizer(instance, finalizer) {
			controllerutil.AddFinalizer(instance, finalizer)
			if err := r.Client.Update(ctx, instance); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	solution := &symphonyv1.Solution{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.Spec.Solution, Namespace: instance.Namespace}, solution)
	if err != nil && apierrors.IsNotFound(err) {
		log.Error(err, "Solution not found so we'd need to wait for it")
		return ctrl.Result{}, nil
	}

	var props models.HelmPropertyConfig
	for _, v := range solution.Spec.Components {
		log.Info("TODO: Need to handle multiple components?")
		err := json.Unmarshal(v.Properties.Raw, &props)
		if err != nil {
			log.Error(err, "Failed to unmarshal helm properties")
		}
	}

	settings := cli.New()
	actionConfig := new(action.Configuration)
	debugLog := func(format string, v ...interface{}) {
		msg := fmt.Sprintf(format, v...)
		log.Info(fmt.Sprint("Debug log: ", msg))
	}

	if err = actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), debugLog); err != nil {
		log.Error(err, "failed initializing actionConfig")
		return ctrl.Result{}, err
	}

	isToBeDeleted := instance.GetDeletionTimestamp() != nil
	if isToBeDeleted {
		log.Info(fmt.Sprintf("Instance %s to be deleted", instance.Name))
		if controllerutil.ContainsFinalizer(instance, finalizer) {
			log.Info("Performing Finalizer Operations for instance before deleting resource")
		}

		log.Info("TODO: Deleting by release name only. Need to figure out how to do it by namespace")
		uninstaller := action.NewUninstall(actionConfig)
		uninstaller.Run(props.Name)

		log.Info("Removing finalizer for instance after successfully perform the operations")
		controllerutil.RemoveFinalizer(instance, finalizer)
		if err := r.Client.Update(ctx, instance); err != nil {
			log.Error(err, "Removing finalizer failed")
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	lister := action.NewList(actionConfig)
	lister.StateMask = action.ListAll
	lister.AllNamespaces = true
	releases, err := lister.Run()
	if err != nil {
		log.Error(err, "Problem running lister")
		return ctrl.Result{}, err
	}

	idx := slices.IndexFunc(releases, func(r *release.Release) bool {
		return r.Name == props.Name && r.Namespace == instance.Spec.Scope
	})

	if idx < 0 {
		log.Info("Installing new helm chart")
		installer := action.NewInstall(actionConfig)
		installer.Namespace = instance.Spec.Scope
		installer.Wait = props.Wait
		installer.CreateNamespace = true
		installer.IsUpgrade = false
		installer.Version = props.Version
		installer.ReleaseName = props.Name
		installer.Timeout = 120 * time.Second

		rcl, err := registry.NewClient(registry.ClientOptDebug(true))
		if err != nil {
			log.Error(err, "error creating registry client")
			return ctrl.Result{}, err
		}

		installer.SetRegistryClient(rcl)
		s, err := installer.ChartPathOptions.LocateChart(props.Repo, settings)
		if err != nil {
			log.Error(err, "Error locating chart")
			return ctrl.Result{}, err
		}

		log.Info(fmt.Sprint("Chart located: ", s))
		crt, err := loader.Load(s)
		if err != nil {
			log.Error(err, "Load failed")
			return ctrl.Result{}, err
		}
		log.Info("Chart loaded")

		log.Info(fmt.Sprintf("Running installer with values: %+v", props.Values))
		runner, err := installer.Run(crt, props.Values)
		if err != nil {
			log.Error(err, "Run failed")
			return ctrl.Result{}, err
		}
		log.Info(fmt.Sprintf("Run result: %+v", runner))
	} else {
		log.Info("Helm chart already installed")
		log.Info("TODO: Handle upgrades")
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *InstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&symphonyv1.Instance{}).
		Complete(r)
}
