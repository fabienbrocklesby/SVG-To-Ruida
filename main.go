package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"svg-to-ruida/converter"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})

	http.HandleFunc("/convert", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		file, _, err := r.FormFile("svg")
		if err != nil {
			http.Error(w, "Failed to read SVG file from form", http.StatusBadRequest)
			return
		}
		defer file.Close()

		svgData, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read SVG data", http.StatusInternalServerError)
			return
		}

		rdData := converter.Convert(svgData)

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=\"output.rd\"")
		w.Write(rdData)
	})

	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
