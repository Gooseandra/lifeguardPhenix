package services

import (
	"fmt"
	"log"
)

type (
	Log     struct{}
	LogFunc string
)

func NewLog() *Log { return &Log{} }

func (l Log) Func(n string) LogFunc { return LogFunc(n) }

func (f LogFunc) BadRequest(m string, a ...any) { f.print("BadRequest", m, a) }

func (f LogFunc) InternalSerer(m string, a ...any) {
	log.Println(string(f) + " InternalSerer: " + fmt.Sprintf(m, a...))
}

func (f LogFunc) NotFound(msg string, args ...any) {
	log.Println(string(f) + " Not Found: " + fmt.Sprintf(msg, args...))
}

func (f LogFunc) OK(msg string, args ...any) {
	log.Println(string(f) + " OK: " + fmt.Sprintf(msg, args...))
}

func (f LogFunc) print(s, m string, a ...any) {
	log.Println(string(f) + " " + s + ": " + fmt.Sprintf(m, a...))
}

func (f LogFunc) Unauthorized(msg string, args ...any) {
	log.Println(string(f) + " Unauthorized: " + fmt.Sprintf(msg, args...))
}
