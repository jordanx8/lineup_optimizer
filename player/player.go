package player

import (
	"fmt"
	"sort"
)

type player struct {
	name      string
	positions []string
	team      string
	points    int
}

//	lineup := make(map[string]player)
//	lineup["PG"] = *p
//	lineup["SG"] = *p2
//	fmt.Println(lineup)

func newPlayer(name string, positions []string, team string, points int) *player {
	p := player{name: name, positions: positions, team: team, points: points}
	return &p
}

func orderPlayers(availableplayers []player) []player {
	sort.SliceStable(availableplayers, func(i, j int) bool {
		return availableplayers[i].points > availableplayers[j].points
	})
	return availableplayers
}

func remove(slice []player, s int) []player {
	return append(slice[:s], slice[s+1:]...)
}

func setUtils(availableplayers []player, lineup map[string]player) (map[string]player, []player) {
	lineup["Util"] = availableplayers[0]
	lineup["Util2"] = availableplayers[1]
	availableplayers = remove(availableplayers, 0)
	availableplayers = remove(availableplayers, 0)
	return lineup, availableplayers
}

func setLineup(availableplayers []player, lineup map[string]player) (map[string]player, []player) {
	i := 0
	j := 0
	removed := false

	for j < len(availableplayers) {
		i = 0
		for i < len(availableplayers[j].positions) {
			fmt.Println("Checking if", availableplayers[j], "can be added to position", availableplayers[j].positions[i])
			if value, ok := lineup[availableplayers[j].positions[i]]; ok {
				fmt.Println(value, "already at position", availableplayers[j].positions[i])
				if i == len(availableplayers[j].positions)-1 {
					fmt.Println("No open spots available for", availableplayers[j])
					k := 0
					fmt.Println("Checking if another player can be moved to make room for", availableplayers[j])
					for k < len(availableplayers[j].positions) {
						if value, ok := lineup[availableplayers[j].positions[k]]; ok {
							fmt.Println("Checking if", value, "at", availableplayers[j].positions[k], "can be moved")
							l := 0
							for l < len(value.positions) {
								if value2, ok := lineup[value.positions[l]]; ok {
									fmt.Println(value2, "already at position", value.positions[l])
								} else {
									fmt.Println(value.positions[l], "IS OPEN")
									lineup[value.positions[l]] = value
									fmt.Println(value, "has been moved from", availableplayers[j].positions[k], "to", value.positions[l])
									lineup[availableplayers[j].positions[k]] = availableplayers[j]
									fmt.Println("Adding", availableplayers[j], "to position", availableplayers[j].positions[k])
									availableplayers = remove(availableplayers, j)
									removed = true
								}
								l++
							}
						}
						k++
					}
				}
			} else {
				fmt.Println(availableplayers[j].positions[i], "position is open, adding", availableplayers[j], "to lineup")
				lineup[availableplayers[j].positions[i]] = availableplayers[j]
				i = len(availableplayers[j].positions) //end loop when player is added to lineup
				availableplayers = remove(availableplayers, j)
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

//	availableplayers := []player{}
//	availableplayers = append(availableplayers, *p, *p2, *p3, *p4)
//	availableplayers = orderPlayers(availableplayers)
