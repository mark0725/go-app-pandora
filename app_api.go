package pandora

import (
	"fmt"
	"net/http"

	"github.com/mark0725/go-app-pandora/entities"

	"github.com/gin-gonic/gin"
	base_db "github.com/mark0725/go-app-base/db"
	base_web "github.com/mark0725/go-app-base/web"
	_ "github.com/mark0725/go-app-base/web/middleware"
)

type AppApi struct{}

var g_AppApi AppApi = AppApi{}

func (api *AppApi) Config(c *gin.Context) {
	appConfig := AppConfig{
		App: AppInfo{
			Logo: g_appConfig.Pandora.Logo,
			// Name:    version.AppName,
			Title: g_appConfig.Pandora.Title,
			// Version: version.Version,
		},
		Menu: MenuConfig{
			Main: []*MenuItem{},
			Nav2: []*MenuItem{},
		},
	}
	if len(g_appConfig.Pandora.User.AuthType) > 0 {
		appConfig.Auth = &UserAuth{
			AuthType:   g_appConfig.Pandora.User.AuthType,
			AuthUrl:    g_appConfig.Pandora.User.AuthUrl,
			SigninUrl:  g_appConfig.Pandora.User.SigninUrl,
			SignoutUrl: g_appConfig.Pandora.User.SignoutUrl,
		}
		userId := ""
		if v, ok := c.Get(base_web.CtxKeyAuthenticatedConsumer); !ok {
			c.JSON(http.StatusOK, ApiReponse{Code: "OK", Message: "OK", Data: appConfig})
			return
		} else {
			if vv, ok := v.(*base_web.AuthenticatedConsumer); ok {
				appConfig.Auth.Authed = true
				userId = vv.Id
			} else {
				logger.Errorf("AuthenticatedConsumer type error: %s", userId)
				c.JSON(http.StatusInternalServerError, ApiReponse{Code: "ERROR", Message: "User not found"})
				return
			}
		}

		sqlParams := map[string]any{
			"ORG_ID":  g_appConfig.Org.OrgId,
			"USER_ID": userId,
		}
		logger.Debug("sqlParams:", sqlParams)

		users, err := base_db.DBQueryEnt2[entities.IdmUser](base_db.DB_CONN_NAME_DEFAULT, entities.DB_TABLE_IDM_USER, base_db.NewDBQueryOptions().Where("ORG_ID={ORG_ID} and USER_ID={USER_ID} and STATUS='active'").Params(sqlParams))
		if err != nil {
			logger.Error("DBQueryEnt fail: ", err)
			c.JSON(http.StatusInternalServerError, ApiReponse{Code: "ERROR", Message: "DBQueryEnt fail"})
			return
		}

		if len(users) == 0 {
			c.JSON(http.StatusBadRequest, ApiReponse{Code: "bad_request", Message: "user not found"})
			return
		}

		userInfo := users[0]
		appConfig.User = &UserInfo{
			Id:     userInfo.UserId,
			Name:   userInfo.UserName,
			Dept:   userInfo.DeptNo,
			Mail:   userInfo.Email,
			Avatar: fmt.Sprintf(g_appConfig.Pandora.User.Avatar, userId),
		}
	}
	sqlParams := map[string]any{
		"ORG_ID": g_appConfig.Org.OrgId,
	}
	logger.Debug("sqlParams:", sqlParams)

	modules, err := base_db.DBQueryEnt2[entities.BaseModuleNav](base_db.DB_CONN_NAME_DEFAULT, entities.DB_TABLE_BASE_MODULE_NAV, base_db.NewDBQueryOptions().Where("ORG_ID={ORG_ID} and MODULE_ID='main' and STATUS='00'").Params(sqlParams).Order("order_no"))
	if err != nil {
		logger.Error("DBQueryEnt fail: ", err)
		c.JSON(http.StatusInternalServerError, ApiReponse{Code: "ERROR", Message: "DBQueryEnt fail"})
		return
	}
	for _, module := range modules {
		menu := MenuItem{
			Id:         module.NavId,
			Type:       module.ViewType,
			Title:      module.NavName,
			TitleShort: module.TitleShort,
			View:       module.NavType,
			Url:        module.Url,
			Ico:        module.NavIcon,
		}
		switch module.NavPosition {
		case "main":
			appConfig.Menu.Main = append(appConfig.Menu.Main, &menu)
		case "nav2":
			appConfig.Menu.Nav2 = append(appConfig.Menu.Nav2, &menu)
		case "nav-user":
			if appConfig.Menu.NavUser == nil {
				appConfig.Menu.NavUser = make([]*MenuItem, 0)
			}
			appConfig.Menu.NavUser = append(appConfig.Menu.NavUser, &menu)

		default:
			logger.Errorf("module: %s nav_position: %s error: %v", module.ModuleId, module.NavPosition, err)
		}

	}

	c.JSON(http.StatusOK, ApiReponse{Code: "OK", Message: "OK", Data: appConfig})
}

