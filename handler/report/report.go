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
	"sync"
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
// @Param payload body Request true "请求数据"
// @Success 200 {object} report.Response
// @Router /evaluation/{id}/report/ [post]
func ReportEvaluation(c *gin.Context) {
	// 获取评课 ID
	eid, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetParam, nil, err.Error())
		return
	}
	// 获取当前用户 ID
	uid, ok := c.Get("id")
	if !ok || uid == nil {
		handler.SendUnauthorized(c, errno.ErrAuthFailed, nil, "API get user id from token failed.")
		return
	}

	var requestPayload Request
	if err := c.BindJSON(&requestPayload); err != nil {
		handler.SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}

	var existed bool = false
	var tot int = 0
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		existed = model.ReportExisted(eid, uid.(uint64))
		wg.Done()
	}()
	go func() {
		// not passed total of reports
		tot = model.CountEvaluationBeReportedTimes(eid)
		wg.Done()
	}()
	wg.Wait()

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
					Subject:     "木犀课栈左边儿胡同儿口的驿站: 举报通知",
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
