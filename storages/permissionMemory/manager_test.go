package permissionMemory

import (
	"github.com/stretchr/testify/require"
	"swagger/storages/permissionModel"
	"testing"
)

func TestPermissionMemoryManagerNew(t *testing.T) {
	type Item struct {
		name permissionModel.EntityName
		sign permissionModel.EntitySign
	}
	actual := manager{
		dictId:   managerDictId{},
		dictName: managerDictName{},
		dictSign: managerDictSign{}}
	row1 := permissionModel.NewEntityDefault(1, "one", "Один")
	row2 := permissionModel.NewEntityDefault(2, "two", "Два")
	row3 := permissionModel.NewEntityDefault(3, "three", "Три")
	expected := manager{
		dictId:   managerDictId{row1.ID(): &row1, row2.ID(): &row2, row3.ID(): &row3},
		dictName: managerDictName{row1.Name(): &row1, row2.Name(): &row2, row3.Name(): &row3},
		dictSign: managerDictSign{row1.Sign(): &row1, row2.Sign(): &row2, row3.Sign(): &row3},
		sequence: 3}
	items := []Item{{name: "one", sign: "Один"}, {name: "two", sign: "Два"}, {name: "three", sign: "Три"}}
	for _, item := range items {
		_, err := actual.New(item.name, item.sign)
		require.NoError(t, err)
	}
	require.Equal(t, expected, actual)
}
