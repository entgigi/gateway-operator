package reconcilers

import (
	"context"

	"github.com/entgigi/gateway-operator.git/api/v1alpha1"
	"github.com/entgigi/gateway-operator.git/utility"

	netv1 "k8s.io/api/networking/v1"
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
	/*
		replicas := cr.Spec.Replicas
		deploymentName := makeDeploymentName(cr)
		containerName := makeContainerName(cr)
		labels := map[string]string{labelKey: containerName}
		port := int32(cr.Spec.Port)
	*/
	ingress := &netv1.Ingress{} // set owner
	ctrl.SetControllerReference(cr, ingress, scheme)
	return ingress
}

func makeContainerName(cr *v1alpha1.EntandoGatewayV2) string {
	return "plugin-" + utility.TruncateString(cr.GetName(), 200) + "-container"
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
