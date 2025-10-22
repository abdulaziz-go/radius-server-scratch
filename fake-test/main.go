package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/rfc2866"
)

const (
	username     = "testusername"
	sharedSecret = "secretkey123"
	acctAddr     = "127.0.0.1:1813"
	nasPort      = 0
	sessionID    = "test-session-123"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	fmt.Println("📥 Accounting-Start:")
	if err := sendAccounting(ctx, rfc2866.AcctStatusType_Value_Start); err != nil {
		log.Fatalf("Accounting-Start failed: %v", err)
	}

	fmt.Println("📤 Accounting-Interim-Update:")
	if err := sendAccounting(ctx, rfc2866.AcctStatusType_Value_InterimUpdate); err != nil {
		log.Fatalf("Accounting-Interim-Update failed: %v", err)
	}

	fmt.Println("📤 Accounting-Stop:")
	if err := sendAccounting(ctx, rfc2866.AcctStatusType_Value_Stop); err != nil {
		log.Fatalf("Accounting-Stop failed: %v", err)
	}
}

func sendAccounting(ctx context.Context, status rfc2866.AcctStatusType) error {
	packet := radius.New(radius.CodeAccountingRequest, []byte(sharedSecret))

	packet.Set(radius.Type(1), []byte(username))

	packet.Set(radius.Type(44), []byte(sessionID))

	rfc2865.NASPort_Set(packet, rfc2865.NASPort(nasPort))

	rfc2866.AcctStatusType_Set(packet, status)

	if status == rfc2866.AcctStatusType_Value_Stop {
		rfc2866.AcctSessionTime_Set(packet, 3600)
		rfc2866.AcctInputOctets_Set(packet, 1024000)
		rfc2866.AcctOutputOctets_Set(packet, 2048000)
	}

	resp, err := radius.Exchange(ctx, packet, acctAddr)
	if err != nil {
		return err
	}
	fmt.Printf("Response: %v\n", resp.Code)
	if resp.Code != radius.CodeAccountingResponse {
		return fmt.Errorf("expected Accounting-Response, got %v", resp.Code)
	}
	return nil
}
