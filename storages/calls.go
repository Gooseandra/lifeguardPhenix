package storages

import "time"

type (
	Call interface {
		ID() CallID
		Description() CallDescription
		SummingUp() *CallSummingUp
		Address() CallAddress
		TimeStart() CallTime
		TimeFinish() *CallTime
		Title() CallTitle
		Crew() CallCrew
	}

	CallCrew = uint64

	CallTitle = string

	CallID = uint64

	CallDescription = string

	CallSummingUp = string

	CallAddress = string

	CallTime = time.Time

	Calls interface {
		ByID(id CallID) (Call, error)
		New(description CallDescription, address CallAddress, time CallTime, title CallTitle, crew CallCrew) (Call, error)
		List(skip uint64, count uint32) ([]Call, error)
		Update(description CallDescription, address CallAddress, time *CallTime, timeFinish CallTime,
			summingUp CallSummingUp, title CallTitle, crew CallCrew, id CallID) (Call, error)
	}
)
