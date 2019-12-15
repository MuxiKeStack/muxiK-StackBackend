package service

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"

	"github.com/lexkong/log"
)

// Get collections which have been processed by userId in table page.
func GetCollectionsList(userId uint32) (*[]model.CourseInfoInCollections, error) {
	var result []model.CourseInfoInCollections
	courseIds, err := model.GetCollectionsByUserId(userId)
	if err != nil {
		log.Error("GetCollectionsByUserId function error", err)
		return nil, err
	}

	hiddenCourseIds, err := GetAllClassIdsInTables(userId)
	if err != nil {
		log.Error("GetAllClassIdsInTables function error", err)
		return nil, err
	}

	errChan := make(chan error, 1)
	wg := &sync.WaitGroup{}
	finished := make(chan bool, 1)
	dataChan := make(chan *model.CourseInfoInCollections, 10)

	for _, courseId := range courseIds {
		// skip the course which has been added into tables
		if ok := hiddenCourseIds[courseId]; ok {
			continue
		}

		wg.Add(1)
		go func(courseId string) {
			defer wg.Done()

			classes, err := model.GetClassesByCourseHash(courseId)
			if err != nil {
				log.Error("GetClassesByCourseHash function error", err)
				errChan <- err
			}

			classInfos, err := GetClassInfoInCollection(classes)
			if err != nil {
				log.Error("GetClassInfoInCollection function error", err)
				errChan <- err
			}

			data := &model.CourseInfoInCollections{
				CourseId:   courseId,
				CourseName: (*classes)[0].Teacher,
				ClassSum:   len(*classes),
				Classes:    classInfos,
			}
			dataChan <- data

		}(courseId)
	}

	go func() {
		wg.Wait()
		close(dataChan)
	}()

	go func() {
		for class := range dataChan {
			result = append(result, *class)
		}
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan:
		return nil, err
	}

	return &result, nil
}

// Get class info by original classes for table image's collection.
func GetClassInfoInCollection(classes *[]model.UsingCourseModel) (*[]model.ClassInfoInCollections, error) {
	var infos []model.ClassInfoInCollections

	for _, class := range *classes {
		// 解析上课地点
		places := []string{class.Place1}
		if class.Place2 != "" {
			places = append(places, class.Place2)
		}
		if class.Place3 != "" {
			places = append(places, class.Place3)
		}

		// 获取上课周次和单双周状态
		weeks := []string{class.Weeks1}
		if class.Weeks2 != "" {
			weeks = append(weeks, class.Weeks2)
		}
		if class.Weeks3 != "" {
			weeks = append(weeks, class.Weeks3)
		}

		// 获取课堂上课时间
		times := []string{class.Time1}
		if class.Time2 != "" {
			times = append(times, class.Time2)
		}
		if class.Time3 != "" {
			times = append(times, class.Time3)
		}

		// 解析上课时间详情
		var timeInfos []model.ClassTimeInfoInCollections
		for i, time := range times {
			split1 := strings.Index(time, "-")
			split2 := strings.Index(time, "#")

			start, err := strconv.ParseInt(time[:split1], 10, 8)
			if err != nil {
				log.Error("strconv.ParseInt function error when parsing start", err)
				return nil, err
			}

			stop, err := strconv.ParseInt(time[split1+1:split2], 10, 8)
			if err != nil {
				log.Error("strconv.ParseInt function error when parsing stop", err)
				return nil, err
			}

			// 上课星期
			day, err := strconv.ParseInt(time[split2+1:], 10, 8)
			if err != nil {
				log.Error("strconv.ParseInt function error when parsing day", err)
				return nil, err
			}

			// 解析上课周次和单双周状态
			week := weeks[i]
			splitWeek := strings.Index(week, "#")

			weekState, err := strconv.ParseInt(week[splitWeek+1:], 10, 8)
			if err != nil {
				log.Error("strconv.ParseInt function error when parsing weekState", err)
				return nil, err
			}

			timeInfos = append(timeInfos, model.ClassTimeInfoInCollections{
				Time:      fmt.Sprintf("%d-%d", start, stop),
				Day:       int8(day),
				Weeks:     week[:splitWeek],
				WeekState: int8(weekState),
			})
		}

		infos = append(infos, model.ClassInfoInCollections{
			ClassId:         class.Hash,
			ClassName:       class.Name,
			TeachingClassId: class.ClassId,
			Teacher:         class.Teacher,
			Times:           &timeInfos,
			Places:          &places,
		})
	}
	return &infos, nil
}
