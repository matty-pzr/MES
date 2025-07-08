package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"manu-node-cli/internal/node"
)

// Storage handles persistence of manufacturing nodes
type Storage struct {
	filePath string
	mu       sync.RWMutex
}

// NewStorage creates a new storage instance
func NewStorage(dataDir string) (*Storage, error) {
	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	return &Storage{
		filePath: filepath.Join(dataDir, "nodes.json"),
	}, nil
}

// Load reads all nodes from storage
func (s *Storage) Load() ([]*node.Node, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check if file exists
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		// Return empty slice if file doesn't exist yet
		return []*node.Node{}, nil
	}

	// Read file
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read nodes file: %w", err)
	}

	// Handle empty file
	if len(data) == 0 {
		return []*node.Node{}, nil
	}

	// Unmarshal JSON
	var nodes []*node.Node
	if err := json.Unmarshal(data, &nodes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal nodes: %w", err)
	}

	return nodes, nil
}

// Save writes all nodes to storage
func (s *Storage) Save(nodes []*node.Node) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Marshal to JSON with indentation for readability
	data, err := json.MarshalIndent(nodes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal nodes: %w", err)
	}

	// Write to file
	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write nodes file: %w", err)
	}

	return nil
}

// SaveNode adds or updates a single node
func (s *Storage) SaveNode(n *node.Node) error {
	nodes, err := s.Load()
	if err != nil {
		return err
	}

	// Check if node exists
	found := false
	for i, existing := range nodes {
		if existing.ID == n.ID {
			nodes[i] = n
			found = true
			break
		}
	}

	// Add new node if not found
	if !found {
		nodes = append(nodes, n)
	}

	return s.Save(nodes)
}

// GetNode retrieves a node by ID
func (s *Storage) GetNode(id string) (*node.Node, error) {
	nodes, err := s.Load()
	if err != nil {
		return nil, err
	}

	for _, n := range nodes {
		if n.ID == id {
			return n, nil
		}
	}

	return nil, fmt.Errorf("node with ID %s not found", id)
}

// GetNodeByTitle retrieves a node by title (case-insensitive)
func (s *Storage) GetNodeByTitle(title string) (*node.Node, error) {
	nodes, err := s.Load()
	if err != nil {
		return nil, err
	}

	titleLower := strings.ToLower(title)
	var matches []*node.Node
	
	for _, n := range nodes {
		if strings.ToLower(n.Title) == titleLower {
			matches = append(matches, n)
		}
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("node with title '%s' not found", title)
	}
	if len(matches) > 1 {
		return nil, fmt.Errorf("multiple nodes found with title '%s'", title)
	}

	return matches[0], nil
}

// GetNodeByIDOrTitle tries to get a node by ID first, then by title
func (s *Storage) GetNodeByIDOrTitle(identifier string) (*node.Node, error) {
	// Try ID first
	n, err := s.GetNode(identifier)
	if err == nil {
		return n, nil
	}

	// Try title
	return s.GetNodeByTitle(identifier)
}

// IsTitleUnique checks if a title is unique (case-insensitive)
// excludeID allows checking uniqueness while updating an existing node
func (s *Storage) IsTitleUnique(title string, excludeID string) (bool, error) {
	nodes, err := s.Load()
	if err != nil {
		return false, err
	}

	titleLower := strings.ToLower(title)
	for _, n := range nodes {
		if n.ID != excludeID && strings.ToLower(n.Title) == titleLower {
			return false, nil
		}
	}

	return true, nil
}

// UpdateNode updates an existing node
func (s *Storage) UpdateNode(id string, updated *node.Node) error {
	nodes, err := s.Load()
	if err != nil {
		return err
	}

	found := false
	for i, n := range nodes {
		if n.ID == id {
			// Preserve original ID and creation time
			updated.ID = id
			updated.CreatedAt = n.CreatedAt
			nodes[i] = updated
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("node with ID %s not found", id)
	}

	return s.Save(nodes)
}

// DeleteNode removes a node by ID
func (s *Storage) DeleteNode(id string) error {
	nodes, err := s.Load()
	if err != nil {
		return err
	}

	// Find and remove node
	found := false
	var filtered []*node.Node
	for _, n := range nodes {
		if n.ID == id {
			found = true
			continue
		}
		filtered = append(filtered, n)
	}

	if !found {
		return fmt.Errorf("node with ID %s not found", id)
	}

	return s.Save(filtered)
}