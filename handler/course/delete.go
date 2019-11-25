package course

import (
	_ "github.com/MuxiKeStack/muxiK-StackBackend/handler"
	_ "github.com/MuxiKeStack/muxiK-StackBackend/model"
	_ "github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	_ "github.com/MuxiKeStack/muxiK-StackBackend/pkg/token"

	"github.com/gin-gonic/gin"
)

//删除一个课程
func DeleteCourse(c *gin.Context) {

}
