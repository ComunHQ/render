# Render

Render is a CLI tool for generating the Kubernetes artifacts for a source Helm Chart and a custom values specification.

We built render at [Comun](https://github.com/ComunHQ) because:
1. At Comun, we manage multiple kubernetes clusters via fluxcd. Many of the clusters share the same core helm chart components with similar values. Render allows us to share code between all clusters.
2. We have found helm deployments to be operationally complex, but still want to take advantage of helm's packaging and templating capabilities.
3. Disaster recovery is important to us and render allows us to store raw kubernetes yaml specifications in our fluxcd state repository instead of pointing to remote charts stored in remote helm or ocr repositories.

## How to use Render

Render has two main inputs, a root yaml specification (`render.yaml`) and multiple values specifications (`values.jsonnet`'s).

The outputs of render are `*.generated.yaml` files which contain all the kubernetes artifacts that are included in the helm chart according to the specified values.

### Inputs

#### render.yaml

`render.yaml` tells the Render CLI how and where to render helm templates. Each entry in the `render.yaml` consists of:
1. A link to either a public remote helm repository chart, or a locally stored repository.
2. A list of locations to render the helm chart. The most important of the fields in each list entry is the `workingDirectory` which tells render where to render the chart. In addition, within the `workingDirectory`, there must be a `values.jsonnet` file, which tells the Render CLI what values to use when rendering the chart.

#### values.jsonnet

Each render workingDirectory should have a jsonnet file called `values.jsonnet`. The Render CLI is responsible for compiling the `values.jsonnet` file into json values that are then fed into helm when rendering the chart. You can structure `values.jsonnet` however you want as long as it can compile into a json file via the jsonnet compiler. You can even import from other jsonnet or libsonnet files.

### Outputs

#### *.generated.yaml files

One `.generated.yaml` is generated per render in the `render.yaml`. 

### Example

You can see a full example of a use case of the the `Render` tool in the [example](https://github.com/ComunHQ/render/tree/main/example) folder.
