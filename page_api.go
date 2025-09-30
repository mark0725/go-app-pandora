package pandora

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/mark0725/go-app-pandora/page_model"

	"github.com/gin-gonic/gin"
	base_db "github.com/mark0725/go-app-base/db"
)

type PanPageQuery struct {
	Name    string
	Mapping PanPageMappingFunc
}

var g_PanPageQueries = map[string]PanPageQueryFunc{}

type PanPageQueryFunc func(ctx context.Context, path string, params map[string]any, dataTable *page_model.DataTable) (string, string, error)

func RegisterPageQuery(module string, fn PanPageQueryFunc) {
	g_PanPageQueries[module] = fn
}

type PageApi struct{}

var g_PageApi = PageApi{}

func (api *PageApi) PageQuery(c *gin.Context) {
	//{"code":"OK","message":"success","data":[]}
	data := ApiReponse{
		Code:    "OK",
		Message: "success",
		Data:    []string{},
	}

	module := c.Param("module")
	path := c.Param("path")
	logger.Debugf("Pages path: %s", path)
	pageIndex := 0
	pageSize := 50

	if s := c.Query("page"); s != "" {
		if num, err := strconv.Atoi(s); err == nil {
			pageIndex = num - 1
		}
	}
	if s := c.Query("size"); s != "" {
		if num, err := strconv.Atoi(s); err == nil {
			pageSize = num
		}
	}

	queryParam, err := GetPageQueryParam(c, module, path)
	if err != nil {
		logger.Errorf("Get page query error: %s", err)
		return
	}

	orderBy := ""
	if s := c.Query("order"); s != "" {
		orderBy = fmt.Sprintf(" order by %s", s)
		if !strings.Contains(queryParam.sql, "order") {
			queryParam.sql += orderBy
		}
	}

	logger.Debugf("Page Index: %d, Page Size: %d", pageIndex, pageSize)
	logger.Debugf("Query Param: %+v", queryParam)
	pageQueryResult, err := base_db.DBPageQuery(queryParam.db, queryParam.sql, queryParam.params, pageSize, pageIndex)
	if err != nil {
		logger.Error("DBPageQuery fail: ", err)
		c.JSON(http.StatusInternalServerError, ApiReponse{Code: "ERROR", Message: "DBPageQuery fail"})
		return
	}

	recs := []map[string]any{}
	for _, row := range pageQueryResult.Content {
		rec := map[string]any{}
		for k, v := range row {
			rec[strings.ToUpper(k)] = v
		}
		recs = append(recs, rec)
	}

	pageQueryResult.Content = recs
	data.Data = pageQueryResult

	c.JSON(http.StatusOK, data)
}

func (api *PageApi) Query(c *gin.Context) {
	//{"code":"OK","message":"success","data":[]}
	data := ApiReponse{
		Code:    "OK",
		Message: "success",
		Data:    []string{},
	}

	module := c.Param("module")
	path := c.Param("path")
	logger.Debugf("Pages path: %s", path)

	queryParam, err := GetPageQueryParam(c, module, path)
	if err != nil {
		logger.Errorf("Get page query error: %s", err)
		return
	}

	sqlBuildResult := base_db.QueryNamedParamsBuilder(queryParam.sql, queryParam.params)
	orderBy := ""
	if s := c.Query("order"); s != "" {
		orderBy = fmt.Sprintf(" order by %s", s)
		if !strings.Contains(queryParam.sql, "order") {
			sqlBuildResult.Sql += orderBy
		}
	}
	logger.Debugf("Query Param: %+v", sqlBuildResult)
	rows, err := base_db.DBQuery(queryParam.db, sqlBuildResult.Sql, sqlBuildResult.Params)
	if err != nil {
		logger.Error("DBPageQuery fail: ", err)
		c.JSON(http.StatusInternalServerError, ApiReponse{Code: "ERROR", Message: "DBPageQuery fail"})
		return
	}

	recs := []map[string]any{}
	for _, row := range rows {
		rec := map[string]any{}
		for k, v := range row {
			rec[strings.ToUpper(k)] = v
		}
		recs = append(recs, rec)
	}

	data.Data = recs

	c.JSON(http.StatusOK, data)
}

