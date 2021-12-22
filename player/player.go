package player

import (
	"sort"
)

type Player struct {
	name      string
	positions []string
	status    string
	points    float32
}

func OptimizeLineup(availableplayers []Player) (map[string]Player, []Player) {
	lineup := make(map[string]Player)
	availableplayers = OrderPlayers(availableplayers)
	lineup, availableplayers = SetLineup(availableplayers, lineup)
	lineup, availableplayers = SetUtils(availableplayers, lineup)
	return lineup, availableplayers
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

func NewPlayer(name string, positions []string, status string, points float32) *Player {
	p := Player{name: name, positions: positions, status: status, points: points}
	return &p
}

func OrderPlayers(availableplayers []Player) []Player {
	sort.SliceStable(availableplayers, func(i, j int) bool {
		return availableplayers[i].points > availableplayers[j].points
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

func SetLineup(availableplayers []Player, lineup map[string]Player) (map[string]Player, []Player) {
	i := 0
	j := 0
	removed := false

	for j < len(availableplayers) {
		i = 0
		for i < len(availableplayers[j].positions) {
			// fmt.Println("Checking if", availableplayers[j], "can be added to position", availableplayers[j].positions[i])
			if _, ok := lineup[availableplayers[j].positions[i]]; ok {
				// fmt.Println(value, "already at position", availableplayers[j].positions[i])
				if i == len(availableplayers[j].positions)-1 {
					// fmt.Println("No open spots available for", availableplayers[j])
					k := 0
					// fmt.Println("Checking if another player can be moved to make room for", availableplayers[j])
					for k < len(availableplayers[j].positions) {
						if value, ok := lineup[availableplayers[j].positions[k]]; ok {
							// fmt.Println("Checking if", value, "at", availableplayers[j].positions[k], "can be moved")
							l := 0
							for l < len(value.positions) {
								if _, ok := lineup[value.positions[l]]; ok {
									// fmt.Println(value2, "already at position", value.positions[l])
								} else {
									// fmt.Println(value.positions[l], "IS OPEN")
									lineup[value.positions[l]] = value
									// fmt.Println(value, "has been moved from", availableplayers[j].positions[k], "to", value.positions[l])
									lineup[availableplayers[j].positions[k]] = availableplayers[j]
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
				lineup[availableplayers[j].positions[i]] = availableplayers[j]
				i = len(availableplayers[j].positions) //end loop when player is added to lineup
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
