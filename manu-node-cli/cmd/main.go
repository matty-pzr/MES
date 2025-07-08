package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"manu-node-cli/internal/node"
	"manu-node-cli/internal/storage"
)

// nodeCompleter generates completion items for node IDs and titles
func createNodeCompleter(store *storage.Storage) func(string) []string {
	return func(line string) []string {
		nodes, err := store.Load()
		if err != nil {
			return nil
		}
		
		var suggestions []string
		for _, n := range nodes {
			// Add ID
			suggestions = append(suggestions, n.ID)
			// Add Title (if not empty)
			if n.Title != "" {
				suggestions = append(suggestions, n.Title)
			}
		}
		return suggestions
	}
}

func main() {
	// Create color printers for nice output
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	fmt.Println(cyan("=== Manufacturing Node Manager CLI ==="))
	fmt.Println("Type 'help' for available commands")
	fmt.Println()

	// Initialize storage
	dataDir := filepath.Join(".", "data")
	store, err := storage.NewStorage(dataDir)
	if err != nil {
		fmt.Printf("%s: Failed to initialize storage: %v\n", red("Error"), err)
		os.Exit(1)
	}

	// Create completer
	nodeCompleter := createNodeCompleter(store)
	completer := readline.NewPrefixCompleter(
		readline.PcItem("create"),
		readline.PcItem("list"),
		readline.PcItem("view", readline.PcItemDynamic(nodeCompleter)),
		readline.PcItem("update", readline.PcItemDynamic(nodeCompleter)),
		readline.PcItem("delete", readline.PcItemDynamic(nodeCompleter)),
		readline.PcItem("clear"),
		readline.PcItem("cls"),
		readline.PcItem("help"),
		readline.PcItem("exit"),
		readline.PcItem("quit"),
	)

	// Configure readline
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[32m>>> \033[0m",
		HistoryFile:     filepath.Join(dataDir, ".history"),
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		fmt.Printf("%s: Failed to initialize readline: %v\n", red("Error"), err)
		os.Exit(1)
	}
	defer rl.Close()

	// Main loop
	for {
		// Read user input
		line, err := rl.Readline()
		if err != nil { // io.EOF or user pressed Ctrl+C
			break
		}
		
		input := strings.TrimSpace(line)
		
		// Split input into command and arguments
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}
		
		command := parts[0]
		
		// Handle commands
		switch command {
		case "help":
			showHelp()
		case "create":
			handleCreate(store, rl)
		case "list":
			handleList(store)
		case "view":
			if len(parts) < 2 {
				fmt.Println(red("Usage: view <node-id or title>"))
				continue
			}
			handleView(store, strings.Join(parts[1:], " "))
		case "update":
			if len(parts) < 2 {
				fmt.Println(red("Usage: update <node-id or title>"))
				continue
			}
			handleUpdate(store, strings.Join(parts[1:], " "))
		case "delete":
			if len(parts) < 2 {
				fmt.Println(red("Usage: delete <node-id or title>"))
				continue
			}
			handleDelete(store, strings.Join(parts[1:], " "))
		case "clear", "cls":
			handleClear()
		case "exit", "quit":
			fmt.Println(yellow("Goodbye!"))
			return
		default:
			fmt.Printf("Unknown command: %s. Type 'help' for available commands.\n", command)
		}
	}
}

func showHelp() {
	fmt.Println("\nAvailable commands:")
	fmt.Println("  create  - Create a new manufacturing node")
	fmt.Println("  list    - List all nodes")
	fmt.Println("  view    - View details of a specific node")
	fmt.Println("  update  - Update a node")
	fmt.Println("  delete  - Delete a node")
	fmt.Println("  clear   - Clear the screen")
	fmt.Println("  help    - Show this help message")
	fmt.Println("  exit    - Exit the program")
	fmt.Println()
}

func handleClear() {
	// ANSI escape codes to clear screen and move cursor to top
	fmt.Print("\033[H\033[2J")
	fmt.Print("\033[H")
}

