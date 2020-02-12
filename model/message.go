package model

func (m *Message) TableName() string {
	return "message"
}

func CreateMessage(pub *MessagePub) error {
	d := DB.Self.Create(&Message{
		PubUserId:       pub.PubUserId,
		SubUserId:       pub.SubUserId,
		Kind:            pub.Kind,
		IsRead:          pub.IsRead,
		Reply:           pub.Reply,
		Time:            pub.Time,
		CourseId:        pub.CourseId,
		CourseName:      pub.CourseName,
		Teacher:         pub.Teacher,
		EvaluationId:    pub.EvaluationId,
		Content:         pub.Content,
		Sid:             pub.Sid,
		ParentCommentId: pub.ParentCommentId,
	})
	return d.Error
}

func GetMessages(page, limit, uid uint32) (*[]Message, error) {
	var messages []Message
	DB.Self.Where("sub_user_id = ?", uid).Order("time desc").Limit(limit).Offset((page - 1) * limit).Find(&messages)
	return &messages, nil
}

func GetCount(uid uint32) (uint32, error) {
	var count uint32
	d := DB.Self.Table("message").Where("sub_user_id = ? AND is_read = ?", uid, 0).Count(&count)
	return count, d.Error
}

func ReadAll(uid uint32) error {
	d := DB.Self.Table("message").Where("sub_user_id = ?", uid).Update("is_read", 1)
	// d = DB.Self.Model(&messages).Update("is_read", 1)
	return d.Error
}
