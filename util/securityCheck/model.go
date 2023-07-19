package securityCheck

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/MuxiKeStack/muxiK-StackBackend/util"

	"github.com/MuxiKeStack/muxiK-StackBackend/log"
	"github.com/spf13/viper"
)

// QQ小程序内容安全检测

// QQ小程序 Access Token 管理
type accessTokenManager struct {
	Token     string
	CreateAt  *time.Time
	ExpiresIn time.Duration
}

var (
	QQAppSecret string
	QQAppID     string

	accessToken = &accessTokenManager{ExpiresIn: 7200 * time.Second}

	imgSecCheckURL    = "https://api.q.qq.com/api/json/security/ImgSecCheck?access_token="
	msgSecCheckURL    = "https://api.q.qq.com/api/json/security/MsgSecCheck?access_token="
	accessTokenGetURL = "https://api.q.qq.com/api/getToken?grant_type=client_credential&appid=%s&secret=%s"
)

func QQSecInit() {
	QQAppID = viper.GetString("QQ_APPID")
	QQAppSecret = viper.GetString("QQ_APP_SECRET")

	accessToken.loadToken()

	// fmt.Println(QQAppID, QQAppSecret, accessToken.Token)

	imgSecCheckURL += accessToken.Token
	msgSecCheckURL += accessToken.Token
}

type QQGetTokenPayload struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int32  `json:"expires_in"`
	ErrCode     int32  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

func (t *accessTokenManager) loadToken() error {
	fmt.Println("Begin to load or refresh access token from QQ server...")

	resp, err := http.Get(fmt.Sprintf(accessTokenGetURL, QQAppID, QQAppSecret))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var obj QQGetTokenPayload
	if err := json.Unmarshal([]byte(body), &obj); err != nil {
		log.Error("json unmarshal to QQGetTokenPayload error", err)
		return err
	}

	fmt.Printf("QQ access token: old token: %s; new token: %s\n", t.Token, obj.AccessToken)

	t.Token = obj.AccessToken
	t.CreateAt = util.GetCurrentTime()

	return nil
}

func (t *accessTokenManager) check() error {
	now := time.Now()
	if t.CreateAt.Add(t.ExpiresIn).Sub(now) <= 0 {
		// 过期，更新 token
		if err := t.loadToken(); err != nil {
			log.Error("Refresh access token failed", err)
			return err
		}
		log.Info("Refresh access token OK")
	}

	fmt.Printf("QQ access token info: createAt=%v, expiresIn=%v, sub time from now=%v\n",
		t.CreateAt, t.ExpiresIn, t.CreateAt.Add(t.ExpiresIn).Sub(now))

	return nil
}

// 定时更新 QQ APP token
func RefreshTokenScheduled() {
	for {
		// 提前10分钟更新
		time.Sleep(accessToken.ExpiresIn - time.Minute*10)

		if err := accessToken.loadToken(); err != nil {
			log.Error("Refresh access token failed", err)
			continue
		}

		log.Info("Refresh QQ access token OK")
	}
}
