package endpoints

import (
	"context"
	"errors"
	"strings"

	"github.com/jiuzhou-zhao/go-fundamental/loge"
	"github.com/sbasestarter/post/internal/config"
	"github.com/sbasestarter/post/pkg/sms/tencent"
)

type tencentSMSEndPoint struct {
	appID  string
	appKey string
}

func NewTencentSMSEndPoint(ctx context.Context, endPoint config.EndPoint) EndPoint {
	argv := strings.Split(endPoint.Argument, ",")
	if len(argv) != 2 {
		loge.Errorf(ctx, "invalid argument: [%v] %v", endPoint.Name, endPoint.Argument)
		return nil
	}
	return &tencentSMSEndPoint{
		appID:  argv[0],
		appKey: argv[1],
	}
}

func (endPoint *tencentSMSEndPoint) Send(ctx context.Context, to []string, templateID string, vars []string) error {
	if len(to) != 1 {
		return errors.New("invalid to args")
	}
	return tencent.SendSMS(ctx, endPoint.appID, endPoint.appKey, to[0], templateID, vars)
}
