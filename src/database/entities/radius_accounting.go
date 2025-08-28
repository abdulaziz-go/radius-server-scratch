package entities

import "time"

const RadiusAccountingTable = "radius_accounting"

type RadiusAccounting struct {
	Id                 int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	NasId              int64      `json:"nas_id" gorm:"index;not null"`
	Username           string     `json:"username" gorm:"type:varchar(64);not null"`
	SessionId          string     `json:"session_id" gorm:"type:varchar(128);not null;index"`
	NasIp              string     `json:"nas_ip" gorm:"type:inet;not null"`
	FramedIp           *string    `json:"framed_ip" gorm:"type:inet"`
	AcctStatusType     string     `json:"acct_status_type" gorm:"type:varchar(32);not null"`
	AcctSessionTime    *int       `json:"acct_session_time"`
	AcctInputOctets    int64      `json:"acct_input_octets" gorm:"default:0"`
	AcctOutputOctets   int64      `json:"acct_output_octets" gorm:"default:0"`
	AcctTerminateCause *string    `json:"acct_terminate_cause" gorm:"type:varchar(64)"`
	StartTime          *time.Time `json:"start_time"`
	StopTime           *time.Time `json:"stop_time"`
	CreatedAt          time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

func (RadiusAccounting) TableName() string {
	return RadiusAccountingTable
}
