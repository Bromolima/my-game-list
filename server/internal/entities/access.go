package entities

type AccessType string

const (
	ReadAccess   AccessType = "read"
	UpdateAccess AccessType = "update"
	CreateAcess  AccessType = "create"
	DeleteAcess  AccessType = "delete"
)

type Access struct {
	ID         uint       `gorm:"primaryKey;type:smallint"`
	AccessType AccessType `gorm:"column:access_type;type:char(100);not null"`
}
