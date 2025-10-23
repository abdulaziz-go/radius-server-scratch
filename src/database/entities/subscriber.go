package entities

import (
	numberUtil "radius-server/src/utils/number"
)

type SubscriberData struct {
	SubscriberID    string `json:"subscriber_id"`
	IP              string `json:"ip"`
	IpVersion       string `json:"ip_version"`
	SessionID       string `json:"session_id"`
	LastUpdatedTime int64  `json:"last_updated_time"`
}

func (s *SubscriberData) NewSubscriberDataFromFieldMap(fieldMap map[string]string) *SubscriberData {
	if id, ok := fieldMap["subscriber_id"]; ok {
		s.SubscriberID = id
	}
	if ip, ok := fieldMap["ip"]; ok {
		s.IP = ip
	}
	if sessionID, ok := fieldMap["session_id"]; ok {
		s.SessionID = sessionID
	}
	if lastUpdated, ok := fieldMap["last_updated_time"]; ok {
		value, _ := numberUtil.ParseToInt64(lastUpdated)
		s.LastUpdatedTime = value
	}
	if ipVersion, ok := fieldMap["ip_version"]; ok {
		s.IpVersion = ipVersion
	}
	
	return s
}

func (s *SubscriberData) ToFieldsMap() map[string]interface{} {
	fields := map[string]interface{}{
		"subscriber_id":     s.SubscriberID,
		"ip":                s.IP,
		"ip_version":        s.IpVersion,
		"session_id":        s.SessionID,
		"last_updated_time": s.LastUpdatedTime,
	}

	return fields
}
