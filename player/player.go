package player

type player struct {
	name      string
	positions []string
	team      string
	points    int
}

func newPlayer(name string, positions []string, team string, points int) *player {
	p := player{name: name, positions: positions, team: team, points: points}
	return &p
}
