package service

import (
	"strconv"
	"strings"
	"sync"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"
)

type classInfoList struct {
	Lock *sync.Mutex
	list []model.ClassInfo
}

// 根据课表model获取课表返回详情
func GetTableInfoByTableModel(table *model.ClassTableModel) (*model.ClassTableInfo, error) {
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
	table := &model.ClassTableModel{Id: id}
	if err := table.GetById(); err != nil {
		return nil, err
	}

	return GetTableInfoByTableModel(table)
}

// 根据id获取课表所用的课堂详情
func GetClassInfoForTableById(id string) (*model.ClassInfo, error) {
	class, err := model.GetClassById(id)
	if err != nil {
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

	// 解析上课时间
	times := []string{class.Time1}
	if class.Time2 != "" {
		times = append(times, class.Time2)
	}
	if class.Time3 != "" {
		times = append(times, class.Time3)
	}

	var timeInfos []model.ClassTimeInfo
	for _, time := range times {
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

		//...

		timeInfos = append(timeInfos, model.ClassTimeInfo{
			Start:     int8(start),
			Duration:  int8(stop - start + 1),
			Day:       0,
			Weeks:     "",
			WeekState: 0,
		})
	}

	info := &model.ClassInfo{
		CourseId:   string(class.CourseId),
		ClassId:    id,
		ClassName:  class.Name,
		Teacher:    class.Teacher,
		Places:     &places,
		Times:      &timeInfos,
	}

	return info, nil
}
