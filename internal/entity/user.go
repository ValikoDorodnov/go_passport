package entity

type User struct {
	CommonId string `db:"common_id"`
	Roles    string `db:"roles"`
}
