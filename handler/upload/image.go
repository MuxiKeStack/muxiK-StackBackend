package upload

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"
	"github.com/gin-gonic/gin"
)

type ImageUrlModel struct {
	Url string `json:"url"`
}

// @Tags upload
// @Summary 上传文件，图片，返回url，即为上传头像
// @Description
// @Param token header string "token"
// @Accept multipart/form-data
// @Produce json
// @Success 200 "OK"
// @Router /upload/image [post]
func Image(c *gin.Context) {
	image, _, err := c.Request.FormFile("image")
	if err != nil {
		handler.SendError(c, errno.ErrGetFile, nil, err.Error())
	}
	id, _ := c.Get("id")
	url, err := service.UploadImage(id.(uint32), image) //TODO hwo to the image parameter to pointer, wo need a io.Reader
	if err != nil {
		handler.SendError(c, errno.ErrUploadFile, nil, err.Error())
	}
	handler.SendResponse(c, nil, ImageUrlModel{Url: url})
}
