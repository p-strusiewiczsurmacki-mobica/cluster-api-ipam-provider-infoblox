/*
Copyright 2023 Deutsche Telekom AG.

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
	"net/netip"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/telekom/cluster-api-ipam-provider-infoblox/api/v1alpha1"
	"github.com/telekom/cluster-api-ipam-provider-infoblox/pkg/infoblox"
)

const (
	// ProtectPoolFinalizer is used to prevent deletion of a Pool object while its addresses have not been deleted.
	ProtectPoolFinalizer = "ipam.cluster.x-k8s.io/ProtectPool"
)

// InfobloxIPPoolReconciler reconciles a InfobloxIPPool object.
type InfobloxIPPoolReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	// operatorNamespace     string
	NewInfobloxClientFunc func(config infoblox.Config) (infoblox.Client, error)
}

//+kubebuilder:rbac:groups=ipam.cluster.x-k8s.io,resources=infobloxippools,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ipam.cluster.x-k8s.io,resources=infobloxippools/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ipam.cluster.x-k8s.io,resources=infobloxippools/finalizers,verbs=update

// SetupWithManager sets up the controller with the Manager.
func (r *InfobloxIPPoolReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Uncomment the following line adding a pointer to an instance of the controlled resource as an argument
		For(&v1alpha1.InfobloxIPPool{}).
		Complete(r)
}

// Reconcile an InfobloxIPPool.
func (r *InfobloxIPPoolReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, reterr error) {
	_ = log.FromContext(ctx)

	pool := &v1alpha1.InfobloxIPPool{}
	if err := r.Client.Get(ctx, req.NamespacedName, pool); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	patchHelper, err := patch.NewHelper(pool, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}
	defer func() {
		if err := patchHelper.Patch(ctx, pool, patch.WithOwnedConditions{}); err != nil {
			reterr = kerrors.NewAggregate([]error{reterr, err})
		}
	}()

	return ctrl.Result{}, r.reconcile(ctx, pool)
}

func (r *InfobloxIPPoolReconciler) reconcile(ctx context.Context, pool *v1alpha1.InfobloxIPPool) error {
	logger := log.FromContext(ctx)

	ibclient, err := getInfobloxClientForInstance(ctx, r.Client, pool.Spec.InstanceRef.Name, pool.Namespace, r.NewInfobloxClientFunc)
	if err != nil {
		conditions.MarkFalse(pool,
			clusterv1.ReadyCondition,
			v1alpha1.InfobloxClientCreationFailedReason,
			clusterv1.ConditionSeverityError, "client creation failed for instance %s", pool.Spec.InstanceRef.Name)
		return err
	}

	// TODO: handle this in a better way
	if ok, err := ibclient.CheckNetworkViewExists(pool.Spec.NetworkView); err != nil || !ok {
		logger.Error(err, "could not find network view", "networkView", pool.Spec.NetworkView)
		conditions.MarkFalse(pool,
			clusterv1.ReadyCondition,
			v1alpha1.InfobloxNetworkViewNotFoundReason,
			clusterv1.ConditionSeverityError,
			"could not find network view: %s", err)
		return nil
	}

	for _, sub := range pool.Spec.Subnets {
		subnet, err := netip.ParsePrefix(sub.CIDR)
		if err != nil {
			// We won't set a condition here since this should be caught by validation
			return fmt.Errorf("failed to parse subnet: %w", err)
		}
		if ok, err := ibclient.CheckNetworkExists(pool.Spec.NetworkView, subnet); err != nil || !ok {
			logger.Error(err, "could not find network", "networkView", pool.Spec.NetworkView, "subnet", subnet)
			conditions.MarkFalse(pool,
				clusterv1.ReadyCondition,
				v1alpha1.InfobloxNetworkNotFoundReason,
				clusterv1.ConditionSeverityError,
				"could not find network: %s", err)
			return nil
		}
	}

	return nil
}
