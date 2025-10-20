package entities

const RadiusNasTable = "radius_nas"

type RadiusNas struct {
	Id           int64   `json:"id" gorm:"primaryKey;autoIncrement"`
	NasName      *string `json:"nas_name" gorm:"type:varchar(128)"`
	IpAddress    string  `json:"ip_address" gorm:"type:inet;unique;not null;index:idx_radius_nas_ip_address"`
	Secret       string  `json:"secret" gorm:"type:varchar(64);not null"`
	SubscriberId *string `json:"subscriber_id" gorm:"type:varchar(64)"`
	SessionId    *string `json:"session_id" gorm:"type:varchar(128)"`
	CreatedAt    int64   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    int64   `json:"updated_at" gorm:"autoUpdateTime"`
}

func (RadiusNas) TableName() string {
	return RadiusNasTable
}
