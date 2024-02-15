package userRequestEntity

import (
	"encoding/json"
	"net/http"

	"github.com/gtkmk/finder_api/infra/EmptyBodyVerification"
	"github.com/gtkmk/finder_api/infra/requestEntityFieldsValidation"
)

type ChangeUserPermission struct {
	Id                string `json:"id"`
	PermissionGroupID string `json:"permission_group_id"`
}

const (
	UserIdFieldConst            = "o usuário"
	PermissionGroupIdFieldConst = "o grupo de permissão"
)

func ChangeUserPermissionRequest(req *http.Request) (*ChangeUserPermission, error) {
	if err := EmptyBodyVerification.ValidateBody(req); err != nil {
		return nil, err
	}

	var request *ChangeUserPermission

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func (changeUserPermission *ChangeUserPermission) Validate() error {
	if err := requestEntityFieldsValidation.IsValidUUID(UserIdFieldConst, changeUserPermission.Id); err != nil {
		return err
	}

	return requestEntityFieldsValidation.IsValidUUID(
		PermissionGroupIdFieldConst,
		changeUserPermission.PermissionGroupID,
	)
}
