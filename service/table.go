package service

import (
	"strconv"
	"strings"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"
)

func GetTableInfoByTableModel(table *model.ClassTableModel) (*model.ClassTableInfo, error) {
	ids := strings.Split(table.Classes, ",")
	var classList []model.ClassInfo
	for _, id := range ids {
		classInfo, err := GetClassInfoById(id)
		if err != nil {
			return nil, nil
		}
		classList = append(classList, *classInfo)
	}

	info := &model.ClassTableInfo{
		TableId:   table.Id,
		TableName: table.Name,
		ClassNum:  uint32(len(ids)),
		ClassList: &classList,
	}

	return info, nil
}

func GetClassTableInfoById(id uint32) (*model.ClassTableInfo, error) {
	table := &model.ClassTableModel{Id: id}
	if err := table.GetById(); err != nil {
		return nil, err
	}

	return GetTableInfoByTableModel(table)
}

func GetClassInfoById(id string) (*model.ClassInfo, error) {
	class, err := model.GetClassById(id)
	if err != nil {
		return nil, err
	}

	places := []string{class.Place1}
	if class.Place2 != "" {
		places = append(places, class.Place2)
	}
	if class.Place3 != "" {
		places = append(places, class.Place3)
	}

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
