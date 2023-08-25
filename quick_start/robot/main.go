package main

import (
	"fmt"
	"net/http"

	oapi_sdk_go_demo "oapi-sdk-go-demo"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/core/httpserverext"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
)

func main() {
	// 创建告警群并拉人入群
	chatId, err := CreateAlertChat()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("chatId: " + chatId)

	// 发送告警通知
	err = SendAlertMessage(chatId)
	if err != nil {
		fmt.Println(err)
		return
	}

	//  注册事件回调
	eventHandler := dispatcher.NewEventDispatcher(oapi_sdk_go_demo.VerificationToken, oapi_sdk_go_demo.EncryptKey)
	eventHandler.OnP2MessageReceiveV1(DoP2ImMessageReceiveV1)

	// 注册卡片回调
	cardHandler := larkcard.NewCardActionHandler(oapi_sdk_go_demo.VerificationToken, oapi_sdk_go_demo.EncryptKey, DoInteractiveCard)

	http.HandleFunc("/event", httpserverext.NewEventHandlerFunc(eventHandler,
		larkevent.WithLogLevel(larkcore.LogLevelDebug)))
	http.HandleFunc("/card", httpserverext.NewCardActionHandlerFunc(cardHandler,
		larkevent.WithLogLevel(larkcore.LogLevelDebug)))

	err = http.ListenAndServe(":7777", nil)
	if err != nil {
		panic(err)
	}
}
