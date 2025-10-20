package radius

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"radius-server/src/common/logger"
	"radius-server/src/config"
	"radius-server/src/database"
	"radius-server/src/radius/handlers"
	"radius-server/src/redis"

	"layeh.com/radius"
)

type RadiusServer struct {
}

func New() *RadiusServer {
	return &RadiusServer{}
}

func (rs *RadiusServer) Start() error {
	log.Println("Starting RADIUS server...")
	secretSource := &SecretSource{}
	// Here you would add the actual RADIUS server initialization and start logic
	errChan := make(chan error)

	go func() {
		accessSrv := radius.PacketServer{
			Addr:         fmt.Sprintf(":%d", config.AppConfig.RadiusServer.AccessHandlerServerPort),
			Handler:      radius.HandlerFunc(handlers.AccessHandler),
			SecretSource: secretSource,
		}
		log.Printf("Access server running on :%d", config.AppConfig.RadiusServer.AccessHandlerServerPort)
		if err := accessSrv.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	err := <-errChan

	return err
}

type SecretSource struct {
}

func (s *SecretSource) RADIUSSecret(ctx context.Context, addr net.Addr) ([]byte, error) {
	udpAddr, ok := addr.(*net.UDPAddr)
	if !ok {
		return nil, fmt.Errorf("invalid addr type: %T", addr)
	}

	ip := udpAddr.IP.String()

	nas, err := redis.GetNASByIP(ip)
	if err == nil {
		return []byte(nas.Secret), err
	}

	nas, err = database.GetNasByIp(ip)
	if err != nil {
		return nil, err
	}
	if nas == nil {
		return nil, errors.New("NAS not found")
	}

	if err = redis.HSetNasClient(map[string]interface{}{
		"id":         nas.Id,
		"nas_name":   nas.NasName,
		"ip_address": nas.IpAddress,
		"secret":     nas.Secret,
	}); err != nil {
		logger.Logger.Error().Msgf("error while writing redis. %s", err.Error())
	}

	return []byte(nas.Secret), nil
}
