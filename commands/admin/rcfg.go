package admin

import (
	"../../cfg"
	"../../fact"
	"../../glob"
	"github.com/bwmarrin/discordgo"
)

//Archive map
func ReloadConfig(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {

	glob.GameMapLock.Lock()
	defer glob.GameMapLock.Unlock()

	//Read global and local configs
	cfg.ReadGCfg()
	cfg.ReadLCfg()

	//Re-Write global and local configs
	cfg.WriteGCfg()
	cfg.WriteLCfg()
	fact.CMS(m.ChannelID, "Config files reloaded.")

	//Config reset-interval
	if cfg.Local.ResetScheduleText != "" {
		fact.WriteFact("/resetint " + cfg.Local.ResetScheduleText)
	}

}
