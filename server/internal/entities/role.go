package entities

const (
	RoleUserID  = 1
	RoleAdminID = 2

	RoleUserName  = "USER"
	RoleAdminName = "ADMIN"
)

type Role struct {
	ID     uint     `gorm:"primaryKey;type:smallint"`
	Name   string   `gorm:"char(50);not null;unique"`
	Access []Access `gorm:"many2many:role_accesses"`
	Users  []User
}
