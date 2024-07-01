package constants

var Enter byte = 13
var Escape byte = 27
var Space byte = 32

var Up byte = 65
var Down byte = 66
var Right byte = 67
var Left byte = 68

var MoveDown byte = 106
var MoveUp byte = 107

var ToggleName byte = 108
var Rename byte = 114
var Create byte = 97
var Delete byte = 100
var Help byte = 104
var Keys = map[byte]bool{
	Up:         true,
	Down:       true,
	Right:      true,
	Left:       true,
	MoveDown:   true,
	MoveUp:     true,
	Delete:     true,
	Rename:     true,
	Create:     true,
	ToggleName: true,
	Help:       true,
}
