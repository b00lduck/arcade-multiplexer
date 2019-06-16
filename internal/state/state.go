package state

type UserInterface interface {
	SelectedGame(string)
	CurrentGame(string)
}

type patch struct {
	name string
}

type state struct {
	patches []patch
	current uint8
	ui      UserInterface
}

func NewState(u UserInterface) *state {
	return &state{
		ui:      u,
		current: 0,
		patches: []patch{
			{
				name: "Turrican"},
			{
				name: "Turrican 2"},
			{
				name: "Lotus II"},
			{
				name: "Marble Madness"}}}
}

func (s *state) Up() {
	s.current++
	if s.current >= s.numPatches() {
		s.current = 0
	}
	s.ui.SelectedGame(s.patches[s.current].name)
}

func (s *state) Down() {
	if s.current == 0 {
		s.current = s.numPatches() - 1
	} else {
		s.current--
	}
	s.ui.SelectedGame(s.patches[s.current].name)
}

func (s *state) Choose() {
	s.ui.CurrentGame(s.patches[s.current].name)
}

func (s *state) numPatches() uint8 {
	return uint8(len(s.patches))
}
