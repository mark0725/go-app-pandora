package pandora

import (
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	pandora_assets "github.com/mark0725/go-pandora"
)

func InitRoutes(group string, r *gin.RouterGroup) {
	switch group {
	case "page-view":
		r.GET("/:module/*path", Pages)
	case "page-api":
		r.GET("/:module/mapping", PageApiMapping)
		r.POST("/:module/action/*path", g_PageApi.ActionCreate)
		r.PUT("/:module/action/*path", g_PageApi.ActionUpdate)
		r.DELETE("/:module/action/*path", g_PageApi.ActionDelete)
		r.GET("/:module/query/*path", g_PageApi.QueryOne)
		r.GET("/:module/query-list/*path", g_PageApi.Query)
		r.GET("/:module/query-page/*path", g_PageApi.PageQuery)
	case "app":
		r.GET("/config", g_AppApi.Config)
		r.GET("/config/:module", g_AppApi.ModuleConfig)
	case "page":
		r.GET("/:module/pages/*path", Pages)
		r.GET("/:module/api/mapping", PageApiMapping)
		r.POST("/:module/api/action/*path", g_PageApi.ActionCreate)
		r.PUT("/:module/api/query/*path", g_PageApi.ActionUpdate)
		r.DELETE("/:module/api/query-page/*path", g_PageApi.ActionDelete)
		r.GET("/:module/api/query/*path", g_PageApi.QueryOne)
		r.GET("/:module/api/query-list/*path", g_PageApi.Query)
		r.GET("/:module/api/query-page/*path", g_PageApi.PageQuery)

	default:

	}
}

func Pages(c *gin.Context) {
	//module := c.Param("module")
	filepath := c.Param("path")
	logger.Debugf("Pages path: %s", filepath)
	props := map[string]string{}

	queryParams := map[string]string{}
	for key, value := range c.Request.URL.Query() {
		if len(value) > 0 {
			queryParams[key] = value[0]
		}
	}

	fullPath := path.Join(g_appConfig.Pandora.PagesPath, "pages", filepath+".ds.xml")
	logger.Infof("GetPageView [%s]", fullPath)
	//判断文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		logger.Errorf("GetPageView [%s] not found", fullPath)
		c.JSON(404, "Not Found")
	}

	if content, err := PageEngine(g_appConfig.Pandora.Title, fullPath, g_appConfig.Pandora.Env, props, queryParams); err == nil {
		c.Data(http.StatusOK, "text/html; charset=utf-8", content)
	} else {
		logger.Errorf("PageEngine error: %v", err)
		c.String(http.StatusInternalServerError, "")
	}
}

func PandoraAssets(c *gin.Context) {
	remote := g_appConfig.Pandora.PanUrl
	if strings.HasPrefix(remote, "http://") || strings.HasPrefix(remote, "https://") {
		PandoraAssetsHttp(c)
		return
	}

	reqPath := c.Request.URL.Path
	if reqPath == "/" {
		reqPath = "/index.html"
	}
	reqPath = path.Clean("/" + strings.Trim(reqPath, "/")) // 防止 .. 注入

	if strings.HasPrefix(remote, "embed://") {
		staticFS := http.FS(pandora_assets.StaticFiles)
		staticFilePath := remote[len("embed://"):] + reqPath
		c.FileFromFS(staticFilePath, staticFS)
		return
	}

	c.File(remote + reqPath)
}
func PandoraAssetsHttp(c *gin.Context) {

	remote := g_appConfig.Pandora.PanUrl
	targetURL := remote + c.Request.RequestURI

	// 创建新请求
	req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// 拷贝 headers
	for k, v := range c.Request.Header {
		req.Header[k] = v
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(502, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// 设置响应头
	for k, v := range resp.Header {
		c.Writer.Header()[k] = v
	}
	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}
