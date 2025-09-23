package page_model

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	base_utils "github.com/mark0725/go-app-base/utils"
)

type IViewObject interface {
	GetObject() string
	GetProps() map[string]any
	SetProp(name string, value any)
	GetChildren() []IViewObject
	AppendChildren(value IViewObject)
	AppendElement(name string, value any)
	UnmarshalXML(d *xml.Decoder, start xml.StartElement) error
}

var typeRegistry = map[string]func() IViewObject{
	"PageView": func() IViewObject { return &PageView{ViewObject: *NewViewObject("PageView")} },
	"Form":     func() IViewObject { return &Form{ViewObject: *NewViewObject("Form")} },
}

func newViewObjectByTag(tag string) IViewObject {
	if fn, ok := typeRegistry[tag]; ok {
		return fn()
	}
	return NewViewObject(tag)
}

// // UnmarshalViewObject 通用的解码方法，嵌入了 ViewObject 的类型可在其 UnmarshalXML 中调用
func UnmarshalViewObject(decoder *xml.Decoder, start xml.StartElement, target IViewObject) error {
	if err := setFieldsFromAttrs(start.Attr, target); err != nil {
		return err
	}

	// 3. 依次读出子节点
	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch se := t.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "children":
				if err := parseChildren(decoder, target); err != nil {
					return err
				}
			case "__list":
				if err := parseList(decoder, se, target); err != nil {
					return err
				}
			default:
				if handled, err := parseToField(decoder, se, target); err != nil {
					return err
				} else if !handled {
					// 如果无法匹配到字段，就作为一个子对象存到 Props 中
					childInst, err := parseToProps(decoder, se, target)
					if err != nil {
						return err
					}
					// childInst 就是该子标签解析出来的对象
					_ = childInst
				}
			}

		case xml.EndElement:
			// 遇到结束标签且与当前标签匹配，则结束
			if se.Name.Local == start.Name.Local {
				return nil
			}
		}
	}
	return nil
}

// parseChildren 专门解析 <children> 标签内部的所有子标签，直到读到 </children>
func parseChildren(decoder *xml.Decoder, parent IViewObject) error {

	for {
		t, err := decoder.Token()
		if err != nil {
			return err
		}
		switch se := t.(type) {
		case xml.StartElement:
			// 每个子标签都当做一个新的 IViewObject
			child := newViewObjectByTag(se.Name.Local)
			if err := child.UnmarshalXML(decoder, se); err != nil {
				return err
			}
			parent.AppendChildren(child)
		case xml.EndElement:
			if se.Name.Local == "children" {
				return nil
			}
		}
	}
}

func parseList(decoder *xml.Decoder, start xml.StartElement, parent IViewObject) error {
	propKey := lowerFirstRune(start.Name.Local)
	for _, attr := range start.Attr {
		if attr.Name.Local == "name" {
			propKey = attr.Value
		}
	}

	for {
		t, err := decoder.Token()
		if err != nil {
			return err
		}
		switch se := t.(type) {
		case xml.StartElement:
			// 每个子标签都当做一个新的 IViewObject
			child := newViewObjectByTag(se.Name.Local)
			if err := child.UnmarshalXML(decoder, se); err != nil {
				return err
			}
			parent.AppendElement(propKey, child)
		case xml.EndElement:
			if se.Name.Local == "__list" {
				return nil
			}
		}
	}
}

// parseToField 尝试将一个子标签解析到 target 已定义的 struct 字段
func parseToField(decoder *xml.Decoder, start xml.StartElement, target IViewObject) (bool, error) {
	// 只针对嵌入 struct 进行字段匹配
	v := reflect.ValueOf(target).Elem()
	tp := v.Type()
	tagName := start.Name.Local

	for i := 0; i < tp.NumField(); i++ {
		field := tp.Field(i)
		xmlTag := field.Tag.Get("xml")
		xmlTagName, _, _ := strings.Cut(xmlTag, ",")

		// 命中同名字段
		if xmlTagName == tagName {
			fv := v.Field(i)
			// 如果是切片类型字段
			if fv.Kind() == reflect.Slice {
				sliceElemType := fv.Type().Elem()
				var newElemValue reflect.Value

				// 建立一个新的元素实例
				if sliceElemType.Kind() == reflect.Ptr {
					newElemValue = reflect.New(sliceElemType.Elem())
				} else {
					newElemValue = reflect.New(sliceElemType)
				}
				newElem := newElemValue.Interface()

				// 如果实现了 IViewObject
				if vo, ok := newElem.(IViewObject); ok {
					if err := vo.UnmarshalXML(decoder, start); err != nil {
						return true, err
					}
				} else {
					if err := decoder.DecodeElement(newElem, &start); err != nil {
						return true, err
					}
				}

				// 加入切片
				if sliceElemType.Kind() == reflect.Ptr {
					fv.Set(reflect.Append(fv, newElemValue))
				} else {
					fv.Set(reflect.Append(fv, newElemValue.Elem()))
				}

			} else {
				// 普通字段
				if fv.CanAddr() {
					// 如果字段本身是 IViewObject
					if vo, ok := fv.Addr().Interface().(IViewObject); ok {
						if err := vo.UnmarshalXML(decoder, start); err != nil {
							return true, err
						}
					} else {
						tmpValue := reflect.New(fv.Type())
						if err := decoder.DecodeElement(tmpValue.Interface(), &start); err != nil {
							return true, err
						}
						fv.Set(tmpValue.Elem())
					}
				}
			}
			return true, nil
		}
	}
	return false, nil
}

