package model

type TagModel struct {
	Id      uint32 `gorm:"column:id"`
	Name string `gorm:"column:name"`
}
