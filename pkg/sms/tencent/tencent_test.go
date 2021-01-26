package tencent

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	err := SendSMS(context.Background(), os.Getenv("UTSmsAppID"), os.Getenv("UTSmsAppKey"), "+8618710019781", "557914", []string{"测试测试测试测试测试测试", "10"})
	assert.Nil(t, err)
}
