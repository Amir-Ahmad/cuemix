#Cuemix: helm: {
	// Helm charts are downloaded to a local path
	destDir: string | *".charts"

	// If untar = false, tgz files are stored and used
	// If untar = true, tgz files are extracted then deleted
	untar: bool | *false

	// Some charts rely on capabilities.KubeVersion
	kubeVersion: string | *""
}

#HelmChart: {
	// repo can be a helm repo, oci repo, or a local directory
	repo:        string
	chart:       string
	release:     string
	namespace:   string
	version:     string
	includeCRDs: bool | *false
	values?: {...}
}

// Target for patches, the implementation is pretty much the same as kustomize
// https://github.com/kubernetes-sigs/kustomize/blob/master/examples/patchMultipleObjects.md
#Target: {
	group?:              string
	version?:            string
	kind?:               string
	name?:               string
	namespace?:          string
	annotationSelector?: string
	labelSelector?:      string
}

#Operation: {
	op:     "add" | "remove" | "replace" | "move" | "copy" | "test"
	path:   string
	from?:  string
	value?: _
}

cuemix: #Cuemix

helm?: [string]: #HelmChart

// Manifests can be loaded from a (yaml) file, directory, or URL
manifests?: [...string]

strategicpatch?: [string]: {
	patch: [string]: _

	// Target is optional for strategic patches,
	// it will be derived from the patches gvk, name, and namespace by default
	target: #Target
}

jsonpatch?: [string]: {
	patch: [...#Operation]

	// Target must be specified for json patches
	target: #Target
}
