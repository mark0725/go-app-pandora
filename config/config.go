package config

type PanConfig struct {
	Title       string         `json:"title" yaml:"title"`
	Logo        string         `json:"logo" yaml:"logo"`
	Description string         `json:"description" yaml:"description"`
	PagesPath   string         `json:"page_path" yaml:"page_path"`
	PanUrl      string         `json:"pandora_url" yaml:"pandora_url"` //embed:///pandora/ http://localhost:8080/ file:///pandora/
	Env         map[string]any `json:"env" yaml:"env"`
	User        PanUserConfig  `json:"user"`
}

type PanUserConfig struct {
	AuthType   string `json:"auth_type"`
	AuthUrl    string `json:"auth_url"`
	SigninUrl  string `json:"signin_url"`
	SignoutUrl string `json:"signout_url"`
	Avatar     string `json:"avatar"`
}
