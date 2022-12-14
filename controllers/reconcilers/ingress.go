package reconcilers

import (
	"context"

	"github.com/entgigi/gateway-operator.git/api/v1alpha1"
	"github.com/entgigi/gateway-operator.git/utility"

	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	ctrl "sigs.k8s.io/controller-runtime"
)

func (d *IngressManager) isIngressUpgrade(ctx context.Context, cr *v1alpha1.EntandoGatewayV2, ingress *netv1.Ingress) (error, bool) {
	err := d.Base.Client.Get(ctx, types.NamespacedName{Name: makeDeploymentName(cr), Namespace: cr.GetNamespace()}, ingress)
	if errors.IsNotFound(err) {
		return nil, false
	}
	return err, true
}

func (d *IngressManager) buildIngress(cr *v1alpha1.EntandoGatewayV2, scheme *runtime.Scheme) *netv1.Ingress {
	pp := netv1.PathTypePrefix
	ingress := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.IngressName,
			Namespace: cr.GetNamespace(),
		},
		Spec: netv1.IngressSpec{
			Rules: []netv1.IngressRule{{
				Host: cr.Spec.IngressHost,
				IngressRuleValue: netv1.IngressRuleValue{
					HTTP: &netv1.HTTPIngressRuleValue{
						Paths: []netv1.HTTPIngressPath{{
							Path:     cr.Spec.IngressPath,
							PathType: &pp,
							Backend: netv1.IngressBackend{
								Service: &netv1.IngressServiceBackend{
									Name: cr.Spec.IngressService,
									Port: netv1.ServiceBackendPort{
										Name: cr.Spec.IngressPort,
									},
								},
							},
						}},
					},
				},
			}},
		},
	} // set owner
	ctrl.SetControllerReference(cr, ingress, scheme)
	return ingress
}

func makeDeploymentName(cr *v1alpha1.EntandoGatewayV2) string {
	return "plugin-" + utility.TruncateString(cr.GetName(), 200) + "-deployment"
}

func (d *IngressManager) ApplyKubeIngress(ctx context.Context, cr *v1alpha1.EntandoGatewayV2, scheme *runtime.Scheme) error {
	baseIngress := d.buildIngress(cr, scheme)
	ingress := &netv1.Ingress{}

	err, isUpgrade := d.isIngressUpgrade(ctx, cr, ingress)
	if err != nil {
		return err
	}

	var applyError error
	if isUpgrade {
		ingress.Spec = baseIngress.Spec
		applyError = d.Base.Client.Update(ctx, ingress)

	} else {
		applyError = d.Base.Client.Create(ctx, baseIngress)
	}

	if applyError != nil {
		return applyError
	}
	return nil
}
