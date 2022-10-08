package entity

type User struct {
	CommonId int    `db:"common_id"`
	Roles    string `db:"roles"`
}

func NewUser(commonId int, roles string) *User {
	return &User{CommonId: commonId, Roles: roles}
}
