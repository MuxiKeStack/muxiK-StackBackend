package model

import "github.com/jinzhu/gorm"

// Add a new course.
func (course *UsingCourseModel) Add() error {
	d := DB.Self.Create(course)
	return d.Error
}

// Delete a course.
func (course *UsingCourseModel) Delete() error {
	d := DB.Self.Delete(evaluation)
	return d.Error
}

// Get course by its id.(course list)
func (course *UsingCourseModel) GetById() error {
	d := DB.Self.First(course, "id = ?", course.Id)
	return d.Error
}

// Get course by its type.(course list)
func (course *UsingCourseModel) GetByType() error {
    d := DB.Self.Find(course, "type = ?", course.type)
    return d.Error
}

// Get course by its teacher.(course list)
func (course *UsingCourseModel) GetByTeacher() error {
    d := DB.Self.Find(course, "teacher = ?", course.teacher)
    return d.Error
}

// Get course by its name.(TODO)(course list)
func (course *UsingCourseModel) GetByTeacher() error {
    d := DB.Self.Find(course, "name = ?", course.name)
    return d.Error
}

// Get course by its name.(TODO)(course assistant)
func (course *UsingCourseModel) GetByName(int time,int place) error {   //int为映射，作为筛选条件
    d := DB.Self.Find(course, "name = ?", course.name)
    return d.Error
}

// Get course by its teacher.(course assistant)
func (course *UsingCourseModel) GetByTeacher(int time,int place) error {   //int为映射，作为筛选条件
    d := DB.Self.Find(course, "teacher = ?", course.teacher)
    return d.Error
}

// Get course by its courseid.(course assistant)
func (course *UsingCourseModel) GetByCourseId(int time,int place) error {   //int为映射，作为筛选条件
    d := DB.Self.Find(course, "courseid = ?", course.courseid)
    return d.Error
}

// Favorite course.(TODO)
func (course *UsingCourseModel) Favorite() error {

}

// Unfavorite course.(TODO)
func (course *UsingCourseModel) Unfavorite() error {

}
