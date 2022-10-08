package entity

type User struct {
	CommonId int    `db:"common_id"`
	Roles    string `db:"roles"`
}
