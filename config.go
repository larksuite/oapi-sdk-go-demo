package oapi_sdk_go_demo

import "github.com/larksuite/oapi-sdk-go/v3"

const (
	AppId             = "cli_xxxx"
	AppSecret         = "xxxx"
	EncryptKey        = "xxxx"
	VerificationToken = "xxxx"
)

var Client = lark.NewClient(AppId, AppSecret)