func handleCreate(store *storage.Storage, rl *readline.Instance) {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	
	fmt.Println("\n" + yellow("Creating new node (Press Ctrl+C to cancel)"))
	
	// Temporarily change prompt for input
	oldPrompt := rl.Config.Prompt
	defer func() { rl.SetPrompt(oldPrompt) }()
	
	// Get title
	rl.SetPrompt("Node title: ")
	title, err := rl.Readline()
	if err != nil { // User pressed Ctrl+C
		fmt.Println(red("\nCancelled"))
		return
	}
	title = strings.TrimSpace(title)
	
	// Validate title
	if title == "" {
		fmt.Printf("%s: Title cannot be empty\n", red("Error"))
		return
	}
	
	// Check for control characters or non-printable characters
	if !isValidInput(title) {
		fmt.Printf("%s: Title contains invalid characters\n", red("Error"))
		return
	}
	
	// Check if title is unique
	unique, err := store.IsTitleUnique(title, "")
	if err != nil {
		fmt.Printf("%s: Failed to check title uniqueness: %v\n", red("Error"), err)
		return
	}
	if !unique {
		fmt.Printf("%s: A node with title '%s' already exists\n", red("Error"), title)
		return
	}
	
	// Get description
	rl.SetPrompt("Description: ")
	description, err := rl.Readline()
	if err != nil {
		fmt.Println(red("\nCancelled"))
		return
	}
	description = strings.TrimSpace(description)
	
	// Validate description if not empty
	if description != "" && !isValidInput(description) {
		fmt.Printf("%s: Description contains invalid characters\n", red("Error"))
		return
	}
	
	// Get operations
	rl.SetPrompt("Operations (comma-separated): ")
	opsInput, err := rl.Readline()
	if err != nil {
		fmt.Println(red("\nCancelled"))
		return
	}
	opsInput = strings.TrimSpace(opsInput)
	var operations []string
	if opsInput != "" {
		operations = strings.Split(opsInput, ",")
		for i := range operations {
			operations[i] = strings.TrimSpace(operations[i])
		}
		// Remove any empty strings from the operations and validate
		var filtered []string
		for _, op := range operations {
			if op != "" {
				if !isValidInput(op) {
					fmt.Printf("%s: Operation '%s' contains invalid characters\n", red("Error"), op)
					return
				}
				filtered = append(filtered, op)
			}
		}
		operations = filtered
	}
	
	// Get UNS address
	rl.SetPrompt("UNS Address (e.g., Site/Area/Line/Cell): ")
	unsAddress, err := rl.Readline()
	if err != nil {
		fmt.Println(red("\nCancelled"))
		return
	}
	unsAddress = strings.TrimSpace(unsAddress)
	
	// Validate UNS address if not empty
	if unsAddress != "" && !isValidInput(unsAddress) {
		fmt.Printf("%s: UNS Address contains invalid characters\n", red("Error"))
		return
	}
	
	// Create the node
	newNode := node.NewNode(title, description, operations, unsAddress)
	
	// Save to storage
	if err := store.SaveNode(newNode); err != nil {
		fmt.Printf("%s: Failed to save node: %v\n", red("Error"), err)
		return
	}
	
	fmt.Printf("\n%s Node created successfully!\n", green("✓"))
	fmt.Printf("ID: %s\n", newNode.ID)
	fmt.Printf("Title: %s\n\n", newNode.Title)
}

func handleList(store *storage.Storage) {
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	
	nodes, err := store.Load()
	if err != nil {
		fmt.Printf("%s: Failed to load nodes: %v\n", red("Error"), err)
		return
	}
	
	if len(nodes) == 0 {
		fmt.Println("\nNo nodes found. Create some nodes first!")
		fmt.Println()
		return
	}
	
	// Display header
	fmt.Println("\n" + cyan("Manufacturing Nodes:"))
	fmt.Println(strings.Repeat("-", 90))
	fmt.Printf("%-20s %-30s %-35s\n", "ID", "Title", "UNS Address")
	fmt.Println(strings.Repeat("-", 90))
	
	// Display nodes
	for _, n := range nodes {
		fmt.Printf("%-20s %-30s %-35s\n", 
			n.ID, 
			truncate(n.Title, 28), 
			truncate(n.UNSAddress, 33))
	}
	fmt.Println()
}

func handleView(store *storage.Storage, identifier string) {
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	
	n, err := store.GetNodeByIDOrTitle(identifier)
	if err != nil {
		fmt.Printf("%s: %v\n", red("Error"), err)
		return
	}
	
	fmt.Println("\n" + cyan("Node Details:"))
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("ID:          %s\n", n.ID)
	fmt.Printf("Title:       %s\n", n.Title)
	fmt.Printf("Description: %s\n", n.Description)
	fmt.Printf("UNS Address: %s\n", n.UNSAddress)
	fmt.Printf("Operations:  %s\n", strings.Join(n.Operations, ", "))
	fmt.Printf("Created:     %s\n", n.CreatedAt.Format(time.RFC3339))
	fmt.Printf("Updated:     %s\n", n.UpdatedAt.Format(time.RFC3339))
	fmt.Println()
}

