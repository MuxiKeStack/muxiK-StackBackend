package model

import (
	"encoding/json"
)

func (m *Message) TableName() string {
	return "message"
}

func CreateMessage(pub *MessagePub) error {
	courseJson, err := json.Marshal(pub.CourseInfo)
	if err != nil {
		return err
	}
	d := DB.Self.Create(&Message{
		PubUserId:  pub.PubUserId,
		SubUserId:  pub.SubUserId,
		IsLike:     pub.IsLike,
		IsRead:     pub.IsRead,
		Reply:      pub.Reply,
		Time:       pub.Time,
		CourseInfo: string(courseJson),
	})
	return d.Error
}

func GetMessages(page, limit, uid uint32) (*[]Message, error) {
	var messages []Message
	DB.Self.Where("sub_user_id = ?", uid).Find(&messages).Limit(limit).Offset((page - 1) * limit).Order("time desc")
	return &messages, nil
}

func GetCount(uid uint32) (uint32, error) {
	var count uint32
	d := DB.Self.Table("message").Where("sub_user_id = ? AND is_read = ?", uid, 0).Count(&count)
	return count, d.Error
}

func ReadAll(uid uint32) error {
	var messages []Message
	d := DB.Self.Table("message").Where("sub_user_id = ? AND is_read = ?", uid, '0').Find(&messages)
	d = DB.Self.Model(&messages).Update("is_read", 1)
	return d.Error
}
