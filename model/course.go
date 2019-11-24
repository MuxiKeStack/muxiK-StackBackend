package model

import (
	"database/sql"
	_ "github.com/jinzhu/gorm"
)

// Add a new course.
func (course *UsingCourseModel) Add() error {
	d := DB.Self.Create(course)
	return d.Error
}

// Delete a course.
// Fixed by shiina orez at 2019.11.24 evaluation =>> course
func (course *UsingCourseModel) Delete() error {
	d := DB.Self.Delete(course)
	return d.Error
}

// Get course by its id.(course list)
func (course *UsingCourseModel) GetById() error {
	d := DB.Self.First(course, "id = ?", course.Id)
	return d.Error
}

// Get course by its type.(course list)
// Fixed by shiina orez at 2019.11.24, type =>> Type
func (course *UsingCourseModel) GetByType() error {
	d := DB.Self.Find(course, "type = ?", course.Type)
	return d.Error
}

// Get course by its teacher.(course list)
// Fixed by shiina orez at 2019.11.24 teacher =>> Teacher
func (course *UsingCourseModel) GetByTeacher() error {
	d := DB.Self.Find(course, "teacher = ?", course.Teacher)
	return d.Error
}

// Get course by its name.(TODO)(course list)
// Fixed by shiina orez at 2019.11.24 GetByTeacher =>> GetByName
func (course *UsingCourseModel) GetByName() error {
	d := DB.Self.Find(course, "name = ?", course.Name)
	return d.Error
}

// Get course by its name.(TODO)(course assistant)
// func (course *UsingCourseModel) GetByName(int time, int place) error {   //int为映射，作为筛选条件
//     d := DB.Self.Find(course, "name = ?", course.name)
//     return d.Error
// }

// Get course by its teacher.(course assistant)
// func (course *UsingCourseModel) GetByTeacher(int time, int place) error {   //int为映射，作为筛选条件
//     d := DB.Self.Find(course, "teacher = ?", course.teacher)
//     return d.Error
// }

// Get course by its courseid.(course assistant)
// Fixed by shiina orez at 2019.11.24 `int time` =>> `time int`, `int place` =>> `place int`
func (course *UsingCourseModel) GetByCourseId(time int, place int) error { //int为映射，作为筛选条件
	d := DB.Self.Find(course, "courseid = ?", course.CourseId)
	return d.Error
}

// Favorite course.(TODO)
// Fixed by shiina orez at 2019.11.24, add default return value in function body
func (course *UsingCourseModel) Favorite() error {
	return nil
}

// Unfavorite course.(TODO)
// Fixed by shiina orez at 2019.11.24, add default return value in function body
func (course *UsingCourseModel) Unfavorite() error {
	return nil
}

// Search course by name, courseId or teacher
// Use fulltext search, against and match
func AgainstAndMatchCourses(kw string, page, limit int) (*sql.Rows, error) {
	rows, err := DB.Self.Exec("SELECT name, course_id, teacher FROM using_course WHERE MATCH (name, courseId, teacher) AGAINST (?) LIMIT ? OFFSET ?;", kw, limit, (page-1)*limit).Rows()
	if err != nil {
		return nil, err
	}
	return rows, nil
}
