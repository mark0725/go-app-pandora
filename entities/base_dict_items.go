package entities

const DB_TABLE_BASE_DICT_ITEMS = "base_dict_items"

type BaseDictItems struct {
	OrgId        string `gorm:"column:org_id;type:VARCHAR(64);primaryKey" json:"org_id,omitempty" yaml:"org_id,omitempty" field-id:"org_id" field-type:"varchar" field-comment:"机构ID"`                      // 机构ID
	DictId       string `gorm:"column:dict_id;type:VARCHAR(64);primaryKey" json:"dict_id,omitempty" yaml:"dict_id,omitempty" field-id:"dict_id" field-type:"varchar" field-comment:"码表"`                    // 码表
	DictName     string `gorm:"column:dict_name;type:VARCHAR(180)" json:"dict_name,omitempty" yaml:"dict_name,omitempty" field-id:"dict_name" field-type:"varchar" field-comment:"码表名"`                     // 码表名
	ItemCode     string `gorm:"column:item_code;type:VARCHAR(64);primaryKey" json:"item_code,omitempty" yaml:"item_code,omitempty" field-id:"item_code" field-type:"varchar" field-comment:"代码"`            // 代码
	ItemValue    string `gorm:"column:item_value;type:VARCHAR(180)" json:"item_value,omitempty" yaml:"item_value,omitempty" field-id:"item_value" field-type:"varchar" field-comment:"码值"`                  // 码值
	ItemDesc     string `gorm:"column:item_desc;type:VARCHAR(300)" json:"item_desc,omitempty" yaml:"item_desc,omitempty" field-id:"item_desc" field-type:"varchar" field-comment:"描述"`                      // 描述
	ItemParent   string `gorm:"column:item_parent;type:VARCHAR(64)" json:"item_parent,omitempty" yaml:"item_parent,omitempty" field-id:"item_parent" field-type:"varchar" field-comment:"上级码"`              // 上级码
	DictParent   string `gorm:"column:dict_parent;type:VARCHAR(64)" json:"dict_parent,omitempty" yaml:"dict_parent,omitempty" field-id:"dict_parent" field-type:"varchar" field-comment:"上级码表"`             // 上级码表
	OrdNo        int64  `gorm:"column:ord_no;type:INTEGER" json:"ord_no,omitempty" yaml:"ord_no,omitempty" field-id:"ord_no" field-type:"integer" field-comment:"顺序"`                                       // 顺序
	ModuleId     string `gorm:"column:module_id;type:VARCHAR(64);primaryKey" json:"module_id,omitempty" yaml:"module_id,omitempty" field-id:"module_id" field-type:"varchar" field-comment:"模块"`            // 模块
	ItemColor    string `gorm:"column:item_color;type:VARCHAR(10)" json:"item_color,omitempty" yaml:"item_color,omitempty" field-id:"item_color" field-type:"varchar" field-comment:"颜色"`                   // 颜色
	ItemIcon     string `gorm:"column:item_icon;type:VARCHAR(300)" json:"item_icon,omitempty" yaml:"item_icon,omitempty" field-id:"item_icon" field-type:"varchar" field-comment:"图标"`                      // 图标
	ItemStyle    string `gorm:"column:item_style;type:VARCHAR(10)" json:"item_style,omitempty" yaml:"item_style,omitempty" field-id:"item_style" field-type:"varchar" field-comment:"样式"`                   // 样式
	ItemIconType string `gorm:"column:item_icon_type;type:VARCHAR(10)" json:"item_icon_type,omitempty" yaml:"item_icon_type,omitempty" field-id:"item_icon_type" field-type:"varchar" field-comment:"图标类型"` // 图标类型
	DataCrtDate  string `gorm:"column:data_crt_date;type:VARCHAR(10)" json:"data_crt_date,omitempty" yaml:"data_crt_date,omitempty" field-id:"data_crt_date" field-type:"varchar" field-comment:"创建日期"`     // 创建日期
	DataCrtTime  string `gorm:"column:data_crt_time;type:VARCHAR(20)" json:"data_crt_time,omitempty" yaml:"data_crt_time,omitempty" field-id:"data_crt_time" field-type:"varchar" field-comment:"创建时间"`     // 创建时间
	DataUpdDate  string `gorm:"column:data_upd_date;type:VARCHAR(10)" json:"data_upd_date,omitempty" yaml:"data_upd_date,omitempty" field-id:"data_upd_date" field-type:"varchar" field-comment:"更新日期"`     // 更新日期
	DataUpdTime  string `gorm:"column:data_upd_time;type:VARCHAR(20)" json:"data_upd_time,omitempty" yaml:"data_upd_time,omitempty" field-id:"data_upd_time" field-type:"varchar" field-comment:"更新时间"`     // 更新时间
}

func (BaseDictItems) TableName() string { return "base_dict_items" }
