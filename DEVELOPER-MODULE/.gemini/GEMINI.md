## Module Goal
Create a standalone app for creating and configuring 'Operation Types' (name, icon, recipe schema, ISA-95 connection details for MQTT).

## Phase 1 Plan (Current)
- Initialize React project with Vite.
- Install dependencies: `react-router-dom`, `tailwindcss`, `lucide-react`, `zustand`.
- Create the two-column layout (sidebar, content).
- Add the "Create New Operation" button.

## Data Structures
- `OperationType`: `{ id, name, icon, recipeSchema: {...}, connectionDetails: {...} }`

## Future Ideas
- Search bar
- Import/export functionality