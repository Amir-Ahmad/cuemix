package examples

// helm charts
helm: "ingress-nginx": {
	repo:      "https://kubernetes.github.io/ingress-nginx"
	chart:     "ingress-nginx"
	release:   "ingress-nginx"
	namespace: "ingress-nginx"
	version:   "4.6.0"
	values: controller: replicaCount: 3
}

// Apply patch to remove minAvailable and set maxUnavailable
strategicpatch: pdb: patch: {
	apiVersion: "policy/v1"
	kind:       "PodDisruptionBudget"
	metadata: name: "ingress-nginx-controller"
	spec: {
		minAvailable:   null
		maxUnavailable: "50%"
	}
}
