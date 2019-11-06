package service

import (
	"encoding/json"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
)

func MessageList(offset, limit, uid uint32) (*[]model.MessageSub, error) {
	messages, err := model.GetMessages(offset, limit, uid)
	if err != nil {
		return nil, nil
	}
	var messageSubs []model.MessageSub
	for _, message := range *messages {
		messageSub := model.MessageSub{
			IsLike: message.IsLike,
			IsRead: message.IsRead,
			Reply:  message.Reply,
			Time:   message.Time,
		}
		userInfo, err := GetUserInfoRById(uid)
		if err != nil {
			return nil, err
		}

		var courseInfo model.CourseInfo
		err = json.Unmarshal([]byte(message.CourseInfo), &courseInfo)
		if err != nil {
			return nil, err
		}
		messageSub.CourseInfo = courseInfo
		messageSub.UserInfo = *userInfo
		messageSubs = append(messageSubs, messageSub)
	}
	return &messageSubs, nil
}
