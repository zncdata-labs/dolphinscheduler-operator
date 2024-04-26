package controller

import (
	dolphinv1alpha1 "github.com/zncdata-labs/dolphinscheduler-operator/api/v1alpha1"
	"github.com/zncdata-labs/dolphinscheduler-operator/internal/common"
	"github.com/zncdata-labs/dolphinscheduler-operator/pkg/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var serviceAccountName = func(instanceName string) string { return common.CreateServiceAccountName(instanceName) }

// NewServiceAccount new a ServiceAccountReconciler
func NewServiceAccount(
	scheme *runtime.Scheme,
	instance *dolphinv1alpha1.DolphinschedulerCluster,
	client client.Client,
	mergedLabels map[string]string,
	mergedCfg any,
) *resource.GenericServiceAccountReconciler[*dolphinv1alpha1.DolphinschedulerCluster, any] {
	return resource.NewServiceAccount[*dolphinv1alpha1.DolphinschedulerCluster](scheme, instance, client, mergedLabels, mergedCfg,
		serviceAccountName(instance.GetName()), instance.GetNamespace())
}
