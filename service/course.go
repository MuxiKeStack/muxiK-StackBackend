package service

import (
	_ "fmt"
	_ "sync"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"

	_ "github.com/lexkong/log"
)

// Get course info.
// Fixed by shiina orez at 2019.11.24, add default return value in function body
func CourseList(id, userId uint32) (*model.CourseInfo, error) {
	return nil, nil
}
