package player

import (
	"sort"
)

type Player struct {
	Name          string   `json:"name"`
	Positions     []string `json:"positions"`
	Status        string   `json:"status"`
	Info          string   `json;"info"`
	Points        float32  `json:"points"`
	FinalPosition string   `json:"finalposition"`
}

var playerGetsMoved bool
var currentlineup = make(map[string]Player)
var playersused []Player

func OptimizeLineup(availableplayers []Player) ([]Player, []Player) {
	availableplayers = orderPlayers(availableplayers)
	setLineup(availableplayers)
	// removes used players from availableplayers array
	for i := range playersused {
		for j := range availableplayers {
			if playersused[i].Name == availableplayers[j].Name {
				availableplayers = remove(availableplayers, j)
				break
			}
		}
	}
	availableplayers = setUtils(availableplayers, currentlineup)
	availableplayers = setIL(availableplayers, currentlineup)
	setBN(availableplayers, currentlineup)
	// create structs used for playertable.html template
	var lineupstruct []Player
	var benchilstruct []Player
	lineupstruct = append(lineupstruct, *NewPlayer(currentlineup["PG"].Name, currentlineup["PG"].Positions, currentlineup["PG"].Status, currentlineup["PG"].Info, currentlineup["PG"].Points, "PG"))
	lineupstruct = append(lineupstruct, *NewPlayer(currentlineup["SG"].Name, currentlineup["SG"].Positions, currentlineup["SG"].Status, currentlineup["SG"].Info, currentlineup["SG"].Points, "SG"))
	lineupstruct = append(lineupstruct, *NewPlayer(currentlineup["G"].Name, currentlineup["G"].Positions, currentlineup["G"].Status, currentlineup["G"].Info, currentlineup["G"].Points, "G"))
	lineupstruct = append(lineupstruct, *NewPlayer(currentlineup["SF"].Name, currentlineup["SF"].Positions, currentlineup["SF"].Status, currentlineup["SF"].Info, currentlineup["SF"].Points, "SF"))
	lineupstruct = append(lineupstruct, *NewPlayer(currentlineup["PF"].Name, currentlineup["PF"].Positions, currentlineup["PF"].Status, currentlineup["PF"].Info, currentlineup["PF"].Points, "PF"))
	lineupstruct = append(lineupstruct, *NewPlayer(currentlineup["F"].Name, currentlineup["F"].Positions, currentlineup["F"].Status, currentlineup["F"].Info, currentlineup["F"].Points, "F"))
	lineupstruct = append(lineupstruct, *NewPlayer(currentlineup["C"].Name, currentlineup["C"].Positions, currentlineup["C"].Status, currentlineup["C"].Info, currentlineup["C"].Points, "C"))
	lineupstruct = append(lineupstruct, *NewPlayer(currentlineup["C2"].Name, currentlineup["C2"].Positions, currentlineup["C2"].Status, currentlineup["C2"].Info, currentlineup["C2"].Points, "C"))
	lineupstruct = append(lineupstruct, *NewPlayer(currentlineup["Util"].Name, currentlineup["Util"].Positions, currentlineup["Util"].Status, currentlineup["Util"].Info, currentlineup["Util"].Points, "Util"))
	lineupstruct = append(lineupstruct, *NewPlayer(currentlineup["Util2"].Name, currentlineup["Util2"].Positions, currentlineup["Util2"].Status, currentlineup["Util2"].Info, currentlineup["Util2"].Points, "Util"))
	_, ok := currentlineup["BN"]
	if ok {
		benchilstruct = append(benchilstruct, *NewPlayer(currentlineup["BN"].Name, currentlineup["BN"].Positions, currentlineup["BN"].Status, currentlineup["BN"].Info, currentlineup["BN"].Points, "BN"))
	}
	_, ok = currentlineup["BN2"]
	if ok {
		benchilstruct = append(benchilstruct, *NewPlayer(currentlineup["BN2"].Name, currentlineup["BN2"].Positions, currentlineup["BN2"].Status, currentlineup["BN2"].Info, currentlineup["BN2"].Points, "BN"))
	}
	_, ok = currentlineup["BN3"]
	if ok {
		benchilstruct = append(benchilstruct, *NewPlayer(currentlineup["BN3"].Name, currentlineup["BN3"].Positions, currentlineup["BN3"].Status, currentlineup["BN3"].Info, currentlineup["BN3"].Points, "BN"))
	}
	_, ok = currentlineup["IL"]
	if ok {
		benchilstruct = append(benchilstruct, *NewPlayer(currentlineup["IL"].Name, currentlineup["IL"].Positions, currentlineup["IL"].Status, currentlineup["IL"].Info, currentlineup["IL"].Points, "IL"))
	}
	_, ok = currentlineup["IL2"]
	if ok {
		benchilstruct = append(benchilstruct, *NewPlayer(currentlineup["IL2"].Name, currentlineup["IL2"].Positions, currentlineup["IL2"].Status, currentlineup["IL2"].Info, currentlineup["IL2"].Points, "IL"))
	}
	_, ok = currentlineup["IL3"]
	if ok {
		benchilstruct = append(benchilstruct, *NewPlayer(currentlineup["IL3"].Name, currentlineup["IL3"].Positions, currentlineup["IL3"].Status, currentlineup["IL3"].Info, currentlineup["IL3"].Points, "IL"))
	}
	return lineupstruct, benchilstruct
}

