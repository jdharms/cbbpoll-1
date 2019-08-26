package db

import models "github.com/r-cbb/cbbpoll/pkg"

type DBClient interface {
	AddTeam(team models.Team) (id int64, err error)
	GetTeam(id int64) (team models.Team, err error)
	GetTeams() (teams []models.Team, err error)
	AddUser(user models.User) (name string, err error)
	GetUser(name string) (user models.User, err error)
}
