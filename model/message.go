package model

import (
	"time"
)

/*
this package is created to avoid passing around anonymous
structs. In dev process due to tight connection of message
with db and chat package, and not creating anonymous structs
on recieve and send would have caused circular imports issue.

anonymous structs are defenitely a solution but i had to
define them in multiple places which leads to harder
maintanablity in larger scale.
*/
type Message struct {
	UserID    uint32    `json:"-"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}
