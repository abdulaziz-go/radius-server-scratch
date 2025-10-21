package handlers

import (
	"errors"
	"fmt"
	"layeh.com/radius"
	"radius-server/src/common/logger"
	"radius-server/src/database"
	"radius-server/src/database/entities"
	"radius-server/src/redis"
	"strconv"
)

func getSessionIDAndSubscriberIDFromNAS(r *radius.Request, nasIP string) (string, string, error) {
	nas, err := redis.GetNASByIP(nasIP)
	if err == nil {
		return extractValuesFromNAS(r, nas)
	}

	nas, err = database.GetNasByIp(nasIP)
	if err != nil {
		return "", "", err
	}
	if nas == nil {
		return "", "", errors.New("NAS not found")
	}

	if err = redis.HSetNasClient(nas); err != nil {
		logger.Logger.Error().Msgf("error while writing redis. %s", err.Error())
	}

	return extractValuesFromNAS(r, nas)
}

func extractValuesFromNAS(r *radius.Request, nas *entities.RadiusNas) (string, string, error) {
	sessionID, err := extractSessionIDFromNAS(r, nas)
	if err != nil {
		return "", "", err
	}

	subscriberID, err := extractSubscriberIDFromNAS(r, nas)
	if err != nil {
		return "", "", err
	}

	return sessionID, subscriberID, nil
}

func extractSessionIDFromNAS(r *radius.Request, nas *entities.RadiusNas) (string, error) {
	if nas.SessionId == nil || *nas.SessionId == "" {
		return "", fmt.Errorf("session ID AVP not configured")
	}

	avpNumber, err := strconv.Atoi(*nas.SessionId)
	if err != nil {
		return "", fmt.Errorf("invalid session AVP number: %s", *nas.SessionId)
	}

	value := extractValueByAVPNumber(r, avpNumber)
	if value == "" {
		return "", fmt.Errorf("session ID not found in packet for AVP %d", avpNumber)
	}

	return value, nil
}

func extractSubscriberIDFromNAS(r *radius.Request, nas *entities.RadiusNas) (string, error) {
	if nas.SubscriberId == nil || *nas.SubscriberId == "" {
		return "", fmt.Errorf("subscriber ID AVP not configured")
	}

	avpNumber, err := strconv.Atoi(*nas.SubscriberId)
	if err != nil {
		return "", fmt.Errorf("invalid subscriber AVP number: %s", *nas.SubscriberId)
	}

	value := extractValueByAVPNumber(r, avpNumber)
	if value == "" {
		return "", fmt.Errorf("subscriber ID not found in packet for AVP %d", avpNumber)
	}

	return value, nil
}

func extractValueByAVPNumber(r *radius.Request, avpNumber int) string {
	attr, ok := r.Packet.Lookup(radius.Type(avpNumber))
	if !ok {
		return ""
	}
	return string(attr)
}
