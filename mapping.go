package pandora

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/mark0725/go-app-pandora/entities"

	"github.com/gin-gonic/gin"
	base_db "github.com/mark0725/go-app-base/db"
)

type CustMapping struct {
	Name    string
	Mapping PanPageMappingFunc
}

var g_CustMapping = map[string]PanPageMappingFunc{}

type PanPageMappingFunc func(list []string, params map[string]string) (map[string]*Dict, error)

func RegisterCustMapping(module string, fn PanPageMappingFunc) {
	g_CustMapping[module] = fn
}

func PageApiMapping(c *gin.Context) {
	moduleId := c.Param("module")
	// gin query

	dictList := c.Query("list")
	if dictList == "" {
		c.JSON(http.StatusBadRequest, ApiReponse{
			Code:    "ERROR",
			Message: "list is required",
		})
		return
	}

	logger.Debugf("PageApiMapping dictList: %s", dictList)

	resp := ApiReponse{
		Code:    "OK",
		Message: "success",
	}

	dicts := strings.Split(dictList, "|")
	if len(dicts) == 0 {
		resp.Data = map[string]any{}
		c.JSON(http.StatusOK, resp)
		return
	}

	lang := "en"
	if l, ok := g_appConfig.Pandora.Env["lang"]; ok {
		lang = l.(string)
	}

	if cookie, err := c.Request.Cookie("app_lang"); err == nil {
		lang = cookie.Value
	}

	queryParams := map[string]string{}
	for key, value := range c.Request.URL.Query() {
		if len(value) > 0 {
			queryParams[key] = value[0]
		}
		if key == "lang" {
			lang = value[0]
		}
	}

	//去除lang中除字母和下划线以外的其他字符
	reg := regexp.MustCompile(`[^a-zA-Z_]`)
	lang = reg.ReplaceAllString(lang, "")

	dictData, err := mapping(moduleId, dicts, queryParams, lang)
	if err != nil {
		logger.Errorf("mapping error: %v", err)
		c.JSON(http.StatusInternalServerError, ApiReponse{
			Code:    "ERROR",
			Message: fmt.Sprintf("mapping error: %v", err),
		})
		return
	}

	resp.Data = map[string]any{
		"dict": dictData,
	}

	c.JSON(http.StatusOK, resp)
}

func mapping(moduleId string, dicts []string, queryParams map[string]string, lang string) (map[string]*Dict, error) {
	dictData := map[string]*Dict{}
	if fn, ok := g_CustMapping[moduleId]; ok {
		dict, err := fn(dicts, queryParams)
		if err != nil {
			logger.Errorf("QueryDict error: %v", err)
			return nil, errors.New("QueryDict error")
		}
		for k, v := range dict {
			dictData[k] = v
		}
	}

	moduleDictList := []string{}
	for _, dictId := range dicts {
		if _, ok := dictData[dictId]; !ok {
			moduleDictList = append(moduleDictList, dictId)
		}
	}
	if len(moduleDictList) == 0 {
		return dictData, nil
	}

	moduleDict, err := QueryDict(moduleDictList, moduleId, lang)
	if err != nil {
		logger.Errorf("QueryDict error: %v", err)
		return nil, errors.New("QueryDict error")
	}
	for k, v := range moduleDict {
		dictData[k] = v
	}

	baseDictList := []string{}
	for _, dictId := range dicts {
		if _, ok := dictData[dictId]; !ok {
			baseDictList = append(baseDictList, dictId)
		}
	}

	if len(baseDictList) == 0 {
		return dictData, nil
	}

	baseDict, err := QueryDict(baseDictList, "base", lang)
	if err != nil {
		logger.Errorf("QueryDict error: %v", err)
		return nil, errors.New("QueryDict error")
	}
	for k, v := range baseDict {
		dictData[k] = v
	}

	return dictData, nil
}

type DictItemI18n struct {
	entities.BaseDictItems
	I18nValue    string `field-id:"i18n_value"`
	DefaultValue string `field-id:"default_value"`
}

func QueryDict(list []string, module string, lang string) (map[string]*Dict, error) {
	params := map[string]any{
		"ORG_ID":    g_appConfig.Org.OrgId,
		"MODULE_ID": module,
		"LIST":      list,
	}

	wheres := "a.ORG_ID={ORG_ID} AND a.MODULE_ID={MODULE_ID} AND a.DICT_ID IN ({LIST})"
	fields := "a.*, b.default_value"
	if lang != "" {
		fields = fmt.Sprintf("a.*, %s i18n_value, b.default_value", lang)
	}
	from := "base_dict_items a left join base_i18n b on a.module_id=b.module_id and a.dict_id=b.group_id and a.item_code=b.item_key and a.org_id=b.org_id and b.ns='dict'"

	recs, err := base_db.DBQueryG[DictItemI18n](base_db.DB_CONN_NAME_DEFAULT, base_db.NewDBQueryOptions().From(from).Fields(strings.Split(fields, ",")).Where(wheres).Params(params).Order("ord_no asc"))
	if err != nil {
		logger.Error("DBQuery fail: ", err)
		return nil, err
	}

	dictsItems := map[string]*Dict{}
	if recs == nil {
		logger.Error("providers  not found")
		return dictsItems, nil
	}
	for _, rec := range recs {
		if _, ok := dictsItems[rec.DictId]; !ok {
			dictsItems[rec.DictId] = &Dict{
				Id:   rec.DictId,
				Name: rec.DictName,
				// Type:    rec.DictType,
				// Style:   rec.ItemStyle,
				Items:   map[string]*DictItem{},
				Options: []*DictItem{},
			}
		}
		v := rec.ItemCode
		if rec.I18nValue != "" {
			v = rec.I18nValue
		}
		item := DictItem{
			Value: rec.ItemCode,
			Label: v,
			Icon:  rec.ItemIcon,
			Color: rec.ItemColor,
			Style: rec.ItemStyle,
		}
		item.Fields = map[string]string{
			"DICT_PARENT": rec.DictParent,
			"ITEM_PARENT": rec.ItemParent,
		}
		dictsItems[rec.DictId].Items[rec.ItemCode] = &item
		dictsItems[rec.DictId].Options = append(dictsItems[rec.DictId].Options, &item)
	}

	return dictsItems, nil
}
