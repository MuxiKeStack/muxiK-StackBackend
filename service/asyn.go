package service

import (
	"encoding/json"
	"fmt"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"

	"github.com/MuxiKeStack/muxiK-StackBackend/log"
)

type AsynGradeMsgModel struct {
	model.LoginModel
	UserId uint32
	New    bool
}

// 异步成绩服务，sub 端
func AsynGradeService() {
	log.Info("Asyn grade service starts...")

	var data = &AsynGradeMsgModel{}

	ch := model.GradeSubClient.Self.Channel()
	for msg := range ch {
		// fmt.Println(msg)
		err := json.Unmarshal([]byte(msg.Payload), data)
		if err != nil {
			log.Error(fmt.Sprintf("asyn grade service unmarshal msg(%s) error", msg.Payload), err)
		}
		// fmt.Println(data)

		if data.New {
			GradeImportService(data.UserId, data.Sid, data.Password)
		} else {
			GradeCrawlHandler(data.UserId, data.Sid, data.Password)
		}
	}
}

// 异步成绩服务，pub 端
func GradeServiceHandler(gMsg *AsynGradeMsgModel) {
	msg, err := json.Marshal(gMsg)
	if err != nil {
		log.Error(fmt.Sprintf("marshal asyn-grade-msg error for (userId=%d, sid=%s, psw=%s)", gMsg.UserId, gMsg.Sid, gMsg.Password), err)
		return
	}

	if err := model.PublishMsg(msg, model.GradeChan); err != nil {
		log.Error(fmt.Sprintf("asyn-grade-msg publish error for (%s)", string(msg)), err)
		return
	}
	log.Info("publish msg OK")
}

// 异步课程缓存
// func AsyncStoreSelfCourse() {

// }
