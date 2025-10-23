package handlers

import (
	"errors"
	"net"
	"radius-server/src/common/logger"
	"radius-server/src/config"
	"radius-server/src/database/entities"
	"radius-server/src/redis"
	"time"

	"layeh.com/radius/rfc3162"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/rfc2866"
)

func AccountingHandler(w radius.ResponseWriter, r *radius.Request) {
	statusType := rfc2866.AcctStatusType_Get(r.Packet)
	username := rfc2865.UserName_GetString(r.Packet)
	nasIP := r.RemoteAddr.(*net.UDPAddr).IP.String()
	var framedIPStr, ipVersion string
	framedIP := rfc2865.FramedIPAddress_Get(r.Packet)
	if framedIP != nil {
		framedIPStr = framedIP.String()
		ipVersion = config.IPv4
	} else {
		ipv6Net := rfc3162.FramedIPv6Prefix_Get(r.Packet)
		if ipv6Net != nil && ipv6Net.IP != nil {
			framedIPStr = ipv6Net.IP.String()
			ipVersion = config.IPv6
		}
	}

	sessionID, subscriberID, err := getSessionIDAndSubscriberIDFromNAS(r, nasIP)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("Failed to get session/subscriber ID from NAS")
		return
	}

	logger.Logger.Info().
		Str("session_id", sessionID).
		Str("subscriber_id", subscriberID).
		Str("username", username).
		Str("nas_ip", nasIP).
		Str("framed_ip", framedIPStr).
		Str("status_type", statusType.String()).
		Msg("Processing accounting request")

	switch statusType {
	case rfc2866.AcctStatusType_Value_Start:
		if err := handleAccountingStart(sessionID, subscriberID, username, nasIP, framedIPStr, ipVersion); err != nil {
			logger.Logger.Error().Err(err).Msg("Failed to handle accounting start")
			return
		}

	case rfc2866.AcctStatusType_Value_Stop:
		if err := handleAccountingStop(subscriberID, framedIPStr); err != nil {
			logger.Logger.Error().Err(err).Msg("Failed to handle accounting stop")
			return
		}

	case rfc2866.AcctStatusType_Value_InterimUpdate:
		if err := handleAccountingInterimUpdate(sessionID, subscriberID, username, nasIP, framedIPStr, ipVersion); err != nil {
			logger.Logger.Error().Err(err).Msg("Failed to handle accounting interim update")
			return
		}

	default:
		logger.Logger.Warn().
			Str("status_type", statusType.String()).
			Msg("Unknown accounting status type")
		return
	}

	w.Write(r.Response(radius.CodeAccountingResponse))
}

func handleAccountingStart(sessionID, subscriberID, username, nasIP, framedIP, ipVersion string) error {
	currentTime := time.Now().Unix()

	subscribers, _ := redis.GetSubscriberBySubscriberID(subscriberID)
	if len(subscribers) == 0 {
		if framedIP != "" {
			if err := redis.DeleteSubscriberByIP(framedIP); err != nil {
				logger.Logger.Warn().Err(err).Str("ip", framedIP).Msg("Failed to delete existing subscriber by IP")
				return err
			}
		} else {
			return errors.New("ip address shouldn't be null")
		}
	} else {
		logger.Logger.Info().
			Str("subscriber_id", subscriberID).
			Msg("Subscriber already exists in Redis")
	}

	logger.Logger.Info().
		Int("existing_count", len(subscribers)).
		Str("subscriber_id", subscriberID).
		Msg("Number of existing subscribers found")

	subscriber := &entities.SubscriberData{
		SubscriberID:    subscriberID,
		IP:              framedIP,
		IpVersion:       ipVersion,
		SessionID:       sessionID,
		LastUpdatedTime: currentTime,
	}

	return redis.CreateOrUpdateSubscriber(subscriber)
}

func handleAccountingStop(subscriberId, ip string) error {
	subscribers, _ := redis.GetSubscriberBySubscriberID(subscriberId)
	if len(subscribers) == 0 {
		return nil
	}

	for _, sub := range subscribers {
		if sub.IP == ip {
			return redis.DeleteSubscriberByIP(ip)
		}
	}
	return nil
}

func handleAccountingInterimUpdate(sessionID, subscriberID, username, nasIP, framedIP, ipVersion string) error {
	currentTime := time.Now().Unix()

	// Check if subscriber ID exists
	existingSubscribers, err := redis.GetSubscriberBySubscriberID(subscriberID)
	if err != nil {
		// Subscriber ID doesn't exist, create new record
		logger.Logger.Info().Str("subscriber_id", subscriberID).Msg("Subscriber not found, creating new subscriber")

		subscriber := &entities.SubscriberData{
			SubscriberID:    subscriberID,
			IP:              framedIP,
			IpVersion:       ipVersion,
			SessionID:       sessionID,
			LastUpdatedTime: currentTime,
		}

		return redis.CreateOrUpdateSubscriber(subscriber)
	}

	// Subscriber ID exists, check if IP address matches
	var matchingSubscriber *entities.SubscriberData
	for _, sub := range existingSubscribers {
		if sub.IP == framedIP {
			matchingSubscriber = sub
			break
		}
	}

	if matchingSubscriber != nil {
		// IP matches, update record timestamp
		logger.Logger.Info().Str("subscriber_id", subscriberID).Str("ip", framedIP).Msg("IP matches, updating timestamp")

		updatedSubscriber := &entities.SubscriberData{
			SubscriberID:    matchingSubscriber.SubscriberID,
			IP:              framedIP,
			IpVersion:       ipVersion,
			SessionID:       sessionID,
			LastUpdatedTime: currentTime,
		}

		return redis.CreateOrUpdateSubscriber(updatedSubscriber)
	} else {
		// IP doesn't match, delete old subscriber with different IP and create new record
		logger.Logger.Info().Str("subscriber_id", subscriberID).Str("new_ip", framedIP).Msg("IP doesn't match, deleting old records and creating new one")

		// Delete all existing records for this subscriber ID
		for _, sub := range existingSubscribers {
			if err = redis.DeleteSubscriberByIP(sub.IP); err != nil {
				logger.Logger.Warn().Err(err).Str("old_ip", sub.IP).Msg("Failed to delete old subscriber record")
				return err
			}
		}

		// Create new record with new IP
		newSubscriber := &entities.SubscriberData{
			SubscriberID:    subscriberID,
			IP:              framedIP,
			IpVersion:       ipVersion,
			SessionID:       sessionID,
			LastUpdatedTime: currentTime,
		}

		return redis.CreateOrUpdateSubscriber(newSubscriber)
	}
}
