package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

var artists []Artist

func getArtistsData(url string) ([]Artist, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erreur : statut de la réponse %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var artists []Artist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		return nil, err
	}
	return artists, nil
}
func handlerIndex(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "templates/index.html", artists)
}
func handlerArtistDetail(w http.ResponseWriter, r *http.Request) {
	// Extraire l'ID de l'URL manuellement
	idStr := strings.TrimPrefix(r.URL.Path, "/artist/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}
	// Rechercher l'artiste par ID
	for _, artist := range artists {
		if artist.ID == id {
			renderTemplate(w, "templates/artist.html", artist)
			return
		}
	}
	http.Error(w, "Artist not found", http.StatusNotFound)
}
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(filepath.Join(tmpl))
	if err != nil {
		http.Error(w, "Template loading error", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}
func main() {
	url := "https://groupietrackers.herokuapp.com/api/artists"
	var err error
	artists, err = getArtistsData(url)
	if err != nil {
		log.Fatalf("Erreur lors de la récupération des données des artistes: %v", err)
	}
	// Configurer les routes manuellement
	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("/artist/", handlerArtistDetail) // Route pour les détails d'un artiste
	srv := &http.Server{
		Addr:              ":8443",
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	log.Println("Serveur démarré sur http://localhost:8443")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
