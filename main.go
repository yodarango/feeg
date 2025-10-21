package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

type Background struct {
	Name  string // Display name (underscores replaced with spaces)
	Value string // Actual filename with extension
	Type  string // "image" or "video"
}

var db *sql.DB

// Supported file extensions
var imageExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".bmp": true,
}

var videoExtensions = map[string]bool{
	".mp4": true, ".webm": true, ".mov": true, ".avi": true, ".mkv": true, ".flv": true,
}

func initDB() error {
	var err error
	db, err = sql.Open("sqlite3", "settings.db")
	if err != nil {
		return err
	}

	// Create table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS settings (
		id INTEGER PRIMARY KEY,
		key TEXT UNIQUE,
		value TEXT
	);
	`
	_, err = db.Exec(createTableSQL)
	return err
}

func getBackgrounds() ([]Background, error) {
	files, err := os.ReadDir("public/bkgs")
	if err != nil {
		return []Background{}, nil
	}

	var backgrounds []Background
	for _, file := range files {
		if !file.IsDir() {
			filename := file.Name()
			ext := strings.ToLower(filepath.Ext(filename))
			nameWithoutExt := strings.TrimSuffix(filename, ext)
			// Replace underscores with spaces for display
			displayName := strings.ReplaceAll(nameWithoutExt, "_", " ")

			var bgType string
			if imageExtensions[ext] {
				bgType = "image"
			} else if videoExtensions[ext] {
				bgType = "video"
			} else {
				continue // Skip unsupported file types
			}

			backgrounds = append(backgrounds, Background{
				Name:  displayName,
				Value: filename,
				Type:  bgType,
			})
		}
	}
	return backgrounds, nil
}

type Sound struct {
	Name string // Display name
	Path string // File path
	Icon string // Ionicon name
}

func getSounds() []Sound {
	return []Sound{
		{Name: "Fire", Path: "fire.mp3", Icon: "flame"},
		{Name: "Rain", Path: "rain.mp3", Icon: "water"},
		{Name: "Wind", Path: "wind.mp3", Icon: "cloud"},
		{Name: "Forest", Path: "forest.mp3", Icon: "leaf"},
		{Name: "Ocean", Path: "ocean.mp3", Icon: "water-outline"},
		{Name: "Thunder", Path: "thunder.mp3", Icon: "flash"},
		// Add more sounds here as needed
	}
}

func saveBackground(background string) error {
	_, err := db.Exec(
		"INSERT OR REPLACE INTO settings (key, value) VALUES (?, ?)",
		"selectedBackground",
		background,
	)
	return err
}

func getBackground() (string, error) {
	var background string
	err := db.QueryRow("SELECT value FROM settings WHERE key = ?", "selectedBackground").Scan(&background)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return background, err
}

func main() {
	// Initialize database
	if err := initDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Serve the public directory
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	// API endpoint to save background
	http.HandleFunc("/api/save-background", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req map[string]string
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		if err := saveBackground(req["background"]); err != nil {
			http.Error(w, "Failed to save background", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// API endpoint to get background
	http.HandleFunc("/api/get-background", func(w http.ResponseWriter, r *http.Request) {
		background, err := getBackground()
		if err != nil {
			http.Error(w, "Failed to get background", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"background": background})
	})

	// API endpoint to save sound
	http.HandleFunc("/api/save-sound", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req map[string]string
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		if err := saveBackground(req["sound"]); err != nil {
			http.Error(w, "Failed to save sound", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// API endpoint to get sound
	http.HandleFunc("/api/get-sound", func(w http.ResponseWriter, r *http.Request) {
		sound, err := getBackground()
		if err != nil {
			http.Error(w, "Failed to get sound", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"sound": sound})
	})

	// Serve the index.html for the root path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		// Get backgrounds (images and videos)
		backgrounds, err := getBackgrounds()
		if err != nil {
			http.Error(w, "Failed to read backgrounds", http.StatusInternalServerError)
			return
		}

		// Get sounds
		sounds := getSounds()

		// Parse and execute template
		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			http.Error(w, "Failed to parse template", http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"Backgrounds": backgrounds,
			"Sounds":      sounds,
		}

		tmpl.Execute(w, data)
	})

	port := ":8012"
	fmt.Printf("Server running on http://localhost:8012\n")
	log.Fatal(http.ListenAndServe(port, nil))
}

