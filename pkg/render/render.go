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

func Render(conf render, base renderBase) {
	var helmChart string

	var artifactUrl = base.RemoteArtifact.Url
	if conf.RemoteArtifactOverride.Url != "" {
		artifactUrl = conf.RemoteArtifactOverride.Url
	}

	var artifactModificationCommand = base.RemoteArtifact.ModificationCommand
	if conf.RemoteArtifactOverride.ModificationCommand != "" {
		artifactModificationCommand = conf.RemoteArtifactOverride.ModificationCommand
	}

	var localChart = base.LocalChart.Directory
	if conf.LocalChartOverride.Directory != "" {
		localChart = conf.LocalChartOverride.Directory
	}

	var remoteChart = base.RemoteChart.Chart
	if conf.RemoteChartOverride.Chart != "" {
		remoteChart = conf.RemoteChartOverride.Chart
	}

	var remoteChartVersion = base.RemoteChart.Version
	if conf.RemoteChartOverride.Version != "" {
		remoteChartVersion = conf.RemoteChartOverride.Version
	}

	if remoteChart != "" && remoteChartVersion != "" {
		remoteChart = fmt.Sprintf("%s --version %s", remoteChart, remoteChartVersion)
	}

	if artifactUrl != "" {
		helmChart = GetRemoteArtifact(artifactUrl, artifactModificationCommand)
	} else if localChart != "" {
		helmChart = localChart
	} else if remoteChart != "" {
		helmChart = remoteChart
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
		Render(render, conf.Base)
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