func AddExtraPositions(positions []string) []string {
	alreadyhasF := false
	alreadyhasG := false
	for _, b := range positions {
		if !alreadyhasG && b == "PG" {
			positions = append(positions, "G")
			alreadyhasG = true
		} else if !alreadyhasG && b == "SG" {
			positions = append(positions, "G")
			alreadyhasG = true
		} else if !alreadyhasF && b == "SF" {
			positions = append(positions, "F")
			alreadyhasF = true
		} else if !alreadyhasF && b == "PF" {
			positions = append(positions, "F")
			alreadyhasF = true
		} else if b == "C" {
			positions = append(positions, "C2")
			alreadyhasF = true
		}
	}
	return positions
}

func NewPlayer(name string, positions []string, status string, info string, points float32, finalposition ...string) *Player {
	if len(finalposition) > 0 {
		p := Player{Name: name, Positions: positions, Status: status, Info: info, Points: points, FinalPosition: finalposition[0]}
		return &p
	} else {
		p := Player{Name: name, Positions: positions, Status: status, Info: info, Points: points, FinalPosition: ""}
		return &p
	}

}

func orderPlayers(availableplayers []Player) []Player {
	sort.SliceStable(availableplayers, func(i, j int) bool {
		return availableplayers[i].Points > availableplayers[j].Points
	})
	return availableplayers
}

func remove(slice []Player, s int) []Player {
	return append(slice[:s], slice[s+1:]...)
}

func setUtils(availableplayers []Player, currentlineup map[string]Player) []Player {
	currentlineup["Util"] = availableplayers[0]
	currentlineup["Util2"] = availableplayers[1]
	availableplayers = remove(availableplayers, 0)
	availableplayers = remove(availableplayers, 0)
	return availableplayers
}

func setIL(availableplayers []Player, currentlineup map[string]Player) []Player {
	a, b, c := 99, 99, 99
	i := 0
	for i < len(availableplayers) {
		if availableplayers[i].Status == "INJ" {
			if _, ok := currentlineup["IL"]; ok {
				if _, ok := currentlineup["IL2"]; ok {
					if _, ok := currentlineup["IL3"]; ok {

					} else {
						currentlineup["IL3"] = availableplayers[i]
						c = i
					}
				} else {
					currentlineup["IL2"] = availableplayers[i]
					b = i
				}
			} else {
				currentlineup["IL"] = availableplayers[i]
				a = i
			}
		}
		i++
	}
	if a != 99 {
		availableplayers = remove(availableplayers, a)
	}
	if b != 99 {
		availableplayers = remove(availableplayers, b-1)
	}
	if c != 99 {
		availableplayers = remove(availableplayers, c-2)
	}
	return availableplayers
}

func setBN(availableplayers []Player, currentlineup map[string]Player) map[string]Player {
	for i, b := range availableplayers {
		if i == 0 {
			currentlineup["BN"] = b
		} else if i == 1 {
			currentlineup["BN2"] = b
		} else if i == 2 {
			currentlineup["BN3"] = b
		}
	}
	return currentlineup
}

func setLineup(availableplayers []Player) {
	for i := 0; (len(currentlineup) < 9) && (i < len(availableplayers)); i++ {
		playerGetsMoved = false
		setLineupRecur(availableplayers[i], 0)
		if playerGetsMoved {
			playersused = append(playersused, availableplayers[i])
		}
	}
}

func setLineupRecur(player Player, depth int) {
	//stop recursion at depth 3
	if depth == 3 {
		return
	}
	//checks to see if player can be added to an open position
	for j := 0; j < len(player.Positions); j++ {
		if _, ok := currentlineup[player.Positions[j]]; !ok {
			currentlineup[player.Positions[j]] = player
			playerGetsMoved = true
			return
		}
	}
	//checks if players at positions that currentplayer can play at are able to be moved
	for k := 0; k < len(player.Positions); k++ {
		setLineupRecur(currentlineup[player.Positions[k]], depth+1)
		if playerGetsMoved {
			currentlineup[player.Positions[k]] = player
			return
		}
	}
}
