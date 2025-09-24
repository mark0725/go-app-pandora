package pandora

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	base_app "github.com/mark0725/go-app-base/app"
	base_log "github.com/mark0725/go-app-base/logger"
	base_web "github.com/mark0725/go-app-base/web"
)

const APP_MODULE_NAME string = "pandora"

type AppModule struct{}

var logger = base_log.GetLogger(APP_MODULE_NAME)
var g_appConfig *AppModuleConfig

func init() {
	base_app.AppModuleRegister(APP_MODULE_NAME, &AppModule{}, []string{"db"},
		base_app.AppModuleRegisterOptionWithConfigType(&AppModuleConfig{}),
	)
}

func (m *AppModule) Init(appConfig any, depends []string) error {
	fmt.Printf("config: %#v\n", appConfig)
	if v, ok := appConfig.(*AppModuleConfig); !ok {
		logger.Error("invalid app config")
		return errors.New("invalid app config")
	} else {
		g_appConfig = v
	}

	logger.Tracef("AppModule %s init ... ", APP_MODULE_NAME)
	if err := InitPageEngine(); err != nil {
		return err
	}

	base_web.ServerConfigRegister(APP_MODULE_NAME, func(name string, params map[string]any, c *gin.Engine) {
		switch name {
		case "assets":
			c.NoRoute(PandoraAssets)
		}
	})

	base_web.EndPointRegister(APP_MODULE_NAME, func(group string, r *gin.RouterGroup) {
		InitRoutes(group, r)
	})
	return nil
}