func (api *PageApi) QueryOne(c *gin.Context) {
	//{"code":"OK","message":"success","data":[]}
	data := ApiReponse{
		Code:    "OK",
		Message: "success",
		Data:    []string{},
	}

	module := c.Param("module")
	path := c.Param("path")
	logger.Debugf("Pages path: %s", path)

	queryParam, err := GetPageQueryParam(c, module, path)
	if err != nil {
		logger.Errorf("Get page query error: %s", err)
		return
	}

	sqlBuildResult := base_db.QueryNamedParamsBuilder(queryParam.sql, queryParam.params)
	rows, err := base_db.DBQuery(queryParam.db, sqlBuildResult.Sql, sqlBuildResult.Params)
	if err != nil {
		logger.Error("DBPageQuery fail: ", err)
		c.JSON(http.StatusInternalServerError, ApiReponse{Code: "ERROR", Message: "DBPageQuery fail"})
		return
	}

	if len(rows) == 0 {
		c.JSON(http.StatusOK, ApiReponse{Code: "SUCCESS", Data: map[string]any{}})
		return
	}

	row := rows[0]
	rec := map[string]any{}
	for k, v := range row {
		rec[strings.ToUpper(k)] = v
	}
	data.Data = rec
	c.JSON(http.StatusOK, data)
}

func (api *PageApi) ActionCreate(c *gin.Context) {
	module := c.Param("module")
	path := c.Param("path")
	logger.Debugf("Pages %s path: %s", module, path)

	c.JSON(http.StatusOK, ApiReponse{Code: "OK", Message: "success"})
}

func (api *PageApi) ActionDelete(c *gin.Context) {
	module := c.Param("module")
	path := c.Param("path")
	logger.Debugf("Pages %s path: %s", module, path)

	c.JSON(http.StatusOK, ApiReponse{Code: "OK", Message: "success"})
}

func (api *PageApi) ActionUpdate(c *gin.Context) {
	module := c.Param("module")
	path := c.Param("path")
	logger.Debugf("Pages %s path: %s", module, path)

	c.JSON(http.StatusOK, ApiReponse{Code: "OK", Message: "success"})
}

type PageQueryParam struct {
	db     string
	sql    string
	params map[string]any
}

func GetPageQueryParam(c *gin.Context, module string, path string) (*PageQueryParam, error) {

	viewParams := map[string]string{}
	for key, value := range c.Request.URL.Query() {
		if len(value) > 0 {
			viewParams[key] = value[0]
		}
	}

	pageview, err := GetPageView(path, viewParams)
	if err != nil {
		logger.Errorf("GetPageView error: %v", err)
		c.JSON(http.StatusNotFound, ApiReponse{Code: "NOT_FOUND", Message: "not found error"})
		return nil, err
	}

	ds := c.Query("ds")
	if len(ds) == 0 {
		c.JSON(http.StatusBadRequest, ApiReponse{Code: "BadRequest", Message: "ds is not found"})
		return nil, err
	}

	if _, ok := pageview.DataSet[ds]; !ok {
		logger.Error("DataSet not found: ", ds)
		c.JSON(http.StatusBadRequest, ApiReponse{Code: "BadRequest", Message: "DataSet not found"})
		return nil, err
	}

	dataTable := pageview.DataSet[ds]
	queryParams := map[string]any{}
	fieldIds := []string{}
	staticParams := map[string]string{}
	for _, f := range dataTable.Fields {
		fieldIds = append(fieldIds, f.Id)
		if f.StaticValue != "" {
			staticParams[f.Id] = f.StaticValue
		}
		if f.IsFilter || f.IsQuery {
			if v, exist := c.GetQuery(f.Id); exist {
				queryParams[f.Id] = v
			}
		}
	}
	logger.Debug("queryParams: ", queryParams)

	ctx := context.Background()
	// ctx = context.WithValue(ctx, userIDKey, 12345)

	querySql := ""
	db := base_db.DB_CONN_NAME_DEFAULT
	if fn, ok := g_PanPageQueries[module]; !ok {
		logger.Errorf("Not found module: %s", module)
		c.JSON(http.StatusNotFound, ApiReponse{Code: "NOT_FOUND", Message: "not found error"})
		return nil, errors.New("Not found module")

	} else {
		db, querySql, err = fn(ctx, path, queryParams, dataTable)
		if err != nil {
			logger.Errorf("Build query error: %v", err)
			c.JSON(http.StatusInternalServerError, ApiReponse{Code: "ERROR", Message: "Build query error"})
			return nil, err
		}

		if querySql == "" {

			if dataTable.TableName != "" {
				queryParams["ORG_ID"] = g_appConfig.Org.OrgId
				for k, v := range staticParams {
					queryParams[k] = v
				}
				sqlParts := base_db.NewQueryParamsBuilder().Params(queryParams).Build()
				querySql = fmt.Sprintf("SELECT %s FROM %s WHERE 1=1 ", strings.Join(fieldIds, ","), dataTable.TableName)
				if sqlParts != "" {
					querySql += " AND " + sqlParts
				}
				logger.Debug("use default query: ", querySql)
			} else {
				logger.Errorf("Build query error: %v", err)
				c.JSON(http.StatusInternalServerError, ApiReponse{Code: "ERROR", Message: "Build query error"})
				return nil, err
			}
		}
	}

	return &PageQueryParam{db: db, sql: querySql, params: queryParams}, nil
}
