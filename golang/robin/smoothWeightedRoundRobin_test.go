package robin

import (
	"testing"
)

func TestSmoothWeightedRoundRobin_GetPeer(t *testing.T) {
	servers := []*SmoothServer{
		{weight: 5, host: "server1"},
		{weight: 1, host: "server2"},
		{weight: 1, host: "server3"},
	}

	swr := NewSmoothWeightedRoundRobin(servers)

	expectedHosts := []string{"server1", "server1", "server2", "server1", "server3", "server1", "server1", "server1", "server1", "server2"}

	for i := 0; i < len(expectedHosts); i++ {
		server := swr.getPeer()
		if server == nil {
			t.Error("Expected a server, but got nil.")
		} else if server.host != expectedHosts[i] {
			t.Errorf("Expected host: %s, but got: %s", expectedHosts[i], server.host)
		}
	}
}

func TestSmoothWeightedRoundRobin_AdjustEffectiveWeight(t *testing.T) {
	servers := []*SmoothServer{
		{weight: 3, host: "server1"},
		{weight: 2, host: "server2"},
		{weight: 5, host: "server3"},
	}

	swr := NewSmoothWeightedRoundRobin(servers)

	// Adjust the effective weight of "server1" by -1
	swr.adjustEffectiveWeight("server1", -1)

	// The order of servers after adjustment should be: server3, server2, server1
	expectedHosts := []string{
		"server3", "server1", "server2", "server3", "server3", "server1", "server3", "server2", "server3",
		"server3", "server1", "server2", "server3", "server3", "server1", "server3", "server2", "server3",
	}

	for i := 0; i < len(expectedHosts); i++ {
		server := swr.getPeer()
		if server == nil {
			t.Error("Expected a server, but got nil.")
		} else if server.host != expectedHosts[i] {
			t.Errorf("Expected host: %s, but got: %s", expectedHosts[i], server.host)
		}
	}
}
