package model

// 评课物理表
type CourseEvaluationModel struct {
	Id 					uint64 `gorm:"column:id; primary_key; AUTO_INCREMENT"`
	CourseId 			string `gorm:"column:course_id"` 				// 课程id
	CourseName      	string `gorm:"column:course_name"`       		// 课程名称
	UserId				uint64 `gorm:"column:user_id"` 					// 用户id
	Rate            	uint8  `gorm:"column:rate"`              		// 评价星级
	AttendanceCheckType uint8  `gorm:"column:attendance_check_type"` 	// 考勤方式
	ExamCheckType       uint8  `gorm:"column:exam_check_type"`      	// 考核方式
	Content         	string `gorm:"column:content"`           		// 评课内容
	LikeNum         	int64  `gorm:"column:like_num"`          		// 点赞数
	CommentNum      	int64  `gorm:"column:comment_num"`       		// 一级评论数
	Tags 				string `gorm:"column:tags"` 					// 标签id列表，逗号分隔
	IsAnonymous     	bool   `gorm:"column:is_anonymous"`      		// 是否匿名
	IsValid 			bool   `gorm:"column:is_valid"`					// 是否有效，未被折叠
	Time 				string `gorm:"column:time"`						// 时间，时间戳
}

// 评论物理表
type CommentModel struct {
	Id				uint64 `gorm:"column:id; primary_key; AUTO_INCREMENT"`
	UserId			uint64 `gorm:"column:user_id"`
	ParentId		uint64 `gorm:"column:parent_id"`
	CommentTargetId uint64 `gorm:"column:comment_target_id"`
	Content			string `gorm:"column:content"`
	LikeNum			uint64 `gorm:"column:like_num"`
	IsRoot			bool   `gorm:"column:is_root"`
	Time			string `gorm:"column:time"`
	SubCommentNum	uint64 `gorm:"column:sub_comment_num"`
}

// 发布评课
type EvaluationPublish struct {

}

// 评课信息
type EvaluationInfo struct {

}

// 评论信息
type CommentInfo struct {

}
