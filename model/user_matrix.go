package model

type UserMatrix struct {
	ID         uint  `gorm:"primaryKey"`
	UserID     int64 `gorm:"not null;unique"`
	User       *User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID"`
	IsCreate   *bool `json:"isCreate" validate:"required"`
	IsRead     *bool `json:"isRead" validate:"required"`
	IsUpdate   *bool `json:"isUpdate" validate:"required"`
	IsDelete   *bool `json:"isDelete" validate:"required"`
	IsUpload   *bool `json:"isUpload" validate:"required"`
	IsDownload *bool `json:"isDownload" validate:"required"`
	IsArchive  *bool `json:"isArchive" validate:"required"`
}

type UserMatrixRequest struct {
	IsCreate   *bool `json:"isCreate" validate:"required"`
	IsRead     *bool `json:"isRead" validate:"required"`
	IsUpdate   *bool `json:"isUpdate" validate:"required"`
	IsDelete   *bool `json:"isDelete" validate:"required"`
	IsUpload   *bool `json:"isUpload" validate:"required"`
	IsDownload *bool `json:"isDownload" validate:"required"`
	IsArchive  *bool `json:"isArchive" validate:"required"`
}

type PermissionActionEnum string

const (
	Create   PermissionActionEnum = "C"
	Read     PermissionActionEnum = "R"
	Update   PermissionActionEnum = "U"
	Delete   PermissionActionEnum = "D"
	Upload   PermissionActionEnum = "A"
	Download PermissionActionEnum = "B"
	Archive  PermissionActionEnum = "AV"
)
