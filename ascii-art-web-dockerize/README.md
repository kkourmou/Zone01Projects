# Ascii-Art-Web-Dockerize

## Description
Ascii-art-web-dockerize is a containerized version of the ascii-art-web project. It consists of an HTTP server written in Go that allows users to submit text and select an ASCII banner style (`shadow`, `standard`, or `thinkertoy`) through a responsive single-page interface, packaged and deployed using Docker.

The frontend operates asynchronously via a JSON REST API, displaying the resulting ASCII art in the browser without a full page reload.

## Authors
- kkourmou
- spetridou

## Usage: How to Build & Run

### With Docker (recommended)
1. Ensure you have Docker installed on your system.
2. Clone this repository and navigate to the project root directory.
3. Build the Docker image (includes OCI metadata labels like `org.opencontainers.image.*`):
   ```bash
   docker image build -f Dockerfile -t ascii-art-web .
   ```
4. Run the container (adds container labels via `--label`):
   ```bash
   docker container run -p 8080:8080 -d --name ascii-art ascii-art-web  
   ```
5. Open your browser at `http://localhost:8080`

### With Docker scripts (build/run/clean)
- PowerShell:
  ```powershell
  pwsh scripts/docker.ps1 build
  ```
- sh:
  ```sh
  ./scripts/docker.sh build
  ```

### Without Docker
1. Ensure you have Go installed on your system.
2. Clone this repository and navigate to the project root directory.
3. Start the server:
   ```bash
   go run .
   ```
4. Open your browser at `http://localhost:8080`

## Docker Implementation

### Dockerfile
The project uses a multi-stage Dockerfile to keep the final image small and secure:
- **Stage 1 (builder)**: Uses `golang:1.24-alpine` to compile the Go binary.
- **Stage 2 (runtime)**: Uses `scratch` and runs as a non-root user, copying only the compiled binary and required assets.

### Docker Good Practices Applied
- Multi-stage build to minimize final image size (no Go toolchain in production image).
- `.dockerignore` to exclude unnecessary files from the build context (tests, docs, git files).
- `EXPOSE` to document the port the container listens on.
- Non-root working directory `/app` inside the container.
- Metadata applied to Docker objects:
  - Image: OCI labels in `Dockerfile` (`org.opencontainers.image.*`)
  - Container: labels in `docker run` / `scripts/docker.*`

### Garbage Collection
To remove unused Docker objects and free up space:
```bash
# Remove stopped containers
docker container prune

# Remove unused images
docker image prune

# Remove all unused objects (containers, images, networks, volumes)
docker system prune
```

## Architecture & Implementation Details
The backend is built entirely using standard Go packages (`net/http`, `html/template`, `encoding/json`), strictly avoiding external dependencies.

### Endpoints
- **`GET /`**: Serves the main HTML Single-Page Application.
- **`POST /api/ascii-art`**: REST JSON API endpoint. Accepts a JSON payload with text and banner, returns generated ASCII art.
- **`POST /ascii-art`**: Legacy form-encoded endpoint maintained for backward compatibility.
- **`GET /static/app.css`** and **`GET /static/app.js`**: Serve SPA assets without exposing directory listings.

### Error Handling
- `200 OK`: Successful generation.
- `400 Bad Request`: Unsupported characters or invalid input.
- `404 Not Found`: Missing banner files or invalid routes.
- `500 Internal Server Error`: Template rendering failures or server-side errors.

### Testing
```bash
# Run the Go native testing suite
go test ./...

# Test the backend API endpoints
bash api_test.sh
```
