# cuemix

Cuemix allows you to source external manifests (via helm or yaml) in cue, and apply patches on top.

It uses helm's golang library to interact with helm charts, and kustomize's api to apply patches.

Cuemix only concerns itself with generating manifests and is meant to be used with a deployment tool like ArgoCD.

## Examples

Check out the [examples](examples/) directory to see some sample definitions.

Try out any of the examples with cuemix build, e.g:

```
cuemix build ./examples/local-manifests/

# show only deployments
cuemix build -k deployment ./examples/local-manifests/
```

## Features

- Load helm charts from a helm repo, OCI repo, or local directory
- Load manifests from a local file, directory, or URL
- Apply [strategic](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-api-machinery/strategic-merge-patch.md) or [json](https://datatracker.ietf.org/doc/html/rfc6902) patches
- Output objects sorted by kind, and optionally filter by kind

Check out [schema.cue](internal/app/schema.cue) for all configuration options.

## Planned work
- Support loading cue objects from a configurable path
- e2e tests
