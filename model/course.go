package model

import (
	"fmt"
)

const (
	typeTemp              = "AND RIGHT(LEFT(using_course.course_id, 4), 1) = %s "
	typeCourseTemp        = "AND using_course.type = %s "
	typeHistoryCourseTemp = "AND history_course.type = %s "
	academyTemp           = "AND using_course.academy = '%s' "
	weekdayTemp           = "AND (RIGHT(`time1`, 1) = %s OR RIGHT(`time2`, 1) = %s OR RIGHT(`time3`, 1) = %s) "
	nPlaceTemp            = "AND LEFT(`place1`, 1) = 'N' AND (`place2` = '' OR LEFT(`place2`, 1) = 'N') AND (`place3` = '' OR LEFT(`place3`, 1) = 'N') "
	bPlaceTemp            = "AND LEFT(`place1`, 1) != 'N' AND (`place2` = '' OR LEFT(`place2`, 1) != 'N') AND (`place3` = '' OR LEFT(`place3`, 1) != 'N') "
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
func (course *HistoryCourseModel) GetHistoryByHash() error {
	d := DB.Self.First(course, "hash = ?", course.Hash)
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

// Get course by its courseId.(course assistant)
func (class *UsingCourseModel) GetByCourseId(time int, place int) error { // int为映射，作为筛选条件
	d := DB.Self.Find(class, "course_id = ?", class.CourseId)
	return d.Error
}

// Update history course.
func (course *HistoryCourseModel) Update() error {
	return DB.Self.Save(course).Error
}

// Judge whether a course has already favorited by the current user,
// return record id and bool type.
func HasFavorited(userId uint32, hash string) (uint32, bool) {
	var data CourseListModel
	d := DB.Self.Where("user_id = ? AND course_hash_id = ? ", userId, hash).First(&data)
	return data.Id, !d.RecordNotFound()
}

// Favorite a course by the current user.
func Favorite(userId uint32, hash string) error {
	var data = CourseListModel{
		CourseHashId: hash,
		UserId:       userId,
	}

	d := DB.Self.Create(&data)
	return d.Error
}

// Cancel a course's favorite by the current user.
func Unfavorite(id uint32) error {
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
	// where := "MATCH (using_course.name, using_course.course_id, using_course.teacher) AGAINST ('" + kw + "') "
	// 全文索引无效（无法作用于course_id和teacher, 且name不是模糊查询），但没找到原因，只能改为like
	kw = "%" + kw + "%"
	where := "using_course.name like ? or using_course.course_id like ? or using_course.teacher like ?"
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
		Where(where, kw, kw, kw).
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
	// where := "MATCH (`name`, `teacher`) AGAINST ('" + kw + "') "
	kw = "%" + kw + "%"
	where := "name like ? or teacher like ?"
	if t != "" {
		where += fmt.Sprintf(typeHistoryCourseTemp, t)
	}

	DB.Self.Table("history_course").Where(where, kw, kw).Limit(limit).Offset((page - 1) * limit).Find(courses)
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
	// fmt.Println(where)
	if where == "" {
		DB.Self.Table("using_course").
			Select("using_course.*, history_course.stars_num, history_course.rate").
			Joins("LEFT JOIN history_course ON using_course.hash = history_course.hash").
			Group("using_course.hash").
			Order("history_course.stars_num  DESC, history_course.rate DESC").
			Limit(limit).Offset((page - 1) * limit).Find(&courses)
	} else {
		whereFix := where[4:]
		// fmt.Println(whereFix)
		DB.Self.Table("using_course").
			Select("using_course.*, history_course.stars_num, history_course.rate").
			Joins("LEFT JOIN history_course ON using_course.hash = history_course.hash").
			Where(whereFix).
			Group("using_course.hash").
			Order("history_course.stars_num  DESC, history_course.rate DESC").
			Limit(limit).Offset((page - 1) * limit).Find(&courses)
	}
	return *courses, nil
}

// Get all history courses
func AllHistoryCourses(page, limit uint64, t string) ([]HistoryCourseModel, error) {
	courses := &[]HistoryCourseModel{}
	if t != "" {
		DB.Self.Table("history_course").Where("type = ?", t).
			Order("history_course.stars_num  DESC, history_course.rate DESC").
			Limit(limit).Offset((page - 1) * limit).Find(courses)
	} else {
		DB.Self.Table("history_course").
			Order("history_course.stars_num  DESC, history_course.rate DESC").
			Limit(limit).Offset((page - 1) * limit).Find(courses)
	}
	return *courses, nil
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

// 根据课程id获取教师名
func GetTeacherByCourseId(id string) (string, error) {
	var data struct{ Teacher string }
	d := DB.Self.Table("history_course").Select("teacher").Where("hash = ?", id).Scan(&data)
	return data.Teacher, d.Error
}

// Get history course by hash id.
func GetHistoryCourseByHashId(id string) (*HistoryCourseModel, error) {
	var course = &HistoryCourseModel{Hash: id}
	err := course.GetHistoryByHash()
	return course, err
}

func GetHistoryCoursePartInfoByHashId(hash string) (*HistoryCourseModel, bool, error) {
	var data = &HistoryCourseModel{}
	d := DB.Self.Table("history_course").
		Select("hash, name, teacher, rate, stars_num").
		Where("hash = ?", hash).Scan(data)

	if d.RecordNotFound() {
		return nil, false, nil
	}
	return data, true, d.Error
}

func GetGradeInfoFromHistoryCourseInfo(hash string) (*HistoryCourseModel, bool, error) {
	var data = &HistoryCourseModel{}
	d := DB.Self.Table("history_course").
		Select("total_grade, usual_grade, grade_sample_size, grade_section_1, grade_section_2, grade_section_3").
		Where("hash = ?", hash).Scan(data)
	if d.RecordNotFound() {
		return nil, false, nil
	}
	return data, true, d.Error
}

// 判断课程是否存在
func IsCourseExisting(hash string) bool {
	var data struct{ Id uint32 }

	DB.Self.Table("history_course").Select("id").Where("hash = ?", hash).Scan(&data)
	return data.Id != 0
}

// 创建新的历史课程
func (course *HistoryCourseModel) New() error {
	return DB.Self.Create(course).Error
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
