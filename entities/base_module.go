package entities

const DB_TABLE_BASE_MODULE = "base_module"

type BaseModule struct {
	OrgId       string `gorm:"column:org_id;type:VARCHAR(100);primaryKey" json:"org_id,omitempty" yaml:"org_id,omitempty" field-id:"org_id" field-type:"varchar" field-comment:"机构ID"`             // 机构ID
	ModuleId    string `gorm:"column:module_id;type:VARCHAR(100);primaryKey" json:"module_id,omitempty" yaml:"module_id,omitempty" field-id:"module_id" field-type:"varchar" field-comment:"模块ID"` // 模块ID
	ModuleName  string `gorm:"column:module_name;type:VARCHAR(180)" json:"module_name,omitempty" yaml:"module_name,omitempty" field-id:"module_name" field-type:"varchar" field-comment:"模块名称"`    // 模块名称
	ModuleType  string `gorm:"column:module_type;type:VARCHAR(100)" json:"module_type,omitempty" yaml:"module_type,omitempty" field-id:"module_type" field-type:"varchar" field-comment:"模块类型"`    // 模块类型
	ModuleDesc  string `gorm:"column:module_desc;type:VARCHAR(300)" json:"module_desc,omitempty" yaml:"module_desc,omitempty" field-id:"module_desc" field-type:"varchar" field-comment:"模块说明"`    // 模块说明
	TitleShort  string `gorm:"column:title_short;type:VARCHAR(60)" json:"title_short,omitempty" yaml:"title_short,omitempty" field-id:"title_short" field-type:"varchar"`
	Url         string `gorm:"column:url;type:VARCHAR(300)" json:"url,omitempty" yaml:"url,omitempty" field-id:"url" field-type:"varchar" field-comment:"模块位置"`                                        // 模块位置
	OrderNo     int64  `gorm:"column:order_no;type:INTEGER" json:"order_no,omitempty" yaml:"order_no,omitempty" field-id:"order_no" field-type:"integer" field-comment:"顺序"`                           // 顺序
	Status      string `gorm:"column:status;type:VARCHAR(10)" json:"status,omitempty" yaml:"status,omitempty" field-id:"status" field-type:"varchar" field-comment:"状态"`                               // 状态
	DataCrtDate string `gorm:"column:data_crt_date;type:VARCHAR(10)" json:"data_crt_date,omitempty" yaml:"data_crt_date,omitempty" field-id:"data_crt_date" field-type:"varchar" field-comment:"创建日期"` // 创建日期
	DataCrtTime string `gorm:"column:data_crt_time;type:VARCHAR(20)" json:"data_crt_time,omitempty" yaml:"data_crt_time,omitempty" field-id:"data_crt_time" field-type:"varchar" field-comment:"创建时间"` // 创建时间
	DataUpdDate string `gorm:"column:data_upd_date;type:VARCHAR(10)" json:"data_upd_date,omitempty" yaml:"data_upd_date,omitempty" field-id:"data_upd_date" field-type:"varchar" field-comment:"更新日期"` // 更新日期
	DataUpdTime string `gorm:"column:data_upd_time;type:VARCHAR(20)" json:"data_upd_time,omitempty" yaml:"data_upd_time,omitempty" field-id:"data_upd_time" field-type:"varchar" field-comment:"更新时间"` // 更新时间
}

func (BaseModule) TableName() string { return "base_module" }
