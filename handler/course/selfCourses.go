package course

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"
	"github.com/MuxiKeStack/muxiK-StackBackend/util"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
)

type selfCoursesResponse struct {
	Sum  int                           `json:"sum"`
	Data *[]service.ProducedCourseItem `json:"data"`
}

// @Summary 获取个人历史课程
// @Tags course
// @Param token header string true "token"
// @Param year query string true "学年，默认获取全部"
// @Param term query string true "学期，1/2/3，默认0表示获取全部"
// @Param data body model.LoginModel true "data"
// @Success 200 {object} course.selfCoursesResponse
// @Router /user/courses/ [post]
func GetSelfCourses(c *gin.Context) {
	userId := c.MustGet("id").(uint32)

	var l model.LoginModel
	if err := c.ShouldBindJSON(&l); err != nil {
		handler.SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}

	// 验证学号是否属于该用户
	if service.GetSidById(userId) != l.Sid {
		handler.SendBadRequest(c, errno.ErrAuthFailed, nil, "sid error for this user")
		return
	}

	// 判断学号密码是否正确
	if err := util.LoginRequest(l.Sid, l.Password); err != nil {
		handler.SendResponse(c, errno.ErrAuthFailed, nil)
		return
	}

	year := c.DefaultQuery("year", "0")
	term := c.DefaultQuery("term", "0")

	data, err := service.GetSelfCourseList(userId, l.Sid, l.Password, year, term)
	/*	if err != nil {
			// 从教务处获取选课课表失败，从本地数据库中获取备份
			log.Error("GetSelfCourseList function error", err)
			//handler.SendError(c, errno.ErrGetSelfCourses, nil, err.Error())
			// return
			log.Info("Enter GetSelfCourseListFromLocal function")

			if data, err = service.GetSelfCourseListFromLocal(userId); err != nil {
				handler.SendError(c, errno.ErrGetSelfCourses, nil, err.Error())
				return
			}

			// 获取成功则将数据备份到本地数据库
		} else if err = service.SavingCourseDataToLocal(userId, data); err != nil {
			handler.SendError(c, errno.ErrSavesDataToLocal, nil, err.Error())
			return
		}*/
	if err != nil {
		handler.SendError(c, errno.ErrGetSelfCourses, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, &selfCoursesResponse{
		Sum:  len(*data),
		Data: data,
	})

	/* ------ 成绩爬取服务 ------ */

	// 环境变量设置，是否爬取成绩
	// export MUXIKSTACK_GRADE_CRAWL=true
	if gradeSwitch := viper.GetBool("grade_crawl"); !gradeSwitch {
		return
	}

	// 检查是否加入成绩共享计划
	if ok, err := model.UserHasLicence(userId); err != nil {
		log.Error("UserHasLicence function error", err)
		return
	} else if !ok {
		return
	}
	log.Info("Crawling grades begins")

	// 有许可，导入成绩至统计样本
	if err := service.GradeImportService(userId, l.Sid, l.Password); err != nil {
		log.Error("Grade import failed", err)
		return
	}
	log.Info("Grade sample imported successfully")
}
