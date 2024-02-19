package storages

import (
	"github.com/stretchr/testify/require"
	"swagger/storages/permissionMemory"
	"swagger/storages/permissionModel"
	"testing"
)

func TestPermissionNew(t *testing.T) {
	const message = "manager %d, test %d"
	type Test struct {
		entity permissionModel.Entity
		name   permissionModel.EntityName
		sign   permissionModel.EntitySign
	}
	managers := []permissionModel.Manager{permissionMemory.New()}
	tests := []Test{{name: "one", sign: "Один"}, {name: "two", sign: "Два"}, {name: "three", sign: "Три"}}
	for managerIndex, managerValue := range managers {
		for testIndex, testValue := range tests {
			entity, err := managerValue.New(testValue.name, testValue.sign)
			require.NoError(t, err, message, managerIndex, testIndex)
			tests[testIndex].entity = entity
		}
		for testIndex, testValue := range tests {
			expected := permissionModel.NewEntityDefault(permissionModel.EntityID(testIndex+1), testValue.name, testValue.sign)
			require.Equal(t, &expected, testValue.entity, message, managerIndex, testIndex)
		}
	}
}