func handleUpdate(store *storage.Storage, identifier string) {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	
	// Get existing node
	existing, err := store.GetNodeByIDOrTitle(identifier)
	if err != nil {
		fmt.Printf("%s: %v\n", red("Error"), err)
		return
	}
	
	scanner := bufio.NewScanner(os.Stdin)
	
	fmt.Printf("\nUpdating node: %s\n", existing.Title)
	fmt.Println(yellow("Press Enter to keep current value"))
	
	// Update title
	fmt.Printf("Title [%s]: ", existing.Title)
	scanner.Scan()
	title := strings.TrimSpace(scanner.Text())
	if title == "" {
		title = existing.Title
	} else if title != existing.Title {
		// Check if new title is unique
		unique, err := store.IsTitleUnique(title, existing.ID)
		if err != nil {
			fmt.Printf("%s: Failed to check title uniqueness: %v\n", red("Error"), err)
			return
		}
		if !unique {
			fmt.Printf("%s: A node with title '%s' already exists\n", red("Error"), title)
			return
		}
	}
	
	// Update description
	fmt.Printf("Description [%s]: ", existing.Description)
	scanner.Scan()
	description := scanner.Text()
	if description == "" {
		description = existing.Description
	}
	
	// Update operations
	fmt.Printf("Operations [%s]: ", strings.Join(existing.Operations, ", "))
	scanner.Scan()
	opsInput := strings.TrimSpace(scanner.Text())
	operations := existing.Operations
	if opsInput != "" {
		operations = strings.Split(opsInput, ",")
		for i := range operations {
			operations[i] = strings.TrimSpace(operations[i])
		}
		// Remove any empty strings from the operations and validate
		var filtered []string
		for _, op := range operations {
			if op != "" {
				if !isValidInput(op) {
					fmt.Printf("%s: Operation '%s' contains invalid characters\n", red("Error"), op)
					return
				}
				filtered = append(filtered, op)
			}
		}
		operations = filtered
	}
	
	// Update UNS address
	fmt.Printf("UNS Address [%s]: ", existing.UNSAddress)
	scanner.Scan()
	unsAddress := scanner.Text()
	if unsAddress == "" {
		unsAddress = existing.UNSAddress
	}
	
	// Create updated node
	updated := &node.Node{
		ID:          existing.ID,
		Title:       title,
		Description: description,
		Operations:  operations,
		UNSAddress:  unsAddress,
		CreatedAt:   existing.CreatedAt,
		UpdatedAt:   time.Now(),
	}
	
	// Save updated node
	if err := store.UpdateNode(existing.ID, updated); err != nil {
		fmt.Printf("%s: Failed to update node: %v\n", red("Error"), err)
		return
	}
	
	fmt.Printf("\n%s Node updated successfully!\n", green("✓"))
}

func handleDelete(store *storage.Storage, identifier string) {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	
	// Get node to confirm
	n, err := store.GetNodeByIDOrTitle(identifier)
	if err != nil {
		fmt.Printf("%s: %v\n", red("Error"), err)
		return
	}
	
	// Confirm deletion
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("\n%s Delete node '%s' (ID: %s)? [y/N]: ", 
		yellow("Warning:"), n.Title, n.ID)
	scanner.Scan()
	confirm := strings.ToLower(strings.TrimSpace(scanner.Text()))
	
	if confirm != "y" && confirm != "yes" {
		fmt.Println("Deletion cancelled.")
		return
	}
	
	// Delete node
	if err := store.DeleteNode(n.ID); err != nil {
		fmt.Printf("%s: Failed to delete node: %v\n", red("Error"), err)
		return
	}
	
	fmt.Printf("\n%s Node deleted successfully!\n", green("✓"))
}

// Helper functions
func isValidInput(s string) bool {
	// Check if string contains only printable characters
	for _, r := range s {
		// Allow printable ASCII and common Unicode characters
		if r < 32 || r == 127 { // Control characters
			return false
		}
	}
	return true
}

func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}