// setFieldsFromAttrs 将属性分配给 target 已定义字段或放入 props
func setFieldsFromAttrs(attrs []xml.Attr, target IViewObject) error {
	tv := reflect.ValueOf(target).Elem()
	tt := tv.Type()
	for _, attr := range attrs {
		match := false
		for i := 0; i < tt.NumField(); i++ {
			sf := tt.Field(i)
			xmlTag := sf.Tag.Get("xml")
			tagName, tagOpt, _ := strings.Cut(xmlTag, ",")
			if (tagOpt == "attr" || strings.Contains(xmlTag, "attr") || tagOpt == "") &&
				(tagName == attr.Name.Local || strings.EqualFold(sf.Name, attr.Name.Local)) {
				fv := tv.Field(i)
				if fv.CanSet() && fv.Kind() == reflect.String {
					fv.SetString(attr.Value)
					match = true
					break
				}
			}
		}

		newValue := parseValue(attr.Value)
		if !match {
			switch attr.Name.Space {
			case "object":
				o := map[string]any{}
				if err := json.Unmarshal([]byte(attr.Value), &o); err != nil {
					target.SetProp(attr.Name.Local, newValue)
				} else {

					target.SetProp(attr.Name.Local, o)
				}
			case "list":
				o := []any{}
				if err := json.Unmarshal([]byte(attr.Value), &o); err != nil {
					target.SetProp(attr.Name.Local, newValue)
				} else {
					target.SetProp(attr.Name.Local, o)
				}
			default:
				target.SetProp(attr.Name.Local, newValue)
			}

		}
	}
	return nil
}

// parseToProps 如果没有找到对应字段，就解析为新的 IViewObject 并存入 target.GetProps()
func parseToProps(decoder *xml.Decoder, start xml.StartElement, target IViewObject) (IViewObject, error) {
	child := newViewObjectByTag(start.Name.Local)
	if err := child.UnmarshalXML(decoder, start); err != nil {
		return nil, err
	}

	// 按照需求，如果子标签有 name 属性，则以其 name 属性为键，否则用标签首字母小写
	nameVal := child.GetProps()["name"]
	var propKey string
	if nameVal != nil {
		propKey = nameVal.(string)
	} else {
		propKey = lowerFirstRune(start.Name.Local)
	}

	target.SetProp(propKey, child)

	return child, nil
}

// setFieldsFromAttrs 遍历所有属性，先看能否匹配到结构体已有字段（以 xml tag 或字段名匹配），否则放 props

// lowerFirstRune 将字符串首字母小写
func lowerFirstRune(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func ParsePageModel(path string, params map[string]string) (*PageModel, error) {
	xmlContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var pm PageModel
	if err := xml.Unmarshal(xmlContent, &pm); err != nil {
		return nil, err
	}

	return &pm, nil
}

func ParsePageModel2Map(path string, params map[string]string) (map[string]any, error) {

	pm, err := ParsePageModel(path, params)
	if err != nil {
		return nil, err
	}

	mapObj := base_utils.StructToMap(pm)
	if pageView, ok := mapObj["pageView"]; ok {
		FlattenProps(pageView)
	}

	return mapObj, nil
}

func FlattenProps(obj any) any {
	switch data := obj.(type) {
	case map[string]any:
		// 如果存在 "props" 并且它是一个 map，则将其中所有键值提升到当前层级
		if propsValue, ok := data["props"]; ok {
			if propsMap, isMap := propsValue.(map[string]any); isMap {
				for k, v := range propsMap {
					data[k] = v
				}
			}
			delete(data, "props")
		}

		// 递归处理 map 的所有字段，并将返回值重新赋回
		for k, v := range data {
			data[k] = FlattenProps(v)
		}
		return data

	case []any:
		// 递归处理切片中的所有元素，并将返回值重新赋回
		for i, v := range data {
			data[i] = FlattenProps(v)
		}
		return data

	default:
		// 如果既不是 map 也不是切片，则直接返回
		return obj
	}
}

func parseValue(val string) any {
	if i, err := strconv.ParseInt(val, 10, 64); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(val, 64); err == nil {
		return f
	}
	if b, err := strconv.ParseBool(strings.ToLower(val)); err == nil {
		return b
	}
	return val
}
