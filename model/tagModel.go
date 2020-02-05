package model

type TagModel struct {
	Id   uint32 `gorm:"column:id" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}

type CourseTagModel struct {
	Id       uint32 `gorm:"column:id"`
	TagId    uint32 `gorm:"column:tag_id"`
	CourseId string `gorm:"column:course_id"` //此处的CourseId实际上是hash
	Num      uint32 `gorm:"column:num"`
}
