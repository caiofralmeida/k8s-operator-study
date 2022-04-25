/*
Copyright 2022.

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

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	webstaticv1alpha1 "github.com/caiofralmeida/k8s-operator-study/api/v1alpha1"
)

// MarkdownContentReconciler reconciles a MarkdownContent object
type MarkdownContentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=webstatic.caiofralmeida.github.io,resources=markdowncontents,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=webstatic.caiofralmeida.github.io,resources=markdowncontents/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=webstatic.caiofralmeida.github.io,resources=markdowncontents/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MarkdownContent object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *MarkdownContentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	mc := &webstaticv1alpha1.MarkdownContent{}
	err := r.Get(ctx, req.NamespacedName, mc)

	if k8serrors.IsNotFound(err) {
		log.Log.Info("being removed")
		return ctrl.Result{}, nil
	}

	wp := &webstaticv1alpha1.WebPage{
		ObjectMeta: v1.ObjectMeta{
			Name:      "jamal",
			Namespace: req.Namespace,
		},
	}
	ctrl.SetControllerReference(mc, wp, r.Scheme)

	if err := r.Create(ctx, wp); err != nil {
		if k8serrors.IsAlreadyExists(err) {
			log.Log.Info("already exists ignoring")
			return ctrl.Result{}, nil
		}

		log.Log.Error(err, "reconciling error", "foo", "bar")
		return ctrl.Result{Requeue: false}, err
	}

	log.Log.Info("reconcilied with success", "foo", "bar")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MarkdownContentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webstaticv1alpha1.MarkdownContent{}).
		WithEventFilter(&predicate.GenerationChangedPredicate{}).
		Complete(r)
}
