package pkgs

import (
	"errors"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

// 向手机发送验证码
func SendMsg(tel string, code string) (*dysmsapi.SendSmsResponse, error) {
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", "LTAI5tMH1obRmTtMboBpJqvU", "bEiRZV34a0CZO6s6MGx3yXW232wTpR")
	if err != nil {
		return nil, err
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = tel             //手机号变量值
	request.SignName = "北京水山识流科技"          //签名
	request.TemplateCode = "SMS_297950499" //模板编码
	request.TemplateParam = "{\"code\":\"" + code + "\"}"
	response, err := client.SendSms(request)

	if err != nil {
		return response, err
	}

	if response.Code == "isv.BUSINESS_LIMIT_CONTROL" {
		return response, errors.New("frequency_limit")
	}

	return response, nil
}
