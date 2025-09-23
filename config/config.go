package config

type PanConfig struct {
	Title       string         `json:"title" yaml:"title"`
	Logo        string         `json:"logo" yaml:"logo"`
	Description string         `json:"description" yaml:"description"`
	PagesPath   string         `json:"page_path" yaml:"page_path"`
	PanUrl      string         `json:"pandora_url" yaml:"pandora_url"`
	Auth        string         `json:"auth" yaml:"auth"`
	Env         map[string]any `json:"env" yaml:"env"`
	User        PanUserConfig  `json:"user"`
}

type PanUserConfig struct {
	Avatar string `json:"avatar"`
}
