package pandora

import (
	"errors"
	"fmt"
	"net/http"
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
	queryParams := map[string]string{}
	for key, value := range c.Request.URL.Query() {
		if len(value) > 0 {
			queryParams[key] = value[0]
		}
	}

	dictData, err := mapping(moduleId, dicts, queryParams)
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

func mapping(moduleId string, dicts []string, queryParams map[string]string) (map[string]*Dict, error) {
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

	moduleDict, err := QueryDict(moduleDictList, moduleId)
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

	baseDict, err := QueryDict(baseDictList, "base")
	if err != nil {
		logger.Errorf("QueryDict error: %v", err)
		return nil, errors.New("QueryDict error")
	}
	for k, v := range baseDict {
		dictData[k] = v
	}

	return dictData, nil
}

func QueryDict(list []string, module string) (map[string]*Dict, error) {
	params := map[string]any{
		"ORG_ID":    g_appConfig.Org.OrgId,
		"MODULE_ID": module,
		"LIST":      list,
	}

	wheres := "ORG_ID={ORG_ID} AND MODULE_ID={MODULE_ID} AND DICT_ID IN ({LIST})"

	recs, err := base_db.DBQueryEnt2[entities.BaseDictItems](base_db.DB_CONN_NAME_DEFAULT, entities.DB_TABLE_BASE_DICT_ITEMS, base_db.NewDBQueryOptions().Where(wheres).Params(params))
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
		item := DictItem{
			Value: rec.ItemCode,
			Label: rec.ItemValue,
			Icon:  rec.ItemIcon,
			Color: rec.ItemColor,
			Style: rec.ItemStyle,
		}
		dictsItems[rec.DictId].Items[rec.ItemCode] = &item
		dictsItems[rec.DictId].Options = append(dictsItems[rec.DictId].Options, &item)
	}

	return dictsItems, nil
}
