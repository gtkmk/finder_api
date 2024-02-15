package permissions

import (
	"github.com/gtkmk/finder_api/core/port"
)

type Permissions struct{}

func NewCheckUserPermissions() port.CheckUserPermissionsInterface {
	return &Permissions{}
}

func (p *Permissions) CheckPermissions(permissions []*port.Permission, resource string, permissionType ...string) bool {

	return p.validePermissionByType(permissions, resource, permissionType...)
}

func (p *Permissions) validePermissionByType(permissions []*port.Permission, resource string, permissionType ...string) bool {
	permissionsTransformed := p.transformPermissions(permissions)

	_, ok := permissionsTransformed[resource]

	if !ok {
		return false
	}

	if len(permissionType) > 1 {
		return permissionsTransformed[resource].OP == permissionType[0] || permissionsTransformed[resource].OP == permissionType[1]
	}

	return permissionsTransformed[resource].OP == permissionType[0]
}

func (p *Permissions) transformPermissions(permissions []*port.Permission) map[string]*port.Permission {
	permissionsReturn := make(map[string]*port.Permission)

	for _, permission := range permissions {
		permissionsReturn[permission.RN] = permission
	}

	return permissionsReturn
}
