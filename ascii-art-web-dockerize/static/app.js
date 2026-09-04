(function () {
    "use strict";

    const form = document.getElementById("asciiForm");
    const textInput = document.getElementById("text");
    const generateBtn = document.getElementById("generateBtn");
    const generateLabel = document.getElementById("generateLabel");
    const clearBtn = document.getElementById("clearBtn");
    const copyBtn = document.getElementById("copyBtn");
    const downloadBtn = document.getElementById("downloadBtn");
    const sampleButtons = document.querySelectorAll(".sample-chip");
    const charCount = document.getElementById("charCount");
    const lineCount = document.getElementById("lineCount");
    const inputHint = document.getElementById("inputHint");
    const validationMessage = document.getElementById("validationMessage");
    const statusCard = document.getElementById("appStatus");
    const statusLabel = document.getElementById("statusLabel");
    const statusDetail = document.getElementById("statusDetail");
    const notice = document.getElementById("notice");
    const noticeTitle = document.getElementById("noticeTitle");
    const noticeText = document.getElementById("noticeText");
    const asciiResult = document.getElementById("asciiResult");
    const outputFrame = document.querySelector(".output-frame");
    const metricInput = document.getElementById("metricInput");
    const metricOutput = document.getElementById("metricOutput");
    const metricBanner = document.getElementById("metricBanner");
    const metricTime = document.getElementById("metricTime");

    const maxLength = Number(form.dataset.maxLength || textInput.maxLength || 1000);
    let latestResult = asciiResult.dataset.empty === "false" ? asciiResult.textContent : "";
    let touched = textInput.value.length > 0;

    function selectedBanner() {
        const checked = form.querySelector("input[name='banner']:checked");
        return checked ? checked.value : "standard";
    }

    function lineTotal(value) {
        if (value.length === 0) {
            return 0;
        }
        return value.replace(/\r\n/g, "\n").split("\n").length;
    }

    function outputLineTotal(value) {
        if (!value) {
            return 0;
        }
        const lines = value.endsWith("\n") ? value.slice(0, -1).split("\n") : value.split("\n");
        return lines.length;
    }

    function firstUnsupportedCharacter(value) {
        const normalized = value.replace(/\r\n/g, "\n");
        for (const char of normalized) {
            const code = char.codePointAt(0);
            if ((code < 32 || code > 126) && code !== 10) {
                return { char, code };
            }
        }
        return null;
    }

    function setStatus(state, label, detail) {
        statusCard.dataset.state = state;
        statusLabel.textContent = label;
        statusDetail.textContent = detail;
    }

    function setNotice(tone, title, message) {
        notice.className = "notice tone-" + tone;
        notice.setAttribute("role", tone === "danger" ? "alert" : "status");
        noticeTitle.textContent = title;
        noticeText.textContent = message;
    }

    function setValidation(state, message) {
        validationMessage.dataset.state = state || "";
        validationMessage.textContent = message || "";
        validationMessage.setAttribute("role", state === "danger" ? "alert" : "status");
    }

    function updateMetrics(duration) {
        const value = textInput.value;
        metricInput.textContent = value.length + (value.length === 1 ? " char" : " chars");
        metricOutput.textContent = outputLineTotal(latestResult) + " lines";
        metricBanner.textContent = selectedBanner();
        metricTime.textContent = duration || "-";
    }

    function updateActionState() {
        const hasResult = latestResult.length > 0;
        copyBtn.disabled = !hasResult;
        downloadBtn.disabled = !hasResult;
    }

    function celebrateOutput() {
        if (!outputFrame) {
            return;
        }
        outputFrame.classList.remove("is-celebrating");
        void outputFrame.offsetWidth;
        outputFrame.classList.add("is-celebrating");
        window.setTimeout(function () {
            outputFrame.classList.remove("is-celebrating");
        }, 900);
    }

    function validateInput() {
        const value = textInput.value;
        const count = value.length;
        const remaining = maxLength - count;
        const unsupported = firstUnsupportedCharacter(value);

        charCount.textContent = count + " / " + maxLength;
        lineCount.textContent = lineTotal(value) + (lineTotal(value) === 1 ? " line" : " lines");
        inputHint.textContent = remaining + " characters remaining.";
        charCount.dataset.state = remaining < 0 ? "danger" : remaining <= Math.ceil(maxLength * 0.1) ? "warning" : "";

        if (count === 0) {
            generateBtn.disabled = true;
            textInput.setAttribute("aria-invalid", touched ? "true" : "false");
            setValidation(touched ? "danger" : "", touched ? "Text is required." : "");
            return false;
        }
        if (count > maxLength) {
            generateBtn.disabled = true;
            textInput.setAttribute("aria-invalid", "true");
            setValidation("danger", "Text is too long.");
            return false;
        }
        if (unsupported) {
            generateBtn.disabled = true;
            textInput.setAttribute("aria-invalid", "true");
            setValidation("danger", "Unsupported character: U+" + unsupported.code.toString(16).toUpperCase().padStart(4, "0"));
            return false;
        }

        generateBtn.disabled = false;
        textInput.setAttribute("aria-invalid", "false");
        setValidation("success", "Input is ready.");
        return true;
    }

    async function generateAscii() {
        if (!validateInput()) {
            setStatus("danger", "Needs attention", "Fix the input before generating.");
            return;
        }

        const payload = {
            text: textInput.value,
            banner: selectedBanner()
        };
        const controller = new AbortController();
        const timeout = window.setTimeout(function () {
            controller.abort();
        }, 12000);
        const started = performance.now();

        generateBtn.disabled = true;
        generateBtn.dataset.loading = "true";
        generateBtn.setAttribute("aria-busy", "true");
        asciiResult.setAttribute("aria-busy", "true");
        generateLabel.textContent = "Generating";
        setStatus("loading", "Generating", "Request in progress.");
        setNotice("neutral", "Generating", "Waiting for the server response.");

        try {
            const response = await fetch("/api/ascii-art", {
                method: "POST",
                headers: {
                    "Accept": "application/json",
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(payload),
                signal: controller.signal
            });
            let data;
            try {
                data = await response.json();
            } catch (err) {
                throw new Error("Server returned an unreadable response.");
            }

            if (!response.ok || data.error) {
                throw new Error(data.error || "Request failed with status " + response.status + ".");
            }

            latestResult = data.result || "";
            asciiResult.textContent = latestResult || "No output returned.";
            asciiResult.dataset.empty = latestResult ? "false" : "true";

            const duration = Math.max(1, Math.round(performance.now() - started)) + " ms";
            setStatus("success", "Generated", "Output is available.");
            setNotice("success", "Generated", "Created with the " + payload.banner + " banner.");
            updateMetrics(duration);
            updateActionState();
            celebrateOutput();
        } catch (err) {
            const message = err.name === "AbortError" ? "Request timed out. Try a shorter input." : err.message || "Network error.";
            setStatus("danger", "Could not generate", "Review the message below.");
            setNotice("danger", "Could not generate", message);
        } finally {
            window.clearTimeout(timeout);
            generateBtn.dataset.loading = "false";
            generateBtn.setAttribute("aria-busy", "false");
            asciiResult.setAttribute("aria-busy", "false");
            generateLabel.textContent = "Generate";
            validateInput();
        }
    }

    async function copyOutput() {
        if (!latestResult) {
            return;
        }

        try {
            await navigator.clipboard.writeText(latestResult);
            setStatus("success", "Copied", "Output copied to clipboard.");
            setNotice("success", "Copied", "ASCII output copied.");
        } catch (err) {
            setStatus("danger", "Copy failed", "Clipboard access was blocked.");
            setNotice("danger", "Copy failed", "Clipboard access was blocked by the browser.");
        }
    }

    function downloadOutput() {
        if (!latestResult) {
            return;
        }

        const blob = new Blob([latestResult], { type: "text/plain;charset=utf-8" });
        const url = URL.createObjectURL(blob);
        const link = document.createElement("a");
        link.href = url;
        link.download = "ascii-art-" + selectedBanner() + ".txt";
        document.body.appendChild(link);
        link.click();
        link.remove();
        URL.revokeObjectURL(url);
        setStatus("success", "Saved", "Text file download started.");
        setNotice("success", "Saved", "ASCII output saved as a text file.");
    }

    function clearForm() {
        textInput.value = "";
        latestResult = "";
        touched = false;
        asciiResult.textContent = "No output yet.";
        asciiResult.dataset.empty = "true";
        form.querySelector("input[name='banner'][value='standard']").checked = true;
        setStatus("ready", "Ready", "Waiting for input.");
        setNotice("neutral", "No output yet", "Submit text to create ASCII art.");
        setValidation("", "");
        updateMetrics("-");
        updateActionState();
        validateInput();
        textInput.focus();
    }

    form.addEventListener("submit", function (event) {
        event.preventDefault();
        touched = true;
        generateAscii();
    });

    textInput.addEventListener("input", function () {
        touched = true;
        validateInput();
        updateMetrics("-");
    });

    form.querySelectorAll("input[name='banner']").forEach(function (input) {
        input.addEventListener("change", function () {
            metricBanner.textContent = selectedBanner();
        });
    });

    sampleButtons.forEach(function (button) {
        button.addEventListener("click", function () {
            textInput.value = button.dataset.sample || "";
            touched = true;
            validateInput();
            updateMetrics("-");
            textInput.focus();
        });
    });

    clearBtn.addEventListener("click", clearForm);
    copyBtn.addEventListener("click", copyOutput);
    downloadBtn.addEventListener("click", downloadOutput);

    document.addEventListener("keydown", function (event) {
        if ((event.ctrlKey || event.metaKey) && event.key === "Enter") {
            event.preventDefault();
            touched = true;
            generateAscii();
        }
    });

    validateInput();
    updateMetrics("-");
    updateActionState();
})();
