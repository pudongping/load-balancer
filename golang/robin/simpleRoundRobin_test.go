package robin

import (
	"fmt"
	"testing"
)

// go test -run TestSimpleRoundRobin_Add
func TestSimpleRoundRobin_Add(t *testing.T) {
	rr := SimpleRoundRobin{}

	// Test case 1: Adding multiple servers
	err := rr.Add("server1", "server2", "server3")
	if err != nil {
		t.Errorf("Error while adding servers: %v", err)
	}

	// Test case 2: Adding an empty server
	err = rr.Add()
	if err == nil {
		t.Error("Expected an error when adding an empty server, but got none.")
	}
}

// go test -run TestSimpleRoundRobin_GetPeer
func TestSimpleRoundRobin_GetPeer(t *testing.T) {
	rr := SimpleRoundRobin{}

	peer := rr.GetPeer()
	if peer != "" {
		t.Errorf("Expected empty peer when there are no servers, but got: %s", peer)
	}

	servers := []string{
		"192.168.0.1",
		"192.168.0.2",
		"192.168.0.3",
	}

	expectedPeers := make([]string, len(servers))
	copy(expectedPeers, servers)
	fmt.Println(expectedPeers)

	if err := rr.Add(servers...); err != nil {
		t.Fatalf("Err ==> %v\n", err)
	}

	for i := 0; i < len(expectedPeers); i++ {
		peer := rr.GetPeer()
		if peer != expectedPeers[i] {
			t.Errorf("Expected peer: %s, but got: %s", expectedPeers[i], peer)
		}
	}

}
