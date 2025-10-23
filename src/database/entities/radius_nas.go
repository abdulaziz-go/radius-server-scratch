package entities

import (
	numberUtil "radius-server/src/utils/number"
)

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

func (nas *RadiusNas) NewRadiusNasFromFieldMap(fieldMap map[string]string) *RadiusNas {
	if id, ok := fieldMap["id"]; ok {
		value, _ := numberUtil.ParseToInt64(id)
		nas.Id = value
	}
	if name, ok := fieldMap["nas_name"]; ok {
		nas.NasName = &name
	}
	if ipAddr, ok := fieldMap["ip_address"]; ok {
		nas.IpAddress = ipAddr
	}
	if secret, ok := fieldMap["secret"]; ok {
		nas.Secret = secret
	}
	if sessionId, ok := fieldMap["session_id"]; ok {
		nas.SessionId = &sessionId
	}
	if subsId, ok := fieldMap["subscriber_id"]; ok {
		nas.SubscriberId = &subsId
	}
	if createdAt, ok := fieldMap["created_at"]; ok {
		value, _ := numberUtil.ParseToInt64(createdAt)
		nas.CreatedAt = value
	}
	if updatedAt, ok := fieldMap["updated_at"]; ok {
		value, _ := numberUtil.ParseToInt64(updatedAt)
		nas.UpdatedAt = value
	}

	return nas
}

func (nas *RadiusNas) ToFieldsMap() map[string]interface{} {
	fields := map[string]interface{}{
		"id":         nas.Id,
		"ip_address": nas.IpAddress,
		"secret":     nas.Secret,
	}

	if nas.NasName != nil {
		fields["nas_name"] = *nas.NasName
	}
	if nas.SubscriberId != nil {
		fields["subscriber_id"] = *nas.SubscriberId
	}
	if nas.SessionId != nil {
		fields["session_id"] = *nas.SessionId
	}
	if nas.CreatedAt != 0 {
		fields["created_at"] = nas.CreatedAt
	}
	if nas.UpdatedAt != 0 {
		fields["updated_at"] = nas.UpdatedAt
	}

	return fields
}
