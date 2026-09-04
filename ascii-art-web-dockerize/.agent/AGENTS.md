# Agents Configuration

This project utilizes a structured agentic coding approach to divide the work among different specialized agent roles. Below are the definitions for the roles in `ascii-art-web`.

## Roles

### Frontend & Web Server (Half 1)
- **Responsibilities:**
  - Initialize the `net/http` Go web server on a valid port.
  - Implement and manage routes for `GET /` and `POST /ascii-art`.
  - Compile the `html/template` structure and manage form inputs securely.
  - Handle appropriate HTTP protocol responses and status codes (`200 OK`, `400 Bad Request`, `404 Not Found`, `500 Internal Server`).

### ASCII Engine Logic (Half 2)
- **Responsibilities:**
  - Re-use the previous text parsing logic for mapping banner files (`shadow.txt`, `standard.txt`, `thinkertoy.txt`).
  - Read user inputs and transform them into multi-line printable blocks.
  - Return formatted strings cleanly to the frontend Server.
  - Specifically tackle all edge cases (properly mapping `\n` and catching invalid unprintable characters).

### Auditor & QA (Shared)
- **Responsibilities:**
  - Test exactly against the strings in the `auditors_tests.md` using different banners.
  - Evaluate client-server pipeline stability.
