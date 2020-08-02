package service

import (
	"errors"
	"strings"
	"sync"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/util"

	"github.com/lexkong/log"
)

// Get table info by table id.
func GetTableInfoById(id uint32) (*model.ClassTableInfo, error) {
	table, err := model.GetTableById(id)
	if err != nil {
		log.Error("GetTableById function error", err)
		return nil, err
	}

	return GetTableInfoByTableModel(table)
}

// Get class info for tables by class hash id.
func GetClassInfoForTableById(hashId string, classId string) (*model.ClassInfo, error) {
	class, err := model.GetClassByHashId(hashId, classId)
	if err != nil {
		log.Error("GetClassByHashId function error", err)
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

	var timeInfos []*model.ClassTimeInfo
	for i := 0; i < len(times); i++ {
		// 解析时间和周次信息
		timeInfoItems, err := util.ParseClassTime(times[i], weeks[i])
		if err != nil {
			log.Error("ParseClassTime function error", err)
			return nil, err
		}

		// 可能有复数个时间段
		for _, item := range timeInfoItems {
			timeInfos = append(timeInfos, &model.ClassTimeInfo{
				Start:     item.Start,
				Duration:  item.End - item.Start + 1,
				Day:       item.Day,
				Weeks:     item.Weeks,
				WeekState: item.WeekState,
			})
		}
	}

	info := &model.ClassInfo{
		CourseId:  hashId,
		ClassId:   class.ClassId,
		ClassName: class.Name,
		Teacher:   class.Teacher,
		Places:    places,
		Times:     timeInfos,
		Type:      int8(class.Type),
	}

	return info, nil
}

// Get table response info by original table model.
func GetTableInfoByTableModel(table *model.ClassTableModel) (*model.ClassTableInfo, error) {
	// return if has no class
	if table.Classes == "" {
		return &model.ClassTableInfo{
			TableId:   table.Id,
			TableName: table.Name,
		}, nil
	}

	ids := strings.Split(table.Classes, ",")

	var classList []*model.ClassInfo

	wg := &sync.WaitGroup{}
	errChan := make(chan error, 1)
	finished := make(chan bool, 1)
	classChan := make(chan *model.ClassInfo, 20)

	// 并发获取课堂列表
	for _, id := range ids {
		wg.Add(1)

		go func(id string) {
			defer wg.Done()

			idSq := strings.Split(id, "#")
			// 分隔出错，存储的数据没有按照hashId#classId存储
			if len(idSq) < 2 {
				errChan <- errors.New("classes split error")
				return
			}

			classInfo, err := GetClassInfoForTableById(idSq[0], idSq[1])
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
			classList = append(classList, class)
		}
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
		ClassList: classList,
	}

	return info, nil
}

// Get all classes' id if in the table, returning a map and error, used by collection.
func GetAllClassIdsByTableId(userId uint32, tableId uint32) (map[string]bool, error) {
	table := &model.ClassTableModel{
		Id:     tableId,
		UserId: userId,
	}
	if err := table.Get(); err != nil {
		log.Error("table.GetById function error", err)
		return nil, err
	}

	result := make(map[string]bool)

	ids := strings.Split(table.Classes, ",")
	for _, id := range ids {
		index := strings.Split(id, "#")
		result[index[0]] = true
	}

	return result, nil
}
