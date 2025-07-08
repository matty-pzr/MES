package node

import (
	"testing"
	"time"
)

func TestNewNode(t *testing.T) {
	title := "Test Node"
	description := "This is a test node"
	operations := []string{"cutting", "welding", "painting"}
	unsAddress := "Factory1/Area2/Line3/Cell4"

	n := NewNode(title, description, operations, unsAddress)

	// Check basic fields
	if n.Title != title {
		t.Errorf("Expected title %s, got %s", title, n.Title)
	}
	if n.Description != description {
		t.Errorf("Expected description %s, got %s", description, n.Description)
	}
	if len(n.Operations) != len(operations) {
		t.Errorf("Expected %d operations, got %d", len(operations), len(n.Operations))
	}
	for i, op := range operations {
		if n.Operations[i] != op {
			t.Errorf("Expected operation %s at index %d, got %s", op, i, n.Operations[i])
		}
	}
	if n.UNSAddress != unsAddress {
		t.Errorf("Expected UNS address %s, got %s", unsAddress, n.UNSAddress)
	}

	// Check defaults
	if n.ID == "" {
		t.Error("Expected ID to be generated, got empty string")
	}
	if n.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
	if n.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}

func TestGenerateID(t *testing.T) {
	id1 := generateID()
	time.Sleep(time.Second) // Ensure different timestamp
	id2 := generateID()

	if id1 == "" {
		t.Error("Expected non-empty ID")
	}
	if id1 == id2 {
		t.Error("Expected unique IDs for different timestamps")
	}
	if len(id1) != 14 {
		t.Errorf("Expected ID length of 14, got %d", len(id1))
	}
}


func TestNodeStructure(t *testing.T) {
	// Test that Node struct can be properly created with all fields
	now := time.Now()
	n := &Node{
		ID:          "test-id",
		Title:       "Test Node",
		Description: "A test node",
		Operations:  []string{"op1", "op2"},
		UNSAddress:  "test/address",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if n.ID != "test-id" {
		t.Errorf("Expected ID test-id, got %s", n.ID)
	}
	if n.Title != "Test Node" {
		t.Errorf("Expected title Test Node, got %s", n.Title)
	}
}