package page_model

import (
	"encoding/xml"
)

func NewViewObject(object string) *ViewObject {
	return &ViewObject{
		Object: object,
	}
}

// 实现 IViewObject 接口
func (vo *ViewObject) GetObject() string {
	return vo.Object
}
func (vo *ViewObject) GetProps() map[string]any {
	return vo.Props
}

func (vo *ViewObject) SetProp(name string, value any) {
	if vo.Props == nil {
		vo.Props = make(map[string]any)
	}

	newValue := value
	vo.Props[name] = newValue
}

func (vo *ViewObject) AppendChildren(value IViewObject) {
	if len(vo.Children) == 0 {
		vo.Children = make([]IViewObject, 0)
	}
	vo.Children = append(vo.Children, value)
}

func (vo *ViewObject) AppendElement(name string, value any) {
	if vo.Props[name] == nil {
		vo.Props[name] = make([]any, 0)
	}
	slice := vo.Props[name].([]any)
	slice = append(slice, value)
	vo.Props[name] = slice
}

func (vo *ViewObject) GetChildren() []IViewObject {
	return vo.Children
}

func (vo *ViewObject) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return UnmarshalViewObject(d, start, vo)
}

func (p *PageView) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// 先提取已知属性
	var known = map[string]bool{
		"name": true,
	}

	// 循环处理全部属性
	p.Props = make(map[string]any)
	for _, attr := range start.Attr {
		lc := attr.Name.Local
		if _, ok := known[lc]; ok {
			// 如果后来需要增加 PageView 自身字段可在此处理
			continue
		}
		// 未定义的属性进 props
		p.Props[lc] = attr.Value
	}
	// 通用解析子标签
	return UnmarshalViewObject(d, start, p)
}

type Form struct {
	XMLName xml.Name `xml:"Form"                json:"-"`
	ViewObject
	Name      string `xml:"name,attr,omitempty"       json:"name,omitempty"`
	DataTable string `xml:"dataTable,attr,omitempty"   json:"dataTable,omitempty"`
	Cols      string `xml:"cols,attr,omitempty"        json:"cols,omitempty"`
	Mode      string `xml:"mode,attr,omitempty"        json:"mode,omitempty"`
}

// 这里覆盖 UnmarshalXML，使它先解析已知字段，再将同行/未知属性放入 Props
func (f *Form) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// 先用一个临时的 struct 用于标准的 xml tag 解码（仅匹配已定义字段）
	type AliasForm Form
	aliasTarget := &AliasForm{}

	// 先提取属性并赋值给 aliasTarget
	if err := setFieldsFromAttrs(start.Attr, aliasTarget); err != nil {
		return err
	}
	// 再将已识别的字段赋值回 f
	f.Name = aliasTarget.Name
	f.DataTable = aliasTarget.DataTable
	f.Cols = aliasTarget.Cols
	f.Mode = aliasTarget.Mode

	// 这里再调用通用的 UnmarshalViewObject 处理 children、props 等
	// 注意要将 f 自身传进去，以便 children 继续挂载到当前实例。
	return UnmarshalViewObject(d, start, f)
}

func (pm *PageModel) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// 先解析 PageModel 自身的已知字段属性 (title, initApi, mainView)，忽略无效属性
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "title":
			pm.Title = attr.Value
		case "initApi":
			pm.InitApi = attr.Value
		case "mainView":
			pm.MainView = attr.Value
		}
	}

	pm.DataSet = make(map[string]*DataTable)
	pm.DataStore = make(map[string]*DataObject)
	pm.Operations = make(map[string]*Operation)

	// 读到 </PageModel> 停止
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}
		switch se := t.(type) {
		case xml.StartElement:
			tag := se.Name.Local
			switch tag {
			case "DataSet":
				if err := pm.parseDataSet(d); err != nil {
					return err
				}
			case "DataStore":
				if err := pm.parseDataStore(d); err != nil {
					return err
				}
			case "Operations":
				if err := pm.parseOperations(d); err != nil {
					return err
				}

			case "PageView":
				if err := pm.PageView.UnmarshalXML(d, se); err != nil {
					return err
				}

			default:
				if err := skipElement(d, se.Name); err != nil {
					return err
				}
			}

		case xml.EndElement:
			if se.Name.Local == "PageModel" {
				return nil
			}
		}
	}
}

func (pm *PageModel) parseDataSet(d *xml.Decoder) error {
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "DataTable" {
				var dt DataTable
				if err := d.DecodeElement(&dt, &se); err != nil {
					return err
				}
				if dt.Id != "" {
					pm.DataSet[dt.Id] = &dt
				}
			} else {
				// 跳过其他标签
				if err := skipElement(d, se.Name); err != nil {
					return err
				}
			}
		case xml.EndElement:
			if se.Name.Local == "DataSet" {
				return nil
			}
		}
	}
}

func (pm *PageModel) parseDataStore(d *xml.Decoder) error {
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "Data" {
				var store DataObject
				if err := d.DecodeElement(&store, &se); err != nil {
					return err
				}
				if store.Id != "" {
					pm.DataStore[store.Id] = &store
				}
			} else {
				// 跳过其他标签
				if err := skipElement(d, se.Name); err != nil {
					return err
				}
			}
		case xml.EndElement:
			if se.Name.Local == "DataStore" {
				return nil
			}
		}
	}
}

func (pm *PageModel) parseOperations(d *xml.Decoder) error {
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "Operation" {
				var op Operation
				if err := d.DecodeElement(&op, &se); err != nil {
					return err
				}
				if op.Id != "" {
					pm.Operations[op.Id] = &op
				}
			} else {
				// 跳过其他标签
				if err := skipElement(d, se.Name); err != nil {
					return err
				}
			}
		case xml.EndElement:
			if se.Name.Local == "Operations" {
				return nil
			}
		}
	}
}

// skipElement 略过整个标签（包含子标签）
func skipElement(decoder *xml.Decoder, name xml.Name) error {
	depth := 1
	for {
		t, err := decoder.Token()
		if err != nil {
			return err
		}
		switch x := t.(type) {
		case xml.StartElement:
			depth++
		case xml.EndElement:
			if x.Name.Local == name.Local {
				depth--
				if depth == 0 {
					return nil
				}
			}
		}
	}
}
