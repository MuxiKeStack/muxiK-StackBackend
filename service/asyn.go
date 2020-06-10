package service

import (
	"encoding/json"
	"fmt"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"

	"github.com/lexkong/log"
)

type AsynGradeMsgModel struct {
	model.LoginModel
	UserId uint32
	New    bool
}

// 异步成绩服务
func AsynGradeService() {
	log.Info("Asyn grade service starts...")

	var data = &AsynGradeMsgModel{}

	ch := model.GradeSubClient.Self.Channel()
	for msg := range ch {
		fmt.Println(msg)
		err := json.Unmarshal([]byte(msg.Payload), data)
		if err != nil {
			log.Errorf(err, "asyn grade service unmarshal msg(%s) error", msg.Payload)
		}
		fmt.Println(data)

		if data.New {
			GradeImportService(data.UserId, data.Sid, data.Password)
		} else {
			GradeCrawlHandler(data.UserId, data.Sid, data.Password)
		}
	}
}
