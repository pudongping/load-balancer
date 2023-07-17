package robin

import (
	"testing"
)

// go test -run TestNewWeightedRoundRobin
func TestNewWeightedRoundRobin(t *testing.T) {
	servers := []Server{
		{weight: 3, host: "server1"},
		{weight: 2, host: "server2"},
		{weight: 5, host: "server3"},
	}

	rr := NewWeightedRoundRobin(servers)

	// Check if the returned rr instance is not nil
	if rr == nil {
		t.Error("Expected a non-nil instance of WeightedRoundRobin, but got nil.")
	}

	// Check if the currentWeight and currentIndex are initialized correctly
	if rr.currentWeight != 0 {
		t.Errorf("Expected currentWeight to be 0, but got: %d", rr.currentWeight)
	}

	if rr.currentIndex != -1 {
		t.Errorf("Expected currentIndex to be -1, but got: %d", rr.currentIndex)
	}

	// Check if maxWeight is set correctly
	expectedMaxWeight := 5
	if rr.maxWeight != expectedMaxWeight {
		t.Errorf("Expected maxWeight to be %d, but got: %d", expectedMaxWeight, rr.maxWeight)
	}

	// Check if total is set correctly
	expectedTotal := len(servers)
	if rr.total != expectedTotal {
		t.Errorf("Expected total to be %d, but got: %d", expectedTotal, rr.total)
	}
}

// go test -run TestWeightedRoundRobin_GetPeer
func TestWeightedRoundRobin_GetPeer(t *testing.T) {
	// Test the case when there are no servers
	rr := NewWeightedRoundRobin([]Server{})
	server := rr.GetPeer()
	if server != nil {
		t.Error("Expected GetPeer to return nil when there are no servers, but got a server.")
	}

	// Test the case when there is only one server
	rr = NewWeightedRoundRobin([]Server{{weight: 1, host: "server1"}})
	server = rr.GetPeer()
	if server == nil {
		t.Error("Expected GetPeer to return a server when there is only one server, but got nil.")
	} else if server.host != "server1" {
		t.Errorf("Expected GetPeer to return server1, but got: %s", server.host)
	}

	// Test the case when there are multiple servers
	servers := []Server{
		{weight: 5, host: "server1"},
		{weight: 1, host: "server2"},
		{weight: 1, host: "server3"},
	}

	rr = NewWeightedRoundRobin(servers)

	expectedHosts := []string{"server1", "server1", "server1", "server1", "server1", "server2", "server3"}

	for i := 0; i < len(expectedHosts); i++ {
		server = rr.GetPeer()
		if server == nil {
			t.Error("Expected a server, but got nil.")
		} else if server.host != expectedHosts[i] {
			t.Errorf("Expected host: %s, but got: %s", expectedHosts[i], server.host)
		}
	}
}

// go test -run TestWeightedRoundRobin_calculateGCD
func TestWeightedRoundRobin_calculateGCD(t *testing.T) {
	servers := []Server{
		{weight: 3},
		{weight: 2},
		{weight: 5},
	}

	rr := NewWeightedRoundRobin(servers)

	// Test the case when the GCD is calculated correctly
	expectedGCD := 1
	gcd := rr.calculateGCD()
	if gcd != expectedGCD {
		t.Errorf("Expected GCD to be %d, but got: %d", expectedGCD, gcd)
	}

	// Test the case when there is only one server
	servers = []Server{{weight: 2}}
	rr = NewWeightedRoundRobin(servers)
	expectedGCD = 2
	gcd = rr.calculateGCD()
	if gcd != expectedGCD {
		t.Errorf("Expected GCD to be %d, but got: %d", expectedGCD, gcd)
	}

	// Test the case when all servers have the same weight
	servers = []Server{{weight: 3}, {weight: 3}, {weight: 3}}
	rr = NewWeightedRoundRobin(servers)
	expectedGCD = 3
	gcd = rr.calculateGCD()
	if gcd != expectedGCD {
		t.Errorf("Expected GCD to be %d, but got: %d", expectedGCD, gcd)
	}
}

// go test -run TestWeightedRoundRobin_calculateGCDRecursive
func TestWeightedRoundRobin_calculateGCDRecursive(t *testing.T) {
	rr := NewWeightedRoundRobin([]Server{})

	// Test the case when GCD of 0 and 5 is calculated correctly
	a, b := 0, 5
	expectedGCD := 5
	gcd := rr.calculateGCDRecursive(a, b)
	if gcd != expectedGCD {
		t.Errorf("Expected GCD to be %d, but got: %d", expectedGCD, gcd)
	}

	// Test the case when GCD of 12 and 6 is calculated correctly
	a, b = 12, 6
	expectedGCD = 6
	gcd = rr.calculateGCDRecursive(a, b)
	if gcd != expectedGCD {
		t.Errorf("Expected GCD to be %d, but got: %d", expectedGCD, gcd)
	}

	// Test the case when GCD of 15 and 10 is calculated correctly
	a, b = 15, 10
	expectedGCD = 5
	gcd = rr.calculateGCDRecursive(a, b)
	if gcd != expectedGCD {
		t.Errorf("Expected GCD to be %d, but got: %d", expectedGCD, gcd)
	}
}
