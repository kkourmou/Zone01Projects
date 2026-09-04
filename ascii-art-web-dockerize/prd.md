# PRD – ascii-art-web-dockerize

## 1. Problem Statement
ascii-art-web-dockerize is a containerized Go HTTP server that provides a web UI for generating ASCII art (based on the previous ascii-art project). Users enter text, choose a banner style, and view the generated ASCII art in the browser.

Compared to the original ascii-art-web, this version includes a responsive single-page interface that calls the backend asynchronously via a JSON REST API, and it can be built and run via Docker.

---

## 2. User / Use Case

### Primary Users
- Users wanting to generate ASCII-art visually without using the command line.
- Web developers learning Go HTTP server implementation, HTML templates, and basic JSON APIs.
- Learners practicing Docker “good practices” with a small Go web service.

### Use Cases
- Access a web page to input a string of text.
- Select from multiple ASCII-art banners (shadow, standard, thinkertoy).
- Submit the data and view the resulting ASCII-art directly in the browser without a full page reload.

---

## 3. Web Interface Requirements

### Main Page Features
- Responsive single-page UI served at `/`.
- A text input area for the string to convert.
- A banner selector for `shadow`, `standard`, and `thinkertoy`.
- A generate action that triggers an asynchronous request to the backend (no full page reload).
- The generated ASCII art is displayed in a preformatted area (`<pre>` or equivalent) that preserves spacing.
- Clear user-facing error messages for invalid input, unsupported characters, or server errors.

---

## 4. Technical Requirements

### 4.1 Endpoints
- **`GET /`**:
  - Serves the main HTML page (SPA shell).
  - Uses Go `html/template` to render the page.
- **`POST /api/ascii-art`**:
  - JSON REST endpoint used by the SPA.
  - Accepts a JSON payload containing the input text and selected banner.
  - Returns JSON containing the generated ASCII art (or a structured error response).
- **`POST /ascii-art`** (legacy/backward compatible):
  - Form-encoded endpoint maintained for compatibility with the earlier ascii-art-web behavior.
  - Accepts form fields for text and banner and returns an HTML response with the result.
- **Static assets**:
  - **`GET /static/app.css`** and **`GET /static/app.js`** are served explicitly.
  - Directory listings for `/static/` must not be exposed.

### 4.2 Error Handling and HTTP Status Codes
Endpoints must return appropriate HTTP status codes:
- **`200 OK`**: Everything went without errors.
- **`400 Bad Request`**: Invalid input (malformed JSON / form data, empty text where disallowed, unsupported characters, invalid banner value).
- **`404 Not Found`**: Resources not found (missing templates or banner files, invalid URL routes).
- **`500 Internal Server Error`**: Unhandled server errors (template rendering failures, unexpected internal errors).

### 4.3 Architecture
- **Language**: Go
- **Packages**: Only standard Go packages are allowed (no external libraries/frameworks).
- **Templates**: HTML templates must be stored in the `templates` directory at the project root.
- **Static files**: Frontend assets are stored under `static/` and served without directory listing.
- **Best Practices**: The code must respect good Go programming practices.

### 4.4 Docker Requirements
- Build and run via a `Dockerfile` using a multi-stage build.
- Final image should be minimal (no Go toolchain) and run as a non-root user.
- Only required runtime assets are included in the final image (server binary, templates, banner files, static assets).

---

## 5. Non-Goals
- No external Go packages or routers (like Gin or Echo).
- No user authentication, persistence, or multi-user state.
- No advanced frontend frameworks (keep the UI lightweight).

---

## 6. Acceptance Criteria
- Running the server and visiting `/` on the configured port displays the SPA UI.
- Submitting valid input with a banner choice returns ASCII art and renders it without a full page reload.
- `POST /api/ascii-art` validates input and returns correct status codes and responses for success and error cases.
- Legacy `POST /ascii-art` continues to work for form submissions.
- Missing banner files/templates or invalid endpoints correctly trigger 404/500 errors as appropriate.
- The application can be built and run using Docker and serves templates/static assets correctly in the container.
