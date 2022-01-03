package player

import (
	"sort"
)

type Player struct {
	Name          string   `json:"name"`
	Positions     []string `json:"positions"`
	Status        string   `json:"status"`
	Points        float32  `json:"points"`
	FinalPosition string   `json:"finalposition"`
}

func OptimizeLineup(availableplayers []Player) ([]Player, []Player) {
	lineup := make(map[string]Player)
	availableplayers = OrderPlayers(availableplayers)
	lineup, availableplayers = SetLineup(availableplayers, lineup)
	lineup, availableplayers = SetUtils(availableplayers, lineup)
	lineup, availableplayers = SetIL(availableplayers, lineup)
	lineup = SetBN(availableplayers, lineup)
	var lineupstruct []Player
	var benchilstruct []Player
	lineupstruct = append(lineupstruct, *NewPlayer(lineup["PG"].Name, lineup["PG"].Positions, lineup["PG"].Status, lineup["PG"].Points, "PG"))
	lineupstruct = append(lineupstruct, *NewPlayer(lineup["SG"].Name, lineup["SG"].Positions, lineup["SG"].Status, lineup["SG"].Points, "SG"))
	lineupstruct = append(lineupstruct, *NewPlayer(lineup["G"].Name, lineup["G"].Positions, lineup["G"].Status, lineup["G"].Points, "G"))
	lineupstruct = append(lineupstruct, *NewPlayer(lineup["SF"].Name, lineup["SF"].Positions, lineup["SF"].Status, lineup["SF"].Points, "SF"))
	lineupstruct = append(lineupstruct, *NewPlayer(lineup["PF"].Name, lineup["PF"].Positions, lineup["PF"].Status, lineup["PF"].Points, "PF"))
	lineupstruct = append(lineupstruct, *NewPlayer(lineup["F"].Name, lineup["F"].Positions, lineup["F"].Status, lineup["F"].Points, "F"))
	lineupstruct = append(lineupstruct, *NewPlayer(lineup["C"].Name, lineup["C"].Positions, lineup["C"].Status, lineup["C"].Points, "C"))
	lineupstruct = append(lineupstruct, *NewPlayer(lineup["C2"].Name, lineup["C2"].Positions, lineup["C2"].Status, lineup["C2"].Points, "C"))
	lineupstruct = append(lineupstruct, *NewPlayer(lineup["Util"].Name, lineup["Util"].Positions, lineup["Util"].Status, lineup["Util"].Points, "Util"))
	lineupstruct = append(lineupstruct, *NewPlayer(lineup["Util2"].Name, lineup["Util2"].Positions, lineup["Util2"].Status, lineup["Util2"].Points, "Util"))
	benchilstruct = append(benchilstruct, *NewPlayer(lineup["BN"].Name, lineup["BN"].Positions, lineup["BN"].Status, lineup["BN"].Points, "BN"))
	benchilstruct = append(benchilstruct, *NewPlayer(lineup["BN2"].Name, lineup["BN2"].Positions, lineup["BN2"].Status, lineup["BN2"].Points, "BN"))
	benchilstruct = append(benchilstruct, *NewPlayer(lineup["BN3"].Name, lineup["BN3"].Positions, lineup["BN3"].Status, lineup["BN3"].Points, "BN"))
	_, ok := lineup["IL"]
	if ok {
		benchilstruct = append(benchilstruct, *NewPlayer(lineup["IL"].Name, lineup["IL"].Positions, lineup["IL"].Status, lineup["IL"].Points, "IL"))
	}
	_, ok = lineup["IL2"]
	if ok {
		benchilstruct = append(benchilstruct, *NewPlayer(lineup["IL2"].Name, lineup["IL2"].Positions, lineup["IL2"].Status, lineup["IL2"].Points, "IL"))
	}
	_, ok = lineup["IL3"]
	if ok {
		benchilstruct = append(benchilstruct, *NewPlayer(lineup["IL3"].Name, lineup["IL3"].Positions, lineup["IL3"].Status, lineup["IL3"].Points, "IL"))
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

func NewPlayer(name string, positions []string, status string, points float32, finalposition ...string) *Player {
	if len(finalposition) > 0 {
		p := Player{Name: name, Positions: positions, Status: status, Points: points, FinalPosition: finalposition[0]}
		return &p
	} else {
		p := Player{Name: name, Positions: positions, Status: status, Points: points, FinalPosition: ""}
		return &p
	}

}

func OrderPlayers(availableplayers []Player) []Player {
	sort.SliceStable(availableplayers, func(i, j int) bool {
		return availableplayers[i].Points > availableplayers[j].Points
	})
	return availableplayers
}

func Remove(slice []Player, s int) []Player {
	return append(slice[:s], slice[s+1:]...)
}

func SetUtils(availableplayers []Player, lineup map[string]Player) (map[string]Player, []Player) {
	lineup["Util"] = availableplayers[0]
	lineup["Util2"] = availableplayers[1]
	availableplayers = Remove(availableplayers, 0)
	availableplayers = Remove(availableplayers, 0)
	return lineup, availableplayers
}

func SetIL(availableplayers []Player, lineup map[string]Player) (map[string]Player, []Player) {
	a, b, c := 99, 99, 99
	i := 0
	for i < len(availableplayers) {
		if availableplayers[i].Status == "INJ" {
			if _, ok := lineup["IL"]; ok {
				if _, ok := lineup["IL2"]; ok {
					if _, ok := lineup["IL3"]; ok {

					} else {
						lineup["IL3"] = availableplayers[i]
						c = i
					}
				} else {
					lineup["IL2"] = availableplayers[i]
					b = i
				}
			} else {
				lineup["IL"] = availableplayers[i]
				a = i
			}
		}
		i++
	}
	if a != 99 {
		availableplayers = Remove(availableplayers, a)
	}
	if b != 99 {
		availableplayers = Remove(availableplayers, b-1)
	}
	if c != 99 {
		availableplayers = Remove(availableplayers, c-2)
	}
	return lineup, availableplayers
}

func SetBN(availableplayers []Player, lineup map[string]Player) map[string]Player {
	for i, b := range availableplayers {
		if i == 0 {
			lineup["BN"] = b
		} else if i == 1 {
			lineup["BN2"] = b
		} else if i == 2 {
			lineup["BN3"] = b
		}
	}
	return lineup
}

func SetLineup(availableplayers []Player, lineup map[string]Player) (map[string]Player, []Player) {
	i := 0
	j := 0
	removed := false

	for j < len(availableplayers) {
		i = 0
		for i < len(availableplayers[j].Positions) {
			// fmt.Println("Checking if", availableplayers[j], "can be added to position", availableplayers[j].positions[i])
			if _, ok := lineup[availableplayers[j].Positions[i]]; ok {
				// fmt.Println(value, "already at position", availableplayers[j].positions[i])
				if i == len(availableplayers[j].Positions)-1 {
					// fmt.Println("No open spots available for", availableplayers[j])
					k := 0
					// fmt.Println("Checking if another player can be moved to make room for", availableplayers[j])
					for k < len(availableplayers[j].Positions) {
						if value, ok := lineup[availableplayers[j].Positions[k]]; ok {
							// fmt.Println("Checking if", value, "at", availableplayers[j].positions[k], "can be moved")
							l := 0
							for l < len(value.Positions) {
								if _, ok := lineup[value.Positions[l]]; ok {
									// fmt.Println(value2, "already at position", value.positions[l])
								} else {
									// fmt.Println(value.positions[l], "IS OPEN")
									lineup[value.Positions[l]] = value
									// fmt.Println(value, "has been moved from", availableplayers[j].positions[k], "to", value.positions[l])
									lineup[availableplayers[j].Positions[k]] = availableplayers[j]
									// fmt.Println("Adding", availableplayers[j], "to position", availableplayers[j].positions[k])
									availableplayers = Remove(availableplayers, j)
									removed = true
								}
								l++
							}
						}
						k++
					}
				}
			} else {
				// fmt.Println(availableplayers[j].positions[i], "position is open, adding", availableplayers[j], "to lineup")
				lineup[availableplayers[j].Positions[i]] = availableplayers[j]
				i = len(availableplayers[j].Positions) //end loop when player is added to lineup
				availableplayers = Remove(availableplayers, j)
				removed = true
			}
			i++
		}
		if removed {
			removed = false
		} else {
			j++
		}
	}
	return lineup, availableplayers
}
