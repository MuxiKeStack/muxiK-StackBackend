package model

// Message represents a json for messaging
type Message struct {
	Id         uint32 `gorm:"column:id; primary_key" `
	PubUserId  uint32 `gorm:"column:pub_user_id"`
	SubUserId  uint32 `gorm:"column:sub_user_id"`
	IsLike     bool   `gorm:"column:is_like"`
	IsRead     bool   `gorm:"column:is_read"`
	Reply      string `gorm:"column:reply"`
	Time       string `gorm:"column:time"`
	CourseInfo string `gorm:"column:course_info"`
}

// CourseInf	o represents evaluation or comment information
type CourseInfo struct {
	EvaluationId    uint32 `json:"evaluation_id"`
	Sid             string `json:"sid"`
	ParentCommentId string `json:"parent_comment_id"`
	CourseName      string `json:"course_name"`
	Teacher         string `json:"teacher"`
	Content         string `json:"content"`
}

//
type MessagePub struct {
	PubUserId  uint32     `json:"pub_user_id"`
	SubUserId  uint32     `json:"sub_user_id"`
	IsLike     bool       `json:"is_like"`
	IsRead     bool       `json:"is_read"`
	Reply      string     `json:"reply"`
	Time       string     `json:"time"`
	CourseInfo CourseInfo `json:"course_info"`
}

//
type MessageSub struct {
	UserInfo   UserInfoRequest `json:"user_info"`
	IsLike     bool            `json:"is_like"`
	IsRead     bool            `json:"is_read"`
	Reply      string          `json:"reply"`
	Time       string          `json:"time"`
	CourseInfo CourseInfo      `json:"course_info"`
}
