package tests

import (
	"context"
	"fmt"
	"net"
	"radius-server/src/common/logger"
	"testing"
	"time"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

func RunTestAuth(t *testing.T) {
	logger.Logger.Info().Msg("Running TestAuth")
	// clean database before starting test
	err := cleanUserAndUserTypes()
	if err != nil {
		t.Error(err)
	}

	err = createUsers()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err := cleanUserAndUserTypes()
		if err != nil {
			t.Error(err)
		}
	}()

	t.Run("SUCCESS_ACCESS_ACCEPT", SuccessFullAccessAccept)
}

func SuccessFullAccessAccept(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	packet := radius.New(radius.CodeAccessRequest, []byte(Nas.Secret))
	rfc2865.UserName_SetString(packet, Users[0].Username)
	rfc2865.UserPassword_SetString(packet, Users[0].Password)
	rfc2865.NASIPAddress_Set(packet, net.IP(Nas.IpAddress))
	fmt.Println("Here is ip address ", Nas.IpAddress)
	resp, err := radius.Exchange(ctx, packet, AccessServerAddress)
	if err != nil {
		t.Error(err)
	}
	if resp == nil {
		t.Fail()
		return
	}

	if resp.Code != radius.CodeAccessAccept {
		t.Fail()
	}
}
