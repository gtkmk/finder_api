package port

type CheckUserPermissionsInterface interface {
	CheckPermissions(permissions []*Permission, resource string, permissionType ...string) bool
}
