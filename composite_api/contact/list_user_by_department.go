/*
 获取部门下所有用户列表，使用到两个OpenAPI：
 1. [获取子部门列表](https://open.feishu.cn/document/server-docs/contact-v3/department/children)
 2. [获取部门直属用户列表](https://open.feishu.cn/document/server-docs/contact-v3/user/find_by_department)
*/

package contact

import (
	"context"
	"fmt"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
)

type ListUserByDepartmentRequest struct {
	DepartmentId string
}

type ListUserByDepartmentResponse struct {
	*larkcore.CodeError
	ChildrenDepartmentResponse   *larkcontact.ChildrenDepartmentRespData
	FindByDepartmentUserResponse []*larkcontact.User
}

// ListUserByDepartment 获取部门下所有用户列表
func ListUserByDepartment(client *lark.Client, request *ListUserByDepartmentRequest) (*ListUserByDepartmentResponse, error) {
	// 获取子部门列表
	childrenDepartmentReq := larkcontact.NewChildrenDepartmentReqBuilder().
		DepartmentIdType("open_department_id").
		DepartmentId(request.DepartmentId).
		Build()

	childrenDepartmentResp, err := client.Contact.Department.Children(context.Background(), childrenDepartmentReq)

	if err != nil {
		return nil, err
	}
	if !childrenDepartmentResp.Success() {
		fmt.Printf("client.Contact.Department.Children failed, code: %d, msg: %s, log_id: %s\n",
			childrenDepartmentResp.Code, childrenDepartmentResp.Msg, childrenDepartmentResp.RequestId())
		return nil, childrenDepartmentResp.CodeError
	}

	// 获取部门直属用户列表
	users := make([]*larkcontact.User, 0)
	openDepartmentIds := []string{request.DepartmentId}
	for _, item := range childrenDepartmentResp.Data.Items {
		openDepartmentIds = append(openDepartmentIds, *item.OpenDepartmentId)
	}

	for _, id := range openDepartmentIds {
		findByDepartmentUserReq := larkcontact.NewFindByDepartmentUserReqBuilder().
			DepartmentId(id).
			Build()

		findByDepartmentUserResp, err := client.Contact.User.FindByDepartment(context.Background(), findByDepartmentUserReq)

		if err != nil {
			return nil, err
		}
		if !findByDepartmentUserResp.Success() {
			fmt.Printf("client.Contact.User.FindByDepartment failed, code: %d, msg: %s, log_id: %s\n",
				findByDepartmentUserResp.Code, findByDepartmentUserResp.Msg, findByDepartmentUserResp.RequestId())
			return nil, findByDepartmentUserResp.CodeError
		}

		users = append(users, findByDepartmentUserResp.Data.Items...)
	}

	// 返回结果
	return &ListUserByDepartmentResponse{
		CodeError: &larkcore.CodeError{
			Code: 0,
			Msg:  "success",
		},
		ChildrenDepartmentResponse:   childrenDepartmentResp.Data,
		FindByDepartmentUserResponse: users,
	}, nil
}
