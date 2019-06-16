package state

import "fmt"

type patch struct {
	name string
}

type state struct {
	patches []patch
	current uint8
}

func NewState() *state {
	return &state{
		current: 0,
		patches: []patch{
			{
				name: "Turrican"},
			{
				name: "Turrican 2"},
		},
	}
}

func (s *state) Up() {
	s.current++
}

func (s *state) Down() {
	s.current--
}

func (s *state) Choose() {
	fmt.Println("choose " + fmt.Sprintf("%d", s.current))
}
