package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/png"
)

type Page struct {
	Title    string
	QRCode   string
	ErrorMsg string
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/generator/", viewCodeHandler)

	// Serve static files (CSS, JS, etc.)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Printf("Listening on port 63342\n")
	log.Fatal(http.ListenAndServe(":63342", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{Title: "QR-code generator"}
	renderTemplate(w, "generator.html", p)
}

func viewCodeHandler(w http.ResponseWriter, r *http.Request) {
	dataString := r.FormValue("dataString")

	if dataString == "" {
		p := Page{
			Title:    "QR-code generator",
			ErrorMsg: "Please enter a valid string.",
		}
		renderTemplate(w, "generator.html", p)
		return
	}

	qrCode, err := qr.Encode(dataString, qr.L, qr.Auto)
	if err != nil {
		log.Printf("Error generating QR code: %v", err)
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		return
	}

	qrCode, err = barcode.Scale(qrCode, 512, 512)
	if err != nil {
		log.Printf("Error scaling QR code: %v", err)
		http.Error(w, "Failed to scale QR code", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	if err := png.Encode(w, qrCode); err != nil {
		log.Printf("Error encoding QR code to PNG: %v", err)
		http.Error(w, "Failed to encode QR code", http.StatusInternalServerError)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p Page) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, p); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}
