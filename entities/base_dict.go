package entities

const DB_TABLE_BASE_DICT = "base_dict"

type BaseDict struct {
	OrgId       string `gorm:"column:org_id;type:VARCHAR(64);primaryKey" json:"org_id,omitempty" yaml:"org_id,omitempty" field-id:"org_id" field-type:"varchar" field-comment:"机构ID"`                  // 机构ID
	DictId      string `gorm:"column:dict_id;type:VARCHAR(64);primaryKey" json:"dict_id,omitempty" yaml:"dict_id,omitempty" field-id:"dict_id" field-type:"varchar" field-comment:"码表"`                // 码表
	DictName    string `gorm:"column:dict_name;type:VARCHAR(180)" json:"dict_name,omitempty" yaml:"dict_name,omitempty" field-id:"dict_name" field-type:"varchar" field-comment:"码表名"`                 // 码表名
	DictType    string `gorm:"column:dict_type;type:VARCHAR(64)" json:"dict_type,omitempty" yaml:"dict_type,omitempty" field-id:"dict_type" field-type:"varchar" field-comment:"码表类型"`                 // 码表类型
	ModuleId    string `gorm:"column:module_id;type:VARCHAR(64)" json:"module_id,omitempty" yaml:"module_id,omitempty" field-id:"module_id" field-type:"varchar" field-comment:"模块"`                   // 模块
	DictDesc    string `gorm:"column:dict_desc;type:VARCHAR(300)" json:"dict_desc,omitempty" yaml:"dict_desc,omitempty" field-id:"dict_desc" field-type:"varchar" field-comment:"描述"`                  // 描述
	DataCrtDate string `gorm:"column:data_crt_date;type:VARCHAR(10)" json:"data_crt_date,omitempty" yaml:"data_crt_date,omitempty" field-id:"data_crt_date" field-type:"varchar" field-comment:"创建日期"` // 创建日期
	DataCrtTime string `gorm:"column:data_crt_time;type:VARCHAR(20)" json:"data_crt_time,omitempty" yaml:"data_crt_time,omitempty" field-id:"data_crt_time" field-type:"varchar" field-comment:"创建时间"` // 创建时间
	DataUpdDate string `gorm:"column:data_upd_date;type:VARCHAR(10)" json:"data_upd_date,omitempty" yaml:"data_upd_date,omitempty" field-id:"data_upd_date" field-type:"varchar" field-comment:"更新日期"` // 更新日期
	DataUpdTime string `gorm:"column:data_upd_time;type:VARCHAR(20)" json:"data_upd_time,omitempty" yaml:"data_upd_time,omitempty" field-id:"data_upd_time" field-type:"varchar" field-comment:"更新时间"` // 更新时间
}

func (BaseDict) TableName() string { return "base_dict" }
