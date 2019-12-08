package report

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/constvar"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/gin-gonic/gin"
	"strconv"
	"sync"
)

// @Summary 折叠评课
// @Tags report
// @Param id path integer true "评课ID"
// @Success 200 "OK"
// @Router /evaluation/{id}/block/ [get]
func BlockEvaluation(c *gin.Context) {
	eid, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetParam, nil, err.Error())
		return
	}
	if model.CountEvaluationBeReportedTimes(eid) >= constvar.AllowRemindAdminLimit {
		// block
		evaluation := &model.CourseEvaluationModel{Id: uint32(eid)}
		err := evaluation.GetById()
		if err != nil {
			handler.SendError(c, errno.ErrGetEvaluation, nil, err.Error())
			return
		}
		if evaluation.IsValid == false {
			handler.SendResponse(c, nil, "Evaluation already be blocked.")
		}
		// block it
		evaluation.IsValid = false
		errChan, done := make(chan error), make(chan struct{})

		defer close(errChan)

		go func() {
			err := evaluation.Block()
			if err != nil {
				errChan <- err
			}
			done <- struct{}{}
		}()
		go func() {
			reports, err := model.GetAllReportOfEvaluation(eid)
			if err != nil {
				errChan <- err
				return
			}
			wg := sync.WaitGroup{}
			for _, report := range reports {
				report.Pass = true
				wg.Add(1)
				go func(r model.ReportModel) {
					err := r.Update()
					if err != nil {
						errChan <- err
					}
					wg.Done()
				}(report)
			}
			wg.Wait()
			done <- struct{}{}
		}()
		for i := 0; i < 2; i++ {
			select {
			case _ = <-done:
			case e := <-errChan:
				{
					handler.SendError(c, e, nil, err.Error())
					i = 3
				}
			}
		}
		handler.SendResponse(c, nil, "Block Evaluation Successful! All about reports be passed.")
	} else {
		handler.SendResponse(c, nil, "Evaluation: "+strconv.Itoa(int(eid))+" be reported less than "+strconv.Itoa(constvar.AllowRemindAdminLimit)+" times.")
	}
}
