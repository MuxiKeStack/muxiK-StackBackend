package model

import (
	"fmt"
)

const (
	typeTemp       = "AND RIGHT(LEFT(using_course.course_id, 4), 1) = %s "
	typeCourseTemp = "AND using_course.type = %s "
	academyTemp    = "AND using_course.academy = '%s' "
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

func (*SelfCourseModel) TableName() string {
	return "self_course"
}

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

// Get course by its hash.
func (class *UsingCourseModel) GetByHash2() (uint8, error) {
	d := DB.Self.First(class, "hash = ?", class.Hash)
	if d.RecordNotFound() {
		return 1, nil
	}
	return 0, d.Error
}

// Get history course by its hash.
func (class *HistoryCourseModel) GetHistoryByHash() error {
	d := DB.Self.First(class, "hash = ?", class.Hash)
	return d.Error
}

func (class *UsingCourseModel) GetClass(courseId string, classId string) error {
	d := DB.Self.Where("course_id = ? AND class_id = ? ", courseId, classId).First(&class)
	if d.RecordNotFound() {
		return nil
	}
	return d.Error
}

func GetAllClass(hash string) ([]UsingCourseModel, error) {
	var c []UsingCourseModel
	d := DB.Self.Table("using_course").Where("hash = ?", hash).Find(&c)
	return c, d.Error
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

// Judge whether a course has already favorited by the current user,
// return record id and bool type.
func (class *UsingCourseModel) HasFavorited(userId uint32) (uint32, bool) {
	var data CourseListModel
	d := DB.Self.Where("user_id = ? AND course_hash_id = ? ", userId, class.Hash).First(&data)
	return data.Id, !d.RecordNotFound()
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

// Cancel a course's favorite by the current user.
func (class *UsingCourseModel) Unfavorite(id uint32) error {
	var data = CourseListModel{Id: id}
	d := DB.Self.Delete(&data)
	return d.Error
}

// Search course by name, courseId or teacher
// Use fulltext search, against and match
// 2020-01-15: Add New Filter: type, academy, weekday, place
// 2020-02-10: Add Join History Course SQL: join hash, get stars_num and rate
func AgainstAndMatchCourses(kw string, page, limit uint64, t, a, w, p string) ([]UsingCourseSearchModel, error) {
	courses := &[]UsingCourseSearchModel{}
	where := "MATCH (using_course.name, using_course.course_id, using_course.teacher) AGAINST ('" + kw + "') "
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

	DB.Self.Debug().Table("using_course").
		Select("using_course.*, history_course.stars_num, history_course.rate").
		Where(where).
		Joins("LEFT JOIN history_course ON using_course.hash = history_course.hash").
		Limit(limit).Offset((page - 1) * limit).
		Find(courses)

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
func AllCourses(page, limit uint64, t, a, w, p string) ([]UsingCourseSearchModel, error) {
	courses := &[]UsingCourseSearchModel{}
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
		DB.Self.Table("using_course").
			Joins("LEFT JOIN history_course ON using_course.hash = history_course.hash").
			Limit(limit).Offset((page - 1) * limit).Find(&courses)
	} else {
		DB.Self.Table("using_course").
			Joins("LEFT JOIN history_course ON using_course.hash = history_course.hash").
			Where(where).
			Limit(limit).Offset((page - 1) * limit).Find(&courses)
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

func (course *HistoryCourseModel) UpdateGradeInfo() error {
	return DB.Self.Save(course).Error
}

func GetCourseTags(hash string) ([]CourseTagModel, error) {
	var tags []CourseTagModel
	d := DB.Self.Where("course_id = ?", hash).Find(&tags)
	return tags, d.Error
}

// Get tag name by id.
func GetTagNameById32(id uint32) (string, error) {
	var tag TagModel
	d := DB.Self.Where("id = ?", id).First(&tag)
	return tag.Name, d.Error
}

/*---------------------------- SelfCourse Operation --------------------------*/

func (data *SelfCourseModel) New() error {
	return DB.Self.Create(data).Error
}

func (data *SelfCourseModel) Update() error {
	return DB.Self.Save(data).Error
}

func (data *SelfCourseModel) GetByUserId() (bool, error) {
	d := DB.Self.Where("user_id = ?", data.UserId).First(data)
	if d.RecordNotFound() {
		return false, nil
	}
	return true, d.Error
}

func GetSelfCoursesByUserId(userId uint32) (string, error) {
	var data SelfCourseModel
	d := DB.Self.Where("user_id = ?", userId).First(&data)
	if d.RecordNotFound() {
		return data.Courses, nil
	}
	return data.Courses, d.Error
}
