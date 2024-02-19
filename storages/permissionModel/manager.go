package permissionModel

type Manager interface {
	New(EntityName, EntityName) (Entity, error)
}
