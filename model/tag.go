package model

func (tag *TagModel) TableName() string {
	return "tags"
}

func (data *CourseTagModel) TableName() string {
	return "course_tag"
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

// New tags for a course when publish a new evaluation.
func NewTagsForCourse(tagId uint32, courseId string) error {
	var data CourseTagModel
	d := DB.Self.Where("course_id = ? AND tag_id = ?", courseId, tagId).First(&data)

	if d.RecordNotFound() {
		data = CourseTagModel{
			TagId:    tagId,
			CourseId: courseId,
			Num:      1,
		}
		d = DB.Self.Create(&data)
		return d.Error
	}

	data.Num++
	d = DB.Self.Save(&data)
	return d.Error
}

func GetTagsNumber(tagId uint32, courseId string) (uint32, error) {
	var data CourseTagModel
	d := DB.Self.Where("course_id = ? AND tag_id = ?", courseId, tagId).First(&data)
	return data.Num, d.Error
}
