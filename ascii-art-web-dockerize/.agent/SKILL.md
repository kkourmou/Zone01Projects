---
description: Go HTTP Web Server Skills
---

# Required Skills

To successfully contribute to the `ascii-art-web` project, agents must apply the following skills and best practices:

## Go Standard Library Web Development
- Master `net/http` package components (`http.HandleFunc`, `http.ListenAndServe`, `http.ResponseWriter`, `http.Request`).
- Capture form requests perfectly via parsing tools like `r.FormValue`.
- Send standard web HTTP error codes like `http.StatusBadRequest`, `http.StatusNotFound`, and `http.StatusInternalServerError` intelligently upon failure.

## HTML Templates Mastery
- Employ the standard `html/template` package to inject struct data into `index.html`. 
- Be aware of executing template renders to the `ResponseWriter` safely.
- No external HTML templating systems or JS frameworks are needed.

## Precise Text Processing
- Continue using efficient `os` and `strings` tools to mathematically jump through 8-line tall representations of characters based on specific `.txt` banners.
- Be able to strictly apply TDD edge cases to mimic literal translation mappings directly from `auditors_tests.md`.
