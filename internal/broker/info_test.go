package broker

import "testing"

func TestBrokerInfo(t *testing.T) {
	b := BrokerInfo{
		ID:   1,
		Host: "localhost",
		Port: 9092,
	}
	if b.ID != 1 {
		t.Error("Expected ID 1, got ", b.ID)
	}
	if b.Host != "localhost" {
		t.Error("Expected Host localhost, got ", b.Host)
	}
	if b.Port != 9092 {
		t.Error("Expected Port 9092, got ", b.Port)
	}
}