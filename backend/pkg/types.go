package pkg

type Team struct {
	ID         int64
	FullName   string
	ShortName  string
	Nickname   string
	Conference string
}

type User struct {
	Nickname string
	IsAdmin  bool
}
