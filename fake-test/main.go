package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/rfc2866"
	"layeh.com/radius/rfc3162"
)

const (
	sharedSecret = "secretkey123"
	acctAddr     = "127.0.0.1:1813"
	numSessions  = 50 // Number of concurrent sessions to test
	numUpdates   = 5  // Number of interim updates per session
)

type TestSession struct {
	Username     string
	SessionID    string
	FramedIP     string
	NASPort      uint32
	SubscriberID string
}

type TestResult struct {
	SessionID string
	Success   bool
	Error     error
	Duration  time.Duration
}

func main() {
	// Check if user wants to run individual tests
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "start":
			testAccountingStart()
			return
		case "interim":
			testAccountingInterim()
			return
		case "stop":
			testAccountingStop()
			return
		}
	}

	fmt.Println("🚀 Starting comprehensive RADIUS accounting test...")
	fmt.Printf("📊 Testing with %d concurrent sessions, %d updates each\n", numSessions, numUpdates)

	// Generate test data
	sessions := generateTestSessions(numSessions)

	// Run tests
	start := time.Now()
	results := runConcurrentTests(sessions)
	totalDuration := time.Since(start)

	// Analyze results
	analyzeResults(results, totalDuration)
}

func generateTestSessions(count int) []TestSession {
	sessions := make([]TestSession, count)
	rand.Seed(time.Now().UnixNano())

	// IPv4 and IPv6 test IPs
	ipv4Pool := []string{
		"192.168.1.10", "192.168.1.11", "192.168.1.12", "192.168.1.13", "192.168.1.14",
		"10.0.0.10", "10.0.0.11", "10.0.0.12", "10.0.0.13", "10.0.0.14",
		"172.16.0.10", "172.16.0.11", "172.16.0.12", "172.16.0.13", "172.16.0.14",
	}

	ipv6Pool := []string{
		"2001:db8::1", "2001:db8::2", "2001:db8::3", "2001:db8::4", "2001:db8::5",
		"fe80::1", "fe80::2", "fe80::3", "fe80::4", "fe80::5",
		"2001:db8:85a3::1", "2001:db8:85a3::2", "2001:db8:85a3::3",
	}

	allIPs := append(ipv4Pool, ipv6Pool...)

	for i := 0; i < count; i++ {
		sessions[i] = TestSession{
			Username:     fmt.Sprintf("testuser_%03d", i+1),
			SessionID:    fmt.Sprintf("session_%d_%d", time.Now().Unix(), i),
			FramedIP:     allIPs[i%len(allIPs)],
			NASPort:      uint32(i + 1000),
			SubscriberID: fmt.Sprintf("sub_%d", i+1),
		}
	}

	return sessions
}

func runConcurrentTests(sessions []TestSession) []TestResult {
	var wg sync.WaitGroup
	resultsChan := make(chan TestResult, len(sessions))

	// Start concurrent session tests
	for _, session := range sessions {
		wg.Add(1)
		go func(s TestSession) {
			defer wg.Done()
			result := testSessionLifecycle(s)
			resultsChan <- result
		}(session)

		// Small delay to avoid overwhelming the server
		time.Sleep(10 * time.Millisecond)
	}

	// Wait for all tests to complete
	wg.Wait()
	close(resultsChan)

	// Collect results
	var results []TestResult
	for result := range resultsChan {
		results = append(results, result)
	}

	return results
}

func testSessionLifecycle(session TestSession) TestResult {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Printf("🔄 Testing session %s (IP: %s)\n", session.SessionID, session.FramedIP)

	// 1. Accounting Start
	if err := sendAccountingRequest(ctx, session, rfc2866.AcctStatusType_Value_Start); err != nil {
		return TestResult{
			SessionID: session.SessionID,
			Success:   false,
			Error:     fmt.Errorf("start failed: %w", err),
			Duration:  time.Since(start),
		}
	}

	// 2. Multiple Interim Updates
	for i := 0; i < numUpdates; i++ {
		time.Sleep(100 * time.Millisecond) // Simulate time between updates
		if err := sendAccountingRequest(ctx, session, rfc2866.AcctStatusType_Value_InterimUpdate); err != nil {
			return TestResult{
				SessionID: session.SessionID,
				Success:   false,
				Error:     fmt.Errorf("interim update %d failed: %w", i+1, err),
				Duration:  time.Since(start),
			}
		}
	}

	// 3. Accounting Stop
	if err := sendAccountingRequest(ctx, session, rfc2866.AcctStatusType_Value_Stop); err != nil {
		return TestResult{
			SessionID: session.SessionID,
			Success:   false,
			Error:     fmt.Errorf("stop failed: %w", err),
			Duration:  time.Since(start),
		}
	}

	return TestResult{
		SessionID: session.SessionID,
		Success:   true,
		Error:     nil,
		Duration:  time.Since(start),
	}
}

