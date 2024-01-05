package render

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
	KubeVersion    string         `yaml:"kubeVersion"`
}

type render struct {
	WorkingDirectory       string         `yaml:"workingDirectory"`
	RemoteArtifactOverride remoteArtifact `yaml:"remoteArtifact"`
	LocalChartOverride     localChart     `yaml:"localChart"`
	RemoteChartOverride    remoteChart    `yaml:"remoteChart"`
	OutputFile             string         `yaml:"outputFile"`
	ReleaseName            string         `yaml:"releaseName"`
	Namespace              string         `yaml:"namespace"`
	KubeVersionOverride    string         `yaml:"kubeVersion"`
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
