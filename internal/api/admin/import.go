package admin

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liony823/tools/a2r"
	"github.com/liony823/tools/apiresp"
	"github.com/liony823/tools/errs"
	"github.com/liony823/tools/utils/encrypt"
	"github.com/openimsdk/chat/pkg/common/config"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/common/xlsx"
	"github.com/openimsdk/chat/pkg/common/xlsx/model"
	"github.com/openimsdk/chat/pkg/protocol/chat"
)

// @Summary		通过Excel导入用户
// @Description	通过Excel文件批量导入用户信息
// @Tags			import
// @Id	importUserByXlsx
// @Accept			multipart/form-data
// @Produce		json
// @Param			data	formData	file	true	"Excel文件"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/user/import/xlsx [post]
func (o *Api) ImportUserByXlsx(c *gin.Context) {
	formFile, err := c.FormFile("data")
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	ip, err := o.GetClientIP(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	file, err := formFile.Open()
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	defer file.Close()
	var users []model.User
	if err := xlsx.ParseAll(file, &users); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("xlsx file parse error "+err.Error()))
		return
	}
	us, err := o.xlsx2user(users)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	ctx := o.WithAdminUser(mctx.WithApiToken(c, imToken))
	apiresp.GinError(c, o.registerChatUser(ctx, ip, us))
}

// @Summary		通过JSON导入用户
// @Description	通过JSON数据批量导入用户信息
// @Tags			import
// @Id	importUserByJson
// @Accept			json
// @Produce		json
// @Param			data	body []chat.RegisterUserInfo	true	"用户信息列表"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/user/import/json [post]
func (o *Api) ImportUserByJson(c *gin.Context) {
	req, err := a2r.ParseRequest[struct {
		Users []*chat.RegisterUserInfo `json:"users"`
	}](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	ip, err := o.GetClientIP(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	ctx := o.WithAdminUser(mctx.WithApiToken(c, imToken))
	apiresp.GinError(c, o.registerChatUser(ctx, ip, req.Users))
}

// @Summary		获取导入模板
// @Description	下载用户导入的Excel模板文件
// @Tags			import
// @Id	importTemplate
// @Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Success		200	{file}	binary	"Excel模板文件"
// @Router			/user/import/xlsx [get]
func (o *Api) BatchImportTemplate(c *gin.Context) {
	md5Sum := md5.Sum(config.ImportTemplate)
	md5Val := hex.EncodeToString(md5Sum[:])
	if c.GetHeader("If-None-Match") == md5Val {
		c.Status(http.StatusNotModified)
		return
	}
	c.Header("Content-Disposition", "attachment; filename=template.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Length", strconv.Itoa(len(config.ImportTemplate)))
	c.Header("ETag", md5Val)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", config.ImportTemplate)
}

func (o *Api) xlsx2user(users []model.User) ([]*chat.RegisterUserInfo, error) {
	chatUsers := make([]*chat.RegisterUserInfo, len(users))
	for i, info := range users {
		if info.Nickname == "" {
			return nil, errs.ErrArgs.WrapMsg("nickname is empty")
		}
		if info.AreaCode == "" || info.PhoneNumber == "" {
			return nil, errs.ErrArgs.WrapMsg("areaCode or phoneNumber is empty")
		}
		if info.Password == "" {
			return nil, errs.ErrArgs.WrapMsg("password is empty")
		}
		if !strings.HasPrefix(info.AreaCode, "+") {
			return nil, errs.ErrArgs.WrapMsg("areaCode format error")
		}
		if _, err := strconv.ParseUint(info.AreaCode[1:], 10, 16); err != nil {
			return nil, errs.ErrArgs.WrapMsg("areaCode format error")
		}
		gender, _ := strconv.Atoi(info.Gender)
		chatUsers[i] = &chat.RegisterUserInfo{
			UserID:      info.UserID,
			Nickname:    info.Nickname,
			FaceURL:     info.FaceURL,
			Birth:       o.xlsxBirth(info.Birth).UnixMilli(),
			Gender:      int32(gender),
			AreaCode:    info.AreaCode,
			PhoneNumber: info.PhoneNumber,
			Email:       info.Email,
			Account:     info.Account,
			Password:    encrypt.Md5(info.Password),
		}
	}
	return chatUsers, nil
}

func (o *Api) xlsxBirth(s string) time.Time {
	if s == "" {
		return time.Now()
	}
	var separator byte
	for _, b := range []byte(s) {
		if b < '0' || b > '9' {
			separator = b
		}
	}
	arr := strings.Split(s, string([]byte{separator}))
	if len(arr) != 3 {
		return time.Now()
	}
	year, _ := strconv.Atoi(arr[0])
	month, _ := strconv.Atoi(arr[1])
	day, _ := strconv.Atoi(arr[2])
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	if t.Before(time.Date(1900, 0, 0, 0, 0, 0, 0, time.Local)) {
		return time.Now()
	}
	return t
}
