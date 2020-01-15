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

// TODO FIX upload to change avatar directly.

// @Tags upload
// @Summary 上传文件，图片，返回url，即为上传头像
// @Description
// @Param token header string true "token"
// @Param image formData file true "二进制图片文件"
// @Accept multipart/form-data
// @Produce json
// @Success 200 {object} ImageUrlModel
// @Router /upload/image/ [post]
func Image(c *gin.Context) {
	image, header, err := c.Request.FormFile("image")
	if err != nil {
		handler.SendError(c, errno.ErrGetFile, nil, err.Error())
		return
	}
	dataLen := header.Size
	id, _ := c.Get("id")

	url, err := service.UploadImage(header.Filename, id.(uint32), image, dataLen)
	if err != nil {
		handler.SendError(c, errno.ErrUploadFile, nil, err.Error())
		return
	}
	handler.SendResponse(c, nil, ImageUrlModel{Url: url})
}
