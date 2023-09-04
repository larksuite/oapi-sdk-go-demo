package contact

import (
	"fmt"
	"testing"

	oapi_sdk_go_demo "oapi-sdk-go-demo"
	"oapi-sdk-go-demo/composite_api/contact"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

func TestListUserByDepartment(t *testing.T) {
	req := &contact.ListUserByDepartmentRequest{
		DepartmentId: "0",
	}

	resp, err := contact.ListUserByDepartment(oapi_sdk_go_demo.Client, req)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(larkcore.Prettify(resp))
}
