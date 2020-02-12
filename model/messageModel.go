package model

// Message represents a json for messaging
type Message struct {
	Id        uint32 `gorm:"column:id; primary_key" `
	PubUserId uint32 `gorm:"column:pub_user_id"`
	SubUserId uint32 `gorm:"column:sub_user_id"`
	Kind      uint8  `gorm:"column:kind"`
	IsRead    bool   `gorm:"column:is_read"`
	Reply     string `gorm:"column:reply"`
	Time      string `gorm:"column:time"`

	// MessageInfo string `gorm:"column:message_info"`
	// 消息提醒的一些信息，显示出来的字段，用于回复评论的id	o represents evaluation or comment information
	CourseId   string `gorm:"column:course_id"`
	CourseName string `gorm:"column:course_name"`
	Teacher    string `gorm:"column:teacher"`

	//点击消息提醒中的内容 跳转到 评课 需要 EvaluationId
	EvaluationId uint32 `gorm:"column:evaluation_id"`

	//即为操作对象的内容，如果是对于评课则是评课，如果是对评论则是原评论内容。
	Content string `gorm:"column:content"`

	//用于对评课==一级评论(只需要EnvaluationID)，评论的回复==二级评论(一级评论的ID ParentCommentId+目标用户 Sid).
	//用来发二级评论
	Sid             string `gorm:"column:sid"`
	ParentCommentId string `gorm:"column:parent_comment_id"`
}

// MessagePub 消息提醒的发送者。
type MessagePub struct {
	PubUserId uint32 `json:"pub_user_id"`
	SubUserId uint32 `json:"sub_user_id"`
	Kind      uint8  `json:"kind"`
	IsRead    bool   `json:"is_read"`
	Reply     string `json:"reply"`
	Time      string `json:"time"`

	// MessageInfo string `gorm:"column:message_info"`
	// 消息提醒的一些信息，显示出来的字段，用于回复评论的id	o represents evaluation or comment information
	CourseId   string `json:"course_id"`
	CourseName string `json:"course_name"`
	Teacher    string `json:"teacher"`

	//点击消息提醒中的内容 跳转到 评课 需要 EvaluationId
	EvaluationId uint32 `json:"evaluation_id"`

	//即为操作对象的内容，如果是对于评课则是评课，如果是对评论则是原评论内容。
	Content string `json:"content"`

	//用于对评课==一级评论(只需要EnvaluationID)，评论的回复==二级评论(一级评论的ID ParentCommentId+目标用户 Sid).
	//用来发二级评论
	Sid             string `json:"sid"`
	ParentCommentId string `json:"parent_comment_id"`
}

// MessageSub 消息提醒的接受model
type MessageSub struct {
	UserInfo UserInfoRequest `json:"user_info"`
	//kind 区分 点赞->0 评论->1 举报->2 系统提醒->3
	Kind   uint8  `json:"kind"`
	IsRead bool   `json:"is_read"`
	Reply  string `json:"reply"`
	Time   string `json:"time"`
	// MessageInfo string `gorm:"column:message_info"`
	// 消息提醒的一些信息，显示出来的字段，用于回复评论的id	o represents evaluation or comment information
	CourseId   string `json:"course_id"`
	CourseName string `json:"course_name"`
	Teacher    string `json:"teacher"`

	//点击消息提醒中的内容 跳转到 评课 需要 EvaluationId
	EvaluationId uint32 `json:"evaluation_id"`

	//即为操作对象的内容，如果是对于评课则是评课，如果是对评论则是原评论内容。
	Content string `json:"content"`

	//用于对评课==一级评论(只需要EnvaluationID)，评论的回复==二级评论(一级评论的ID ParentCommentId+目标用户 Sid).
	//用来发二级评论
	Sid             string `json:"sid"`
	ParentCommentId string `json:"parent_comment_id"`
}