func sendAccountingRequest(ctx context.Context, session TestSession, status rfc2866.AcctStatusType) error {
	packet := radius.New(radius.CodeAccountingRequest, []byte(sharedSecret))

	// Set basic attributes
	rfc2865.UserName_SetString(packet, session.Username)
	packet.Set(radius.Type(44), []byte(session.SessionID)) // Acct-Session-Id
	rfc2865.NASPort_Set(packet, rfc2865.NASPort(session.NASPort))
	rfc2866.AcctStatusType_Set(packet, status)

	// Set Framed-IP-Address using proper RFC functions
	if session.FramedIP != "" {
		if ip := net.ParseIP(session.FramedIP); ip != nil {
			if ipv4 := ip.To4(); ipv4 != nil {
				// IPv4 - use RFC 2865 Framed-IP-Address
				rfc2865.FramedIPAddress_Set(packet, ipv4)
			} else {
				// IPv6 - use RFC 3162 Framed-IPv6-Prefix
				// Create a /128 prefix for the IPv6 address
				ipv6Net := &net.IPNet{
					IP:   ip,
					Mask: net.CIDRMask(128, 128), // /128 means single host
				}
				rfc3162.FramedIPv6Prefix_Set(packet, ipv6Net)
			}
		}
	}

	// Add session-specific data for Stop requests
	if status == rfc2866.AcctStatusType_Value_Stop {
		// Random session time (1-7200 seconds)
		sessionTime := rand.Uint32()%7200 + 1
		rfc2866.AcctSessionTime_Set(packet, rfc2866.AcctSessionTime(sessionTime))

		// Random data usage
		inputOctets := rand.Uint32()%10000000 + 1000
		outputOctets := rand.Uint32()%20000000 + 2000
		rfc2866.AcctInputOctets_Set(packet, rfc2866.AcctInputOctets(inputOctets))
		rfc2866.AcctOutputOctets_Set(packet, rfc2866.AcctOutputOctets(outputOctets))
	}

	// Send request
	resp, err := radius.Exchange(ctx, packet, acctAddr)
	if err != nil {
		return fmt.Errorf("exchange failed: %w", err)
	}

	if resp.Code != radius.CodeAccountingResponse {
		return fmt.Errorf("expected Accounting-Response, got %v", resp.Code)
	}

	return nil
}

func analyzeResults(results []TestResult, totalDuration time.Duration) {
	fmt.Println("\n📈 Test Results Analysis:")
	fmt.Println("========================")

	successCount := 0
	failureCount := 0
	totalSessionDuration := time.Duration(0)
	var failures []TestResult

	for _, result := range results {
		if result.Success {
			successCount++
			totalSessionDuration += result.Duration
		} else {
			failureCount++
			failures = append(failures, result)
		}
	}

	fmt.Printf("✅ Successful sessions: %d/%d (%.1f%%)\n",
		successCount, len(results), float64(successCount)/float64(len(results))*100)
	fmt.Printf("❌ Failed sessions: %d/%d (%.1f%%)\n",
		failureCount, len(results), float64(failureCount)/float64(len(results))*100)
	fmt.Printf("⏱️  Total test duration: %v\n", totalDuration)

	if successCount > 0 {
		avgSessionDuration := totalSessionDuration / time.Duration(successCount)
		fmt.Printf("📊 Average session duration: %v\n", avgSessionDuration)
	}

	fmt.Printf("🔢 Total requests sent: %d\n", len(results)*(1+numUpdates+1)) // Start + Updates + Stop
	fmt.Printf("📈 Requests per second: %.1f\n", float64(len(results)*(1+numUpdates+1))/totalDuration.Seconds())

	// Show failures
	if len(failures) > 0 {
		fmt.Println("\n❌ Failed Sessions:")
		for _, failure := range failures {
			fmt.Printf("  - %s: %v\n", failure.SessionID, failure.Error)
		}
	}

	// Test summary
	fmt.Println("\n🎯 Test Summary:")
	if successCount == len(results) {
		fmt.Println("🎉 ALL TESTS PASSED! The accounting service is working correctly with multiple concurrent sessions.")
	} else if successCount > len(results)/2 {
		fmt.Printf("⚠️  PARTIAL SUCCESS: %d/%d sessions completed successfully. Check failed sessions above.\n", successCount, len(results))
	} else {
		fmt.Printf("🚨 MAJOR ISSUES: Only %d/%d sessions completed successfully. Service needs investigation.\n", successCount, len(results))
	}
}

