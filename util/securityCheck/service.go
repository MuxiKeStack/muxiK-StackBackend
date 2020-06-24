package securityCheck

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/lexkong/log"
)

type imgCheckReq struct {
	AppID string `json:"appid"`
}

type msgCheckReq struct {
	AppID   string `json:"appid"`
	Content string `json:"content"`
}

type checkResponse struct {
	ErrCode int32  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

// 消息文本检测
func MsgSecCheck(content string) (bool, error) {
	// accessToken.check()

	data, err := json.Marshal(msgCheckReq{
		AppID:   QQAppID,
		Content: content,
	})
	if err != nil {
		return false, err
	}

	// fmt.Println(string(data))

	resp, err := http.Post(msgSecCheckURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Error("QQ msg security check err", err)
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	// fmt.Println(string(body))

	var rp checkResponse
	if err := json.Unmarshal(body, &rp); err != nil {
		return false, err
	}

	// fmt.Println(rp)
	if rp.ErrCode != 0 {
		log.Info(fmt.Sprintf("msg security check failed. code: %d; msg: %s.", rp.ErrCode, rp.ErrMsg))
		return false, nil
	}
	log.Info("msg security check OK.")

	return true, nil
}

// 图片检测
func ImgSecCheck(image string) (bool, error) {
	return true, nil
}
