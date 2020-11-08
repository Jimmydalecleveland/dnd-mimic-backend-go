package datasources

type CharSkillProficiency struct {
	CharID  int32 `gorm:"column:charID"`
	SkillID int32 `gorm:"column:skillID"`
}

func (CharSkillProficiency) TableName() string {
	return "CharSkillProficiency"
}
