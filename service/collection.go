package service

import (
	"fmt"
	"sync"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/util"

	"github.com/MuxiKeStack/muxiK-StackBackend/log"
)

type CourseInfoForCollectionsList struct {
	Lock  *sync.Mutex
	IdMap map[uint32]*model.CourseInfoForCollections
}

// Get collections which have been processed by userId in table page.
func GetCollectionListForTables(userId uint32, tableId uint32) ([]*model.CourseInfoInTableCollection, error) {
	var result []*model.CourseInfoInTableCollection
	courseIds, err := model.GetCourseHashIdsFromCollection(userId)
	if err != nil {
		log.Error("GetCollectionsByUserId function error", err)
		return nil, err
	}

	// Get all course ids in user's table, as skip items when gets collections
	hiddenCourseIds, err := GetAllClassIdsByTableId(userId, tableId)
	if err != nil {
		log.Error("GetAllClassIdsInTables function error", err)
		return nil, err
	}

	errChan := make(chan error, 1)
	wg := &sync.WaitGroup{}
	finished := make(chan bool, 1)
	dataChan := make(chan *model.CourseInfoInTableCollection)

	for _, courseId := range courseIds {
		// skip the course which has been added into tables
		if exist := hiddenCourseIds[courseId]; exist {
			continue
		}

		wg.Add(1)
		go func(courseId string) {
			defer wg.Done()

			// Get all classes of the course by its hash id.
			classes, err := model.GetClassesByCourseHash(courseId)
			if err != nil {
				log.Error(fmt.Sprintf("GetClassesByCourseHash function error; [Hash: %s]", courseId), err)
				errChan <- err
			}

			// No classes to choose.
			if classes == nil || len(classes) == 0 {
				log.Info(fmt.Sprintf("Table Service Error: [Hash: %s] get classes is nil.", courseId))
				return
			}

			// Parse and produce classes, to get avaliable class infos.
			classInfos, err := GetClassInfoInCollection(classes)
			if err != nil {
				log.Error(fmt.Sprintf("GetClassInfoInCollection function error; [Hash: %s].", courseId), err)
				errChan <- err
			}

			data := &model.CourseInfoInTableCollection{
				CourseId:   courseId,
				CourseName: classes[0].Name,
				ClassSum:   len(classes),
				Classes:    classInfos,
				Type:       int8(classes[0].Type),
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
			result = append(result, class)
		}
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan:
		return nil, err
	}

	return result, nil
}

// Get class info by original classes for table image's collection.
func GetClassInfoInCollection(classes []*model.UsingCourseModel) ([]*model.ClassInfoInCollections, error) {
	var infos []*model.ClassInfoInCollections

	// 选课手册课程，复数个课堂
	for _, class := range classes {
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

		// 一个课堂，复数个时间段
		var timeInfos []*model.ClassTimeInfoInCollections
		for i := 0; i < len(times); i++ {
			// 解析时间和周次信息
			timeInfoItems, err := util.ParseClassTime(times[i], weeks[i])
			if err != nil {
				log.Error("ParseClassTime function error", err)
				return nil, err
			}

			// 可能有复数个时间段
			// 2020-8-2 fix: time 出现个例 3-4,9-10#4
			for _, item := range timeInfoItems {
				timeInfos = append(timeInfos, &model.ClassTimeInfoInCollections{
					Time:      fmt.Sprintf("%d-%d", item.Start, item.End),
					Day:       item.Day,
					Weeks:     item.Weeks,
					WeekState: item.WeekState,
				})
			}
		}

		infos = append(infos, &model.ClassInfoInCollections{
			ClassId:   class.ClassId,
			ClassName: class.Name,
			Teacher:   class.Teacher,
			Times:     timeInfos,
			Places:    places,
		})
	}
	return infos, nil
}

// Get collections in course list page.
func GetCollectionList(userId uint32, lastId, limit int32) ([]*model.CourseInfoForCollections, error) {
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

			course, ok, err := model.GetHistoryCoursePartInfoByHashId(record.CourseHashId)
			if err != nil {
				log.Error("GetHistoryCoursePartInfoByHashId function error", err)
				errChan <- err
			} else if !ok {
				log.Info(fmt.Sprintf("No this history course; hash = %s", record.CourseHashId))
				return
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

	var result []*model.CourseInfoForCollections
	for _, id := range ids {
		// 2020-6-9: fix: history course may not exist
		c, ok := courseInfoList.IdMap[id]
		if !ok {
			continue
		}
		result = append(result, c)
	}

	return result, nil
}
