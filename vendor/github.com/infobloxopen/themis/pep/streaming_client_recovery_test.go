package pep

import (
	"fmt"
	"testing"

	"github.com/infobloxopen/themis/pdp"
	pb "github.com/infobloxopen/themis/pdp-service"
)

const (
	fakeServerAddress    = "127.0.0.1:5555"
	fakeServerAltAddress = "127.0.0.1:5556"
)

func TestStreamClientRecovery(t *testing.T) {
	singleClientRecovery(5, t)
}

func TestConnectionClientRecovery(t *testing.T) {
	singleClientRecovery(1, t)
}

func TestStreamClientRecoveryWithHotSpotBalancer(t *testing.T) {
	hotSotBalancedClientRecovery(10, t)
}

func TestConnectionClientRecoveryWithHotSpotBalancer(t *testing.T) {
	hotSotBalancedClientRecovery(1, t)
}

func singleClientRecovery(streams int, t *testing.T) {
	s, err := newFailServer(fakeServerAddress)
	if err != nil {
		t.Fatalf("couldn't start fake server: %s", err)
	}

	defer s.Stop()

	msgs := make(chan string, 1)

	c := NewClient(
		WithStreams(streams),
		WithConnectionStateNotification(func(addr string, state int, err error) {
			if streams > 1 && state == StreamingConnectionBroken {
				msg := fmt.Sprintf("unexpected connection failure when number of streams set to %d "+
					"(expected only stream failure)", streams)
				select {
				default:
				case msgs <- msg:
				}
			}
		}),
	)

	err = c.Connect(fakeServerAddress)
	if err != nil {
		t.Fatalf("can't connect to fake server: %s", err)
	}
	defer c.Close()

	in := []pdp.AttributeAssignment{
		pdp.MakeIntegerAssignment(IDID, 1),
		pdp.MakeStringAssignment(failID, thisRequest),
	}

	var out pb.Msg
	err = c.Validate(in, &out)
	if err != nil {
		t.Fatalf("can't send first request: %s", err)
	}

	for len(msgs) > 0 {
		msg, ok := <-msgs
		if !ok {
			break
		}

		t.Error(msg)
	}

	var attempts uint64 = 2
	if s.ID != attempts {
		t.Errorf("Expected %d attempts but got %d", attempts, s.ID)
	}
}

func hotSotBalancedClientRecovery(streams int, t *testing.T) {
	s1, err := newFailServer(fakeServerAddress)
	if err != nil {
		t.Fatalf("couldn't start fake server: %s", err)
	}

	defer s1.Stop()

	s2, err := newFailServer(fakeServerAltAddress)
	if err != nil {
		t.Fatalf("couldn't start fake server: %s", err)
	}

	defer s2.Stop()

	msgs := make(chan string, 1)

	c := NewClient(
		WithStreams(streams),
		WithHotSpotBalancer(
			fakeServerAddress,
			fakeServerAltAddress,
		),
		WithConnectionStateNotification(func(addr string, state int, err error) {
			if streams > 1 && state == StreamingConnectionBroken {
				msg := fmt.Sprintf("unexpected connection failure when number of streams set to %d "+
					"(expected only stream failure)", streams)
				select {
				default:
				case msgs <- msg:
				}
			}
		}),
	)

	err = c.Connect(fakeServerAddress)
	if err != nil {
		t.Fatalf("can't connect to fake server: %s", err)
	}
	defer c.Close()

	in := []pdp.AttributeAssignment{
		pdp.MakeIntegerAssignment(IDID, 1),
		pdp.MakeStringAssignment(failID, thisRequest),
	}

	var out pb.Msg
	err = c.Validate(in, &out)
	if err != nil {
		t.Fatalf("can't send first request: %s", err)
	}

	for len(msgs) > 0 {
		msg, ok := <-msgs
		if !ok {
			break
		}

		t.Error(msg)
	}

	var attempts uint64 = 2
	total := s1.ID + s2.ID
	if total != attempts {
		t.Errorf("Expected %d attempts but got %d", attempts, total)
	}
}
