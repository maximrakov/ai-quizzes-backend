package model

import "errors"

type Role string

const (
	Student Role = "student"
	Mentor  Role = "mentor"
)

var InvalidRole = errors.New("invalid role")

func NewRole(value string) (Role, error) {
	role := Role(value)

	switch role {
	case Student, Mentor:
		return role, nil
	default:
		return "", InvalidRole
	}
}
