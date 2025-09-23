package page_model

import "encoding/xml"

type PageModel struct {
	XMLName    xml.Name               `xml:"PageModel"        json:"-"`
	Title      string                 `xml:"title,attr,omitempty"    json:"title,omitempty"`
	InitApi    string                 `xml:"initApi,attr,omitempty"  json:"initApi,omitempty"`
	MainView   string                 `xml:"mainView,attr,omitempty"  json:"mainView,omitempty"`
	DataSet    map[string]*DataTable  `xml:"DataSet"                 json:"dataSet"`
	DataStore  map[string]*DataObject `xml:"DataStore"                 json:"dataStore"`
	Operations map[string]*Operation  `xml:"Operations"              json:"operations"`
	PageView   PageView               `xml:"PageView"                json:"pageView"`
}

type DataTable struct {
	XMLName    xml.Name `xml:"DataTable"                json:"-"`
	Id         string   `xml:"id,attr,omitempty"        json:"id,omitempty"`
	Label      string   `xml:"label,attr,omitempty"     json:"label,omitempty"`
	MappingApi string   `xml:"mappingApi,attr,omitempty" json:"mappingApi,omitempty"`
	Fields     []Field  `xml:"Field"                    json:"fields"`
	TableName  string   `xml:"tableName,attr,omitempty" json:"tableName,omitempty"`
	OrderBy    string   `xml:"orderBy,attr,omitempty" json:"orderBy,omitempty"`
}

type DataObject struct {
	XMLName     xml.Name `xml:"Data"                json:"-"`
	Id          string   `xml:"id,attr,omitempty"        json:"id,omitempty"`
	Type        string   `xml:"type,attr,omitempty"     json:"type,omitempty"`
	Description string   `xml:"description,attr,omitempty" json:"description,omitempty"`
	Api         string   `xml:"api,attr,omitempty"       json:"api,omitempty"`
	Load        string   `xml:"load,attr,omitempty"     json:"load,omitempty"`
	Cache       string   `xml:"cache,attr,omitempty"     json:"cache,omitempty"`
	DataTable   string   `xml:"dataTable,attr,omitempty" json:"dataTable,omitempty"`
}

type Field struct {
	XMLName      xml.Name `xml:"Field"                                        json:"-"`
	Id           string   `xml:"id,attr,omitempty"                            json:"id,omitempty"`
	Component    string   `xml:"component,attr,omitempty"                     json:"component,omitempty"`
	Label        string   `xml:"label,attr,omitempty"                         json:"label,omitempty"`
	Name         string   `xml:"name,attr,omitempty"                          json:"name,omitempty"`
	Source       string   `xml:"source,attr,omitempty"                        json:"source,omitempty"`
	Required     bool     `xml:"required,attr,omitempty"                      json:"required,omitempty"`
	IsFilter     bool     `xml:"isFilter,attr,omitempty"                      json:"isFilter,omitempty"`
	IsQuery      bool     `xml:"isQuery,attr,omitempty"                       json:"isQuery,omitempty"`
	Searchable   bool     `xml:"searchable,attr,omitempty"                    json:"searchable,omitempty"`
	Clearable    bool     `xml:"clearable,attr,omitempty"                     json:"clearable,omitempty"`
	Format       string   `xml:"format,attr,omitempty"                        json:"format,omitempty"`
	InputFormat  string   `xml:"inputFormat,attr,omitempty"                   json:"inputFormat,omitempty"`
	Multiple     bool     `xml:"multiple,attr,omitempty"                      json:"multiple,omitempty"`
	DefaultValue string   `xml:"defaultValue,attr,omitempty"                  json:"defaultValue,omitempty"`
	StaticValue  string   `xml:"staticValue,attr,omitempty"                  json:"staticValue,omitempty"`
	Disabled     bool     `xml:"disabled,attr,omitempty"                      json:"disabled,omitempty"`
}

type Operation struct {
	XMLName    xml.Name `xml:"Operation"                json:"-"`
	Id         string   `xml:"id,attr,omitempty"        json:"id,omitempty"`
	Type       string   `xml:"type,attr,omitempty"      json:"type,omitempty"`
	Label      string   `xml:"label,attr,omitempty"     json:"label,omitempty"`
	Level      string   `xml:"level,attr,omitempty"     json:"level,omitempty"`
	ActionType string   `xml:"actionType,attr,omitempty" json:"actionType,omitempty"`
	Icon       string   `xml:"icon,attr,omitempty"      json:"icon,omitempty"`
	Title      string   `xml:"title,attr,omitempty"     json:"title,omitempty"`
	Feature    string   `xml:"feature,attr,omitempty"   json:"feature,omitempty"`
	Api        string   `xml:"api,attr,omitempty"       json:"api,omitempty"`
	Method     string   `xml:"method,attr,omitempty"       json:"method,omitempty"`
	View       string   `xml:"view,attr,omitempty"       json:"view,omitempty"`
	Confirm    string   `xml:"confirm,attr,omitempty"       json:"confirm,omitempty"`
	Effects    string   `xml:"effects,attr,omitempty"       json:"effects,omitempty"`
}

// PageView
type PageView struct {
	XMLName xml.Name `xml:"PageView" json:"-"`
	ViewObject
}

type ViewObject struct {
	Object   string         `json:"object,omitempty"`
	Props    map[string]any `json:"props,omitempty"`
	Children []IViewObject  `json:"children,omitempty"`
}
