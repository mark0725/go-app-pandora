package i18n

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	i18nXML := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<i18n>
   <Common>
      <Item id="APP_TITLE" en="Application" zh="应用程序" ja="アプリ" />
   </Common>
   <Dataset>
      <DataTable id="AI_SERVICE">
        <Item id="SERVICE_ID" en="Service Id" zh="服务ID" zh_cn="" zh_tw="" ja="" ko="" ar="" ru="" pt="" es="" fr="" de="" it="" he="" hi="" />
        <Item id="SERVICE_NAME" en="Service Name" zh="服务名称" />
      </DataTable>
      <DataTable id="USER_TABLE">
        <Item id="USER_ID" en="User Id" zh="用户ID" />
      </DataTable>
   </Dataset>
   <DataStore>
      <Item id="STORE_KEY" en="Store Key" zh="存储键" />
   </DataStore>
   <Operations>
      <Item id="BTN_ADD" en="Create" zh="创建" zh_cn="" zh_tw="" ja="" ko="" ar="" ru="" pt="" es="" fr="" de="" it="" he="" hi="" />
      <Item id="BTN_DEL" en="Delete" zh="删除" />
   </Operations>
   <PageView>
      <Item id="addView" en="Create" zh="创建" />
      <Item id="editView" en="Edit" zh="编辑" />
   </PageView>
</i18n>`)

	result, err := Parse(i18nXML)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// 测试英文结果
	enData := result["en"]

	// 检查 Common
	if enData.Common["APP_TITLE"] != "Application" {
		t.Errorf("Expected Common.APP_TITLE='Application', got '%s'", enData.Common["APP_TITLE"])
	}

	// 检查 Dataset (嵌套结构)
	if enData.Dataset["AI_SERVICE"]["SERVICE_ID"] != "Service Id" {
		t.Errorf("Expected Dataset.AI_SERVICE.SERVICE_ID='Service Id', got '%s'", enData.Dataset["AI_SERVICE"]["SERVICE_ID"])
	}
	if enData.Dataset["AI_SERVICE"]["SERVICE_NAME"] != "Service Name" {
		t.Errorf("Expected Dataset.AI_SERVICE.SERVICE_NAME='Service Name', got '%s'", enData.Dataset["AI_SERVICE"]["SERVICE_NAME"])
	}
	if enData.Dataset["USER_TABLE"]["USER_ID"] != "User Id" {
		t.Errorf("Expected Dataset.USER_TABLE.USER_ID='User Id', got '%s'", enData.Dataset["USER_TABLE"]["USER_ID"])
	}

	// 检查 DataStore
	if enData.DataStore["STORE_KEY"] != "Store Key" {
		t.Errorf("Expected DataStore.STORE_KEY='Store Key', got '%s'", enData.DataStore["STORE_KEY"])
	}

	// 检查 Operations
	if enData.Operations["BTN_ADD"] != "Create" {
		t.Errorf("Expected Operations.BTN_ADD='Create', got '%s'", enData.Operations["BTN_ADD"])
	}
	if enData.Operations["BTN_DEL"] != "Delete" {
		t.Errorf("Expected Operations.BTN_DEL='Delete', got '%s'", enData.Operations["BTN_DEL"])
	}

	// 检查 PageView
	if enData.PageView["addView"] != "Create" {
		t.Errorf("Expected PageView.addView='Create', got '%s'", enData.PageView["addView"])
	}

	// 测试中文结果
	zhData := result["zh"]
	if zhData.Dataset["AI_SERVICE"]["SERVICE_ID"] != "服务ID" {
		t.Errorf("Expected zh Dataset.AI_SERVICE.SERVICE_ID='服务ID', got '%s'", zhData.Dataset["AI_SERVICE"]["SERVICE_ID"])
	}
	if zhData.Operations["BTN_ADD"] != "创建" {
		t.Errorf("Expected zh Operations.BTN_ADD='创建', got '%s'", zhData.Operations["BTN_ADD"])
	}

	// 测试日语结果
	jaData := result["ja"]
	if jaData.Common["APP_TITLE"] != "アプリ" {
		t.Errorf("Expected ja Common.APP_TITLE='アプリ', got '%s'", jaData.Common["APP_TITLE"])
	}
}

func TestParseEmptyContent(t *testing.T) {
	i18nXML := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<i18n>
   <Common></Common>
   <Dataset></Dataset>
   <DataStore></DataStore>
   <Operations></Operations>
   <PageView></PageView>
</i18n>`)

	result, err := Parse(i18nXML)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// 验证所有语言的数据结构都已初始化
	for _, lang := range supportedLanguages {
		if result[lang].Common == nil {
			t.Errorf("Expected Common map initialized for lang %s", lang)
		}
		if result[lang].Dataset == nil {
			t.Errorf("Expected Dataset map initialized for lang %s", lang)
		}
		if result[lang].DataStore == nil {
			t.Errorf("Expected DataStore map initialized for lang %s", lang)
		}
		if result[lang].Operations == nil {
			t.Errorf("Expected Operations map initialized for lang %s", lang)
		}
		if result[lang].PageView == nil {
			t.Errorf("Expected PageView map initialized for lang %s", lang)
		}
	}
}

