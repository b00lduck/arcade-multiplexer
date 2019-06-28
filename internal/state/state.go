package state

import (
	"github.com/b00lduck/arcade-multiplexer/internal/cores"
	"github.com/b00lduck/arcade-multiplexer/internal/games"
)

type UserInterface interface {
	SelectedGame(string)
	CurrentGame(string)
}

type state struct {
	cores    map[string]cores.Core
	selected uint8
	current  games.Game
	ui       UserInterface
}

func NewState(ui UserInterface) *state {
	return &state{
		ui:       ui,
		selected: 0}
}

func (s *state) Up() {
	s.selected++
	if s.selected >= s.numPatches() {
		s.selected = 0
	}
	s.ui.SelectedGame(games.Games[s.selected].Name)
}

func (s *state) Down() {
	if s.selected == 0 {
		s.selected = s.numPatches() - 1
	} else {
		s.selected--
	}
	s.ui.SelectedGame(games.Games[s.selected].Name)
}

func (s *state) Choose() {
	oldCore := s.current.Core
	s.current = games.Games[s.selected]
	newCore := s.current.Core
	if oldCore != newCore {
		cores.ChangeCore(oldCore, newCore)
	}
	cores.LoadGame(newCore, games.Games[s.selected].PrgIndex)
	s.ui.CurrentGame(games.Games[s.selected].Name)
}

func (s *state) numPatches() uint8 {
	return uint8(len(games.Games))
}
