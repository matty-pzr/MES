package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"manu-node-cli/internal/node"
)

func setupTestStorage(t *testing.T) (*Storage, func()) {
	// Create temporary directory for tests
	tempDir, err := os.MkdirTemp("", "storage_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	store, err := NewStorage(tempDir)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Return cleanup function
	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return store, cleanup
}

func TestNewStorage(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "storage_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store, err := NewStorage(tempDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	if store.filePath != filepath.Join(tempDir, "nodes.json") {
		t.Errorf("Expected file path %s, got %s", 
			filepath.Join(tempDir, "nodes.json"), store.filePath)
	}
}

func TestLoadEmptyStorage(t *testing.T) {
	store, cleanup := setupTestStorage(t)
	defer cleanup()

	nodes, err := store.Load()
	if err != nil {
		t.Fatalf("Failed to load empty storage: %v", err)
	}

	if len(nodes) != 0 {
		t.Errorf("Expected 0 nodes, got %d", len(nodes))
	}
}

func TestSaveAndLoad(t *testing.T) {
	store, cleanup := setupTestStorage(t)
	defer cleanup()

	// Create test nodes
	nodes := []*node.Node{
		{
			ID:          "node1",
			Title:       "Test Node 1",
			Description: "First test node",
			Operations:  []string{"op1", "op2"},
			UNSAddress:  "test/address/1",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "node2",
			Title:       "Test Node 2",
			Description: "Second test node",
			Operations:  []string{"op3", "op4"},
			UNSAddress:  "test/address/2",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Save nodes
	err := store.Save(nodes)
	if err != nil {
		t.Fatalf("Failed to save nodes: %v", err)
	}

	// Load nodes
	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Failed to load nodes: %v", err)
	}

	if len(loaded) != len(nodes) {
		t.Errorf("Expected %d nodes, got %d", len(nodes), len(loaded))
	}

	// Verify loaded data
	for i, n := range loaded {
		if n.ID != nodes[i].ID {
			t.Errorf("Expected ID %s, got %s", nodes[i].ID, n.ID)
		}
		if n.Title != nodes[i].Title {
			t.Errorf("Expected title %s, got %s", nodes[i].Title, n.Title)
		}
	}
}

func TestSaveNode(t *testing.T) {
	store, cleanup := setupTestStorage(t)
	defer cleanup()

	// Create and save first node
	node1 := &node.Node{
		ID:          "node1",
		Title:       "Test Node 1",
		Description: "First test node",
		Operations:  []string{"op1"},
		UNSAddress:  "test/1",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := store.SaveNode(node1)
	if err != nil {
		t.Fatalf("Failed to save node: %v", err)
	}

	// Load and verify
	nodes, err := store.Load()
	if err != nil {
		t.Fatalf("Failed to load nodes: %v", err)
	}
	if len(nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(nodes))
	}

	// Save second node
	node2 := &node.Node{
		ID:          "node2",
		Title:       "Test Node 2",
		Description: "Second test node",
		Operations:  []string{"op2"},
		UNSAddress:  "test/2",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = store.SaveNode(node2)
	if err != nil {
		t.Fatalf("Failed to save second node: %v", err)
	}

	// Load and verify both nodes exist
	nodes, err = store.Load()
	if err != nil {
		t.Fatalf("Failed to load nodes: %v", err)
	}
	if len(nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(nodes))
	}
}

func TestGetNode(t *testing.T) {
	store, cleanup := setupTestStorage(t)
	defer cleanup()

	// Save a node
	testNode := &node.Node{
		ID:          "test-id",
		Title:       "Test Node",
		Description: "Test description",
		Operations:  []string{"op1"},
		UNSAddress:  "test/address",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := store.SaveNode(testNode)
	if err != nil {
		t.Fatalf("Failed to save node: %v", err)
	}

	// Get existing node
	retrieved, err := store.GetNode("test-id")
	if err != nil {
		t.Fatalf("Failed to get node: %v", err)
	}

	if retrieved.ID != testNode.ID {
		t.Errorf("Expected ID %s, got %s", testNode.ID, retrieved.ID)
	}
	if retrieved.Title != testNode.Title {
		t.Errorf("Expected title %s, got %s", testNode.Title, retrieved.Title)
	}

	// Try to get non-existent node
	_, err = store.GetNode("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent node, got nil")
	}
}

func TestUpdateNode(t *testing.T) {
	store, cleanup := setupTestStorage(t)
	defer cleanup()

	// Save initial node
	originalNode := &node.Node{
		ID:          "update-test",
		Title:       "Original Title",
		Description: "Original description",
		Operations:  []string{"op1"},
		UNSAddress:  "original/address",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := store.SaveNode(originalNode)
	if err != nil {
		t.Fatalf("Failed to save node: %v", err)
	}

	// Update node
	updatedNode := &node.Node{
		ID:          "different-id", // Should be ignored
		Title:       "Updated Title",
		Description: "Updated description",
		Operations:  []string{"op1", "op2"},
		UNSAddress:  "updated/address",
		CreatedAt:   time.Now().Add(time.Hour), // Should be preserved
		UpdatedAt:   time.Now().Add(time.Hour),
	}

	err = store.UpdateNode("update-test", updatedNode)
	if err != nil {
		t.Fatalf("Failed to update node: %v", err)
	}

	// Verify update
	retrieved, err := store.GetNode("update-test")
	if err != nil {
		t.Fatalf("Failed to get updated node: %v", err)
	}

	if retrieved.ID != "update-test" {
		t.Errorf("Expected ID to remain %s, got %s", "update-test", retrieved.ID)
	}
	if retrieved.Title != updatedNode.Title {
		t.Errorf("Expected title %s, got %s", updatedNode.Title, retrieved.Title)
	}
	if !retrieved.CreatedAt.Equal(originalNode.CreatedAt) {
		t.Error("Expected CreatedAt to be preserved from original")
	}

	// Try to update non-existent node
	err = store.UpdateNode("non-existent", updatedNode)
	if err == nil {
		t.Error("Expected error for non-existent node, got nil")
	}
}

func TestDeleteNode(t *testing.T) {
	store, cleanup := setupTestStorage(t)
	defer cleanup()

	// Save multiple nodes
	nodes := []*node.Node{
		{
			ID:        "node1",
			Title:     "Node 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "node2",
			Title:     "Node 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "node3",
			Title:     "Node 3",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, n := range nodes {
		err := store.SaveNode(n)
		if err != nil {
			t.Fatalf("Failed to save node: %v", err)
		}
	}

	// Delete middle node
	err := store.DeleteNode("node2")
	if err != nil {
		t.Fatalf("Failed to delete node: %v", err)
	}

	// Verify deletion
	remaining, err := store.Load()
	if err != nil {
		t.Fatalf("Failed to load nodes: %v", err)
	}

	if len(remaining) != 2 {
		t.Errorf("Expected 2 nodes after deletion, got %d", len(remaining))
	}

	// Verify correct nodes remain
	for _, n := range remaining {
		if n.ID == "node2" {
			t.Error("Deleted node still exists")
		}
	}

	// Try to delete non-existent node
	err = store.DeleteNode("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent node, got nil")
	}
}

func TestConcurrentAccess(t *testing.T) {
	store, cleanup := setupTestStorage(t)
	defer cleanup()

	// Test concurrent saves
	done := make(chan bool, 3)

	go func() {
		for i := 0; i < 10; i++ {
			n := &node.Node{
				ID:        "concurrent1",
				Title:     "Concurrent Node 1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			store.SaveNode(n)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 10; i++ {
			store.Load()
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 10; i++ {
			store.GetNode("concurrent1")
		}
		done <- true
	}()

	// Wait for all goroutines to complete
	for i := 0; i < 3; i++ {
		<-done
	}

	// If we get here without deadlock or panic, concurrent access is safe
	t.Log("Concurrent access test passed")
}