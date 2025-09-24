package pandora

import (
	"encoding/json"
	"os"
	"path"

	"github.com/mark0725/go-app-pandora/page_model"
)

type PageParamBase struct {
	PanUrl string
}

func InitPageEngine() error {

	logger.Info("InitPageEngine end")
	return nil
}

func GetPageView(pagePath string, params map[string]string) (*page_model.PageModel, error) {
	fullPath := path.Join(g_appConfig.Pandora.PagesPath, "pages", pagePath+".ds.xml")
	logger.Infof("GetPageView [%s]", fullPath)
	//判断文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		logger.Errorf("GetPageView [%s] not found", fullPath)
	}

	PageView, err := page_model.ParsePageModel(fullPath, params)
	if err != nil {
		logger.Errorf("PageEngine [%s] parse error", fullPath)
		return nil, err
	}

	return PageView, nil
}

func PageEngine(title string, pagePath string, env any, props any, params map[string]string) ([]byte, error) {
	mapObj, err := page_model.ParsePageModel2Map(pagePath, params)
	if err != nil {
		logger.Errorf("GetPageView [%s] error: %v", pagePath, err)
		return nil, err
	}

	schemaStr, _ := json.Marshal(mapObj)
	return schemaStr, nil
}
