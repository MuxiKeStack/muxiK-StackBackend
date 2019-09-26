package model

// 评课
type CourseCommentTable struct {
	CourseCommentId string   `json:"course_comment_id"` // 评课id
	CourseName      string   `json:"course_name"`       // 课程名称
	Attendance      string   `json:"attendance"`        // 考勤方式
	ExamWay         string   `json:"exam_way"`          // 考核方式
	Content         string   `json:"content"`           // 评课内容
	IsAnonymous     bool     `json:"is_anonymous"`      // 是否匿名
	LikeNum         int64    `json:"like_num"`          // 点赞数
	Star            int64    `json:"star"`              // 评价星级
	CommentNum      int64    `json:"comment_num"`       // 一级评论数
	LikeUsers       []string `json:"like_users"`        // 点赞用户（用户id）
}

// 评论
type CommentTable struct {

}
