package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Sound represents an ambient sound
type Sound struct {
	Name string // Display name
	Path string // File path
	Icon string // Ionicon name
}

// Background represents a background image or video
type Background struct {
	Name string // Display name
	Path string // File path (filename with extension)
	Icon string // Ionicon name
	Type string // "image" or "video"
}

// GetSounds returns the list of available sounds
func GetSounds() []Sound {
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

// GetBackgrounds returns the list of available backgrounds
func GetBackgrounds() []Background {
	return []Background{
		{Name: "Vermont Cozziness", Path: "autum_cozy_vermont_bethroom.mov", Icon: "image", Type: "video"},
		{Name: "Milan Lib", Path: "autum_eropean_library.mov", Icon: "image", Type: "video"},
		{Name: "Boring street", Path: "autum_rain_sad_european_street.mov", Icon: "image", Type: "video"},
		{Name: "River delight", Path: "autum_river_fire.mov", Icon: "image", Type: "video"},
		{Name: "The Porch", Path: "rain_fire_summer_morning.mov", Icon: "image", Type: "video"},
	}
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
		backgrounds := GetBackgrounds()

		// Get sounds
		sounds := GetSounds()

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

