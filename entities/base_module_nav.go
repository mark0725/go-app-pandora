package entities

const DB_TABLE_BASE_MODULE_NAV = "base_module_nav"

type BaseModuleNav struct {
	OrgId       string `gorm:"column:org_id;type:VARCHAR(100);primaryKey" json:"org_id,omitempty" yaml:"org_id,omitempty" field-id:"org_id" field-type:"varchar" field-comment:"机构ID"`             // 机构ID
	ModuleId    string `gorm:"column:module_id;type:VARCHAR(100);primaryKey" json:"module_id,omitempty" yaml:"module_id,omitempty" field-id:"module_id" field-type:"varchar" field-comment:"模块ID"` // 模块ID
	NavId       string `gorm:"column:nav_id;type:VARCHAR(100);primaryKey" json:"nav_id,omitempty" yaml:"nav_id,omitempty" field-id:"nav_id" field-type:"varchar" field-comment:"导航ID"`             // 导航ID
	NavName     string `gorm:"column:nav_name;type:VARCHAR(180)" json:"nav_name,omitempty" yaml:"nav_name,omitempty" field-id:"nav_name" field-type:"varchar" field-comment:"导航名称"`                // 导航名称
	NavType     string `gorm:"column:nav_type;type:VARCHAR(60)" json:"nav_type,omitempty" yaml:"nav_type,omitempty" field-id:"nav_type" field-type:"varchar" field-comment:"导航类型"`                 // 导航类型
	NavIcon     string `gorm:"column:nav_icon;type:VARCHAR(300)" json:"nav_icon,omitempty" yaml:"nav_icon,omitempty" field-id:"nav_icon" field-type:"varchar" field-comment:"导航图标"`                // 导航图标
	NavDesc     string `gorm:"column:nav_desc;type:VARCHAR(300)" json:"nav_desc,omitempty" yaml:"nav_desc,omitempty" field-id:"nav_desc" field-type:"varchar" field-comment:"导航说明"`                // 导航说明
	TitleShort  string `gorm:"column:title_short;type:VARCHAR(60)" json:"title_short,omitempty" yaml:"title_short,omitempty" field-id:"title_short" field-type:"varchar"`
	ViewType    string `gorm:"column:view_type;type:VARCHAR(60)" json:"view_type,omitempty" yaml:"view_type,omitempty" field-id:"view_type" field-type:"varchar" field-comment:"视图类型"` // 视图类型
	ViewSize    string `gorm:"column:view_size;type:VARCHAR(30)" json:"view_size,omitempty" yaml:"view_size,omitempty" field-id:"view_size" field-type:"varchar" field-comment:"视图大小"` // 视图大小
	Url         string `gorm:"column:url;type:VARCHAR(300)" json:"url,omitempty" yaml:"url,omitempty" field-id:"url" field-type:"varchar" field-comment:"导航位置"`                        // 导航位置
	MenuType    string `gorm:"column:menu_type;type:VARCHAR(30)" json:"menu_type,omitempty" yaml:"menu_type,omitempty" field-id:"menu_type" field-type:"varchar" field-comment:"菜单类型"` // 菜单类型
	MenuApi     string `gorm:"column:menu_api;type:VARCHAR(300)" json:"menu_api,omitempty" yaml:"menu_api,omitempty" field-id:"menu_api" field-type:"varchar" field-comment:"菜单数据API"` // 菜单数据API
	IconSize    string `gorm:"column:icon_size;type:VARCHAR(30)" json:"icon_size,omitempty" yaml:"icon_size,omitempty" field-id:"icon_size" field-type:"varchar"`
	OrderNo     int64  `gorm:"column:order_no;type:INTEGER" json:"order_no,omitempty" yaml:"order_no,omitempty" field-id:"order_no" field-type:"integer"`
	DataCrtDate string `gorm:"column:data_crt_date;type:VARCHAR(10)" json:"data_crt_date,omitempty" yaml:"data_crt_date,omitempty" field-id:"data_crt_date" field-type:"varchar"`
	DataCrtTime string `gorm:"column:data_crt_time;type:VARCHAR(20)" json:"data_crt_time,omitempty" yaml:"data_crt_time,omitempty" field-id:"data_crt_time" field-type:"varchar"`
	DataUpdDate string `gorm:"column:data_upd_date;type:VARCHAR(10)" json:"data_upd_date,omitempty" yaml:"data_upd_date,omitempty" field-id:"data_upd_date" field-type:"varchar"`
	DataUpdTime string `gorm:"column:data_upd_time;type:VARCHAR(20)" json:"data_upd_time,omitempty" yaml:"data_upd_time,omitempty" field-id:"data_upd_time" field-type:"varchar"`
	Status      string `gorm:"column:status;type:VARCHAR(30)" json:"status,omitempty" yaml:"status,omitempty" field-id:"status" field-type:"varchar"`
	AppModule   string `gorm:"column:app_module;type:VARCHAR(100)" json:"app_module,omitempty" yaml:"app_module,omitempty" field-id:"app_module" field-type:"varchar"`
	NavPosition string `gorm:"column:nav_position;type:VARCHAR(100)" json:"nav_position,omitempty" yaml:"nav_position,omitempty" field-id:"nav_position" field-type:"varchar"`
	ParamName   string `gorm:"column:param_name;type:VARCHAR(100)" json:"param_name,omitempty" yaml:"param_name,omitempty" field-id:"param_name" field-type:"varchar"`
	PNavId      string `gorm:"column:p_nav_id;type:VARCHAR(100)" json:"p_nav_id,omitempty" yaml:"p_nav_id,omitempty" field-id:"p_nav_id" field-type:"varchar"`
}

func (BaseModuleNav) TableName() string { return "base_module_nav" }
