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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	knativeclient "knative.dev/serving/pkg/client/clientset/versioned"
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

// +kubebuilder:rbac:groups=slo.ibm.com,resources=slodescs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=slo.ibm.com,resources=slodescs/status,verbs=get;update;patch

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
	if len(slo.Status.Metrics) == 0 {
		log.Printf("New SLODesc created, waiting for performance inputs")
	} else {
		latestMetric := slo.Status.Metrics[len(slo.Status.Metrics)-1]
		inputrate := latestMetric.InputRate
		meanlat := latestMetric.MeanLatency
		taillat := latestMetric.TailLatency
		stddev := latestMetric.StddevLatency

		log.Printf("Input rate: %d Mean lat %d stddev %d Tail lat %d\n", inputrate, meanlat, stddev, taillat)

		expected_inputrate := slo.Spec.RateLimit
		expected_taillat := slo.Spec.Tail.Latency

		if inputrate < expected_inputrate && inputrate > 0 {
			opts := metav1.GetOptions{}
			// Load related functions
			for _, serviceName := range slo.Spec.Workflow.Tasks {
				service, err := kc.ServingV1().ServicesGetter.Services("kcontainer").Get(serviceName, opts)
				if err != nil {
					log.Printf("Unable to get service %s", serviceName)
				} else {
					reservation := 1
					if reservation, ok := service.Annotations["autoscaling.knative.dev/minScale"]; ok {
						if tailat > expected_taillat {
							reservation *= 2
						} else {
							if reservation > 0 {
								reservation--
							}
						}
					}
					service.Annotations["autoscaling.knative.dev/minScale"] = reservation
					log.Printf("Set reservation for %s to %d", serviceName, reservation)
				}
			}
		}

	}

	return ctrl.Result{}, nil
}

// func newRestClient(cfg *rest.Config) (*rest.RESTClient, error) {
// 	scheme := runtime.NewScheme()
// 	SchemeBuilder := runtime.NewSchemeBuilder(knativescheme.AddKnownTypes)
// 	if err := SchemeBuilder.AddToScheme(scheme); err != nil {
// 		return nil, err
// 	}
// 	config := *cfg
// 	config.GroupVersion = &SchemeGroupVersion
// 	config.APIPath = "/apis"
// 	config.ContentType = runtime.ContentTypeJSON
// 	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme)
// 	client, err := rest.RESTClientFor(&config)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return client, nil
// }

var kc knativeclient.Clientset

func (r *SLODescReconciler) SetupWithManager(mgr ctrl.Manager) error {

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	kc, err = knativeclient.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&sloibmcomv1alpha1.SLODesc{}).
		Complete(r)
}
