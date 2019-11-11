package model

type TagModel struct {
	Id      uint32 `gorm:"column:id"`
	TagName string `gorm:"column:tag_name"`
}
