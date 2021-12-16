package player

type player struct {
	name      string
	positions []string
	team      string
	points    int
}

type lineup struct {
	PG    player
	SG    player
	G     player
	SF    player
	PF    player
	F     player
	C     player
	C2    player
	Util  player
	Util2 player
}

func newPlayer(name string, positions []string, team string, points int) *player {
	p := player{name: name, positions: positions, team: team, points: points}
	return &p
}
