package handlers

import (
	"net"
	"radius-server/src/common/logger"
	"radius-server/src/redis"
	"time"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/rfc2866"
)

func AccountingHandler(w radius.ResponseWriter, r *radius.Request) {
	statusType := rfc2866.AcctStatusType_Get(r.Packet)

	sessionID := rfc2866.AcctSessionID_GetString(r.Packet)
	username := rfc2865.UserName_GetString(r.Packet)
	nasIP := r.RemoteAddr.(*net.UDPAddr).IP.String()
	framedIP := rfc2865.FramedIPAddress_Get(r.Packet)

	var framedIPStr string
	if framedIP != nil {
		framedIPStr = framedIP.String()
	}

	logger.Logger.Info().
		Str("session_id", sessionID).
		Str("username", username).
		Str("nas_ip", nasIP).
		Str("framed_ip", framedIPStr).
		Str("status_type", statusType.String()).
		Msg("Processing accounting request")

	switch statusType {
	case rfc2866.AcctStatusType_Value_Start:
		if err := handleAccountingStart(sessionID, username, nasIP, framedIPStr); err != nil {
			logger.Logger.Error().Err(err).Msg("Failed to handle accounting start")
		}

	case rfc2866.AcctStatusType_Value_Stop:
		if err := handleAccountingStop(sessionID); err != nil {
			logger.Logger.Error().Err(err).Msg("Failed to handle accounting stop")
		}

	case rfc2866.AcctStatusType_Value_InterimUpdate:
		if err := handleAccountingInterimUpdate(sessionID, username, nasIP, framedIPStr); err != nil {
			logger.Logger.Error().Err(err).Msg("Failed to handle accounting interim update")
		}

	default:
		logger.Logger.Warn().
			Str("status_type", statusType.String()).
			Msg("Unknown accounting status type")
	}

	w.Write(r.Response(radius.CodeAccountingResponse))
}

func handleAccountingStart(sessionID, username, nasIP, framedIP string) error {
	currentTime := time.Now().Unix()

	if framedIP != "" {
		if err := redis.DeleteSubscriberByIP(framedIP); err != nil {
			logger.Logger.Warn().Err(err).Str("ip", framedIP).Msg("Failed to delete existing subscriber by IP")
		}
	}

	subscriberID := time.Now().UnixNano()

	subscriber := &redis.SubscriberData{
		SubscriberID:    subscriberID,
		IP:              framedIP,
		SessionID:       sessionID,
		LastUpdatedTime: currentTime,
	}

	return redis.CreateOrUpdateSubscriber(subscriber)
}

func handleAccountingStop(sessionID string) error {
	return redis.DeleteSubscriberBySessionID(sessionID)
}

func handleAccountingInterimUpdate(sessionID, username, nasIP, framedIP string) error {
	currentTime := time.Now().Unix()

	existingSubscriber, err := redis.GetSubscriberBySessionID(sessionID)
	if err != nil {
		logger.Logger.Info().Str("session_id", sessionID).Msg("Subscriber not found, creating new subscriber")

		subscriberID := time.Now().UnixNano()

		subscriber := &redis.SubscriberData{
			SubscriberID:    subscriberID,
			IP:              framedIP,
			SessionID:       sessionID,
			LastUpdatedTime: currentTime,
		}

		return redis.CreateOrUpdateSubscriber(subscriber)
	}

	if existingSubscriber.IP != "" && existingSubscriber.IP != framedIP {
		if err := redis.DeleteSubscriberByIP(existingSubscriber.IP); err != nil {
			logger.Logger.Warn().Err(err).Str("old_ip", existingSubscriber.IP).Msg("Failed to delete old IP mapping")
		}
	}

	updatedSubscriber := &redis.SubscriberData{
		SubscriberID:    existingSubscriber.SubscriberID,
		IP:              framedIP,
		SessionID:       sessionID,
		LastUpdatedTime: currentTime,
	}

	return redis.CreateOrUpdateSubscriber(updatedSubscriber)
}
