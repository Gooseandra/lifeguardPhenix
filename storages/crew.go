package storages

import "time"

type (
	Crew interface {
		ID() CrewID
		Start() CrewTime
		Finish() CrewTime
		Leader() CrewLeader
		Comment() CrewComment
		Roaster() CrewRoster
		Calls() CrewCalls
	}

	CrewID = uint64

	CrewTime = time.Time

	CrewLeader = uint64

	CrewComment = string

	CrewCalls = []uint64

	CrewRoster = []uint64

	Crews interface {
		ByID(id CrewID) (Crew, error)
		New(start CrewTime, leader CrewLeader, comment CrewComment, roster CrewRoster) (Crew, error)
		List(skip uint64, count uint32) ([]Crew, error)
		Update(id CrewID, start CrewTime, finish CrewTime, leader CrewLeader, comment CrewComment,
			roster CrewRoster) (Crew, error)
	}
)
