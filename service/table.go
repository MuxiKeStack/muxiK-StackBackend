package service

import (
	"strconv"
	"strings"
	"sync"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"

	"github.com/lexkong/log"
)

type classInfoList struct {
	Lock *sync.Mutex
	list []model.ClassInfo
}

// 根据课表model获取课表返回详情
func GetTableInfoByTableModel_2(table *model.ClassTableModel) (*model.ClassTableInfo, error) {
	// return if has no class
	if table.Classes == "" {
		return &model.ClassTableInfo{
			TableId: table.Id,
			TableName: table.Name,
		}, nil
	}

	ids := strings.Split(table.Classes, ",")

	var classList classInfoList

	wg := &sync.WaitGroup{}
	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	// 并发获取课堂列表
	for _, id := range ids {
		wg.Add(1)

		go func(id string) {
			defer wg.Done()

			classInfo, err := GetClassInfoForTableById(id)
			if err != nil {
				errChan <- err
				return
			}
			classList.Lock.Lock()
			defer classList.Lock.Unlock()
			classList.list = append(classList.list, *classInfo)

		}(id)
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

	info := &model.ClassTableInfo{
		TableId:   table.Id,
		TableName: table.Name,
		ClassNum:  uint32(len(ids)),
		ClassList: &classList.list,
	}

	return info, nil
}

// 根据id获取课表详情
func GetTableInfoById(id uint32) (*model.ClassTableInfo, error) {
	log.Info("GetTableInfoById function is called")

	table := &model.ClassTableModel{Id: id}
	if err := table.GetById(); err != nil {
		log.Error("table.GetById function error", err)
		return nil, err
	}

	return GetTableInfoByTableModel(table)
}

// 根据id获取课表所用的课堂详情
func GetClassInfoForTableById(id string) (*model.ClassInfo, error) {
	class, err := model.GetClassByHashId(id)
	if err != nil {
		log.Error("GetClassByHashId function error", err)
		return nil, err
	}

	// Get course's hash id
	couseId, err := model.GetCourseHashIdById(class.CourseId)
	if err != nil {
		log.Error("GetCourseHashIdById function err", err)
		return nil, err
	}

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
	var timeInfos []model.ClassTimeInfo
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

		timeInfos = append(timeInfos, model.ClassTimeInfo{
			Start:     int8(start),
			Duration:  int8(stop - start + 1),
			Day:       int8(day),
			Weeks:     week[:splitWeek],
			WeekState: int8(weekState),
		})
	}

	info := &model.ClassInfo{
		CourseId:  couseId,
		ClassId:   id,
		ClassName: class.Name,
		Teacher:   class.Teacher,
		Places:    &places,
		Times:     &timeInfos,
	}

	return info, nil
}

// 根据课表model获取课表返回详情
func GetTableInfoByTableModel(table *model.ClassTableModel) (*model.ClassTableInfo, error) {
	// return if has no class
	if table.Classes == "" {
		return &model.ClassTableInfo{
			TableId: table.Id,
			TableName: table.Name,
		}, nil
	}

	ids := strings.Split(table.Classes, ",")

	var classList []model.ClassInfo

	wg := &sync.WaitGroup{}
	errChan := make(chan error, 1)
	finished := make(chan bool, 1)
	classChan := make(chan *model.ClassInfo, 20)

	// 并发获取课堂列表
	for _, id := range ids {
		wg.Add(1)

		go func(id string) {
			defer wg.Done()

			classInfo, err := GetClassInfoForTableById(id)
			if err != nil {
				errChan <- err
				return
			}
			classChan <- classInfo

		}(id)
	}

	go func() {
		wg.Wait()
		close(classChan)
	}()

	go func() {
		for class := range classChan {
			classList = append(classList, *class)
		}
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan:
		// 会不会goroutine泄露
		//for range classChan {}
		return nil, err
	}

	info := &model.ClassTableInfo{
		TableId:   table.Id,
		TableName: table.Name,
		ClassNum:  uint32(len(ids)),
		ClassList: &classList,
	}

	return info, nil
}
