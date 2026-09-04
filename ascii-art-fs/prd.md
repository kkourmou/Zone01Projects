# PRD – ascii‑art‑color (Week 2 Deliverable)

## 1. Problem Statement
The program must generate ASCII‑art from a given string and optionally color either the entire output or specific substrings using a `--color=<color>` flag. The system must correctly parse color flags, validate formats, apply coloring only to matching substrings, and remain compatible with other optional flags (e.g., banner selection) if implemented.

---

## 2. User / Use Case

### Primary Users
- Students learning Go, CLI tools, and terminal styling  
- Developers needing stylized terminal banners

### Use Cases
- Convert text into ASCII‑art  
- Highlight specific parts of the output using colors  
- Create decorative CLI headers or readable logs

---

## 3. CLI Contract

### Valid Commands
- `go run . <string>`
- `go run . --color=<color> "<string>"`
- `go run . --color=<color> <substring> "<string>"`
- `go run . [OPTION] [BANNER] <string>` (if additional flags exist)

### Inputs
- **STRING:** text to convert into ASCII‑art  
- **OPTION:**  
  - `--color=<color>` where `<color>` may be:
    - Named color (`red`, `blue`, `yellow`, etc.)
    - Hex (`#ff0000`)
    - RGB (`rgb(255,0,0)`)
    - HSL (`hsl(0,100%,50%)`)
- **Substring:** optional  
  - If omitted → entire output is colored  
  - If provided → only exact matches are colored (case‑sensitive)

### Output
- ASCII‑art representation of the input string  
- Colored segments according to the flag  
- Correct handling of spaces, newlines, and special characters

### Error Handling
Invalid flag formats must print:

~~~
Usage: go run . [OPTION] [STRING]
EX: go run . --color=<color> <substring to be colored> "something"
~~~

---

## 4. Functional Requirements

### 4.1 ASCII‑Art Rendering
- Characters rendered using banner files (8 lines per character)
- Support all printable ASCII characters
- Handle `\n` inside the input string
- Empty string → no output  
- `"\n"` → one empty line

### 4.2 Coloring Rules
- If substring provided:
  - Only exact matches are colored
  - All occurrences must be colored
  - Case‑sensitive
- If substring not provided:
  - Entire ASCII‑art output is colored
- Coloring must use ANSI escape codes (or conversions from hex/rgb/hsl)

### 4.3 Flag Validation
- `--color=<color>` → valid  
- `--color red` → invalid  
- `--colour=red` → invalid  
- Unknown flags → invalid  

### 4.4 Compatibility
- Other ascii‑art flags (e.g., `--banner=shadow`) must still work  
- Program must still run with only a string argument

---

## 5. Non‑Goals
- No external packages  
- No multiple colors in one run  
- No Unicode support  
- No interactive mode  
- No dynamic banner editing  
- No gradients or animations  

---

## 6. Acceptance Criteria

### 6.1 Required Audit Tests
- `--color red "banana"` → invalid → usage message  
- `--color=red "hello world"` → whole output red  
- `--color=green "1 + 1 = 2"` → whole output green  
- `--color=yellow "(%&) ??"` → whole output yellow  
- `--color=red ell "hello"` → only “ell” colored  
- `--color=blue B 'RGB()'` → only “B” colored  
- `--color=orange GuYs "HeY GuYs"` → only “GuYs” colored  
- Random strings + random colors → must work  
- Special characters + substring → must work  
- Mixed case + substring → must work  

### 6.2 Additional Tests
- Empty string  
- Only spaces  
- Only special characters  
- Mixed alphanumeric  
- Multiple occurrences of substring  
- Substring at beginning/end/middle  

---

## 7. Implementation Approach

### Architecture Components
- **Argument parser**
  - Detects `--color=<color>`
  - Extracts substring (if provided)
  - Validates flag format

- **Color parser**
  - Accepts named, hex, rgb(), hsl()
  - Converts to ANSI escape code

- **ASCII‑art renderer**
  - Loads banner file
  - Maps characters to 8‑line blocks
  - Builds output line by line

- **Color applier**
  - If substring provided → replace occurrences  
  - If not → wrap entire output in ANSI codes  

### High‑Level Flow

~~~
Parse CLI args
|
Validate flag format
|
Parse color → convert to ANSI
|
Load banner file
|
Render ASCII-art for input string
|
If substring provided:
- Find all occurrences
- Apply color only to those segments
Else:
- Apply color to entire output
|
Print final result
~~~

---

## 8. Milestones

1. Argument parsing  
2. Color parsing + ANSI conversion  
3. ASCII‑art rendering  
4. Substring detection + coloring  
5. Error handling + usage message  
6. Full audit test coverage  
7. Optional: support multiple banner types  

---

## 9. Risks / Open Questions

### Risks
- Substring coloring must align with ASCII‑art output  
- Color conversion accuracy  
- Overlapping substrings (e.g., "aaa" with substring "aa")

### Open Questions
- Should substring matching optionally support case‑insensitive mode?  
- Should multiple colors be allowed in future versions?  
