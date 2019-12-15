package model

func (*CourseListModel) TableName() string {
	return "course_list"
}

func GetCollectionsByUserId(userId uint32) ([]string, error) {
	var data []CourseListModel
	var result []string
	d := DB.Self.Where("user_id = ?", userId).Find(&data)
	if d.RecordNotFound() {
		return nil, nil
	}
	for _, i := range data {
		result = append(result, i.CourseHashId)
	}
	return result, d.Error
}

func GetClassesByCourseHash(id string) (*[]UsingCourseModel, error) {
	var classes []UsingCourseModel
	d := DB.Self.Where("hash = ?", id).Find(&classes)
	if d.RecordNotFound() {
		return nil, nil
	}
	return &classes, d.Error
}
