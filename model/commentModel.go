package model

// 评课物理表
type CourseEvaluationModel struct {
	Id                  uint64 `gorm:"column:id; primary_key; AUTO_INCREMENT"`
	CourseId            string `gorm:"column:course_id"`             // 课程id
	CourseName          string `gorm:"column:course_name"`           // 课程名称
	UserId              uint64 `gorm:"column:user_id"`               // 用户id
	Rate                uint8  `gorm:"column:rate"`                  // 评价星级
	AttendanceCheckType uint8  `gorm:"column:attendance_check_type"` // 考勤方式
	ExamCheckType       uint8  `gorm:"column:exam_check_type"`       // 考核方式
	Content             string `gorm:"column:content"`               // 评课内容
	LikeNum             uint64 `gorm:"column:like_num"`              // 点赞数
	CommentNum          uint64 `gorm:"column:comment_num"`           // 一级评论数
	Tags                string `gorm:"column:tags"`                  // 标签id列表，逗号分隔
	IsAnonymous         bool   `gorm:"column:is_anonymous"`          // 是否匿名
	IsValid             bool   `gorm:"column:is_valid"`              // 是否有效，未被折叠
	Time                string `gorm:"column:time"`                  // 时间，时间戳
}

// 评论物理表
type CommentModel struct {
	Id              uint64 `gorm:"column:id; primary_key; AUTO_INCREMENT"`
	UserId          uint64 `gorm:"column:user_id"`
	ParentId        uint64 `gorm:"column:parent_id"`
	CommentTargetId uint64 `gorm:"column:comment_target_id"`
	Content         string `gorm:"column:content"`
	LikeNum         uint64 `gorm:"column:like_num"`
	IsRoot          bool   `gorm:"column:is_root"`
	Time            string `gorm:"column:time"`
	SubCommentNum   uint64 `gorm:"column:sub_comment_num"`
}

type EvaluationLikeModel struct {
	Id           uint64 `gorm:"column:id; primary_key; AUTO_INCREMENT"`
	EvaluationId uint64 `gorm:"column:evaluation_id"`
	UserId       uint64 `gorm:"column:user_id"`
}

// 评论点赞中间表
type CommentLikeModel struct {
	Id        uint64 `gorm:"column:id; primary_key; AUTO_INCREMENT"`
	CommentId uint64 `gorm:"column:comment_id"`
	UserId    uint64 `gorm:"column:user_id"`
}

// 发布评课
type EvaluationPublish struct {
	CourseId            string  `json:"course_id"`
	CourseName          string  `json:"course_name"`
	Rate                uint8   `json:"rate"`
	AttendanceCheckType uint8   `json:"attendance_check_type"`
	ExamCheckType       uint8   `json:"exam_check_type"`
	Content             string  `json:"content"`
	IsAnonymous         bool    `json:"is_anonymous"`
	Tags                []uint8 `json:"tags"`
}

// 评课信息
type EvaluationInfo struct {
	CourseId            string
	CourseName          string
	teacher             string
	Rate                uint8
	AttendanceCheckType uint8
	ExamCheckType       uint8
	Content             string
	Time                string
	IsAnonymous         bool
	IsLike              bool
	LikeNum             uint64
	CommentNum          uint64
	Tags                []uint8
	UserInfo            UserInfo
}

// 评论信息
type CommentInfo struct {
	Content           string
	LikeNum           uint64
	IsLike            bool
	Time              string
	UserInfo          UserInfo
	CommentTargetInfo UserInfo
}

// 新增评论请求模型
type NewCommentRequest struct {
	Content     string `json:"content,omitempty"`
	IsAnonymous bool   `json:"is_anonymous,omitempty"`
}
