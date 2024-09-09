package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
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

	srv := &http.Server{
		Addr:              ":8443",
		Handler:           http.HandlerFunc(handlerIndex),
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
