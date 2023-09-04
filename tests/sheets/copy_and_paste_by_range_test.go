package sheets

import (
	"fmt"
	"testing"

	oapi_sdk_go_demo "oapi-sdk-go-demo"
	"oapi-sdk-go-demo/composite_api/sheets"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

func TestCopyAndPasteByRange(t *testing.T) {
	req := &sheets.CopyAndPasteByRangeRequest{
		SpreadsheetToken: "T90VsUqrYhrnGCtBKS3cLCgQnih",
		SrcRange:         "53988e!A1:B5",
		DstRange:         "53988e!C1:D5",
	}

	resp, err := sheets.CopyAndPasteRange(oapi_sdk_go_demo.Client, req)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(larkcore.Prettify(resp))
}
