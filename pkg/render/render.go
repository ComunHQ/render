package render

import (
	"fmt"
	"os"

	"github.com/google/go-jsonnet"
	"gopkg.in/yaml.v3"
	sigsyaml "sigs.k8s.io/yaml"
)

func ValuesYaml(workingDirectory string) string {
	jsonnetFile := fmt.Sprintf("%s/values.jsonnet", workingDirectory)
	vm := jsonnet.MakeVM()
	jsonValues, err := vm.EvaluateFile(jsonnetFile)
	if err != nil {
		panic(err)
	}

	yaml, err := sigsyaml.JSONToYAML([]byte(jsonValues))
	if err != nil {
		panic(err)
	}

	return string(yaml)
}

func Render(renderType string, conf render, base renderBase) {
	chartAndVersion := GetChartAndVersion(renderType, conf, base)

	var outputFilename = conf.OutputFile
	if outputFilename == "" {
		outputFilename = "release"
	}

	var kubeVersion = base.KubeVersion
	if conf.KubeVersionOverride != "" {
		kubeVersion = conf.KubeVersionOverride
	}

	valuesYaml := ValuesYaml(conf.WorkingDirectory)

	HelmTemplate(
		conf.WorkingDirectory,
		*chartAndVersion,
		outputFilename,
		valuesYaml,
		conf.Namespace,
		conf.ReleaseName,
		base.IncludeCrds,
		kubeVersion,
	)
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
