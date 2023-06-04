package k8s

// Taken from https://github.com/helm/helm/blob/main/pkg/releaseutil/kind_sorter.go and slightly modified

import "sort"

var KindOrder = []string{
	"Namespace",
	"NetworkPolicy",
	"ResourceQuota",
	"LimitRange",
	"PodSecurityPolicy",
	"PodDisruptionBudget",
	"ServiceAccount",
	"Secret",
	"SecretList",
	"ConfigMap",
	"StorageClass",
	"PersistentVolume",
	"PersistentVolumeClaim",
	"CustomResourceDefinition",
	"ClusterRole",
	"ClusterRoleList",
	"ClusterRoleBinding",
	"ClusterRoleBindingList",
	"Role",
	"RoleList",
	"RoleBinding",
	"RoleBindingList",
	"Service",
	"DaemonSet",
	"Pod",
	"ReplicationController",
	"ReplicaSet",
	"Deployment",
	"HorizontalPodAutoscaler",
	"StatefulSet",
	"Job",
	"CronJob",
	"IngressClass",
	"Ingress",
	"APIService",
}

// sort manifests by kind.
//
// Results are sorted by 'ordering', keeping order of items with equal kind/priority
func sortManifestsByKind(manifests ObjectStore) ObjectStore {
	sort.SliceStable(manifests, func(i, j int) bool {
		return lessByKind(manifests[i], manifests[j], manifests[i].GetKind(), manifests[j].GetKind())
	})

	return manifests
}

func lessByKind(a interface{}, b interface{}, kindA string, kindB string) bool {
	ordering := make(map[string]int, len(KindOrder))
	for v, k := range KindOrder {
		ordering[k] = v
	}

	first, aok := ordering[kindA]
	second, bok := ordering[kindB]

	if !aok && !bok {
		if kindA != kindB {
			return kindA < kindB
		}
		return first < second
	}

	if !aok {
		return false
	}
	if !bok {
		return true
	}

	return first < second
}
