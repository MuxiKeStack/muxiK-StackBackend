package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/util"

	"github.com/lexkong/log"
)

type ProducedCourseItem struct {
	CourseId     string `json:"course_id"`
	Name         string `json:"name"`
	Teacher      string `json:"teacher"`
	HasEvaluated bool   `json:"has_evaluated"`
	Year         string `json:"year"` // 学期，2018
	Term         string `json:"term"` // 学年，1/2/3
}

// Get one's all courses from XK.
func GetSelfCourseList(userId uint32, sid, pwd, year, term string) ([]*ProducedCourseItem, error) {
	originalCourses, err := util.GetSelfCoursesFromXK(sid, pwd, year, term)
	if err != nil {
		log.Error("GetSelfCoursesFromXK function error", err)
		return nil, err
	}

	wg := sync.WaitGroup{}
	infoChan := make(chan *ProducedCourseItem, 5)
	var list []*ProducedCourseItem
	//(*originalCourses.Items)[0].Jsxx = "2006982627/葛非,2006982646/彭熙,2006982670/刘明,2007980066/姚华雄"

	for _, item := range *originalCourses.Items {
		wg.Add(1)
		go func(item util.OriginalCourseItem) {
			defer wg.Done()

			teacher := util.GetTeachersSqStrBySplitting(item.Jsxx)
			hashId := util.HashCourseId(item.Kch, teacher)
			info := &ProducedCourseItem{
				CourseId:     hashId,
				Name:         item.Kcmc,
				Teacher:      teacher,
				HasEvaluated: model.HasEvaluated(userId, hashId),
				Year:         item.Xnm,
				Term:         item.Xqmc,
			}
			infoChan <- info
		}(item)
	}

	go func() {
		wg.Wait()
		close(infoChan)
	}()

	for info := range infoChan {
		list = append(list, info)
	}

	return list, nil
}

// 暂时废弃...
// 获取个人课程备份
func GetSelfCourseListFromLocal(userId uint32) ([]*ProducedCourseItem, error) {
	hashIdStr, err := model.GetSelfCoursesByUserId(userId)
	if err != nil {
		log.Error("GetSelfCoursesByUserId function error", err)
		return nil, err
	}

	hashIds := strings.Split(hashIdStr, ",")
	var list []*ProducedCourseItem

	for _, hashId := range hashIds {
		course := &model.UsingCourseModel{Hash: hashId}
		if err := course.GetByHash(); err != nil {
			log.Error("GetByHash function error", err)
			return nil, err
		}
		item := &ProducedCourseItem{
			CourseId:     hashId,
			Name:         course.Name,
			Teacher:      course.Teacher,
			HasEvaluated: model.HasEvaluated(userId, hashId),
		}
		list = append(list, item)
	}

	return list, nil
}

// 暂时废弃...
// 存入备份数据至本地数据库
func SavingCourseDataToLocal(userId uint32, list []*ProducedCourseItem) error {
	var record = &model.SelfCourseModel{UserId: userId}
	ok, err := record.GetByUserId()
	if err != nil {
		log.Error("SelfCourseModel.GetByUserId error", err)
		return err
	}

	curNum := len(list)

	// 记录存在且课程无变化，无需更新
	if ok && record.Num == uint32(curNum) {
		return nil
	}

	// 无记录则新添
	//var hashIds []string
	var hashIdStr string
	for i, item := range list {
		//hashIds = append(hashIds, item.CourseId)
		if i > 0 {
			hashIdStr = fmt.Sprintf("%s,%s", hashIdStr, item.CourseId)
			continue
		}
		hashIdStr += item.CourseId
	}
	//hashIdStr := strings.Join(hashIds, ",")

	record.Courses = hashIdStr
	record.Num = uint32(curNum)

	// 不存在记录则新添记录
	if !ok {
		err = record.New()
	} else {
		// 若存在且课程变化则更新
		err = record.Update()
	}

	if err != nil {
		return err
	}
	return nil
}

var (
	CourseDataKey  = "course-data"
	CourseCountKey = "course-count"
)

// 个人课程数据缓存，Redis
func SelfCoursesCacheStoreToRedis(userId uint32, list []*ProducedCourseItem) error {
	userIdStr := strconv.Itoa(int(userId))

	// 获取课程数量
	s, ok, err := model.HashGet(CourseCountKey, userIdStr)
	if err != nil {
		log.Error("Redis hashGet course count error", err)
		return err
	} else if ok {
		// 课程数量存在，根据数量判断是否有新课程

		count, err := strconv.Atoi(s)
		if err != nil {
			return err
		}

		// 无新课程，无需更新缓存数据
		if count >= len(list) {
			return nil
		}
	}

	// 序列化
	coursesBytes, err := json.Marshal(list)
	if err != nil {
		return err
	}

	// 更新课程数据
	if err = model.HashSet(CourseDataKey, userIdStr, string(coursesBytes)); err != nil {
		log.Error("Redis hashSet course data error", err)
		return err
	}

	// 更新课程数量
	if err = model.HashSet(CourseCountKey, userIdStr, len(list)); err != nil {
		log.Error("Redis hashSet course count error", err)
		return err
	}

	log.Info("Store self-courses successfully.")
	return nil
}

// 从 Redis 中获取课程缓存数据
func SelfCoursesCacheGetFromRedis(userId uint32) ([]*ProducedCourseItem, error) {
	userIdStr := strconv.Itoa(int(userId))

	s, ok, err := model.HashGet(CourseDataKey, userIdStr)
	if err != nil {
		return nil, err
	} else if !ok {
		// 无缓存数据
		return nil, errors.New("Can not find self courses cache in redis")
	}

	var courses []*ProducedCourseItem
	if err := json.Unmarshal([]byte(s), &courses); err != nil {
		return nil, err
	}

	return courses, nil
}

// 从本地缓存中获取课程数据
func GetSelfCoursesFromLocalCache(userId uint32, year, term string) ([]*ProducedCourseItem, error) {
	// 从 redis 获取全部课程数据
	list, err := SelfCoursesCacheGetFromRedis(userId)
	if err != nil {
		log.Error("SelfCoursesCacheGetFromRedis function error", err)
		return nil, err
	}

	// 不筛选，返回全部课程数据
	if year == "0" && term == "0" {
		return list, nil
	}

	var result []*ProducedCourseItem
	for _, item := range list {
		// 筛选符合学年和学期的课程
		// year, term == "0" 为默认全部
		if (year == "0" || item.Year == year) && (term == "0" || item.Term == term) {
			result = append(result, item)
		}
	}

	return result, nil
}
