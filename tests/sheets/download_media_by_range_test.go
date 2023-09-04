package sheets

import (
	"fmt"
	"testing"

	oapi_sdk_go_demo "oapi-sdk-go-demo"
	"oapi-sdk-go-demo/composite_api/sheets"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

func TestDownloadMediaByRange(t *testing.T) {
	req := &sheets.DownloadMediaByRangeRequest{
		SpreadsheetToken: "T90VsUqrYhrnGCtBKS3cLCgQnih",
		Range:            "53988e!A1:A7",
	}

	resp, err := sheets.DownloadMediaByRange(oapi_sdk_go_demo.Client, req)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(larkcore.Prettify(resp))
}
