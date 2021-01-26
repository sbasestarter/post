package endpoints

import "context"

type EndPoint interface {
	Send(ctx context.Context, to []string, templateID string, vars []string) error
}
