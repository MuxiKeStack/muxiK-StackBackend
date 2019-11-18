package model

import "time"

// 父评论物理表
type ParentCommentModel struct {
	Id            string     `gorm:"column:id; primary_key"` // uuid
	UserId        uint32     `gorm:"column:user_id"`
	EvaluationId  uint32     `gorm:"column:evaluation_id"`
	Content       string     `gorm:"column:content"`
	Time          *time.Time `gorm:"column:time"`
	SubCommentNum uint32     `gorm:"column:sub_comment_num"`
	IsAnonymous   bool       `gorm:"column:is_anonymous"`
	IsValid       bool       `gorm:"column:is_valid"`
}

// 子评论物理表
type SubCommentModel struct {
	Id           string     `gorm:"column:id; primary_key"` // uuid
	UserId       uint32     `gorm:"column:user_id"`
	TargetUserId uint32     `gorm:"column:target_user_id"` // 回复的用户id
	ParentId     string     `gorm:"column:parent_id"`
	Content      string     `gorm:"column:content"`
	Time         *time.Time `gorm:"column:time"`
	IsAnonymous  bool       `gorm:"column:is_anonymous"`
	IsValid      bool       `gorm:"column:is_valid"`
}

// 评论点赞中间表
type CommentLikeModel struct {
	Id        uint32 `gorm:"column:id; primary_key; AUTO_INCREMENT"`
	UserId    uint32 `gorm:"column:user_id"`
	CommentId string `gorm:"column:comment_id"`
}

// 评论信息
type CommentInfo struct {
	Id             string            `json:"id"`
	Content        string            `json:"content"`
	LikeSum        uint32            `json:"like_sum"`
	IsLike         bool              `json:"is_like"`
	Time           int64             `json:"time"`
	IsAnonymous    bool              `json:"is_anonymous"`
	UserInfo       *UserInfoResponse `json:"user_info"`
	TargetUserInfo *UserInfoResponse `json:"target_user_info"`
}

// 返回的评论列表，一级评论模型
type ParentCommentInfo struct {
	Id              string            `json:"id"` // 父评论id
	Content         string            `json:"content"`
	LikeSum         uint32            `json:"like_sum"`
	IsLike          bool              `json:"is_like"`
	Time            int64             `json:"time"`
	IsAnonymous     bool              `json:"is_anonymous"`
	UserInfo        *UserInfoResponse `json:"user_info"`
	SubCommentsNum  uint32            `json:"sub_comments_num"`
	SubCommentsList *[]CommentInfo    `json:"sub_comments_list"`
}
