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

func setLineup(availableplayers []player, lineup map[string]player) map[string]player {
	i := 0
	j := 0

	for j < len(availableplayers) {
		i = 0
		for i < len(availableplayers[j].positions) {
			if value, ok := lineup[availableplayers[j].positions[i]]; ok {
				fmt.Println("Player already at this position: ", availableplayers[j].positions[i])
				fmt.Println("value: ", value)
			} else {
				fmt.Println(availableplayers[j].positions[i], "position is open, adding player to lineup")
				lineup[availableplayers[j].positions[i]] = availableplayers[j]
				i = len(availableplayers[j].positions) //end loop when player is added to lineup
			}
			i++
		}
		j++
	}
	return lineup
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
