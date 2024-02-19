package sessions

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"swagger/services"
	"swagger/storages"
	"swagger/storages/userMemory"
	"testing"
	"time"
)

func TestSessionsNew(t *testing.T) {
	const oneName, twoName, threeName = "oneUser", "twoName", "threeName"
	const onePassword, twoPassword, threePassword = "onePassword", "twoPassword", "threePassword"
	one := func(t *testing.T) {
		var fail error
		type iteration struct {
			session        *SessionEntity
			name, password string
			user           storages.User
		}
		duration, log, users := time.Second, services.NewLog(), services.NewUsers(userMemory.NewStorage())
		actual := Sessions{dict: map[uuid.UUID]*SessionEntity{}, duration: duration, log: log, users: users}
		go actual.routine()
		iterations := []iteration{
			{name: oneName, password: onePassword},
			{name: twoName, password: twoPassword},
			{name: threeName, password: threePassword}}
		for index, item := range iterations {
			iterations[index].user, fail = users.New(item.name, item.password)
			require.NoError(t, fail)
			iterations[index].session, fail = actual.New(iterations[index].user)
			require.NoError(t, fail)
		}
		dict := map[uuid.UUID]*SessionEntity{
			iterations[0].session.ID(): iterations[0].session,
			iterations[1].session.ID(): iterations[1].session,
			iterations[2].session.ID(): iterations[2].session}
		require.Equal(t, duration, actual.duration)
		require.Equal(t, dict, actual.dict)
		require.Equal(t, iterations[0].session, actual.first)
		require.Equal(t, iterations[2].session, actual.last)
		time.Sleep(5 * duration)
	}
	t.Run("one", one)
}
