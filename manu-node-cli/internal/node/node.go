package node

import (
	"time"
)

// Node represents a manufacturing node in the system
type Node struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Operations  []string  `json:"operations"`
	UNSAddress  string    `json:"uns_address"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewNode creates a new manufacturing node
func NewNode(title, description string, operations []string, unsAddress string) *Node {
	return &Node{
		ID:          generateID(),
		Title:       title,
		Description: description,
		Operations:  operations,
		UNSAddress:  unsAddress,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// generateID creates a simple ID for the node
func generateID() string {
	return time.Now().Format("20060102150405")
}