package alismsproxyapi

import (
	"fmt"
	"strings"

	"github.com/lfhy/alismsproxyapi/client"
	openapi "github.com/lfhy/alismsproxyapi/src/github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/lfhy/alismsproxyapi/src/github.com/alibabacloud-go/tea/tea"
)

type SMSClinet struct {
	c      *client.Client
	option *SMSClinetOptions
}

type SMSClinetOptions struct {
	// 鉴权AK
	AccessKeyId string
	// 鉴权SK
	AccessKeySecret string
	// 代理链接支持HTTP HTTPS SOCKS5
	ProxyURL string
	// 短信签名
	SignName string
	// 模版码
	TemplateCode string
	// 访问域名 默认:dysmsapi.aliyuncs.com
	Endpoint string
	// 访问协议 默认:HTTP
	Protocol string
}

// 创建阿里云验证码服务客户端
func NewSMSClient(option SMSClinetOptions) (*SMSClinet, error) {
	if option.AccessKeyId == "" || option.AccessKeySecret == "" {
		return nil, fmt.Errorf("Client 'AccessKeyId' or 'AccessKeySecret' is empty")
	}

	if option.SignName == "" || option.TemplateCode == "" {
		return nil, fmt.Errorf("SMSTemplate 'SignName' or 'TemplateCode' is empty")
	}
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: tea.String(option.AccessKeyId),
		// 必填，您的 AccessKey Secret
		AccessKeySecret: tea.String(option.AccessKeySecret),
	}
	//协议
	if option.Protocol == "" {
		config.SetProtocol("http")
	} else {
		config.SetProtocol(option.Protocol)
	}

	if option.ProxyURL != "" {
		if strings.HasPrefix(option.ProxyURL, "http:") {
			config.SetHttpProxy(option.ProxyURL)
		} else if strings.HasPrefix(option.ProxyURL, "https:") {
			config.SetHttpsProxy(option.ProxyURL)
		} else if strings.HasPrefix(option.ProxyURL, "socks5:") {
			config.SetSocks5Proxy(option.ProxyURL).SetSocks5NetWork("tcp")
		}

	}
	// 访问的域名
	if option.Endpoint != "" {
		config.SetEndpoint(option.Endpoint)
	} else {
		config.SetEndpoint("dysmsapi.aliyuncs.com")
	}

	c, err := client.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &SMSClinet{
		c:      c,
		option: &option,
	}, nil
}

// 发送验证码
// 参数1:发送的手机号
// 参数2:发送模板的参数 如：{\"code\":12345\"}
func (c *SMSClinet) SendSMS(phone string, templateParam string) (*client.SendSmsResponse, error) {
	sendSmsRequest := &client.SendSmsRequest{
		PhoneNumbers:  tea.String(phone),
		SignName:      tea.String(c.option.SignName),
		TemplateCode:  tea.String(c.option.TemplateCode),
		TemplateParam: tea.String(templateParam),
	}

	return c.c.SendSms(sendSmsRequest)
}
