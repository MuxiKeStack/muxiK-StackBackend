package model

import (
<<<<<<< HEAD
	// _"database/sql"
=======
>>>>>>> 9c76ad1583a96f445e34ac7d1144368131ac906e
	"fmt"
	_ "github.com/jinzhu/gorm"
	// _ "log"
)

const (
	typeTemp       = "AND RIGHT(LEFT(`course_id`, 4), 1) = %s "
	typeCourseTemp = "AND `type` = %s "
	academyTemp    = "AND `academy` = %s "
	weekdayTemp    = "AND (RIGHT(`time1`, 1) = %s OR RIGHT(`time2`, 1) = %s OR RIGHT(`time3`, 1) = %s) "
	nPlaceTemp     = "AND LEFT(`place1`, 1) = 'N' AND (`place2` = '' OR LEFT(`place2`, 1) = 'N') AND (`place3` = '' OR LEFT(`place3`, 1) = 'N') "
	bPlaceTemp     = "AND LEFT(`place1`, 1) != 'N' AND (`place2` = '' OR LEFT(`place2`, 1) != 'N') AND (`place3` = '' OR LEFT(`place3`, 1) != 'N') "
)

func (UsingCourseModel) TableName() string {
	return "using_course"
}

func (HistoryCourseModel) TableName() string {
	return "history_course"
}

func (class *UsingCourseModel) T() {
	fmt.Println("---")
}

/*
func (CourseLikeModel) TableName() string {
	return "course_like"
}*/

// Add a new course.
func (class *UsingCourseModel) Add() error {
	d := DB.Self.Create(class)
	return d.Error
	//return nil
}

// Delete a course.
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
	fmt.Println(data)

	d := DB.Self.Delete(&data)
	return d.Error
}

// Search course by name, courseId or teacher
// Use fulltext search, against and match
// 2020-01-15: Add New Filter: type, academy, weekday, place
func AgainstAndMatchCourses(kw string, page, limit uint64, t, a, w, p string) ([]UsingCourseModel, error) {
	courses := &[]UsingCourseModel{}
	where := "MATCH (`name`, `course_id`, `teacher`) AGAINST ('" + kw + "') "
	if t != "" {
		where += fmt.Sprintf(typeTemp, t)
	}
	if a != "" {
		where += fmt.Sprintf(academyTemp, a)
	}
	if w != "" {
		where += fmt.Sprintf(weekdayTemp, w, w, w)
	}
	if p == "本校区" {
		where += bPlaceTemp
	}
	if p == "南湖校区" {
		where += nPlaceTemp
	}
	DB.Self.Debug().Table("using_course").Where(where).Limit(limit).Offset((page - 1) * limit).Find(courses)
	return *courses, nil
}

// Search history course by name or teacher
// Use fulltext search, against and match
// 2020-01-15: Add New Filter: type
func AgainstAndMatchHistoryCourses(kw string, page, limit uint64, t string) ([]HistoryCourseModel, error) {
	courses := &[]HistoryCourseModel{}
	where := "MATCH (`name`, `teacher`) AGAINST ('" + kw + "') "
	if t != "" {
		where += fmt.Sprintf(typeCourseTemp, t)
	}
	DB.Self.Table("history_course").Where(where).Limit(limit).Offset((page - 1) * limit).Find(courses)
	return *courses, nil
}

// Get all courses
func AllCourses(page, limit uint64, t, a, w, p string) ([]UsingCourseModel, error) {
	courses := &[]UsingCourseModel{}
	where := ""
	if t != "" {
		where += fmt.Sprintf(typeTemp, t)
	}
	if a != "" {
		where += fmt.Sprintf(academyTemp, a)
	}
	if w != "" {
		where += fmt.Sprintf(weekdayTemp, w, w, w)
	}
	if p == "本校区" {
		where += bPlaceTemp
	}
	if p == "南湖校区" {
		where += nPlaceTemp
	}
	if where == "" {
		DB.Self.Table("using_course").Limit(limit).Offset((page - 1) * limit).Find(&courses)
	} else {
		DB.Self.Table("using_course").Where(where).Limit(limit).Offset((page - 1) * limit).Find(&courses)
	}
	return *courses, nil
}

// Get all history courses
func AllHistoryCourses(page, limit uint64, t string) ([]HistoryCourseModel, error) {
	courses := &[]HistoryCourseModel{}
	if t != "" {
		DB.Self.Table("history_course").Where("type = ?", t).Limit(limit).Offset((page - 1) * limit).Find(courses)
	} else {
		DB.Self.Table("history_course").Limit(limit).Offset((page - 1) * limit).Find(courses)
	}
	return *courses, nil
}
