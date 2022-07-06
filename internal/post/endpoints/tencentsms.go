package endpoints

import (
	"context"
	"errors"
	"strings"

	"github.com/sbasestarter/post/internal/config"
	"github.com/sbasestarter/post/pkg/sms/tencent"
	"github.com/sgostarter/i/l"
)

type tencentSMSEndPoint struct {
	appID  string
	appKey string
	logger l.WrapperWithContext
}

func NewTencentSMSEndPoint(ctx context.Context, endPoint config.EndPoint, logger l.Wrapper) EndPoint {
	if logger == nil {
		logger = l.NewNopLoggerWrapper()
	}

	lc := logger.WithFields(l.StringField(l.ClsKey, "tencentSMSEndPoint")).GetWrapperWithContext()

	argv := strings.Split(endPoint.Argument, ",")
	if len(argv) != 2 {
		lc.Errorf(ctx, "invalid argument: [%v] %v", endPoint.Name, endPoint.Argument)

		return nil
	}
	return &tencentSMSEndPoint{
		appID:  argv[0],
		appKey: argv[1],
		logger: lc,
	}
}

func (endPoint *tencentSMSEndPoint) Send(ctx context.Context, to []string, templateID string, vars []string) error {
	if len(to) != 1 {
		return errors.New("invalid to args")
	}

	return tencent.SendSMS(ctx, endPoint.appID, endPoint.appKey, to[0], templateID, vars)
}
