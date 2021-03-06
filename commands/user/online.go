package user

import (
	"../../fact"
	"github.com/bwmarrin/discordgo"
)

// online executes /p o on the server.
func PlayersOnline(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {

	if fact.IsFactRunning() {

		if fact.GetNumPlayers() == 0 {
			fact.CMS(m.ChannelID, "No players are currently online.")
		} else {
			fact.WriteFact("/online")
			fact.CMS(m.ChannelID, "Players online:")
		}
	} else {
		fact.CMS(m.ChannelID, "Factorio is currently offline.")
	}
}
