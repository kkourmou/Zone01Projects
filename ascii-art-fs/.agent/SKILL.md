---
description: Golang CLI Agent Skills
---

# Required Skills

To successfully contribute to the `ascii-art-color` project, agents must apply the following skills and best practices:

## Test Driven Development (TDD)
- Define test files `<file>_test.go` and table-driven test cases before finalizing implementation logic.
- Ensure strict coverage for edge cases, missing arguments, and formatting errors.

## Go Standard Library Mastery
- Use `strings` (e.g., `HasPrefix`, `TrimPrefix`, `Split`, `ReplaceAll`) for text manipulation.
- Use `fmt` (e.g., `Printf`, `Println`, `Sprint`) for formatted output.
- Use `os` (e.g., `os.Args`, `os.Open`) for command line interaction and file system operations.
- Avoid any external dependencies (e.g., `flag` or `github.com/...`). Only standard library packages.

## ANSI Terminal Output
- Render output with ANSI escape codes for basic colors (`\033[31m` ... `\033[37m`, etc.).
- Ensure that color output is bounded correctly with reset sequences (`\033[0m`) after styled regions.
