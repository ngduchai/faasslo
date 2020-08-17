/*


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
	"log"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	sloibmcomv1alpha1 "github.com/ngduchai/faasslo/api/v1alpha1"
)

// SLODescReconciler reconciles a SLODesc object
type SLODescReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=slo.ibm.com.slo.ibm.com,resources=slodescs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=slo.ibm.com.slo.ibm.com,resources=slodescs/status,verbs=get;update;patch

func (r *SLODescReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("slodesc", req.NamespacedName)

	// your logic here
	rst := ctrl.Result{}
	slo := &sloibmcomv1alpha1.SLODesc{}
	err := r.Get(ctx, req.NamespacedName, slo)
	if err != nil {
		return rst, err
	}
	latestMetric := slo.Status.Metrics[len(slo.Status.Metrics)-1]
	inputrate := latestMetric.InputRate
	meanlat := latestMetric.MeanLatency
	taillat := latestMetric.TailLatency
	stddev := latestMetric.StddevLatency
	log.Printf("Input rate: %d Mean lat %d stddev %d Tail lat %d\n", inputrate, meanlat, stddev, taillat)

	return ctrl.Result{}, nil
}

func (r *SLODescReconciler) SetupWithManager(mgr ctrl.Manager) error {

	return ctrl.NewControllerManagedBy(mgr).
		For(&sloibmcomv1alpha1.SLODesc{}).
		Complete(r)
}