func (api *AppApi) ModuleConfig(c *gin.Context) {
	moduleId := c.Param("module")
	if moduleId == "" {
		logger.Error("moduleid is required")
		c.JSON(http.StatusBadRequest, ApiReponse{Code: "BadRequest", Message: "moduleid is required"})
		return
	}

	sqlParams := map[string]any{
		"ORG_ID":    g_appConfig.Org.OrgId,
		"MODULE_ID": moduleId,
	}
	logger.Debug("sqlParams:", sqlParams)
	recs, err := base_db.DBQueryEnt[entities.BaseModuleNav](base_db.DB_CONN_NAME_DEFAULT, entities.DB_TABLE_BASE_MODULE_NAV, "ORG_ID={ORG_ID} and MODULE_ID='main' and NAV_ID={MODULE_ID} and STATUS='00'", sqlParams)
	if err != nil {
		logger.Error("DBQueryEnt fail: ", err)
		c.JSON(http.StatusInternalServerError, ApiReponse{Code: "ERROR", Message: "DBQueryEnt fail"})
		return
	}

	if len(recs) == 0 {
		c.JSON(http.StatusBadRequest, ApiReponse{Code: "DATA_EXIST", Message: "not found module " + moduleId})
		return
	}

	moduleInfo := recs[0]
	pageConfig := PageConfig{
		Title:      moduleInfo.NavName,
		TitleShort: moduleInfo.TitleShort,
		Type:       moduleInfo.ViewType,
		Menu:       []*MenuItem{},
	}
	if moduleInfo.ViewType == "select-nav-page" {
		items, err := mapping(moduleInfo.AppModule, []string{moduleInfo.MenuApi}, nil)
		if err != nil {
			logger.Error("QueryDict error: ", err)
			c.JSON(http.StatusInternalServerError, ApiReponse{Code: "ERROR", Message: "QueryDict error"})
		}

		if len(items[moduleInfo.MenuApi].Options) > 0 {
			pageConfig.Select = &MenuSelect{
				Param: moduleInfo.ParamName,
				Value: items[moduleInfo.MenuApi].Options[0].Value,
				Items: items[moduleInfo.MenuApi].Options,
			}
		}
	}
	navs, err := base_db.DBQueryEnt2[entities.BaseModuleNav](base_db.DB_CONN_NAME_DEFAULT, entities.DB_TABLE_BASE_MODULE_NAV, base_db.NewDBQueryOptions().Where("ORG_ID={ORG_ID} and MODULE_ID={MODULE_ID} and STATUS='00'").Params(sqlParams).Order("order_no"))
	if err != nil {
		logger.Error("DBQueryEnt fail: ", err)
		c.JSON(http.StatusInternalServerError, ApiReponse{Code: "ERROR", Message: "DBQueryEnt fail"})
		return
	}

	for _, nav := range navs {
		menu := MenuItem{
			Id:         nav.NavId,
			Type:       nav.ViewType,
			Title:      nav.NavName,
			TitleShort: nav.TitleShort,
			Url:        nav.Url,
			Ico:        nav.NavIcon,
			View:       nav.NavType,
		}
		pageConfig.Menu = append(pageConfig.Menu, &menu)

	}

	c.JSON(http.StatusOK, ApiReponse{Code: "OK", Message: "OK", Data: pageConfig})
}
