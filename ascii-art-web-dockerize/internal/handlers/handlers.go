package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"ascii-art/internal/ascii"
)

const (
	maxRequestBytes = 8 * 1024
	maxTextBytes    = 1000
)

// PageData is the structure passed to the HTML template.
type PageData struct {
	Text           string
	SelectedBanner string
	AsciiResult    string
	ErrorMsg       string
	MaxTextBytes   int
	Banners        []BannerOption
}

// BannerOption describes one selectable ASCII banner font.
type BannerOption struct {
	Value       string
	Label       string
	Description string
}

var bannerOptions = []BannerOption{
	{Value: "standard", Label: "Standard", Description: "Balanced letterforms for general use."},
	{Value: "shadow", Label: "Shadow", Description: "Heavier output with dimensional edges."},
	{Value: "thinkertoy", Label: "Thinkertoy", Description: "Compact mechanical styling."},
}

// Parse templates from the project root "templates" directory.
var tmpl = template.Must(template.ParseFiles(projectFilePath("templates/index.html")))

// SecurityHeaders applies browser-facing hardening headers to every response.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; base-uri 'none'; form-action 'self'; frame-ancestors 'none'")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}

// HomeHandler serves the GET / endpoint.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		methodNotAllowed(w, http.MethodGet)
		return
	}

	renderPage(w, defaultPageData())
}

// AsciiArtHandler serves the legacy POST /ascii-art form endpoint.
func AsciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}

	req, status, msg := readFormRequest(w, r)
	if msg != "" {
		data := defaultPageData()
		data.Text = req.Text
		data.SelectedBanner = safeSelectedBanner(req.Banner)
		data.ErrorMsg = msg
		renderPageWithStatus(w, status, data)
		return
	}

	result, err := ascii.GenerateAscii(req.Text, req.Banner)
	if err != nil {
		status, msg = classifyGenerateError(err)
		data := defaultPageData()
		data.Text = req.Text
		data.SelectedBanner = safeSelectedBanner(req.Banner)
		data.ErrorMsg = msg
		renderPageWithStatus(w, status, data)
		return
	}

	data := defaultPageData()
	data.Text = req.Text
	data.SelectedBanner = req.Banner
	data.AsciiResult = result
	renderPage(w, data)
}

// StaticHandler serves the SPA stylesheet and script without exposing directory listings.
func StaticHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		methodNotAllowed(w, http.MethodGet)
		return
	}

	files := map[string]struct {
		path        string
		contentType string
	}{
		"/static/app.css": {path: "static/app.css", contentType: "text/css; charset=utf-8"},
		"/static/app.js":  {path: "static/app.js", contentType: "text/javascript; charset=utf-8"},
	}

	file, ok := files[r.URL.Path]
	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", file.contentType)
	w.Header().Set("Cache-Control", "public, max-age=3600")
	http.ServeFile(w, r, projectFilePath(file.path))
}

// ApiRequest represents the expected JSON payload format.
type ApiRequest struct {
	Text   string `json:"text"`
	Banner string `json:"banner"`
}

// ApiResponse represents the returned JSON payload.
type ApiResponse struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

// ApiAsciiArtHandler serves the POST /api/ascii-art JSON endpoint.
func ApiAsciiArtHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSON(w, http.StatusMethodNotAllowed, ApiResponse{Error: "405 Method Not Allowed: use POST"})
		return
	}

	req, status, msg := readJSONRequest(w, r)
	if msg != "" {
		writeJSON(w, status, ApiResponse{Error: msg})
		return
	}

	result, err := ascii.GenerateAscii(req.Text, req.Banner)
	if err != nil {
		status, msg = classifyGenerateError(err)
		writeJSON(w, status, ApiResponse{Error: msg})
		return
	}

	writeJSON(w, http.StatusOK, ApiResponse{Result: result})
}

func defaultPageData() PageData {
	return PageData{
		SelectedBanner: "standard",
		MaxTextBytes:   maxTextBytes,
		Banners:        bannerOptions,
	}
}

func renderPage(w http.ResponseWriter, data PageData) {
	renderPageWithStatus(w, http.StatusOK, data)
}

