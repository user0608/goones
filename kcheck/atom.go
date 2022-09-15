package kcheck

import "fmt"

type Atom struct {
	Name  string
	Value string
}

func (a Atom) String() string {
	return fmt.Sprintf("Atom:{Name: %s, Value: %s}", a.Name, a.Value)
}
