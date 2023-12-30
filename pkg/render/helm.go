package render

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/bitfield/script"
	"github.com/codeclysm/extract/v3"
)

type helmChartAndVersion struct {
	Chart   string
	Version string
}

func GetRemoteArtifact(artifactUrl string, artifactModificationCommand string) string {
	var dir, err = os.MkdirTemp("", "")
	if err != nil {
		panic(err)
	}

	resp, err := http.Get(artifactUrl)
	if err != nil {
		panic(err)
	}

	stripFirstComponent := func(filepath string) string {
		parts := strings.Split(filepath, string(os.PathSeparator))
		return strings.Join(parts[1:], string(os.PathSeparator))
	}

	err = extract.Archive(context.TODO(), resp.Body, dir, stripFirstComponent)
	if err != nil {
		panic(err)
	}

	if artifactModificationCommand != "" {
		if _, err = script.Exec(fmt.Sprintf("bash -c 'cd %s && %s'", dir, artifactModificationCommand)).Stdout(); err != nil {
			panic(err)
		}
	}
	return dir
}

func GetChartAndVersionForRemoteArtifact(artifact remoteArtifact) *helmChartAndVersion {
	result := helmChartAndVersion{}
	result.Chart = GetRemoteArtifact(artifact.Url, artifact.ModificationCommand)
	return &result
}

func GetChartAndVersionForLocalChart(chart localChart) *helmChartAndVersion {
	result := helmChartAndVersion{}
	result.Chart = chart.Directory
	return &result
}

func GetChartAndVersionForRemoteChart(chart remoteChart) *helmChartAndVersion {
	result := helmChartAndVersion{}
	result.Chart = chart.Chart
	if chart.Version != "" {
		result.Version = chart.Version
	}
	return &result
}

func GetChartAndVersion(renderType string, conf render, base renderBase) *helmChartAndVersion {
	if renderType == "remoteArtifact" {
		artifact := base.RemoteArtifact
		if conf.RemoteArtifactOverride.Url != "" {
			artifact.Url = conf.RemoteArtifactOverride.Url
		}
		if conf.RemoteArtifactOverride.ModificationCommand != "" {
			artifact.ModificationCommand = conf.RemoteArtifactOverride.ModificationCommand
		}
		return GetChartAndVersionForRemoteArtifact(artifact)
	} else if renderType == "localChart" {
		chart := base.LocalChart
		if conf.LocalChartOverride.Directory != "" {
			chart.Directory = conf.LocalChartOverride.Directory
		}
		return GetChartAndVersionForLocalChart(chart)
	} else if renderType == "remoteChart" {
		chart := base.RemoteChart
		if conf.RemoteChartOverride.Chart != "" {
			chart.Chart = conf.RemoteChartOverride.Chart
		}
		if conf.RemoteChartOverride.Version != "" {
			chart.Version = conf.RemoteChartOverride.Version
		}
		return GetChartAndVersionForRemoteChart(chart)
	} else {
		panic(fmt.Sprintf("renderType `%s` not recognized", renderType))
	}
}

func HelmTemplate(
	workingDirectory string,
	chartAndVersion helmChartAndVersion,
	outputFilename string,
	valuesFile string,
	namespace string,
	releaseName string,
	includeCrds bool,
) {
	helmChart := chartAndVersion.Chart
	if chartAndVersion.Version != "" {
		helmChart += " --version " + chartAndVersion.Version
	}

	var command = fmt.Sprintf("helm template %s", helmChart)
	command = command + fmt.Sprintf(" --values %s", valuesFile)
	if namespace != "" {
		command = command + fmt.Sprintf(" --namespace %s", namespace)
	}
	if releaseName != "" {
		command = command + fmt.Sprintf(" --name-template %s", releaseName)
	}
	if includeCrds {
		command = command + " --include-crds"
	}
	slog.Debug("Running command: " + command)
	var pipe = script.Exec(command).WithStderr(os.Stderr)

	var outputFile = path.Join(workingDirectory, outputFilename+".generated.yaml")
	pipe.WriteFile(outputFile)
}
