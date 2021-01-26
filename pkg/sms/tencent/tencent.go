package tencent

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/sbasestarter/post/pkg"
)

const (
	urlTemplate = "https://yun.tim.qq.com/v5/tlssmssvr/sendsms?sdkappid=%v&random=%v"
)

type tcTel struct {
	Mobile     string `json:"mobile"`
	NationCode string `json:"nationcode"`
}

type tcResquest struct {
	Ext    string   `json:"ext"`
	Extend string   `json:"extend"`
	Params []string `json:"params"`
	Sig    string   `json:"sig"`
	Sign   string   `json:"sign"`
	Tel    tcTel    `json:"tel"`
	Time   int64    `json:"time"`
	TplID  string   `json:"tpl_id"`
}

type tcResponse struct {
	Result int    `json:"result"`
	ErrMsg string `json:"errmsg"`
	Ext    string `json:"ext"`
	Fee    int    `json:"fee"`
	Sid    string `json:"sid"`
}

func calcSMSSign(appKey, rand, time, mobile string) string {
	s := fmt.Sprintf("appkey=%v&random=%v&time=%v&mobile=%v", appKey, rand, time, mobile)

	h := sha256.New()
	h.Write([]byte(s))
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}

func SendSMS(ctx context.Context, appID, appKey, phone, templateID string, vars []string) error {
	number, _, regionCode, valid := pkg.ParsePhone(phone)
	if !valid {
		return errors.Errorf("parsePhone %v failed", phone)
	}

	req := &tcResquest{}
	req.Params = vars
	req.Tel.Mobile = number
	req.Tel.NationCode = fmt.Sprintf("%v", regionCode)
	req.Time = time.Now().Unix()
	req.TplID = templateID
	randValue := fmt.Sprintf("%v", rand.Int63())
	req.Sig = calcSMSSign(appKey, randValue, fmt.Sprintf("%v", req.Time), req.Tel.Mobile)

	bytesData, err := json.Marshal(req)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(bytesData)
	url := fmt.Sprintf(urlTemplate, appID, randValue)
	request, err := http.NewRequestWithContext(ctx, "POST", url, reader)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var res tcResponse
	err = json.Unmarshal(respBytes, &res)
	if err != nil {
		return err
	}

	if res.Result == 0 {
		return nil
	}
	return errors.Errorf("%v", res.ErrMsg)
}
