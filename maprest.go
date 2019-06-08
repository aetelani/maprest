package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

type Continent struct {
	Name string `json:"name,omitempty"`
}

type Country struct {
	Name      string     `json:"name,omitempty"`
	Continent *Continent `json:"continent,omitempty"`
}

type City struct {
	Name    string   `json:"name,omitempty"`
	Country *Country `json:"country,omitempty"`
}

type DataBaseI interface {
	getCity(string) City
	getContinent(string) Continent
	getCountry(string) Country
	getCities() *Cities
	getContinents() *Continents
	getCountries() *Countries
	addCity(City)
	addContinent(Continent)
	addCountry(Country)
	deleteCity(string)
	deleteContinent(string)
	deleteCountry(string)
}

type Cities map[string]City
type Continents map[string]Continent
type Countries map[string]Country

type MemDataBaseImp struct {
	info       string `default:"InMemory DataBase implementation"`
	continents Continents
	cities     Cities
	countries  Countries
}

func (db MemDataBaseImp) getCity(name string) City           { return db.cities[name] }
func (db MemDataBaseImp) getContinent(name string) Continent { return db.continents[name] }
func (db MemDataBaseImp) getCountry(name string) Country     { return db.countries[name] }
func (db *MemDataBaseImp) getCities() *Cities                { return &db.cities }
func (db *MemDataBaseImp) getContinents() *Continents        { return &db.continents }
func (db *MemDataBaseImp) getCountries() *Countries          { return &db.countries }

func (db *MemDataBaseImp) addCity(city City)                { db.cities[city.Name] = city }
func (db *MemDataBaseImp) addContinent(continent Continent) { db.continents[continent.Name] = continent }
func (db *MemDataBaseImp) addCountry(country Country)       { db.countries[country.Name] = country }
func (db *MemDataBaseImp) deleteCity(city string)           { delete(db.cities, city) }
func (db *MemDataBaseImp) deleteContinent(continent string) { delete(db.continents, continent) }
func (db *MemDataBaseImp) deleteCountry(country string)     { delete(db.countries, country) }

var db DataBaseI

func createCity(w http.ResponseWriter, req *http.Request) {
	var city City
	err := json.NewDecoder(req.Body).Decode(&city)
	log.Println(err)
	db.addCity(city)
	log.Println("city", City{Name: "Test", Country: &Country{Name: "fi"}})
	json.NewEncoder(os.Stdout).Encode(city)
}

func getCity(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(os.Stdout, "%s %s", req.Method, req.RequestURI)
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	name := params[NAMEID]
	city := db.getCity(name)
	err := json.NewEncoder(w).Encode(city)
	log.Println(err)
}
func getCities(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stdout, "%s %s", r.Method, r.RequestURI)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(db.getCities())
	if err != nil {
		log.Fatal(err)
	}
}
func deleteCity(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(os.Stdout, "%s %s", req.Method, req.RequestURI)
	params := mux.Vars(req)
	name := params[NAMEID]
	db.deleteCity(name)
}
func createCountry(w http.ResponseWriter, req *http.Request) {
	var country Country
	err := json.NewDecoder(req.Body).Decode(&country)
	log.Println(err)
	db.addCountry(country)
	json.NewEncoder(os.Stdout).Encode(country)
}
func getCountry(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(os.Stdout, "%s %s", req.Method, req.RequestURI)
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	name := params[NAMEID]
	city := db.getCountry(name)
	err := json.NewEncoder(w).Encode(city)
	log.Println(err)
}
func deleteCountry(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(os.Stdout, "%s %s", req.Method, req.RequestURI)
	params := mux.Vars(req)
	name := params[NAMEID]
	db.deleteCountry(name)
}
func createContinent(w http.ResponseWriter, req *http.Request) {
	var continent Continent
	err := json.NewDecoder(req.Body).Decode(&continent)
	log.Println(err)
	db.addContinent(continent)
	json.NewEncoder(os.Stdout).Encode(continent)
}
func getContinent(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(os.Stdout, "%s %s", req.Method, req.RequestURI)
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	name := params[NAMEID]
	continent := db.getContinent(name)
	err := json.NewEncoder(w).Encode(continent)
	log.Println(err)
}
func deleteContinent(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(os.Stdout, "%s %s", req.Method, req.RequestURI)
	params := mux.Vars(req)
	name := params[NAMEID]
	db.deleteContinent(name)
}
func searchCities(w http.ResponseWriter, req *http.Request) {
	results := make(Cities)
	fmt.Fprintf(os.Stdout, "SE %s %s", req.Method, req.RequestURI)
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)

	continents := params[CONTINENTSID]
	continentList := strings.Split(continents, ",")
	log.Println("params", params)
	countries := params[COUNTRIESID]
	countriesList := strings.Split(countries, ",")
	log.Println(countriesList)
	for _, v := range *db.getCities() {
		log.Println(v)
		for _, n := range countriesList {
			if v.Country != nil && v.Country.Name == n {
				results[v.Name] = v
			}
		}
		for _, n := range continentList {
			if v.Country != nil && v.Country.Continent != nil && v.Country.Continent.Name == n {
				results[v.Name] = v
			}
		}
	}
	err := json.NewEncoder(w).Encode(results)
	log.Println(err)
	log.Println(results)
}
func getCountriesByContinent(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(os.Stdout, "%s %s", req.Method, req.RequestURI)
	params := mux.Vars(req)
	name := params[NAMEID]
	for _, v := range *db.getCountries() {
		if v.Continent.Name == name {
			log.Println(v)
		}
	}
}

