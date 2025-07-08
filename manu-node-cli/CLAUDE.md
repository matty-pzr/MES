# Important Things to Remember

## Manufacturing Node Definition
A manufacturing node is a discrete unit of production capability - any machine, workstation, or process point that performs operations to transform materials or add value in a manufacturing workflow.

## Project Context
- Building a Manufacturing Node Manager CLI in Go
- Matty is learning Go for the first time
- Taking a simple, incremental approach

## Development Approach
1. CLI first (currently working on this)
2. Then MCP server for AI integration
3. API only if needed later for web/mobile

## Architecture Decisions
- Using JSON files for persistence (human readable, simple)
- Interactive CLI style (not command-line arguments)
- Single JSON file for all nodes (not separate files)

## Key Reminders
- Keep explanations beginner-friendly
- Small steps, don't overwhelm
- Test everything before moving forward
- Matty prefers concise responses
- Matty wants you to make suggestions, brainstorm and plan first before execution
- Testing approach: Both Claude and Matty run tests together (collaborative testing)

## Current Status
- Basic CLI structure complete
- Working on data persistence
- Node creation works (but doesn't save yet)

## Tech Stack
- Go 1.24.4
- JSON for data storage
- Will add MCP server functionality later