package page_model

import (
	"encoding/json"
	"testing"
)

// 写单元测试
func TestParsePageModel(t *testing.T) {
	// 读取文件
	path := "../../pan-pages/pages/models/model.ds.xml"
	mapObj, err := ParsePageModel2Map(path, nil)
	if err != nil {
		t.Fatal("Error parsing XML:", err)
	}
	jsonData, err := json.MarshalIndent(mapObj, "", "  ")
	if err != nil {
		t.Fatal("Error marshalling to JSON:", err)
	}
	t.Logf("JSON: %s\n", jsonData)
}
