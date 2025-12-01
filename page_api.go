package pandora

import (
	"context"
	"encoding/base64"
	"encoding/json"
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

type PanPageQueryFunc func(ctx context.Context, path string, params map[string]any, dataTable *page_model.DataTable) (string, *base_db.QueryParamsOptions, error)

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
		return nil, errors.New("dataset param not set")
	}

	if _, ok := pageview.DataSet[ds]; !ok {
		logger.Error("DataSet not found: ", ds)
		c.JSON(http.StatusBadRequest, ApiReponse{Code: "BadRequest", Message: "DataSet not found"})
		return nil, errors.New("dataSet not found")
	}

	dataTable := pageview.DataSet[ds]
	queryParams := map[string]any{}
	fieldIds := []string{}
	staticParams := map[string]string{}
	for _, f := range dataTable.Fields {
		fieldIds = append(fieldIds, f.Id)

		if f.IsFilter {
			if v, exist := c.GetQuery(f.Id); exist {
				queryParams[f.Id] = v
			}
		}
		if f.IsQuery {
			if v, exist := c.GetQuery(f.Id); exist {
				queryParams[f.Id] = v
			} else {
				logger.Error("Query param not found: ", f.Id)
				c.JSON(http.StatusBadRequest, ApiReponse{Code: "BadRequest", Message: fmt.Sprintf("Query param %s not found", f.Id)})
				return nil, fmt.Errorf("Query param %s not found", f.Id)
			}
		}

		if f.StaticValue != "" {
			staticParams[f.Id] = f.StaticValue
		}
		// if f.StaticValue != "" {
		// 	queryParams[f.Id] = f.StaticValue
		// }
	}

	spOps := map[string]string{}
	filterParams := FilterParam{}
	if filterEnc, exist := c.GetQuery("filter"); exist {
		rawURLDec := base64.RawURLEncoding
		dst1 := make([]byte, rawURLDec.DecodedLen(len(filterEnc)))
		_, err := rawURLDec.Decode(dst1, []byte(filterEnc))
		if err != nil {
			logger.Error(err)
			c.JSON(http.StatusBadRequest, ApiReponse{Code: "BadRequest", Message: "filter decode error"})
		}
		logger.Debug("filter: ", string(dst1))

		if err := json.Unmarshal(dst1, &filterParams); err != nil {
			logger.Error(err)
			c.JSON(http.StatusBadRequest, ApiReponse{Code: "BadRequest", Message: "filter json decode error"})
		}

		for _, item := range filterParams.Items {
			if _, ok := queryParams[item.Key]; ok {
				continue
			}
			queryParams[item.Key] = item.Value
			switch item.TypeValue {
			case "is-null":
				spOps[item.Key] = "is null"
			case "is-not-null":
				spOps[item.Key] = "is not null"
			case "not-in":
				spOps[item.Key] = "not in"
			default:
				spOps[item.Key] = item.TypeValue
			}
		}
	}

	logger.Debug("queryParams: ", queryParams)
	logger.Debugf("filter: %+v", filterParams)

	ctx := context.Background()
	// ctx = context.WithValue(ctx, userIDKey, 12345)

	db := base_db.DB_CONN_NAME_DEFAULT
	fn, ok := g_PanPageQueries[module]
	if !ok {
		logger.Errorf("Not found module: %s", module)
		c.JSON(http.StatusNotFound, ApiReponse{Code: "NOT_FOUND", Message: "not found error"})
		return nil, errors.New("not found module")

	}

	queryBuilder := base_db.NewQueryParamsBuilder()
	db, sqlOption, err := fn(ctx, path, queryParams, dataTable)
	if err != nil {
		logger.Errorf("Build query error: %v", err)
		c.JSON(http.StatusInternalServerError, ApiReponse{Code: "ERROR", Message: "Build query error"})
		return nil, err
	}

	querySql := ""
	var sqlParams map[string]any
	if sqlOption == nil {
		if dataTable.TableName != "" {
			queryParams["ORG_ID"] = g_appConfig.Org.OrgId
			for k, v := range staticParams {
				queryParams[k] = v
			}
			querySql, sqlParams = queryBuilder.From(dataTable.TableName).Columns(fieldIds...).WhereMap(queryParams).SpOps(spOps).Build()
			logger.Debug("use default query: ", querySql)
		} else {
			logger.Errorf("Build query error: %v", err)
			c.JSON(http.StatusNotFound, ApiReponse{Code: "ERROR", Message: "no query found"})
			return nil, errors.New("no query found")
		}
	} else {
		querySql, sqlParams = sqlOption.SpOps(spOps).Build()
	}

	return &PageQueryParam{db: db, sql: querySql, params: sqlParams}, nil
}
