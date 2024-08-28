package internal

import (
	"fmt"
	"strings"
	"time"
)

type Person struct {
	Name    string
	Aliases []string
}

func (p Person) String() string {
	aliasString := ""
	if len(p.Aliases) > 0 {
		aliasString = " (" + strings.Join(p.Aliases, ", ") + ")"
	}
	return fmt.Sprintf("  %s %s", p.Name, aliasString)
}

type Note struct {
	Name      string    `json:"name"`
	Note      string    `json:"note"`
	Timestamp time.Time `json:"timestamp"`
}

func (n Note) String() string {
	return n.Name + "\t" + n.Timestamp.Format("Mon 01/02/2006 03:04:05 PM MST") + "\n" + n.Note
}
