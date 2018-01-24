package amagi

import (
	"github.com/fatih/color"
)

// FgColorizer logger fg color switch
func FgColorizer(msg, logLevel string) string {
	var colorVal color.Attribute

	switch logLevel {
	case "i":
		colorVal = color.FgGreen
	case "w":
		colorVal = color.FgYellow
	case "e":
		colorVal = color.FgRed
	case "f":
		colorVal = color.FgHiRed
	default:
		colorVal = color.FgHiBlack
	}

	return color.New(colorVal).SprintFunc()(msg)
}