// Individual test functions for manual testing

// testAccountingStart tests only the Accounting-Start functionality
func testAccountingStart() {
	fmt.Println("🚀 Testing Accounting-Start only...")
	rand.Seed(time.Now().UnixNano())

	// Use fixed timestamp for consistent session IDs
	fixedTimestamp := int64(1729600000) // Fixed timestamp for consistent session IDs

	// Test with both IPv4 and IPv6
	testIPs := []string{
		"192.168.1.100", // IPv4
		"2001:db8::100", // IPv6
	}

	for i, ip := range testIPs {
		session := TestSession{
			Username:     fmt.Sprintf("start_test_%d", i+1),
			SessionID:    fmt.Sprintf("start_session_%d_%d", fixedTimestamp, i),
			FramedIP:     ip,
			NASPort:      uint32(2000 + i),
			SubscriberID: fmt.Sprintf("start_sub_%d", i+1),
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		fmt.Printf("📝 Testing Accounting-Start for %s (IP: %s)\n", session.Username, session.FramedIP)

		start := time.Now()
		err := sendAccountingRequest(ctx, session, rfc2866.AcctStatusType_Value_Start)
		duration := time.Since(start)

		if err != nil {
			fmt.Printf("❌ FAILED: %v (duration: %v)\n", err, duration)
		} else {
			fmt.Printf("✅ SUCCESS: Accounting-Start sent successfully (duration: %v)\n", duration)
		}
	}
}

// testAccountingInterim tests only the Accounting-Interim-Update functionality
func testAccountingInterim() {
	fmt.Println("🚀 Testing Accounting-Interim-Update only...")
	rand.Seed(time.Now().UnixNano())

	// Test with both IPv4 and IPv6
	testIPs := []string{
		"192.168.1.101", // IPv4
		"2001:db8::101", // IPv6
	}

	for i, ip := range testIPs {
		session := TestSession{
			Username:     fmt.Sprintf("interim_test_%d", i+1),
			SessionID:    fmt.Sprintf("interim_session_%d_%d", time.Now().Unix(), i),
			FramedIP:     ip,
			NASPort:      uint32(3000 + i),
			SubscriberID: fmt.Sprintf("interim_sub_%d", i+1),
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		fmt.Printf("📝 Testing Accounting-Interim-Update for %s (IP: %s)\n", session.Username, session.FramedIP)

		start := time.Now()
		err := sendAccountingRequest(ctx, session, rfc2866.AcctStatusType_Value_InterimUpdate)
		duration := time.Since(start)

		if err != nil {
			fmt.Printf("❌ FAILED: %v (duration: %v)\n", err, duration)
		} else {
			fmt.Printf("✅ SUCCESS: Accounting-Interim-Update sent successfully (duration: %v)\n", duration)
		}
	}
}

// testAccountingStop tests only the Accounting-Stop functionality
// Uses the exact same session data as start test for proper lifecycle testing
func testAccountingStop() {
	fmt.Println("🚀 Testing Accounting-Stop only...")
	rand.Seed(time.Now().UnixNano())

	// Use exact same timestamp as start test for identical session IDs
	fixedTimestamp := int64(1729600000) // Same fixed timestamp as start test

	// Use same IPs and session data as start test
	testIPs := []string{
		"192.168.1.100", // IPv4 - same as start test
		"2001:db8::100", // IPv6 - same as start test
	}

	for i, ip := range testIPs {
		session := TestSession{
			Username:     fmt.Sprintf("start_test_%d", i+1), // Same username as start test
			SessionID:    fmt.Sprintf("start_session_%d_%d", fixedTimestamp, i), // Exact same session ID as start test
			FramedIP:     ip,
			NASPort:      uint32(2000 + i), // Same NAS port as start test
			SubscriberID: fmt.Sprintf("start_sub_%d", i+1), // Same subscriber ID as start test
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		fmt.Printf("📝 Testing Accounting-Stop for %s (IP: %s)\n", session.Username, session.FramedIP)

		start := time.Now()
		err := sendAccountingRequest(ctx, session, rfc2866.AcctStatusType_Value_Stop)
		duration := time.Since(start)

		if err != nil {
			fmt.Printf("❌ FAILED: %v (duration: %v)\n", err, duration)
		} else {
			fmt.Printf("✅ SUCCESS: Accounting-Stop sent successfully (duration: %v)\n", duration)
		}
	}
}
