package slack

import (
	"math/rand"
	"time"

	"github.com/fatih/color"
)

var (
	colorList = []string{
		"red",
		"green",
		"yellow",
		"blue",
		"magenta",
		"cyan",
		"white",
	}
)

func setHostColor(host *Host) error {
	rand.Seed(time.Now().Unix())
	host.Color = colorList[rand.Intn(len(colorList))]

	return nil
}

// ColorizedHost colorized hostname
func ColorizedHost(host *Host) string {
	var colorStr string
	switch host.Color {
	case colorList[0]:
		colorStr = color.RedString(host.MicroAppName)
	case colorList[1]:
		colorStr = color.GreenString(host.MicroAppName)
	case colorList[2]:
		colorStr = color.YellowString(host.MicroAppName)
	case colorList[3]:
		colorStr = color.BlueString(host.MicroAppName)
	case colorList[4]:
		colorStr = color.MagentaString(host.MicroAppName)
	case colorList[5]:
		colorStr = color.CyanString(host.MicroAppName)
	case colorList[6]:
		colorStr = color.WhiteString(host.MicroAppName)
	}
	return colorStr
}
