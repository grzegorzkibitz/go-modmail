package helpers

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"strconv"
)

// Color converts a hexadecimal string to a Color type.
func Color(hexColor string) discord.Color {
	colorInt, err := strconv.ParseInt(hexColor, 16, 32)
	if err == nil {
		return discord.DefaultEmbedColor
	}
	return discord.Color(colorInt)
}
