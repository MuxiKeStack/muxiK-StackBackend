package model

import (
	// _"database/sql"
	_ "github.com/jinzhu/gorm"
	// _ "log"
)

const (
	thSQL = " AND LOCATE('5', `course_id`, 3) = 1 "
)

func (UsingCourseModel) TableName() string {
	return "using_course"
}

func (HistoryCourseModel) TableName() string {
	return "history_course"
}

/*
func (CourseLikeModel) TableName() string {
	return "course_like"
}*/

// Add a new course.
func (class *UsingCourseModel) Add() error {
	d := DB.Self.Create(class)
	return d.Error
}

// Delete a course.
// Fixed by shiina orez at 2019.11.24 evaluation =>> course
func (class *UsingCourseModel) Delete() error {
	d := DB.Self.Delete(class)
	return d.Error
}

// Get course by its id.(course list)
func (class *UsingCourseModel) GetById() error {
	d := DB.Self.First(class, "id = ?", class.Id)
	return d.Error
}

// Get course by its hash.
func (class *UsingCourseModel) GetByHash() error {
	d := DB.Self.First(class, "hash = ?", class.Hash)
	return d.Error
}

// Get history course by its hash.
func (class *HistoryCourseModel) GetHistoryByHash() error {
	d := DB.Self.First(class, "hash = ?", class.Hash)
	return d.Error
}

func (class *UsingCourseModel) GetClass(courseId string, classId uint64) error {
	d := DB.Self.Where("course_id = ? AND class_id = ? ", courseId, classId).First(&class)
	if d.RecordNotFound() {
		return nil
	}
	return d.Error
}

// Get course by its type.(course list)
func (class *UsingCourseModel) GetByType() error {
	d := DB.Self.Find(class, "type = ?", class.Type)
	return d.Error
}

// Get course by its teacher.(course list)
func (class *UsingCourseModel) GetByTeacher() error {
	d := DB.Self.Find(class, "teacher = ?", class.Teacher)
	return d.Error
}

/*
func GetCourseList(userId uint32) ([]string, error) {
	var data []CourseLikeModel
	var result []string
	d := DB.Self.Where("user_id = ?", userId).Find(&data)
	if d.RecordNotFound() {
		return nil, nil
	}
	for _, i := range data {
		result = append(result, i.CourseHash)
	}
	return result, d.Error
}*/

// Get course by its courseid.(course assistant)
func (class *UsingCourseModel) GetByCourseId(time int, place int) error { //int为映射，作为筛选条件
	d := DB.Self.Find(class, "course_id = ?", class.CourseId)
	return d.Error
}

// Judge whether a course has already favorited by the current user.
func (class *UsingCourseModel) HasFavorited(userId uint32) bool {
	var data CourseListModel
	var count int
	DB.Self.Where("user_id = ? AND course_hash_id = ? ", userId, class.Hash).First(&data).Count(&count)
	return count > 0
}

// Favorite a course by the current user.
func (class *UsingCourseModel) Favorite(userId uint32) error {
	var data = CourseListModel{
		CourseHashId: class.Hash,
		UserId:       userId,
	}

	d := DB.Self.Create(&data)
	return d.Error
}

func (class *UsingCourseModel) Unfavorite(userId uint32) error {
	var data = CourseListModel{
		CourseHashId: class.Hash,
		UserId:       userId,
	}

	d := DB.Self.Delete(&data)
	return d.Error
}

// Search course by name, courseId or teacher
// Use fulltext search, against and match
func AgainstAndMatchCourses(kw string, page, limit uint64, th bool) ([]UsingCourseModel, error) {
	courses := &[]UsingCourseModel{}
	// log.Println("Query:", kw, page, limit, th)
	if !th {
		DB.Self.Debug().Table("using_course").Where("MATCH (`name`, `course_id`, `teacher`) AGAINST ('" + kw + "') ").Find(courses).Limit(limit).Offset((page - 1) * limit)
	} else {
		DB.Self.Debug().Table("using_course").Where("MATCH (`name`, `course_id`, `teacher`) AGAINST ('" + kw + "')" + thSQL).Find(courses).Limit(limit).Offset((page - 1) * limit)
	}
	return *courses, nil
}

// Search history course by name or teacher
// Use fulltext search, against and match
func AgainstAndMatchHistoryCourses(kw string, page, limit uint64) ([]HistoryCourseModel, error) {
	courses := &[]HistoryCourseModel{}
	DB.Self.Table("history_course").Where("MATCH (`name`, `teacher`) AGAINST ('" + kw + "') ").Find(courses).Limit(limit).Offset((page - 1) * limit)
	return *courses, nil
}

// Get all courses
func AllCourses(page, limit uint64, th bool) ([]UsingCourseModel, error) {
	courses := &[]UsingCourseModel{}
	if th {
		DB.Self.Table("using_course").Where("LOCATE ('5', `course_id`, 3) = 1").Find(courses).Limit(limit).Offset((page - 1) * limit)
	} else {
		DB.Self.Table("using_course").Find(&courses).Limit(limit).Offset((page - 1) * limit)
	}
	return *courses, nil
}

// Get all history courses
func AllHistoryCourses(page, limit uint64) ([]HistoryCourseModel, error) {
	courses := &[]HistoryCourseModel{}
	DB.Self.Table("history_course").Find(courses).Limit(limit).Offset((page - 1) * limit)
	return *courses, nil
}
