package render

import (
	"fmt"
	"os"

	"github.com/google/go-jsonnet"
	"gopkg.in/yaml.v3"
)

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

	valuesFile := ValuesFile(conf.WorkingDirectory)

	HelmTemplate(
		conf.WorkingDirectory,
		*chartAndVersion,
		outputFilename,
		valuesFile,
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
