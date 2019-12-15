package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"

	"github.com/lexkong/log"
)

func GetCollectionsList(userId uint32) (*[]model.CourseInfoInCollections, error) {
	var result []model.CourseInfoInCollections
	courseIds, err := model.GetCollectionsByUserId(userId)
	if err != nil {
		log.Info("GetCollectionsByUserId function error")
		return nil, err
	}

	hiddenCourseIds, err := GetAllClassIdsInTables(userId)
	if err != nil {
		return nil, err
	}

	for _, courseId := range courseIds {
		// skip the course which has been added into tables
		if ok := hiddenCourseIds[courseId]; ok {
			continue
		}

		classes, err := model.GetClassesByCourseHash(courseId)
		if err != nil {
			return nil, err
		}

		classInfos, err := GetClassInfoInCollection(classes)
		if err != nil {
			return nil, err
		}

		result = append(result, model.CourseInfoInCollections{
			CourseId:   courseId,
			CourseName: (*classes)[0].Teacher,
			ClassSum:   len(*classes),
			Classes:    classInfos,
		})
	}

	return &result, nil
}

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
			weeks = append(weeks, class.Place2)
		}
		if class.Weeks3 != "" {
			weeks = append(weeks, class.Place3)
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
				return nil, err
			}

			stop, err := strconv.ParseInt(time[split1+1:split2], 10, 8)
			if err != nil {
				return nil, err
			}

			// 上课星期
			day, err := strconv.ParseInt(time[split2+1:], 10, 8)
			if err != nil {
				return nil, err
			}

			// 解析上课周次和单双周状态
			week := weeks[i]
			splitWeek := strings.Index(week, "#")

			weekState, err := strconv.ParseInt(week[splitWeek+1:], 10, 8)
			if err != nil {
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

func GetAllClassIdsInTables(userId uint32) (map[string]bool, error) {
	tables, err := model.GetTablesByUserId(userId)
	if err != nil {
		return nil, err
	}

	result := make(map[string]bool)
	for _, table := range *tables {
		ids := strings.Split(table.Classes, ",")
		for _, id := range ids {
			result[id] = true
		}
	}
	return result, nil
}