func TestParseInvalidXML(t *testing.T) {
	invalidXML := []byte(`this is not valid xml`)

	_, err := Parse(invalidXML)
	if err == nil {
		t.Error("Expected error for invalid XML, got nil")
	}
}

func TestGetI18nForLang(t *testing.T) {
	i18nXML := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<i18n>
   <Common></Common>
   <Dataset>
      <DataTable id="TEST">
        <Item id="FIELD1" en="Field 1" zh="字段1" />
      </DataTable>
   </Dataset>
   <DataStore></DataStore>
   <Operations>
      <Item id="BTN_TEST" en="Test" zh="测试" />
   </Operations>
   <PageView></PageView>
</i18n>`)

	// 测试获取英语
	enData, err := GetI18nForLang(i18nXML, "en")
	if err != nil {
		t.Fatalf("GetI18nForLang failed: %v", err)
	}
	if enData.Dataset["TEST"]["FIELD1"] != "Field 1" {
		t.Errorf("Expected 'Field 1', got '%s'", enData.Dataset["TEST"]["FIELD1"])
	}

	// 测试不支持的语言
	_, err = GetI18nForLang(i18nXML, "unsupported_lang")
	if err == nil {
		t.Error("Expected error for unsupported language")
	}
}

func TestItemWithoutID(t *testing.T) {
	// 测试没有id属性的Item，应该使用标签名作为key
	i18nXML := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<i18n>
   <Common>
      <Title en="Title" zh="标题" />
      <Description en="Description" zh="描述" />
   </Common>
   <Dataset></Dataset>
   <DataStore></DataStore>
   <Operations></Operations>
   <PageView></PageView>
</i18n>`)

	result, err := Parse(i18nXML)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	enData := result["en"]
	if enData.Common["Title"] != "Title" {
		t.Errorf("Expected Common.Title='Title', got '%s'", enData.Common["Title"])
	}
	if enData.Common["Description"] != "Description" {
		t.Errorf("Expected Common.Description='Description', got '%s'", enData.Common["Description"])
	}

	zhData := result["zh"]
	if zhData.Common["Title"] != "标题" {
		t.Errorf("Expected zh Common.Title='标题', got '%s'", zhData.Common["Title"])
	}
}

func TestAllSupportedLanguages(t *testing.T) {
	expectedLangs := []string{"en", "zh", "zh_cn", "zh_tw", "ja", "ko", "ar", "ru", "pt", "es", "fr", "de", "it", "he", "hi"}

	if !reflect.DeepEqual(supportedLanguages, expectedLangs) {
		t.Errorf("Supported languages mismatch. Expected %v, got %v", expectedLangs, supportedLanguages)
	}
}
