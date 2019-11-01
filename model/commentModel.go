package model

// 评课物理表
type CourseEvaluationModel struct {
	Id                  uint32 `gorm:"column:id; primary_key; AUTO_INCREMENT"`
	CourseId            string `gorm:"column:course_id"`             // 课程id
	CourseName          string `gorm:"column:course_name"`           // 课程名称
	UserId              uint32 `gorm:"column:user_id"`               // 用户id
	Rate                uint8  `gorm:"column:rate"`                  // 评价星级
	AttendanceCheckType uint8  `gorm:"column:attendance_check_type"` // 考勤方式
	ExamCheckType       uint8  `gorm:"column:exam_check_type"`       // 考核方式
	Content             string `gorm:"column:content"`               // 评课内容
	LikeNum             uint32 `gorm:"column:like_num"`              // 点赞数
	CommentNum          uint32 `gorm:"column:comment_num"`           // 一级评论数
	Tags                string `gorm:"column:tags"`                  // 标签id列表，逗号分隔
	IsAnonymous         bool   `gorm:"column:is_anonymous"`          // 是否匿名
	IsValid             bool   `gorm:"column:is_valid"`              // 是否有效，未被折叠
	Time                string `gorm:"column:time"`                  // 时间，时间戳
}

// 评论物理表
type CommentModel struct {
	Id              uint32 `gorm:"column:id; primary_key; AUTO_INCREMENT"`
	UserId          uint32 `gorm:"column:user_id"`
	ParentId        uint32 `gorm:"column:parent_id"`
	CommentTargetId uint32 `gorm:"column:comment_target_id"`
	Content         string `gorm:"column:content"`
	LikeNum         uint32 `gorm:"column:like_num"`
	IsRoot          bool   `gorm:"column:is_root"`
	Time            string `gorm:"column:time"`
	SubCommentNum   uint32 `gorm:"column:sub_comment_num"`
	IsAnonymous     bool   `gorm:"column:is_anonymous"`
	IsValid         bool   `gorm:"column:is_valid"`
}

type EvaluationLikeModel struct {
	Id           uint32 `gorm:"column:id; primary_key; AUTO_INCREMENT"`
	EvaluationId uint32 `gorm:"column:evaluation_id"`
	UserId       uint32 `gorm:"column:user_id"`
}

// 评论点赞中间表
type CommentLikeModel struct {
	Id        uint32 `gorm:"column:id; primary_key; AUTO_INCREMENT"`
	CommentId uint32 `gorm:"column:comment_id"`
	UserId    uint32 `gorm:"column:user_id"`
}

// 评课信息
type EvaluationInfo struct {
	Id                  uint32            `json:"id"`
	CourseId            string            `json:"course_id"`
	CourseName          string            `json:"course_name"`
	Teacher             string            `json:"teacher"`
	Rate                uint8             `json:"rate"`
	AttendanceCheckType uint8             `json:"attendance_check_type"`
	ExamCheckType       uint8             `json:"exam_check_type"`
	Content             string            `json:"content"`
	Time                string            `json:"time"`
	IsAnonymous         bool              `json:"is_anonymous"`
	IsLike              bool              `json:"is_like"`
	LikeNum             uint32            `json:"like_num"`
	CommentNum          uint32            `json:"comment_num"`
	Tags                []uint8           `json:"tags"`
	UserInfo            *UserInfoResponse `json:"user_info"`
}

// 评论信息
type CommentInfo struct {
	Id             uint32            `json:"id"`
	Content        string            `json:"content"`
	LikeNum        uint32            `json:"like_num"`
	IsLike         bool              `json:"is_like"`
	Time           string            `json:"time"`
	IsAnonymous    bool              `json:"is_anonymous"`
	UserInfo       *UserInfoResponse `json:"user_info"`
	TargetUserInfo *UserInfoResponse `json:"target_user_info"`
}

// 新增评论请求模型
type NewCommentRequest struct {
	Content     string `json:"content"`
	IsAnonymous bool   `json:"is_anonymous"`
}

// 返回的评论列表，一级评论模型
type ParentCommentInfo struct {
	Id              uint32            `json:"id"`
	CommentId       uint32            `json:"comment_id"`
	Content         string            `json:"content"`
	LikeNum         uint32            `json:"like_num"`
	IsLike          bool              `json:"is_like"`
	Time            string            `json:"time"`
	IsAnonymous     bool              `json:"is_anonymous"`
	UserInfo        *UserInfoResponse `json:"user_info"`
	SubCommentsNum  uint32            `json:"sub_comments_num"`
	SubCommentsList *[]CommentInfo    `json:"sub_comments_list"`
}
