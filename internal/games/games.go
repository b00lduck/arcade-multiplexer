package games

import "github.com/b00lduck/arcade-multiplexer/internal/cores"

type Game struct {
	Name     string
	Core     *cores.Core
	PrgIndex int
}

var Games = []Game{
	{
		Name:     "Lazy Jones",
		Core:     cores.C64,
		PrgIndex: 7},
	{
		Name:     "Lotus II",
		Core:     cores.Amiga,
		PrgIndex: 17},
	{
		Name:     "Marble Madness",
		Core:     cores.Amiga,
		PrgIndex: 18}}
