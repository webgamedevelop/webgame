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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	webgamev1 "github.com/webgamedevelop/webgame/api/v1"
)

// WebGameReconciler reconciles a WebGame object
type WebGameReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=webgame.webgame.tech,resources=webgames,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=webgame.webgame.tech,resources=webgames/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=webgame.webgame.tech,resources=webgames/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *WebGameReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("webgame event received")
	defer func() { logger.Info("webgame event handling completed") }()

	var webgame webgamev1.WebGame
	if err := r.Get(ctx, req.NamespacedName, &webgame); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("webgame not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request
		logger.Error(err, "unable to fetch webgame, requeue")
		return ctrl.Result{}, err
	}

	logger.Info("show webgame name", "name", webgame.GetName())
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WebGameReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webgamev1.WebGame{}).
		Complete(r)
}
