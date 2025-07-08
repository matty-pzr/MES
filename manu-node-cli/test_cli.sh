#!/bin/bash

# Test script for Manufacturing Node Manager CLI

echo "=== Testing Manufacturing Node Manager CLI ==="
echo ""
echo "This script will demonstrate all CRUD operations."
echo "Press Ctrl+C to exit at any time."
echo ""

# Create a test input file
cat > test_input.txt << 'EOF'
help
create
CNC Machine 1
High precision CNC milling machine
drilling, milling, cutting
Factory1/Area2/Line1/Cell3
create
Welding Robot
Automated welding station
welding, spot welding, seam welding
Factory1/Area2/Line1/Cell4
list
view 20250108130742
update 20250108130742
CNC Machine 1 - Updated

drilling, milling, cutting, polishing
running

list
delete 20250108130743
y
list
exit
EOF

echo "Running test commands..."
echo ""

# Run the CLI with test input
go run cmd/main.go < test_input.txt

echo ""
echo "Test completed!"