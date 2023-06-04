package examples

helm: grafana: {
	repo:      "vendored/grafana/"
	chart:     "grafana"
	release:   "grafana"
	namespace: "monitoring"
	version:   "6.57.0"
	values: {}
}
