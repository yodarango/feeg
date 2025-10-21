package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Sound represents an ambient sound
type Sound struct {
	Name string // Display name
	Path string // File path
	Icon string // Ionicon name
}

// Background represents a background image or video
type Background struct {
	Name      string // Display name
	Path      string // File path (filename with extension)
	Thumbnail string // Thumbnail image path (.webp)
	Icon      string // Ionicon name
	Type      string // "image" or "video"
}

// GetSounds returns the list of available sounds
func GetSounds() []Sound {
	return []Sound{
		{Name: "Fire", Path: "fire.mp3", Icon: "flame"},
		{Name: "Rain", Path: "rain.mp3", Icon: "water"},
		{Name: "Wind", Path: "wind.mp3", Icon: "cloud"},
		{Name: "Forest", Path: "forest.mp3", Icon: "leaf"},
		{Name: "Ocean", Path: "ocean_waves.mp3", Icon: "water-outline"},
		{Name: "Thunder", Path: "thunderstorm.mp3", Icon: "flash"},
		{Name: "Coffee Shop", Path: "coffeeshop.mp3", Icon: "cafe"},
		{Name: "Coffee Shop Loud", Path: "coffee_shop_loud.mp3", Icon: "cafe"},
		{Name: "Cricket", Path: "cricket.mp3", Icon: "bug"},
		{Name: "Heavy Rain", Path: "heavy_rain.mp3", Icon: "water"},
		{Name: "Night Field", Path: "night_field.mp3", Icon: "moon"},
		{Name: "Owl", Path: "owl.mp3", Icon: "egg"},
		{Name: "Rain", Path: "rain_a.mp3", Icon: "water"},
		{Name: "River", Path: "river.mp3", Icon: "water"},
		{Name: "Siren", Path: "siren.mp3", Icon: "alert"},
		{Name: "Traffic Outside", Path: "traffic_outside.mp3", Icon: "car"},
		{Name: "Wolf", Path: "wolf.mp3", Icon: "paw"},
	}
}

// GetBackgrounds returns the list of available backgrounds
func GetBackgrounds() []Background {
	files, err := os.ReadDir("public/bkgs")
	if err != nil {
		return []Background{}
	}

	var backgrounds []Background
	seenNames := make(map[string]bool)

	for _, file := range files {
		if !file.IsDir() {
			filename := file.Name()
			ext := strings.ToLower(filepath.Ext(filename))
			nameWithoutExt := strings.TrimSuffix(filename, ext)

			// Skip if we've already processed this background
			if seenNames[nameWithoutExt] {
				continue
			}
			seenNames[nameWithoutExt] = true

			// Replace underscores with spaces for display
			displayName := strings.ReplaceAll(nameWithoutExt, "_", " ")

			backgrounds = append(backgrounds, Background{
				Name:      displayName,
				Path:      nameWithoutExt + ".mov",
				Thumbnail: nameWithoutExt + ".webp",
				Icon:      "play-circle",
				Type:      "video",
			})
		}
	}
	return backgrounds
}



func main() {
	// Serve the public directory
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

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

