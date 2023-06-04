package examples

manifests: [
	"https://raw.githubusercontent.com/argoproj/argo-cd/v2.6.7/manifests/install.yaml",
]

// Patch argocd-rbac-cm configmap
strategicpatch: "argocd-rbac-cm": patch: {
	apiVersion: "v1"
	kind:       "ConfigMap"
	metadata: name: "argocd-rbac-cm"
	data: {
		"policy.default": ""
		"policy.csv": """
			g, admins, role:admin
			"""
	}
}
