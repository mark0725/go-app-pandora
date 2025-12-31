package i18n

import (
	"encoding/xml"
	"fmt"
)

// I18nData 存储单个语言的所有翻译内容
type I18nData struct {
	Common     map[string]string            `json:"Common"`
	Dataset    map[string]map[string]string `json:"Dataset"` // DataTable.id -> Item.id -> value
	DataStore  map[string]string            `json:"DataStore"`
	Operations map[string]string            `json:"Operations"`
	PageView   map[string]string            `json:"PageView"`
}

// 支持的语言列表
var supportedLanguages = []string{
	"en", "zh", "zh_cn", "zh_tw", "ja", "ko", "ar", "ru", "pt", "es", "fr", "de", "it", "he", "hi",
}

// XML 解析结构体
type xmlI18n struct {
	XMLName    xml.Name   `xml:"i18n"`
	Common     xmlSection `xml:"Common"`
	Dataset    xmlDataset `xml:"Dataset"`
	DataStore  xmlSection `xml:"DataStore"`
	Operations xmlSection `xml:"Operations"`
	PageView   xmlSection `xml:"PageView"`
}

type xmlSection struct {
	Items []xmlItem `xml:",any"`
}

type xmlDataset struct {
	DataTables []xmlDataTable `xml:"DataTable"`
}

type xmlDataTable struct {
	ID    string    `xml:"id,attr"`
	Items []xmlItem `xml:"Item"`
}

type xmlItem struct {
	XMLName xml.Name
	ID      string `xml:"id,attr"`
	En      string `xml:"en,attr"`
	Zh      string `xml:"zh,attr"`
	ZhCn    string `xml:"zh_cn,attr"`
	ZhTw    string `xml:"zh_tw,attr"`
	Ja      string `xml:"ja,attr"`
	Ko      string `xml:"ko,attr"`
	Ar      string `xml:"ar,attr"`
	Ru      string `xml:"ru,attr"`
	Pt      string `xml:"pt,attr"`
	Es      string `xml:"es,attr"`
	Fr      string `xml:"fr,attr"`
	De      string `xml:"de,attr"`
	It      string `xml:"it,attr"`
	He      string `xml:"he,attr"`
	Hi      string `xml:"hi,attr"`
}

// getLangValue 根据语言代码获取对应的翻译值
func (item *xmlItem) getLangValue(lang string) string {
	switch lang {
	case "en":
		return item.En
	case "zh":
		return item.Zh
	case "zh_cn":
		return item.ZhCn
	case "zh_tw":
		return item.ZhTw
	case "ja":
		return item.Ja
	case "ko":
		return item.Ko
	case "ar":
		return item.Ar
	case "ru":
		return item.Ru
	case "pt":
		return item.Pt
	case "es":
		return item.Es
	case "fr":
		return item.Fr
	case "de":
		return item.De
	case "it":
		return item.It
	case "he":
		return item.He
	case "hi":
		return item.Hi
	default:
		return ""
	}
}

// getItemKey 获取Item的key：有id用id，没有id用标签名
func (item *xmlItem) getItemKey() string {
	if item.ID != "" {
		return item.ID
	}
	return item.XMLName.Local
}

// Parse 解析i18n.xml内容，返回不同语言的翻译数据
func Parse(i18nContent []byte) (map[string]*I18nData, error) {
	var xmlData xmlI18n
	if err := xml.Unmarshal(i18nContent, &xmlData); err != nil {
		return nil, fmt.Errorf("failed to parse i18n xml: %w", err)
	}

	result := make(map[string]*I18nData)

	// 为每种语言初始化数据结构
	for _, lang := range supportedLanguages {
		result[lang] = &I18nData{
			Common:     make(map[string]string),
			Dataset:    make(map[string]map[string]string),
			DataStore:  make(map[string]string),
			Operations: make(map[string]string),
			PageView:   make(map[string]string),
		}
	}

	// 解析 Common 部分
	for _, item := range xmlData.Common.Items {
		key := item.getItemKey()
		for _, lang := range supportedLanguages {
			value := item.getLangValue(lang)
			if value != "" {
				result[lang].Common[key] = value
			}
		}
	}

	// 解析 Dataset 部分 (嵌套结构: DataTable -> Item)
	for _, dataTable := range xmlData.Dataset.DataTables {
		tableID := dataTable.ID
		for _, lang := range supportedLanguages {
			if result[lang].Dataset[tableID] == nil {
				result[lang].Dataset[tableID] = make(map[string]string)
			}
		}
		for _, item := range dataTable.Items {
			key := item.getItemKey()
			for _, lang := range supportedLanguages {
				value := item.getLangValue(lang)
				if value != "" {
					result[lang].Dataset[tableID][key] = value
				}
			}
		}
	}

	// 解析 DataStore 部分
	for _, item := range xmlData.DataStore.Items {
		key := item.getItemKey()
		for _, lang := range supportedLanguages {
			value := item.getLangValue(lang)
			if value != "" {
				result[lang].DataStore[key] = value
			}
		}
	}

	// 解析 Operations 部分
	for _, item := range xmlData.Operations.Items {
		key := item.getItemKey()
		for _, lang := range supportedLanguages {
			value := item.getLangValue(lang)
			if value != "" {
				result[lang].Operations[key] = value
			}
		}
	}

	// 解析 PageView 部分
	for _, item := range xmlData.PageView.Items {
		key := item.getItemKey()
		for _, lang := range supportedLanguages {
			value := item.getLangValue(lang)
			if value != "" {
				result[lang].PageView[key] = value
			}
		}
	}

	return result, nil
}

// GetI18nForLang 获取指定语言的翻译数据（便捷方法）
func GetI18nForLang(i18nContent []byte, lang string) (*I18nData, error) {
	allLangs, err := Parse(i18nContent)
	if err != nil {
		return nil, err
	}

	data, exists := allLangs[lang]
	if !exists {
		return nil, fmt.Errorf("unsupported language: %s", lang)
	}

	return data, nil
}
