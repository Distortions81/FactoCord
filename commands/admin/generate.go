package admin

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"../../cfg"
	"../../fact"
	"../../glob"
	"../../logs"
	"github.com/bwmarrin/discordgo"
)

//Generate map
func Generate(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {

	if fact.IsFactRunning() {
		fact.CMS(m.ChannelID, "Stop server first! ($stop)")
		return
	}

	glob.FactorioLaunchLock.Lock()
	defer glob.FactorioLaunchLock.Unlock()

	t := time.Now()
	ourseed := uint64(t.UnixNano())
	argnum := len(args)

	MapPreset := cfg.Local.MapPreset

	if glob.LastMapSeed > 0 {
		ourseed = glob.LastMapSeed
		glob.LastMapSeed = 0
	}

	if argnum > 0 {
		arg := args[0]
		pnum := 0

		//Remove leading zero
		if arg[0] == '0' {
			pnum, _ = strconv.Atoi(fmt.Sprintf("%c", arg[1]))
		} else {
			pnum, _ = strconv.Atoi(arg[0:1])
		}

		MapPreset = fact.GetMapTypeName(pnum)
		decoded, _ := base64.RawURLEncoding.DecodeString(arg[2:])
		ourseed = binary.BigEndian.Uint64(decoded)
	}

	if ourseed <= 0 {
		fact.CMS(m.ChannelID, "Error, no seed.")
		return
	}

	if MapPreset == "Error" {
		fact.CMS(m.ChannelID, "Invalid map preset.")
		return
	}

	fact.CMS(m.ChannelID, "Generating map...")

	//Generate code to make filename
	buf := new(bytes.Buffer)

	_ = binary.Write(buf, binary.BigEndian, ourseed)
	ourcode := fmt.Sprintf("%02d%v", fact.GetMapTypeNum(MapPreset), base64.RawURLEncoding.EncodeToString(buf.Bytes()))
	filename := cfg.Global.PathData.FactorioServersRoot + cfg.Global.PathData.FactorioHomePrefix + cfg.Local.ServerCallsign + "/" + cfg.Global.PathData.SaveFilePath + "/" + ourcode + ".zip"

	factargs := []string{"--preset", MapPreset, "--map-gen-seed", fmt.Sprintf("%v", ourseed), "--create", filename}

	//Append map gen if set
	if cfg.Local.MapGenPreset != "" {
		factargs = append(factargs, "--map-gen-settings")
		factargs = append(factargs, cfg.Global.PathData.FactorioServersRoot+cfg.Global.PathData.MapGenPath+"/"+cfg.Local.MapGenPreset+"-gen.json")

		factargs = append(factargs, "--map-settings")
		factargs = append(factargs, cfg.Global.PathData.FactorioServersRoot+cfg.Global.PathData.MapGenPath+"/"+cfg.Local.MapGenPreset+"-set.json")
	} else {
		factargs = append(factargs, "--preset", MapPreset)
	}

	cmd := exec.Command(cfg.Global.PathData.FactorioServersRoot+cfg.Global.PathData.FactorioHomePrefix+cfg.Local.ServerCallsign+cfg.Global.PathData.FactorioBinary, factargs...)
	out, aerr := cmd.CombinedOutput()

	if aerr != nil {
		logs.Log(fmt.Sprintf("An error occurred attempting to generate the map. Details: %s", aerr))
	}

	lines := strings.Split(string(out), "\n")

	for _, l := range lines {
		if strings.Contains(l, "Creating new map") {
			fact.CMS(m.ChannelID, "New map saved as: "+filename)
			return
		}
	}

	fact.CMS(m.ChannelID, "Unknown error.")
}
