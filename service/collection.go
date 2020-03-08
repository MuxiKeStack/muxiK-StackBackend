package service

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"

	"github.com/lexkong/log"
)

type CourseInfoForCollectionsList struct {
	Lock  *sync.Mutex
	IdMap map[uint32]*model.CourseInfoForCollections
}

// Get collections which have been processed by userId in table page.
func GetCollectionListForTables(userId uint32, tableId uint32) (*[]model.CourseInfoInTableCollection, error) {
	var result []model.CourseInfoInTableCollection
	courseIds, err := model.GetCourseHashIdsFromCollection(userId)
	if err != nil {
		log.Error("GetCollectionsByUserId function error", err)
		return nil, err
	}

	hiddenCourseIds, err := GetAllClassIdsByTableId(userId, tableId)
	if err != nil {
		log.Error("GetAllClassIdsInTables function error", err)
		return nil, err
	}

	errChan := make(chan error, 1)
	wg := &sync.WaitGroup{}
	finished := make(chan bool, 1)
	dataChan := make(chan *model.CourseInfoInTableCollection, 10)

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

			data := &model.CourseInfoInTableCollection{
				CourseId:   courseId,
				CourseName: (*classes)[0].Name,
				ClassSum:   len(*classes),
				Classes:    classInfos,
				Type:       int8((*classes)[0].Type),
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
				log.Error("ParseInt function error when parsing start", err)
				return nil, err
			}

			stop, err := strconv.ParseInt(time[split1+1:split2], 10, 8)
			if err != nil {
				log.Error("ParseInt function error when parsing stop", err)
				return nil, err
			}

			// 上课星期
			day, err := strconv.ParseInt(time[split2+1:], 10, 8)
			if err != nil {
				log.Error("ParseInt function error when parsing day", err)
				return nil, err
			}

			// 解析上课周次和单双周状态
			week := weeks[i]
			splitWeek := strings.Index(week, "#")

			weekState, err := strconv.ParseInt(week[splitWeek+1:], 10, 8)
			if err != nil {
				log.Error("ParseInt function error when parsing weekState", err)
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
			ClassId:   class.ClassId,
			ClassName: class.Name,
			//TeachingClassId: class.ClassId,
			Teacher: class.Teacher,
			Times:   &timeInfos,
			Places:  &places,
		})
	}
	return &infos, nil
}

// Get collections in course list page.
func GetCollectionList(userId uint32, lastId, limit int32) (*[]model.CourseInfoForCollections, error) {
	records, err := model.GetCollectionsByUserId(userId, lastId, limit)
	if err != nil {
		log.Error("GetCollectionsByUserId function error", err)
		return nil, err
	}

	var ids []uint32
	for _, record := range *records {
		ids = append(ids, record.Id)
	}

	courseInfoList := CourseInfoForCollectionsList{
		Lock:  new(sync.Mutex),
		IdMap: make(map[uint32]*model.CourseInfoForCollections, len(*records)),
	}

	wg := sync.WaitGroup{}
	finished := make(chan bool, 1)
	errChan := make(chan error, 1)

	for _, record := range *records {
		wg.Add(1)

		go func(record model.CourseListModel) {
			defer wg.Done()

			course, err := model.GetHistoryCourseByHashId(record.CourseHashId)
			if err != nil {
				log.Error("GetHistoryCourseByHashId function error", err)
				errChan <- err
			}

			attendanceTypeNum, err := model.GetTheMostAttendanceCheckType(record.CourseHashId)
			if err != nil {
				log.Error("GetTheMostAttendanceCheckType function error", err)
				errChan <- err
			}

			examCheckTypeNum, err := model.GetTheMostExamCheckType(record.CourseHashId)
			if err != nil {
				log.Error("GetTheMostExamCheckType function error", err)
				errChan <- err
			}

			tags, err := GetTwoMostTagsOfOneCourse(record.CourseHashId)
			if err != nil {
				log.Error("GetTwoMostTagsOfOneCourse function error", err)
				errChan <- err
			}

			info := &model.CourseInfoForCollections{
				Id:                  record.Id,
				CourseId:            course.Hash,
				CourseName:          course.Name,
				Teacher:             course.Teacher,
				EvaluationNum:       course.StarsNum,
				Rate:                course.Rate,
				AttendanceCheckType: GetAttendanceCheckTypeByCode(attendanceTypeNum),
				ExamCheckType:       GetExamCheckTypeByCode(examCheckTypeNum),
				Tags:                &tags,
			}

			courseInfoList.Lock.Lock()
			defer courseInfoList.Lock.Unlock()

			courseInfoList.IdMap[info.Id] = info
		}(record)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan:
		return nil, err
	}

	var result []model.CourseInfoForCollections
	for _, id := range ids {
		result = append(result, *courseInfoList.IdMap[id])
	}

	return &result, nil
}
