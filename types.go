package pandora

type ApiReponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type Dict struct {
	Id      string               `json:"id"`
	Name    string               `json:"name"`
	Type    string               `json:"type,omitempty"`
	Struct  string               `json:"struct,omitempty"` //options,tree,table
	Style   string               `json:"style,omitempty"`
	Items   map[string]*DictItem `json:"items,omitempty"`
	Options []*DictItem          `json:"options,omitempty"`
}

type DictItem struct {
	Value  string            `json:"value,omitempty"`
	Label  string            `json:"label,omitempty"`
	Icon   string            `json:"icon,omitempty"`
	Style  string            `json:"style,omitempty"`
	Color  string            `json:"color,omitempty"`
	Fields map[string]string `json:"fields,omitempty"`
}

type AppConfig struct {
	App  AppInfo    `json:"app"`
	Auth *UserAuth  `json:"auth,omitempty"`
	User *UserInfo  `json:"user,omitempty"`
	Menu MenuConfig `json:"menu"`
}

type AppInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Logo    string `json:"logo"`
	Title   string `json:"title"`
}

type UserAuth struct {
	Authed     bool   `json:"authed"`
	AuthType   string `json:"auth_type"`
	AuthUrl    string `json:"auth_url"`
	SigninUrl  string `json:"signin_url"`
	SignoutUrl string `json:"signout_url"`
}

type UserInfo struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Dept   string `json:"dept"`
	Mail   string `json:"mail"`
	Avatar string `json:"avatar"`
}

type MenuConfig struct {
	Main    []*MenuItem `json:"main"`
	Nav2    []*MenuItem `json:"nav2,omitempty"`
	NavUser []*MenuItem `json:"navuser,omitempty"`
}

type MenuItem struct {
	Type       string      `json:"type"`
	Id         string      `json:"id"`
	Title      string      `json:"title"`
	TitleShort string      `json:"title_short,omitempty"`
	Url        string      `json:"url,omitempty"`
	Ico        string      `json:"ico,omitempty"`
	View       string      `json:"view,omitempty"`
	ViewSize   string      `json:"view_size,omitempty"`
	Count      int         `json:"count,omitempty"`
	Children   []*MenuItem `json:"children,omitempty"`
}

type MenuSelect struct {
	Param string      `json:"param"`
	Value string      `json:"value"`
	Items []*DictItem `json:"items"`
}

type PageConfig struct {
	Title      string      `json:"title"`
	Type       string      `json:"type"`
	TitleShort string      `json:"title_short,omitempty"`
	Select     *MenuSelect `json:"select,omitempty"`
	Menu       []*MenuItem `json:"menu"`
}

type FilterParamItem struct {
	Key       string `json:"key"`
	Type      string `json:"type"`
	TypeValue string `json:"type_value"`
	Value     any    `json:"value"`
}

type FilterParam struct {
	Opr   string            `json:"opr"`
	Items []FilterParamItem `json:"items"`
}
