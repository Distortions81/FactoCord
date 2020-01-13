package admin

import (
	"../../glob"
	"../../support"
	"github.com/bwmarrin/discordgo"
	"io"
)

// StopServer saves and stops the server.
func StopServer(s *discordgo.Session, m *discordgo.MessageCreate) {

	glob.Refresh = true
	_, err := s.ChannelMessageSend(support.Config.FactorioChannelID, "Server shutting down.")
	if err != nil {
		support.ErrorLog(err)
	}
	if glob.Pipe != nil {
		_, err = io.WriteString(glob.Pipe, "/quit\n")
		if err != nil {
			support.ErrorLog(err)
		}
	}
	glob.Shutdown = true
}
