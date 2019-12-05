package model

import "time"

// 评课物理表
type CourseEvaluationModel struct {
	Id                  uint32     `gorm:"column:id; primary_key; AUTO_INCREMENT"`
	CourseId            string     `gorm:"column:course_id"`             // 课程id
	CourseName          string     `gorm:"column:course_name"`           // 课程名称
	UserId              uint32     `gorm:"column:user_id"`               // 用户id
	Rate                float32    `gorm:"column:rate"`                  // 评价星级
	AttendanceCheckType uint8      `gorm:"column:attendance_check_type"` // 考勤方式
	ExamCheckType       uint8      `gorm:"column:exam_check_type"`       // 考核方式
	Content             string     `gorm:"column:content"`               // 评课内容
	LikeNum             uint32     `gorm:"colum:like_num"`               // 点赞数
	CommentNum          uint32     `gorm:"column:comment_num"`           // 一级评论数
	Tags                string     `gorm:"column:tags"`                  // 标签id列表，逗号分隔
	IsAnonymous         bool       `gorm:"column:is_anonymous"`          // 是否匿名
	IsValid             bool       `gorm:"column:is_valid"`              // 是否有效，未被折叠
	Time                *time.Time `gorm:"column:time"`                  // 时间，时间戳
	DeletedAt           *time.Time `gorm:"column:deleted_at"`
}

// 评课点赞中间表
type EvaluationLikeModel struct {
	Id           uint32 `gorm:"column:id; primary_key; AUTO_INCREMENT"`
	EvaluationId uint32 `gorm:"column:evaluation_id"`
	UserId       uint32 `gorm:"column:user_id"`
}

// 评课信息
type EvaluationInfo struct {
	Id                  uint32            `json:"id"`
	CourseId            string            `json:"course_id"`
	CourseName          string            `json:"course_name"`
	Teacher             string            `json:"teacher"`
	Rate                float32           `json:"rate"`
	AttendanceCheckType string            `json:"attendance_check_type"`
	ExamCheckType       string            `json:"exam_check_type"`
	Content             string            `json:"content"`
	Time                int64             `json:"time"`
	IsAnonymous         bool              `json:"is_anonymous"`
	IsLike              bool              `json:"is_like"`
	LikeNum             uint32            `json:"like_num"`
	CommentNum          uint32            `json:"comment_num"`
	Tags                []string          `json:"tags"`
	UserInfo            *UserInfoResponse `json:"user_info"`
	IsValid             bool              `json:"is_valid"`
	CanDelete           bool              `json:"can_delete"`
}
