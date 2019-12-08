package course

import (
	//	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	//	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	_ "github.com/MuxiKeStack/muxiK-StackBackend/pkg/token"
	_ "github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

//获取课程信息
func GetCourseInfo(c *gin.Context) {
	/*	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		id32 := uint32(id)
		if err != nil {
			handler.SendBadRequest(c, errno.ErrGetParam, nil, err.Error())
			return
		}
		course := &model.UsingCourseModel{Id: id32}
		if err := course.GetById(); err != nil {
			log.Info("course.GetById() error.")
			return
		}

		handler.SendResponse(c, nil, course)*/

	hash := c.Param("hash")
	if hash == "" {
		log.Info("Get Param error")
		return
	}
	
	course := &model.UsingCourseModel{Hash: hash}
	if err := course.GetByHash(); err != nil {
		log.Info("course.GetByHash() error.")
		return
	}

	handler.SendResponse(c, nil, course)
}
