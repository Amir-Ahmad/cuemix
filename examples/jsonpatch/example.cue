package examples

manifests: [
	"deploy.yaml",
]

jsonpatch: rm: {
	patch: [{
		op:   "remove"
		path: "/spec/template/spec/containers/0"
	}]
	target: {
		kind: "Deployment"
		name: "hello2"
	}
}
