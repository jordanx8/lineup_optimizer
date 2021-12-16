package player

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

func orderPlayers(availableplayers []player) []player {
	sort.SliceStable(availableplayers, func(i, j int) bool {
		return availableplayers[i].points > availableplayers[j].points
	})
	return availableplayers
}

//	availableplayers := []player{}
//	availableplayers = append(availableplayers, *p, *p2, *p3, *p4)
//	availableplayers = orderPlayers(availableplayers)
