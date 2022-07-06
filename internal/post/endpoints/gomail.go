package endpoints

import (
	"context"
	"strconv"
	"strings"

	"github.com/sbasestarter/post/internal/config"
	"github.com/sbasestarter/post/pkg/email/gomail"
	"github.com/sbasestarter/post/pkg/post"
	"github.com/sgostarter/i/l"
)

type goMailEndPoint struct {
	from         string
	fromName     string
	smtpHost     string
	smtpPort     int
	authUsername string
	authPass     string

	logger l.WrapperWithContext
}

func NewGoMailEndPoint(ctx context.Context, endPoint config.EndPoint, logger l.Wrapper) EndPoint {
	if logger == nil {
		logger = l.NewNopLoggerWrapper()
	}

	lc := logger.WithFields(l.StringField(l.ClsKey, "goMailEndPoint")).GetWrapperWithContext()

	// from, fromName, host, port, user, pass
	argv := strings.Split(endPoint.Argument, ",")
	if len(argv) != 6 {
		lc.Errorf(ctx, "invalid argument: [%v] %v", endPoint.Name, endPoint.Argument)

		return nil
	}

	port, err := strconv.ParseInt(argv[3], 10, 32)
	if err != nil {
		lc.Errorf(ctx, "invalid argument: %v [%v] %v", argv[3], endPoint.Name, endPoint.Argument)

		return nil
	}

	return &goMailEndPoint{
		from:         argv[0],
		fromName:     argv[1],
		smtpHost:     argv[2],
		smtpPort:     int(port),
		authUsername: argv[4],
		authPass:     argv[5],
		logger:       lc,
	}
}

func (endPoint *goMailEndPoint) Send(ctx context.Context, to []string, _ string, vars []string) error {
	// vars title subject fromName
	subject, body, fromName, err := post.ParseEmailVars(vars)
	if err != nil {
		endPoint.logger.Errorf(ctx, "parse email vars failed: %v", err)

		return err
	}
	if fromName == "" {
		fromName = endPoint.fromName
	}

	return gomail.SendEmail(endPoint.smtpHost, endPoint.smtpPort, endPoint.authUsername, endPoint.authPass, endPoint.from,
		fromName, to, subject, body, true)
}
