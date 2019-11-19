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
	d := DB.Self.Create(Message{
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

func GetMessages(offset, limit, uid uint32) (*[]Message, error) {
	var messages []Message
	DB.Self.Where("sub_user_id = ?", uid).Find(&messages).Limit(limit).Offset(offset * limit).Order("time desc")
	return &messages, nil
}

func GetCount(uid uint32) uint32 {
	var count uint32
	DB.Self.Where("sub_user_id = ? AND is_read = ?", uid, '0').Count(&count)
	return count
}

func ReadAll(uid uint32) error {
	var messages []Message
	var count, i uint32
	d := DB.Self.Where("sub_user_id = ? AND is_read = ?", uid, '0').Find(&messages).Count(&count)
	for i = 0; i < count; i++ {
		d = DB.Self.Model(messages[i]).Update("is_read", 1)
	}
	return d.Error
}
