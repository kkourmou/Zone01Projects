# Agents Configuration

This project utilizes a structured agentic coding approach to divide the work among different specialized agent roles. Below are the definitions for the roles in `ascii-art-color`.

## Roles

### Coder 1: Orchestrator & CLI Parser
- **Responsibilities:**
  - Parse CLI arguments (`--color`, `<substring>`, `<text>`).
  - Validate flags.
  - Wire components together in `main.go`.

### Coder 2: Color Parser
- **Responsibilities:**
  - Convert color names/formats (CSS named colors, hex, rgb, hsl) to ANSI escape codes.
  - Maintain the dictionary of supported colors.

### Coder 3: Colorizer (Renderer Logic)
- **Responsibilities:**
  - Apply generated ANSI codes to the ASCII art blocks.
  - Handle both full string coloring and specific substring matching.
  - Wrap segments properly to ensure correct terminal rendering.

### Coder 4: Tester & Auditor
- **Responsibilities:**
  - Write independent integration tests.
  - Ensure compatibility with standard project constraints.
  - Run the full audit suite defined in the `README.md`.
