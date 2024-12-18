package util

import (
	"context"
	"fmt"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/liony823/tools/errs"
	"github.com/openimsdk/chat/pkg/common/mctx"
)

type Api struct {
	ImUserID        string
	ProxyHeader     string
	ChatAdminUserID string
	BasicAuthUser   string
	BasicAuthPass   string
}

func (o *Api) WithAdminUser(ctx context.Context) context.Context {
	return mctx.WithAdminUser(ctx, o.ChatAdminUserID)
}

func (o *Api) GetClientIP(c *gin.Context) (string, error) {
	if o.ProxyHeader == "" {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		return ip, err
	}
	ip := c.Request.Header.Get(o.ProxyHeader)
	if ip == "" {
		return "", errs.ErrInternalServer.Wrap()
	}
	if ip := net.ParseIP(ip); ip == nil {
		return "", errs.ErrInternalServer.WrapMsg(fmt.Sprintf("parse proxy ip header %s failed", ip))
	}
	return ip, nil
}

func (o *Api) GetDefaultIMAdminUserID() string {
	return o.ImUserID
}
