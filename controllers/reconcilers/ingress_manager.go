package reconcilers

import (
	"context"
	"time"

	"github.com/entgigi/gateway-operator.git/api/v1alpha1"

	"github.com/entgigi/gateway-operator.git/common"
	"github.com/entgigi/gateway-operator.git/controllers/services"

	"k8s.io/apimachinery/pkg/runtime"
)

type IngressManager struct {
	Base       *common.BaseK8sStructure
	Conditions *services.ConditionService
}

func NewIngressManager(base *common.BaseK8sStructure, conditions *services.ConditionService) *IngressManager {
	return &IngressManager{
		Base:       base,
		Conditions: conditions,
	}
}

func (d *IngressManager) IsIngressApplied(ctx context.Context, cr *v1alpha1.EntandoGatewayV2) bool {

	return d.Conditions.IsIngressApplied(ctx, cr)
}

func (d *IngressManager) IsIngressReady(ctx context.Context, cr *v1alpha1.EntandoGatewayV2) bool {

	return d.Conditions.IsIngressReady(ctx, cr)
}

func (d *IngressManager) ApplyIngress(ctx context.Context, cr *v1alpha1.EntandoGatewayV2, scheme *runtime.Scheme) error {
	applyError := d.ApplyKubeIngress(ctx, cr, scheme)
	if applyError != nil {
		return applyError
	}

	return d.Conditions.SetConditionIngressApplied(ctx, cr)
}

func (d *IngressManager) CheckIngress(ctx context.Context, cr *v1alpha1.EntandoGatewayV2) (bool, error) {
	time.Sleep(time.Second * 10)
	ready := true

	if ready {
		return ready, d.Conditions.SetConditionIngressReady(ctx, cr)
	}

	return ready, nil

}
