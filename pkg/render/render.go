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
	"github.com/google/go-jsonnet"
	"gopkg.in/yaml.v3"
)

type renderConf struct {
	Type    string     `yaml:"type"`
	Base    renderBase `yaml:"base"`
	Renders []render   `yaml:"renders"`
}

type renderBase struct {
	RemoteArtifact remoteArtifact `yaml:"remoteArtifact"`
	LocalChart     localChart     `yaml:"localChart"`
	RemoteChart    remoteChart    `yaml:"remoteChart"`
	IncludeCrds    bool           `yaml:"includeCrds"`
}

type render struct {
	WorkingDirectory       string         `yaml:"workingDirectory"`
	RemoteArtifactOverride remoteArtifact `yaml:"remoteArtifact"`
	LocalChartOverride     localChart     `yaml:"localChart"`
	RemoteChartOverride    remoteChart    `yaml:"remoteChart"`
	OutputFile             string         `yaml:"outputFile"`
	ReleaseName            string         `yaml:"releaseName"`
	Namespace              string         `yaml:"namespace"`
}

type remoteArtifact struct {
	Url                 string `yaml:"url"`
	ModificationCommand string `yaml:"modificationCommand"`
}

type localChart struct {
	Directory string `yaml:"directory"`
}

type remoteChart struct {
	Chart   string `yaml:"chart"`
	Version string `yaml:"version"`
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

func ValuesFile(workingDirectory string) string {
	var dir, err = os.MkdirTemp("", "")
	if err != nil {
		panic(err)
	}

	jsonnetFile := fmt.Sprintf("%s/values.jsonnet", workingDirectory)
	vm := jsonnet.MakeVM()
	jsonValues, err := vm.EvaluateFile(jsonnetFile)
	if err != nil {
		panic(err)
	}

	valuesFile := fmt.Sprintf("%s/values.json", dir)
	err = os.WriteFile(valuesFile, []byte(jsonValues), 0644)
	if err != nil {
		panic(err)
	}
	return valuesFile
}

type helmChartAndVersion struct {
	Chart   string
	Version string
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

func Render(renderType string, conf render, base renderBase) {
	chartAndVersion := GetChartAndVersion(renderType, conf, base)
	helmChart := chartAndVersion.Chart
	if chartAndVersion.Version != "" {
		helmChart += " --version " + chartAndVersion.Version
	}

	var outputFilename = conf.OutputFile
	if outputFilename == "" {
		outputFilename = "release"
	}

	valuesFile := ValuesFile(conf.WorkingDirectory)

	var command = fmt.Sprintf("helm template %s", helmChart)
	command = command + fmt.Sprintf(" --values %s", valuesFile)
	if conf.Namespace != "" {
		command = command + fmt.Sprintf(" --namespace %s", conf.Namespace)
	}
	if conf.ReleaseName != "" {
		command = command + fmt.Sprintf(" --name-template %s", conf.ReleaseName)
	}
	if base.IncludeCrds {
		command = command + " --include-crds"
	}
	slog.Debug("Running command: " + command)
	var pipe = script.Exec(command).WithStderr(os.Stderr)

	var outputFile = path.Join(conf.WorkingDirectory, outputFilename+".generated.yaml")
	pipe.WriteFile(outputFile)
}

func Run(conf renderConf) {
	for _, render := range conf.Renders {
		Render(conf.Type, render, conf.Base)
	}
}

func GetConfigs(file string) map[string]renderConf {
	confs := make(map[string]renderConf)
	bytes, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(bytes, &confs)
	if err != nil {
		panic(err)
	}
	return confs
}
