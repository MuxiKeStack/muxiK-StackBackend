package model

import (
	"encoding/json"
)

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
	DB.Self.Where("user_id = ?", uid).Find(&messages).Limit(limit).Offset(offset)
	return &messages, nil
}
