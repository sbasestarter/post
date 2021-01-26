package gomail

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendEmail(t *testing.T) {
	err := SendEmail("smtp.163.com", 465, "nobodyinword@163.com", os.Getenv("UTEmailPass"),
		"nobodyinword@163.com", "YmiPro", []string{"stwstw0123@163.com"},
		"测试", "测试，只是一个测试", true)
	assert.Nil(t, err)
}
