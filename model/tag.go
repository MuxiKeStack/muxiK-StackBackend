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

// Get tag amount of one course.
func GetTagsNumber(tagId uint32, courseId string) (uint32, error) {
	var data CourseTagModel
	d := DB.Self.Where("course_id = ? AND tag_id = ?", courseId, tagId).First(&data)
	return data.Num, d.Error
}

// 获取课程的前二的tag的名字
func GetTwoMostTagNamesOfCourseByHashId(courseID string) ([]string, error) {
	var names []struct{ Name string }
	d := DB.Self.Table("course_tag").
		Select("tags.name").
		Joins("JOIN tags ON tags.id = course_tag.tag_id").
		Where("course_tag.course_id = ?", courseID).
		Order("num desc").Limit(2).Scan(&names)
	res := make([]string, len(names))
	for i, name := range names {
		res[i] = name.Name
	}

	if d.RecordNotFound() {
		return res, nil
	}
	return res, d.Error
}

// Get two most tags' ids of a course by its hash id.
func GetTwoMostTagIdsOfCourseByHashId(courseId string) ([]int, error) {
	var tags []struct{ Id int }
	d := DB.Self.Table("course_tag").Select("tag_id AS id").Where("course_id = ?", courseId).Order("num desc").Limit(2).Scan(&tags)

	var ids []int
	for _, tag := range tags {
		ids = append(ids, tag.Id)
	}

	if d.RecordNotFound() {
		return ids, nil
	}
	return ids, d.Error
}

// Get four most tags' ids of a course by its hash id.
func GetFourMostTagIdsOfCourseByHashId(courseId string) ([]int, error) {
	var tags []struct{ Id int }
	d := DB.Self.Table("course_tag").Select("tag_id AS id").Where("course_id = ?", courseId).Order("num desc").Limit(4).Scan(&tags)

	var ids []int
	for _, tag := range tags {
		ids = append(ids, tag.Id)
	}

	if d.RecordNotFound() {
		return ids, nil
	}
	return ids, d.Error
}
