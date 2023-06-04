package helm

type HelmChart struct {
	Repo        string                 `json:"repo"`
	Chart       string                 `json:"chart"`
	Release     string                 `json:"release"`
	Namespace   string                 `json:"namespace"`
	Version     string                 `json:"version"`
	IncludeCRDs bool                   `json:"includeCRDs"`
	Values      map[string]interface{} `json:"values"`
}
