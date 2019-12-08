package report

import (
	"fmt"
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/constvar"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/util/smtpMail"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

type Request struct {
	Reason string `json:"reason"`
}

type Response struct {
	Fail   bool   `json:"fail"`
	Reason string `json:"reason"`
}

// @Summary 举报评课
// @Tags report
// @Param token header string true "token"
// @Param id path integer true "评课ID"
// @Param payload body {object} report.Request true "请求数据"
// @Success 200 {object} report.Response
// @Router /evaluation/{id}/report/ [post]
func ReportEvaluation(c *gin.Context) {
	eid, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetParam, nil, err.Error())
		return
	}
	uid, _ := c.Get("id")
	var requestPayload Request
	if err := c.BindJSON(&requestPayload); err != nil {
		handler.SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}
	reportExistedCh := make(chan bool)
	beReportedTotCh := make(chan int)
	defer close(reportExistedCh)
	defer close(beReportedTotCh)

	go func() {
		reportExistedCh <- model.ReportExisted(eid, uid.(uint64))
	}()
	go func() {
		// not passed total of reports
		beReportedTotCh <- model.CountEvaluationBeReportedTimes(eid)
	}()
	var existed bool
	var tot int
	for i := 0; i < 2; i++ {
		select {
		case tot = <-beReportedTotCh:
		case existed = <-reportExistedCh:
		}
	}
	if existed {
		handler.SendResponse(c, nil, Response{
			Fail:   true,
			Reason: "You have been reported this evaluation!",
		})
		return
	} else {
		tot += 1
		newReport := model.ReportModel{
			EvaluationId: eid,
			UserId:       uid.(uint64),
			Pass:         false,
			Reason:       requestPayload.Reason,
		}
		err = newReport.Create()
		if err != nil {
			handler.SendError(c, errno.ErrCreateReport, nil, err.Error())
			return
		}
		// send a remind to admin
		if tot >= constvar.AllowRemindAdminLimit {
			// goroutine to send main
			go func() {
				mailContent := strings.Replace(constvar.EmailTemp, "REPORT_TOT", strconv.Itoa(tot), 1)
				mailContent = strings.Replace(mailContent, "EVALUATION_ID", strconv.Itoa(int(eid)), 1)
				reasonTemp := "<p>用户（%d）：%s</p>"
				allReasons := []string{}
				allReports, err := model.GetAllReportOfEvaluation(eid)
				if err != nil {
					fmt.Println("Send Report Mail to Admin Failed! Reason:", err.Error())
					return
				} else {
					for _, report := range allReports {
						allReasons = append(allReasons, fmt.Sprintf(reasonTemp, report.UserId, report.Reason))
					}
				}
				mailContent = strings.Replace(mailContent, "ALL_REASON", strings.Join(allReasons, "\n"), 1)
				log.Infof("Start to send email to: %s", constvar.DefaultAdminEmailAddr)

				err = smtpMail.SendMail("muxistudio@qq.com", viper.GetString("authcode"), []string{constvar.DefaultAdminEmailAddr}, smtpMail.Content{
					NickName:    "木犀课栈: 评课举报通知",
					User:        "muxistudio@qq.com",
					Subject:     "木犀课栈左边儿胡同儿的驿站: 举报通知",
					Body:        mailContent,
					ContentType: "Content-Type: text/html; charset=UTF-8",
				})
				if err != nil {
					fmt.Println("Send Report Mail to Admin Failed! Reason:", err.Error())
					return
				}
				log.Info("Email sent successful.")
			}()
		}
		handler.SendResponse(c, nil, Response{
			Fail:   false,
			Reason: "",
		})
	}
}
