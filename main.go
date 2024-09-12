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
	"time"
)

// Structures des données

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	DatesURL  string   `json:"dates"`
}

// Structure pour encapsuler les données des lieux
type LocationsAPIResponse struct {
	Index []Location `json:"index"`
}

type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

// Structure pour encapsuler les données des dates
type DatesAPIResponse struct {
	Index []Date `json:"index"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// Structure pour encapsuler les données des relations
type RelationsAPIResponse struct {
	Index []Relation `json:"index"`
}

var artists []Artist
var locations []Location
var dates []Date
var relations []Relation

func getAPIData[T any](url string, target *[]T) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erreur : statut de la réponse %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, target)
	if err != nil {
		return err
	}
	return nil
}

func getLocationsData(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erreur : statut de la réponse %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var locationsResponse LocationsAPIResponse
	err = json.Unmarshal(body, &locationsResponse)
	if err != nil {
		return err
	}

	// Assignez les locations récupérées à la variable globale `locations`
	locations = locationsResponse.Index
	return nil
}

func getDatesData(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erreur : statut de la réponse %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var datesResponse DatesAPIResponse
	err = json.Unmarshal(body, &datesResponse)
	if err != nil {
		return err
	}

	// Assignez les dates récupérées à la variable globale `dates`
	dates = datesResponse.Index
	return nil
}

func getRelationsData(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erreur : statut de la réponse %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var relationsResponse RelationsAPIResponse
	err = json.Unmarshal(body, &relationsResponse)
	if err != nil {
		return err
	}

	// Assignez les relations récupérées à la variable globale `relations`
	relations = relationsResponse.Index
	return nil
}

func loadAllData() error {
	// URLs des APIs
	urls := map[string]string{
		"artists":   "https://groupietrackers.herokuapp.com/api/artists",
		"locations": "https://groupietrackers.herokuapp.com/api/locations",
		"dates":     "https://groupietrackers.herokuapp.com/api/dates",
		"relation":  "https://groupietrackers.herokuapp.com/api/relation",
	}

	// Récupération des données des artistes
	if err := getAPIData(urls["artists"], &artists); err != nil {
		return fmt.Errorf("Erreur lors de la récupération des artistes: %v", err)
	}
	// Récupération des données des localisations
	if err := getLocationsData(urls["locations"]); err != nil {
		return fmt.Errorf("Erreur lors de la récupération des localisations: %v", err)
	}

	// Récupération des données des dates
	if err := getDatesData(urls["dates"]); err != nil {
		return fmt.Errorf("Erreur lors de la récupération des dates: %v", err)
	}

	// Récupération des données des relations
	if err := getRelationsData(urls["relation"]); err != nil {
		return fmt.Errorf("Erreur lors de la récupération des relations: %v", err)
	}

	return nil
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFoundHandler(w, r)
		return
	}
	renderTemplate(w, "index.html", artists)
}

func handlerArtistDetail(w http.ResponseWriter, r *http.Request) {
	// Extraire l'ID depuis les paramètres de la requête GET
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "ID d'artiste invalide", http.StatusBadRequest)
		return
	}
	// Rechercher l'artiste par ID
	for _, artist := range artists {
		if artist.ID == id {
			// Optionnel: Récupérer et afficher les détails supplémentaires comme les localisations, dates, etc.
			renderTemplate(w, "artist.html", artist)
			return
		}
	}

	http.Error(w, "Artiste non trouvé", http.StatusNotFound)
}

func handlerConcerts(w http.ResponseWriter, r *http.Request) {
	// Extraire l'ID de l'artiste depuis les paramètres de la requête GET
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "ID d'artiste invalide", http.StatusBadRequest)
		return
	}

	// Rechercher les détails des concerts pour l'artiste en utilisant l'ID
	var artistConcerts Relation
	for _, relation := range relations {
		if relation.ID == id {
			artistConcerts = relation
			break
		}
	}

	// Si aucun concert n'est trouvé pour cet artiste
	if artistConcerts.ID == 0 {
		http.Error(w, "Aucun concert trouvé pour cet artiste", http.StatusNotFound)
		return
	}

	// Rendre le template des concerts avec les données récupérées
	renderTemplate(w, "concerts.html", artistConcerts)
}

// Handler pour les pages non trouvées (Not Found)
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	renderTemplate(w, "error404.html", nil)
}

// Handler pour les erreurs serveur (Internal Server Error)
func internalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "error500.html", nil)
}

// Middleware pour gérer les erreurs et éviter que le serveur plante
func errorHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Println("Captured error:", rec)
				internalServerErrorHandler(w, r)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmplPath := filepath.Join("templates", tmpl)
	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Erreur de chargement du template %s: %v", tmplPath, err)
		http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Printf("Erreur lors de l'exécution du template %s: %v", tmplPath, err)
		http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
	}
}

func main() {
	// Charger toutes les données des APIs
	if err := loadAllData(); err != nil {
		log.Fatalf("Erreur lors du chargement des données: %v", err)
	}

	// Configurer les routes
	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("/artist", handlerArtistDetail)
	http.HandleFunc("/artist/", handlerArtistDetail)
	http.HandleFunc("/concerts", handlerConcerts)
	http.HandleFunc("/404", notFoundHandler)
	http.HandleFunc("/500", internalServerErrorHandler)

	// Configurer et démarrer le serveur
	srv := &http.Server{
		Addr:              ":8440",
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	log.Println("Serveur démarré sur http://localhost:8440")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
