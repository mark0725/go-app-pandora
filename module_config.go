package pandora

import (
	base_config "github.com/mark0725/go-app-base/config"
	base_logger "github.com/mark0725/go-app-base/logger"
	pan_config "github.com/mark0725/go-app-pandora/config"
)

type AppModuleConfig struct {
	Org     base_config.OrgConfig                            `json:"org"`
	App     base_config.AppConfigApp                         `json:"app"`
	Serves  map[string]map[string]base_config.AppServeConfig `json:"serves"`
	Log     base_logger.LogConfig                            `json:"log"`
	Pandora pan_config.PanConfig                             `json:"pandora"`
}
