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
	ArtifactUrl                 string `yaml:"artifactUrl"`
	ArtifactModificationCommand string `yaml:"artifactModificationCommand"`
	Chart                       string `yaml:"chart"`
	IncludeCrds                 bool   `yaml:"includeCrds"`
}

type render struct {
	WorkingDirectory                    string `yaml:"workingDirectory"`
	ArtifactUrlOverride                 string `yaml:"artifactUrl"`
	ArtifactModificationCommandOverride string `yaml:"artifactModificationCommand"`
	ChartOverride                       string `yaml:"chart"`
	OutputFile                          string `yaml:"outputFile"`
	ReleaseName                         string `yaml:"releaseName"`
	Namespace                           string `yaml:"namespace"`
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

	var artifactUrl = base.ArtifactUrl
	if conf.ArtifactUrlOverride != "" {
		artifactUrl = conf.ArtifactUrlOverride
	}

	var artifactModificationCommand = base.ArtifactModificationCommand
	if conf.ArtifactModificationCommandOverride != "" {
		artifactModificationCommand = conf.ArtifactModificationCommandOverride
	}

	var chart = base.Chart
	if conf.ChartOverride != "" {
		chart = conf.ChartOverride
	}

	if artifactUrl != "" {
		helmChart = GetRemoteArtifact(artifactUrl, artifactModificationCommand)
	} else if chart != "" {
		helmChart = chart
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
	if conf.OutputFile != "" {
		var outputFile = path.Join(conf.WorkingDirectory, conf.OutputFile+".generated.yaml")
		pipe.WriteFile(outputFile)
	} else {
		pipe.Stdout()
	}
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