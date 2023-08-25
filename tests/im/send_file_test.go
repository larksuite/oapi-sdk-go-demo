package im

import (
	"fmt"
	"os"
	"testing"

	oapi_sdk_go_demo "oapi-sdk-go-demo"
	"oapi-sdk-go-demo/composite_api/im"

	"github.com/larksuite/oapi-sdk-go/v3/core"
)

func TestSendFile(t *testing.T) {
	file, err := os.Open("/Users/bytedance/Desktop/demo.pdf")
	if err != nil {
		t.Error(err)
		return
	}

	req := &im.SendFileRequest{
		FileType:      "pdf",
		FileName:      "demo.pdf",
		File:          file,
		ReceiveIdType: "open_id",
		ReceiveId:     "ou_a79a0f82add14976e3943f4deb17c3fa",
	}

	resp, err := im.SendFile(oapi_sdk_go_demo.Client, req)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(larkcore.Prettify(resp))
}
