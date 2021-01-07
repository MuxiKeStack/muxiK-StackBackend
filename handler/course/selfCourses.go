package course

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

type selfCoursesResponse struct {
	Sum  int                           `json:"sum"`
	Data []*service.ProducedCourseItem `json:"data"`
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

	var loginRequest model.LoginModel
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		handler.SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}

	// 验证学号是否属于该用户
	if service.GetSidById(userId) != loginRequest.Sid {
		handler.SendBadRequest(c, errno.ErrAuthFailed, nil, "sid error for this user")
		return
	}

	// 判断学号密码是否正确
	// 放到获取教务课程中判断，减少一次模拟登陆的时间
	// if err := util.LoginRequest(loginRequest.Sid, loginRequest.Password); err != nil {
	// 	handler.SendResponse(c, errno.ErrAuthFailed, nil)
	// 	return
	// }

	year := c.DefaultQuery("year", "0")
	term := c.DefaultQuery("term", "0")

	var data []*service.ProducedCourseItem
	var err error

	// 先从教务系统获取个人课程,
	// 如获取成功，则更新 Redis 中的缓存数据,
	// 如获取失败，则直接从 Redis 获取缓存数据。
	// To do: redis 异步更新

	// 另一个想法：
	// 以缓存数据为主，先从缓存中读取数据，如有则直接返回，
	// 并且异步爬取教务课表数据，更新缓存数据；
	// 如从缓存数据中读取失败（无数据），则从教务获取，再更新缓存。
	// 但是这样就有一个问题，学号密码的正确性就检验不了了：
	// 先检查 redis，有则直接返回，然后爬取课表失败，因为账号错误，
	// 然后就一直查询一直错误，缓存中数据永远得不到更新。
	// 所以要采用这种方式，就需要确保学号密码的验证。

	data, err = service.GetSelfCourseList(userId, loginRequest.Sid, loginRequest.Password, year, term)
	if err != nil {
		// 首先判断是否是用户名密码错误
		if err == errno.ErrAuthFailed {
			handler.SendResponse(c, errno.ErrAuthFailed, nil)
			return
		}

		// 从教务处获取选课课表失败，获取缓存数据
		log.Error("GetSelfCourseList function error", err)
		log.Info("Try to get courses cache data from redis...")

		data, err = service.GetSelfCoursesFromLocalCache(userId, year, term)
		if err != nil {
			log.Error("Getting courses from cache failed", err)
			handler.SendError(c, errno.ErrGetSelfCourses, nil, "getting courses from xk and cache failed")
			return
		}
	} else {
		// 获取教务课程成功则将数据备份到 redis
		if err := service.SelfCoursesCacheStoreToRedis(userId, data); err != nil {
			log.Error("Storing courses into redis failed", err)
		}
		log.Info("Storing courses into redis succeed")
	}

	handler.SendResponse(c, nil, &selfCoursesResponse{
		Sum:  len(data),
		Data: data,
	})

	/* ------ 成绩服务 ------ */

	gradeMsg := &service.AsynGradeMsgModel{
		LoginModel: model.LoginModel{
			Sid:      loginRequest.Sid,
			Password: loginRequest.Password,
		},
		UserId: userId,
		New:    false,
	}
	service.GradeServiceHandler(gradeMsg)
}
