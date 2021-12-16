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
