package db

import "github.com/r-cbb/cbbpoll/pkg"

type DBClient interface {
	AddTeam(team pkg.Team) (id int64, err error)
	GetTeam(id int64) (team pkg.Team, err error)
	GetTeams() (teams []pkg.Team, err error)
}
