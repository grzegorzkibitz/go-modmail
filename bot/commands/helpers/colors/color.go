package colors

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"strconv"
)

// ColorName is a custom type for color names
type ColorName string

const (
	Default           ColorName = "DEFAULT"
	Aqua              ColorName = "AQUA"
	DarkAqua          ColorName = "DARK_AQUA"
	Green             ColorName = "GREEN"
	DarkGreen         ColorName = "DARK_GREEN"
	Blue              ColorName = "BLUE"
	DarkBlue          ColorName = "DARK_BLUE"
	Purple            ColorName = "PURPLE"
	DarkPurple        ColorName = "DARK_PURPLE"
	LuminousVividPink ColorName = "LUMINOUS_VIVID_PINK"
	DarkVividPink     ColorName = "DARK_VIVID_PINK"
	Gold              ColorName = "GOLD"
	DarkGold          ColorName = "DARK_GOLD"
	Orange            ColorName = "ORANGE"
	DarkOrange        ColorName = "DARK_ORANGE"
	Red               ColorName = "RED"
	DarkRed           ColorName = "DARK_RED"
	Grey              ColorName = "GREY"
	DarkGrey          ColorName = "DARK_GREY"
	DarkerGrey        ColorName = "DARKER_GREY"
	LightGrey         ColorName = "LIGHT_GREY"
	Navy              ColorName = "NAVY"
	DarkNavy          ColorName = "DARK_NAVY"
	Yellow            ColorName = "YELLOW"
)

var Colors = map[ColorName]discord.Color{
	Default:           0,
	Aqua:              1752220,
	DarkAqua:          1146986,
	Green:             5763719,
	DarkGreen:         2067276,
	Blue:              3447003,
	DarkBlue:          2123412,
	Purple:            10181046,
	DarkPurple:        7419530,
	LuminousVividPink: 15277667,
	DarkVividPink:     11342935,
	Gold:              15844367,
	DarkGold:          12745742,
	Orange:            15105570,
	DarkOrange:        11027200,
	Red:               15548997,
	DarkRed:           10038562,
	Grey:              9807270,
	DarkGrey:          9936031,
	DarkerGrey:        8359053,
	LightGrey:         12370112,
	Navy:              3426654,
	DarkNavy:          2899536,
	Yellow:            16776960,
}

// ToDiscordColor converts a hexadecimal string to a Color type.
func ToDiscordColor(hexColor string) discord.Color {
	colorInt, err := strconv.ParseInt(hexColor, 16, 32)
	if err != nil {
		return discord.DefaultEmbedColor
	}
	return discord.Color(colorInt)
}

func GetColor(colour ColorName) discord.Color {
	if c, ok := Colors[colour]; ok {
		return c
	}
	return discord.DefaultEmbedColor
}
