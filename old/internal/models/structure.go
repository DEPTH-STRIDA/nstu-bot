package models

import "gorm.io/gorm"

// StructureOfGroup представляет таблицу structure_of_group
type StructureOfGroup struct {
	gorm.Model
	GroupID   uint    `gorm:"column:group_id"`
	Group     Group   `gorm:"foreignKey:GroupID"`
	Structure []int64 `gorm:"column:structure;type:bigint[]"`
}

// TableName определяет имя таблицы для StructureOfGroup
func (StructureOfGroup) TableName() string {
	return "structure_of_group"
}
