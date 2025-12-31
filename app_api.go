package pandora

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/mark0725/go-app-pandora/entities"

	"github.com/gin-gonic/gin"
	base_db "github.com/mark0725/go-app-base/db"
	base_web "github.com/mark0725/go-app-base/web"
)

type AppApi struct{}

var g_AppApi AppApi = AppApi{}

type I18nItem struct {
	FieldId   string `field-id:"field_id"`
	ItemKey   string `field-id:"item_key"`
	ItemValue string `field-id:"item_value"`
}

type ModuleNavI18nData struct {
	NavName   *I18nItem
	ShortName *I18nItem
}

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

	i18nMap, err := api.GetModuleNavi18n(c, g_appConfig.Org.OrgId, "main")
	if err != nil {
		logger.Error("GetModuleNavi18n error: ", err)
		c.JSON(http.StatusInternalServerError, ApiReponse{Code: "ERROR", Message: "GetModuleNavi18n error"})
		return
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

		title := module.NavName
		titleShort := module.TitleShort
		if i18n, ok := i18nMap[module.NavId]; ok {
			if i18n.NavName != nil {
				title = i18n.NavName.ItemValue
			}
			if i18n.ShortName != nil {
				titleShort = i18n.ShortName.ItemValue
			}
		}

		menu := MenuItem{
			Id:         module.NavId,
			Type:       module.ViewType,
			Title:      title,
			TitleShort: titleShort,
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

	i18nMainMap, err := api.GetModuleNavi18n(c, g_appConfig.Org.OrgId, "main")
	if err != nil {
		logger.Error("GetModuleNavi18n error: ", err)
		c.JSON(http.StatusInternalServerError, ApiReponse{Code: "ERROR", Message: "GetModuleNavi18n error"})
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
	moduleTitle := moduleInfo.NavName
	moduleTitleShort := moduleInfo.TitleShort
	if i18n, ok := i18nMainMap[moduleId]; ok {
		if i18n.NavName != nil {
			moduleTitle = i18n.NavName.ItemValue
		}
		if i18n.ShortName != nil {
			moduleTitleShort = i18n.ShortName.ItemValue
		}
	}

	pageConfig := PageConfig{
		Title:      moduleTitle,
		TitleShort: moduleTitleShort,
		Type:       moduleInfo.ViewType,
		Menu:       []*MenuItem{},
	}
	if moduleInfo.ViewType == "select-nav-page" {
		items, err := mapping(moduleInfo.AppModule, []string{moduleInfo.MenuApi}, nil, "en")
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
	i18nMap, err := api.GetModuleNavi18n(c, g_appConfig.Org.OrgId, moduleId)
	if err != nil {
		logger.Error("GetModuleNavi18n error: ", err)
		c.JSON(http.StatusInternalServerError, ApiReponse{Code: "ERROR", Message: "GetModuleNavi18n error"})
		return
	}
	navs, err := base_db.DBQueryEnt2[entities.BaseModuleNav](base_db.DB_CONN_NAME_DEFAULT, entities.DB_TABLE_BASE_MODULE_NAV, base_db.NewDBQueryOptions().Where("ORG_ID={ORG_ID} and MODULE_ID={MODULE_ID} and STATUS='00'").Params(sqlParams).Order("order_no"))
	if err != nil {
		logger.Error("DBQueryEnt fail: ", err)
		c.JSON(http.StatusInternalServerError, ApiReponse{Code: "ERROR", Message: "DBQueryEnt fail"})
		return
	}

	for _, nav := range navs {
		title := nav.NavName
		titleShort := nav.TitleShort
		if i18n, ok := i18nMap[nav.NavId]; ok {
			if i18n.NavName != nil {
				title = i18n.NavName.ItemValue
			}
			if i18n.ShortName != nil {
				titleShort = i18n.ShortName.ItemValue
			}
		}
		menu := MenuItem{
			Id:         nav.NavId,
			Type:       nav.ViewType,
			Title:      title,
			TitleShort: titleShort,
			Url:        nav.Url,
			Ico:        nav.NavIcon,
			View:       nav.NavType,
		}
		pageConfig.Menu = append(pageConfig.Menu, &menu)

	}

	c.JSON(http.StatusOK, ApiReponse{Code: "OK", Message: "OK", Data: pageConfig})
}

func (api *AppApi) GetModuleNavi18n(c *gin.Context, orgid string, group string) (map[string]*ModuleNavI18nData, error) {
	lang := "en"
	if l, ok := g_appConfig.Pandora.Env["lang"]; ok {
		lang = l.(string)
	}

	if cookie, err := c.Request.Cookie("app_lang"); err == nil {
		lang = cookie.Value
	}

	if l := c.Query("lang"); l != "" {
		lang = l
	}

	i18nParams := map[string]any{
		"ORG_ID":   orgid,
		"GROUP_ID": group,
	}
	reg := regexp.MustCompile(`[^a-zA-Z_]`)
	lang = reg.ReplaceAllString(lang, "")
	fields := fmt.Sprintf("field_id, item_key, %s item_value", lang)
	i18nItems, err := base_db.DBQueryG[I18nItem](base_db.DB_CONN_NAME_DEFAULT, base_db.NewDBQueryOptions().Table("base_i18n").Fields(strings.Split(fields, ",")).Where("ORG_ID={ORG_ID} and NS='menu' AND GROUP_ID={GROUP_ID} and STATUS='active'").Params(i18nParams))
	if err != nil {
		logger.Error("DBQueryEnt fail: ", err)
		c.JSON(http.StatusInternalServerError, ApiReponse{Code: "ERROR", Message: "DBQueryEnt fail"})
		return nil, err
	}

	i18nMap := map[string]*ModuleNavI18nData{}
	for _, item := range i18nItems {
		if i18nMap[item.ItemKey] == nil {
			i18nMap[item.ItemKey] = &ModuleNavI18nData{}
		}

		switch item.FieldId {
		case "nav_name":
			i18nMap[item.ItemKey].NavName = item
		case "title_short":
			i18nMap[item.ItemKey].ShortName = item
		}

	}

	return i18nMap, nil
}
