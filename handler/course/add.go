package course

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/token"

	"github.com/gin-gonic/gin"
)

//增加一个课程
func AddCourse(c *gin.Context) {
	var data UsingCourseInfo
	if err := c.ShouldBindJSON(&data); err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	if err := data.Import(); err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}
	
	var course = &model.UsingCourseModel{
		Hash          data.hash
		Name          data.name
		Credit        data.credit
		Teacher       data.teacher
		CourseId      data.courseId
		ClassId       data.classId
		Type          data.type
		CreditType    data.creditType
		TotalScore    data.totalScore
		OrdinaryScore data.ordinaryScore
		Time1         data.time1
		Place1        data.place1
		Time2         data.time2
		Place2        data.place2
		Time3         data.time3
		Place3        data.place1
		Weeks1        data.weeks1
		Weeks2        data.weeks1
		Weeks3        data.weeks1
		Region        data.region
	}
	
	handler.SendResponse(c, nil, course})
}
