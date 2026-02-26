package model

type User struct {
	Id       int
	Nickname string
	Password string
	Role     Role
}

func NewUser(nickname, password string, role Role) *User {
	return &User{
		Nickname: nickname,
		Password: password,
		Role:     role,
	}
}
