package radius

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"radius-server/src/config"
	"radius-server/src/database"
	"radius-server/src/radius/handlers"

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
	nas, err := database.GetNasByIp(ip)
	if err != nil {
		return nil, err
	}
	if nas == nil {
		return nil, errors.New("NAS not found")
	}

	return []byte(nas.Secret), nil
}