const (
	CITIES     string = "/cities"
	COUNTRIES  string = "/countries"
	CONTINENTS string = "/continents"

	// Serach terms
	SEARCH               string = "/search"
	COUNTRIESBYCONTINENT string = COUNTRIES + CONTINENTS
	SEARCHCITIES         string = SEARCH + CITIES

	// Keywords
	NAMEID         string = "name"
	NAMEPART       string = "/{name}"
	CONTINENTSID   string = "continents"
	CONTINENTSPART string = "{continents}"
	COUNTRIESID    string = "countries"
	COUNTRIESPART  string = "{countries}"

	// Literals
	POST   string = "POST"
	GET    string = "GET"
	DELETE string = "DELETE"
)

func getMemDataBaseV1() DataBaseI {
	return &MemDataBaseImp{
		continents: make(Continents),
		cities:     make(Cities),
		countries:  make(Countries),
	}
}

func NewDataBase(db func() DataBaseI) DataBaseI {
	return db()
}

func main() {

	db = NewDataBase(getMemDataBaseV1)

	log.Println("Starting...")
	router := mux.NewRouter()
	router.HandleFunc(CITIES, createCity).Methods(POST)
	router.HandleFunc(CITIES, getCities).Methods(GET)
	router.HandleFunc(CITIES+NAMEPART, getCity).Methods(GET)
	router.HandleFunc(CITIES+NAMEPART, deleteCity).Methods(DELETE)
	router.HandleFunc(COUNTRIES, createCountry).Methods(POST)
	router.HandleFunc(COUNTRIES+NAMEPART, getCountry).Methods(GET)
	router.HandleFunc(COUNTRIES+NAMEPART, deleteCountry).Methods(DELETE)
	router.HandleFunc(CONTINENTS, createContinent).Methods(POST)
	router.HandleFunc(CONTINENTS+NAMEPART, getContinent).Methods(GET)
	router.HandleFunc(CONTINENTS+NAMEPART, deleteContinent).Methods(DELETE)

	router.HandleFunc(COUNTRIESBYCONTINENT+NAMEPART, getCountriesByContinent)
	router.HandleFunc(SEARCHCITIES, searchCities).Queries(CONTINENTSID, CONTINENTSPART, COUNTRIESID, COUNTRIESPART)

	// Index page
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<html><ul>")
		fmt.Fprintf(w, "<li><a href=%s>CITIES</a>", CITIES)
		fmt.Fprintf(w, "<li><a href=%s>COUNTRIES</a>", COUNTRIES)
		fmt.Fprintf(w, "<li><a href=%s>CONTINENTS</a>", CONTINENTS)
		fmt.Fprintln(w, "</ul></html>")
	})

	// Listen and serve
	log.Fatal(http.ListenAndServe(":8181", router))
}
