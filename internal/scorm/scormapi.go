package scorm

import (
	"os"
	"path/filepath"
)

func (g *Generator) generateSCORMAPI(dir string) error {
	jsContent := `// SCORM 1.2 API Wrapper
var scormAPI = null;
var scormAPIInitialized = false;

function findSCORMAPI(win) {
    var attempts = 0;
    var maxAttempts = 500;

    while (!win.API && win.parent && win.parent != win && attempts < maxAttempts) {
        attempts++;
        win = win.parent;
    }

    if (win.API) {
        return win.API;
    }

    // Check opener window
    if (window.opener) {
        attempts = 0;
        win = window.opener;
        while (!win.API && win.parent && win.parent != win && attempts < maxAttempts) {
            attempts++;
            win = win.parent;
        }
        if (win.API) {
            return win.API;
        }
    }

    return null;
}

function initSCORM() {
    scormAPI = findSCORMAPI(window);

    if (scormAPI) {
        var result = scormAPI.LMSInitialize("");
        if (result == "true") {
            scormAPIInitialized = true;
            console.log("SCORM API initialized successfully");

            // Set initial status
            var status = scormAPI.LMSGetValue("cmi.core.lesson_status");
            if (status == "not attempted" || status == "") {
                setSCORMValue("cmi.core.lesson_status", "incomplete");
            }
        } else {
            console.error("Failed to initialize SCORM API");
        }
    } else {
        console.warn("SCORM API not found - running in standalone mode");
    }
}

function setSCORMValue(element, value) {
    if (scormAPIInitialized && scormAPI) {
        var result = scormAPI.LMSSetValue(element, value);
        if (result == "true") {
            scormAPI.LMSCommit("");
            console.log("Set " + element + " = " + value);
            return true;
        } else {
            console.error("Failed to set " + element);
            return false;
        }
    }
    return false;
}

function getSCORMValue(element) {
    if (scormAPIInitialized && scormAPI) {
        var value = scormAPI.LMSGetValue(element);
        console.log("Get " + element + " = " + value);
        return value;
    }
    return "";
}

function finishSCORM() {
    if (scormAPIInitialized && scormAPI) {
        scormAPI.LMSCommit("");
        var result = scormAPI.LMSFinish("");
        if (result == "true") {
            console.log("SCORM session finished successfully");
        } else {
            console.error("Failed to finish SCORM session");
        }
    }
}
`

	return os.WriteFile(filepath.Join(dir, "scormapi.js"), []byte(jsContent), 0644)
}
