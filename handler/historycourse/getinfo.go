package course

import (
	_ "github.com/MuxiKeStack/muxiK-StackBackend/handler"
	_ "github.com/MuxiKeStack/muxiK-StackBackend/model"
	_ "github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	_ "github.com/MuxiKeStack/muxiK-StackBackend/pkg/token"

	"github.com/gin-gonic/gin"
)

//获取历史课程信息
func GetHistoryCourseInfo(c *gin.Context) {

}
