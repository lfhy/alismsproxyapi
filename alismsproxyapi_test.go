package alismsproxyapi_test

import (
	"alismsproxyapi"
	"fmt"
	"testing"
)

func TestSendCode(t *testing.T) {
	// 创建Option
	option := alismsproxyapi.SMSClinetOptions{
		AccessKeyId:     "your_access_key_id",
		AccessKeySecret: "your_access_key_secret",
		SignName:        "sign_name",
		TemplateCode:    "template_code",
		ProxyURL:        "socks5://user:pass@socks5.server:2333",
	}

	// 建立Client
	client, err := alismsproxyapi.NewSMSClient(option)
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.SendSMS("12345678901", "{\"code\":12345\"}")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}
