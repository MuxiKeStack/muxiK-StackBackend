package model

import "github.com/jinzhu/gorm"

func (*CourseListModel) TableName() string {
	return "course_list"
}

// Get all courses' hash ids from collection by userId, return array of course's hash id and error.
func GetCourseHashIdsFromCollection(userId uint32) ([]string, error) {
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

// Get collections' records by userId.
func GetCollectionsByUserId(userId uint32, lastId, limit int32) (*[]CourseListModel, error) {
	var data []CourseListModel
	var d *gorm.DB
	if lastId != 0 {
		d = DB.Self.Where("id < ? AND user_id = ?", lastId, userId).Order("id DESC").Limit(limit).Find(&data)
	} else {
		d = DB.Self.Where("user_id = ?", userId).Order("id DESC").Limit(limit).Find(&data)
	}

	if d.RecordNotFound() {
		return nil, nil
	}
	return &data, d.Error
}

// Get classes by course's hash id.
func GetClassesByCourseHash(id string) (*[]UsingCourseModel, error) {
	var classes []UsingCourseModel
	d := DB.Self.Where("hash = ?", id).Find(&classes)
	if d.RecordNotFound() {
		return nil, nil
	}
	return &classes, d.Error
}

func GetTheMostAttendanceCheckType(courseId string) (uint8, error) {
	var data struct{ AttendanceCheckType uint8 }
	sql := "SELECT * FROM (SELECT attendance_check_type FROM course_evaluation WHERE course_id = ? GROUP BY attendance_check_type) AS a LIMIT 1"
	d := DB.Self.Raw(sql, courseId).Scan(&data)
	if d.RecordNotFound() {
		return data.AttendanceCheckType, nil
	}
	return data.AttendanceCheckType, d.Error
}

func GetTheMostExamCheckType(courseId string) (uint8, error) {
	var data struct{ ExamCheckType uint8 }
	sql := "SELECT * FROM (SELECT exam_check_type FROM course_evaluation WHERE course_id = ? GROUP BY exam_check_type) AS a LIMIT 1"
	d := DB.Self.Raw(sql, courseId).Scan(&data)
	if d.RecordNotFound() {
		return data.ExamCheckType, nil
	}
	return data.ExamCheckType, d.Error
}
