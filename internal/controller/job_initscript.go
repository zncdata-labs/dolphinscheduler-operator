package controller

import (
	"context"
	"fmt"
	dolphinv1alpha1 "github.com/zncdata-labs/dolphinscheduler-operator/api/v1alpha1"
	"github.com/zncdata-labs/dolphinscheduler-operator/internal/common"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ common.ResourceBuilder = &JobInitScriptReconciler{}

type JobInitScriptReconciler struct {
	common.GeneralResourceStyleReconciler[*dolphinv1alpha1.DolphinschedulerCluster, *dolphinv1alpha1.RoleGroupSpec]
}

func (j *JobInitScriptReconciler) Build(_ context.Context) (client.Object, error) {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      j.Instance.Name,
			Namespace: j.Instance.Namespace,
			Labels:    j.Labels,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						j.InitDbContainer(),
					},
					InitContainers: []corev1.Container{
						j.waitDbContainer(),
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}, nil
}

func (j *JobInitScriptReconciler) InitDbContainer() corev1.Container {
	return corev1.Container{
		Name:  string(dbInitJob),
		Image: dolphinv1alpha1.DbInitImage,
		Args: []string{
			"tools/bin/upgrade-schema.sh",
		},
		Env: j.InitDbContainerEnvs(),
		EnvFrom: []corev1.EnvFromSource{
			{
				ConfigMapRef: &corev1.ConfigMapEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: common.EnvsConfigMapName(j.Instance.GetName(), j.GroupName),
					},
				},
			},
		},
	}
}

func (j *JobInitScriptReconciler) InitDbContainerEnvs() []corev1.EnvVar {
	envs := []corev1.EnvVar{
		{
			Name:  "TZ",
			Value: "Asia/Shanghai",
		},
		{
			Name:  "SPRING_JACKSON_TIME_ZONE",
			Value: "Asia/Shanghai",
		},
		{
			Name:  "REGISTRY_TYPE",
			Value: "zookeeper",
		},
		{
			Name: "REGISTRY_ZOOKEEPER_CONNECT_STRING",
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "dolphinscheduler-common",
					},
					Key: common.ZookeeperDiscoveryKey,
				},
			},
		},
	}
	envs = append(envs, common.MakeDataBaseEnvs(j.Instance.Spec.ClusterConfigSpec.Database)...)
	return envs
}

func (j *JobInitScriptReconciler) waitDbContainer() corev1.Container {
	_, params := common.ExtractDataBaseReference(j.Instance.Spec.ClusterConfigSpec.Database)
	dbHost := params.Host
	return corev1.Container{
		Name:  string(waitForDb),
		Image: "busybox:1.30.1",
		Command: []string{
			"sh",
			"-xc",
			fmt.Sprintf("for i in $(seq 1 180); do nc -z -w3 %s 5432 && exit 0 || sleep 5; done; exit 1", dbHost),
		},
	}
}

const (
	dbInitJob common.ContainerComponent = "dolphinscheduler-db-init-job"
	waitForDb common.ContainerComponent = "wait-for-database"
)