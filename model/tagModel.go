package model

type TagModel struct {
	Id   uint32 `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

type CourseTagModel struct {
	Id       uint32 `gorm:"column:id"`
	TagId    uint32 `gorm:"column:tag_id"`
	CourseId string `gorm:"column:course_id"`
}