func renderPageWithStatus(w http.ResponseWriter, status int, data PageData) {
	if data.SelectedBanner == "" {
		data.SelectedBanner = "standard"
	}
	if data.MaxTextBytes == 0 {
		data.MaxTextBytes = maxTextBytes
	}
	if len(data.Banners) == 0 {
		data.Banners = bannerOptions
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		http.Error(w, "500 Internal Server Error: failed to render template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = buf.WriteTo(w)
}

func readJSONRequest(w http.ResponseWriter, r *http.Request) (ApiRequest, int, string) {
	if msg := validateJSONContentType(r); msg != "" {
		return ApiRequest{}, http.StatusUnsupportedMediaType, msg
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBytes)
	defer r.Body.Close()

	var req ApiRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		return ApiRequest{}, requestDecodeStatus(err), decodeErrorMessage(err)
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return ApiRequest{}, http.StatusBadRequest, "400 Bad Request: JSON body must contain a single object"
	}

	return validateRequest(req)
}

func readFormRequest(w http.ResponseWriter, r *http.Request) (ApiRequest, int, string) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBytes)
	defer r.Body.Close()

	if err := r.ParseForm(); err != nil {
		return ApiRequest{}, requestDecodeStatus(err), decodeErrorMessage(err)
	}

	return validateRequest(ApiRequest{
		Text:   r.PostForm.Get("text"),
		Banner: r.PostForm.Get("banner"),
	})
}

func validateRequest(req ApiRequest) (ApiRequest, int, string) {
	req.Banner = strings.TrimSpace(req.Banner)
	req.Text = strings.ReplaceAll(req.Text, "\r\n", "\n")
	req.Text = strings.ReplaceAll(req.Text, "\r", "\n")

	if req.Text == "" {
		return req, http.StatusBadRequest, "400 Bad Request: text is required"
	}
	if len(req.Text) > maxTextBytes {
		return req, http.StatusRequestEntityTooLarge, fmt.Sprintf("413 Payload Too Large: text must be %d bytes or fewer", maxTextBytes)
	}
	if !ascii.IsValidBanner(req.Banner) {
		return req, http.StatusBadRequest, "400 Bad Request: choose a supported banner"
	}
	for _, ch := range req.Text {
		if (ch < ' ' || ch > '~') && ch != '\n' {
			return req, http.StatusBadRequest, fmt.Sprintf("400 Bad Request: unsupported character %q", ch)
		}
	}

	return req, http.StatusOK, ""
}

func validateJSONContentType(r *http.Request) string {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return "415 Unsupported Media Type: Content-Type must be application/json"
	}

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil || mediaType != "application/json" {
		return "415 Unsupported Media Type: Content-Type must be application/json"
	}

	return ""
}

func requestDecodeStatus(err error) int {
	var maxBytesError *http.MaxBytesError
	if errors.As(err, &maxBytesError) {
		return http.StatusRequestEntityTooLarge
	}
	return http.StatusBadRequest
}

func decodeErrorMessage(err error) string {
	var maxBytesError *http.MaxBytesError
	if errors.As(err, &maxBytesError) {
		return "413 Payload Too Large: request body is too large"
	}
	return "400 Bad Request: invalid request payload"
}

func classifyGenerateError(err error) (int, string) {
	errStr := err.Error()
	if strings.Contains(errStr, "unsupported character") || strings.Contains(errStr, "invalid banner") {
		return http.StatusBadRequest, errStr
	}
	if strings.Contains(errStr, "banner not found") {
		return http.StatusNotFound, "404 Not Found: selected banner is not available"
	}
	return http.StatusInternalServerError, "500 Internal Server Error: could not generate ASCII art"
}

func safeSelectedBanner(banner string) string {
	if ascii.IsValidBanner(banner) {
		return banner
	}
	return "standard"
}

func writeJSON(w http.ResponseWriter, status int, res ApiResponse) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "500 Internal Server Error: failed to encode response", http.StatusInternalServerError)
	}
}

func methodNotAllowed(w http.ResponseWriter, allowed string) {
	w.Header().Set("Allow", allowed)
	http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
}

func projectFilePath(relativePath string) string {
	cwd, err := os.Getwd()
	if err != nil {
		return relativePath
	}

	for {
		candidate := filepath.Join(cwd, relativePath)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}

		parent := filepath.Dir(cwd)
		if parent == cwd {
			return relativePath
		}
		cwd = parent
	}
}
