package userMemory

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"swagger/models"
)

const emailOne, emailTwo = "emailOne", "emailTwo"
const firstOne, firstTwo = "firstOne", "firstTwo"
const passwordOne, passwordTwo = "passwordOne", "passwordTwo"
const phoneOne, phoneTwo = "phoneOne", "phoneTwo"
const iDOne, iDTwo = 1, 2
const lastOne, lastTwo = "lastOne", "lastTwo"
const middleOne, middleTwo = "middleOne", "middleTwo"
const nickZero, nickOne, nickTwo = "nickZero", "nickOne", "nickTwo"
const tgOne, tgTwo = "emailOne", "emailTwo"
const vkOne, vkTwo = "phoneOne", "phoneTwo"

var null = (*time.Time)(nil)

func initTest(result *[2]entity) manager {
	contactsOne := models.NewUserContactsDefault(emailOne, phoneOne, tgOne, vkOne)
	contactsTwo := models.NewUserContactsDefault(emailTwo, phoneTwo, tgTwo, vkTwo)
	nameOne := models.NewUserNameDefault(firstOne, lastOne, middleOne, nickOne)
	nameTwo := models.NewUserNameDefault(firstTwo, lastTwo, middleTwo, nickTwo)
	modelOne := models.NewUserDefault(iDOne, contactsOne, nameOne, passwordOne, null, time.Now())
	modelTwo := models.NewUserDefault(iDTwo, contactsTwo, nameTwo, passwordTwo, null, time.Now())
	result[0] = entity{UserDefault: modelOne}
	result[1] = entity{UserDefault: modelTwo}
	return manager{
		dictById:   managerId{iDOne: &result[0], iDTwo: &result[1]},
		dictByName: managerName{nickOne: &result[0], nickTwo: &result[1]}}
}

func TestStorageByName(t *testing.T) {
	notFoundName := func(t *testing.T) {
		var nme models.UserNameMissingError
		var entities [2]entity
		_, fail := initTest(&entities).ByName(nickZero)
		require.Error(t, fail)
		require.ErrorAs(t, fail, &nme)
		require.Equal(t, nme, models.UserNameMissingError(nickZero))
	}
	success := func(t *testing.T) {
		var entities [2]entity
		actual := initTest(&entities)
		_, fail := actual.ByName(nickOne)
		require.NoError(t, fail)
	}
	t.Run("not found name", notFoundName)
	t.Run("success", success)
}

func TestManagerList(t *testing.T) {
	var entities [2]entity
	a := initTest(&entities)
	r, e := a.List(0, 0, 9)
	require.NoError(t, e)
	require.Len(t, r, 2)
	require.Equal(t, r[0], &entities[0])
	require.Equal(t, r[1], &entities[1])
}

func TestStorageNew(t *testing.T) {
	duplicateId := func(t *testing.T) {
		//var a model.IdExistError
		//var entities [2]entity
		//test := initTest(&entities)
		//test.sequence = 1
		//_, r := test.New()
		//require.ErrorAs(t, r, &a)
		//require.Equal(t, a, model.IdExistError(1))
	}
	duplicateName := func(t *testing.T) {
		//var e storages.UserNameExistError
		//a := manager{dictByName: storageRowsByName{n2: nil}}
		//_, f := a.New(n2, p2)
		//require.ErrorAs(t, f, &e)
		//require.Equal(t, e, storages.UserNameExistError(n2))
	}
	success := func(t *testing.T) {
		//type New struct {
		//	entity   storages.User
		//	name     storages.UserName
		//	password storages.UserPassword
		//}
		//a := manager{dictById: storageRowsById{}, dictByName: storageRowsByName{}}
		//n := []New{
		//	{name: n1, password: "one"},
		//	{name: n2, password: "two"},
		//	{name: n3, password: "three"},
		//	{name: n4, password: "four"}}
		//for i, v := range n {
		//	var err error
		//	n[i].entity, err = a.New(v.name, v.password)
		//	require.NoError(t, err)
		//	require.Equal(t, storages.UserID(i+1), n[i].entity.ID())
		//	require.Equal(t, v.name, n[i].entity.Name())
		//	require.Equal(t, v.password, n[i].entity.Password())
		//}
		//require.Equal(t, storages.User(a.first), n[3].entity)
		//require.NotNil(t, a.first.next)
		//require.Equal(t, storages.User(a.first.next), n[0].entity)
		//require.NotNil(t, a.first.next.next)
		//require.Equal(t, storages.User(a.first.next.next), n[2].entity)
		//require.NotNil(t, a.first.next.next.next)
		//require.Equal(t, storages.User(a.first.next.next.next), n[1].entity)
		//require.Nil(t, a.first.next.next.next.next)
		//require.Equal(t, storages.User(a.dictById[1]), n[0].entity)
		//require.Equal(t, storages.User(a.dictById[2]), n[1].entity)
		//require.Equal(t, storages.User(a.dictById[3]), n[2].entity)
		//require.Equal(t, storages.User(a.dictById[4]), n[3].entity)
		//require.Equal(t, storages.User(a.dictByName[n1]), n[0].entity)
		//require.Equal(t, storages.User(a.dictByName[n2]), n[1].entity)
		//require.Equal(t, storages.User(a.dictByName[n3]), n[2].entity)
		//require.Equal(t, storages.User(a.dictByName[n4]), n[3].entity)
	}
	t.Run("duplicateId", duplicateId)
	t.Run("duplicateName", duplicateName)
	t.Run("success", success)
}
