package model

func (tag *TagModel) TableName() string {
	return "tags"
}

func (data *CourseTagModel) TableName() string  {
	return "course_tags"
}

// Get tag name by id.
func GetTagNameById(id int) (string, error) {
	var tag TagModel
	d := DB.Self.Where("id = ?", id).First(&tag)
	return tag.Name, d.Error
}

// Get tag total number.
func GetTagSum() int {
	var count int
	var tag TagModel
	DB.Self.Find(&tag).Count(&count)
	return count
}

// Get all tags including ids and names.
func GetTags() *[]TagModel {
	var tags []TagModel
	DB.Self.Find(&tags)
	return &tags
}
