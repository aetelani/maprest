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

type CountryId string
type ContinentId string
type CityId string

type Continent struct {
	Name ContinentId `json:"name,omitempty"`
}

type Country struct {
	Name      CountryId  `json:"name,omitempty"`
	Continent *Continent `json:"continent,omitempty"`
}

func (country Country) getContinent() Continent { return dataBase.getContinent(country.Continent.Name) }

type City struct {
	Name    CityId   `json:"name,omitempty"`
	Country *Country `json:"country,omitempty"`
}

func (city City) getCountry() Country { return dataBase.getCountry(city.Country.Name) }

func (city City) getContinent() Continent { return city.getCountry().getContinent() }

type DataBaseI interface {
	getCity(CityId) City
	getContinent(ContinentId) Continent
	getCountry(CountryId) Country
	getCities() *Cities
	getContinents() *Continents
	getCountries() *Countries
	addCity(City)
	addContinent(Continent)
	addCountry(Country)
	deleteCity(CityId)
	deleteContinent(ContinentId)
	deleteCountry(CountryId)
}

type Cities map[CityId]City
type Continents map[ContinentId]Continent
type Countries map[CountryId]Country

type MemDataBaseImp struct {
	info       string `default:"In-Memory DataBase implementation"`
	continents Continents
	cities     Cities
	countries  Countries
}

func (db MemDataBaseImp) getCity(name CityId) City                { return db.cities[name] }
func (db MemDataBaseImp) getContinent(name ContinentId) Continent { return db.continents[name] }
func (db MemDataBaseImp) getCountry(name CountryId) Country       { return db.countries[name] }
func (db *MemDataBaseImp) getCities() *Cities                     { return &db.cities }
func (db *MemDataBaseImp) getContinents() *Continents             { return &db.continents }
func (db *MemDataBaseImp) getCountries() *Countries               { return &db.countries }

func (db *MemDataBaseImp) addCity(city City)                     { db.cities[city.Name] = city }
func (db *MemDataBaseImp) addContinent(continent Continent)      { db.continents[continent.Name] = continent }
func (db *MemDataBaseImp) addCountry(country Country)            { db.countries[country.Name] = country }
func (db *MemDataBaseImp) deleteCity(city CityId)                { delete(db.cities, city) }
func (db *MemDataBaseImp) deleteContinent(continent ContinentId) { delete(db.continents, continent) }
func (db *MemDataBaseImp) deleteCountry(country CountryId)       { delete(db.countries, country) }

var dataBase DataBaseI

func createCity(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(os.Stdout, "%s %s\n", req.Method, req.RequestURI)
	var city City
	err := json.NewDecoder(req.Body).Decode(&city)
	log.Println(err)
	dataBase.addCity(city)
	json.NewEncoder(os.Stdout).Encode(city)
}

func getCity(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(os.Stdout, "%s %s\n", req.Method, req.RequestURI)
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	id := CityId(params[NAMEID])
	city := dataBase.getCity(id)
	err := json.NewEncoder(w).Encode(city)
	log.Println(err)
}
func getCities(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stdout, "%s %s\n", r.Method, r.RequestURI)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(dataBase.getCities())
	if err != nil {
		log.Fatal(err)
	}
}
func deleteCity(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	name := CityId(params[NAMEID])
	fmt.Fprintf(os.Stdout, "%s %s %s\n", req.Method, req.RequestURI, name)
	dataBase.deleteCity(name)
}
func createCountry(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(os.Stdout, "%s %s\n", req.Method, req.RequestURI)
	var country Country
	err := json.NewDecoder(req.Body).Decode(&country)
	log.Println(err)
	dataBase.addCountry(country)
	json.NewEncoder(os.Stdout).Encode(country)
}
func getCountry(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(os.Stdout, "%s %s\n", req.Method, req.RequestURI)
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	id := CountryId(params[NAMEID])
	city := dataBase.getCountry(id)
	err := json.NewEncoder(w).Encode(city)
	log.Println(err)
}
func deleteCountry(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(os.Stdout, "%s %s\n", req.Method, req.RequestURI)
	params := mux.Vars(req)
	id := CountryId(params[NAMEID])
	dataBase.deleteCountry(id)
}
func createContinent(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(os.Stdout, "%s %s\n", req.Method, req.RequestURI)
	var continent Continent
	err := json.NewDecoder(req.Body).Decode(&continent)
	log.Println(err)
	dataBase.addContinent(continent)
	json.NewEncoder(os.Stdout).Encode(continent)
}
func getContinent(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(os.Stdout, "%s %s\n", req.Method, req.RequestURI)
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	id := ContinentId(params[NAMEID])
	continent := dataBase.getContinent(id)
	err := json.NewEncoder(w).Encode(continent)
	log.Println(err)
}
func deleteContinent(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(os.Stdout, "%s %s\n", req.Method, req.RequestURI)
	params := mux.Vars(req)
	id := ContinentId(params[NAMEID])
	dataBase.deleteContinent(id)
}
func searchCities(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(os.Stdout, "%s %s\n", req.Method, req.RequestURI)
	results := make(Cities)
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)

	continents := params[CONTINENTSID]
	continentList := strings.Split(continents, ",")
	log.Println("params", params)
	countries := params[COUNTRIESID]
	countriesList := strings.Split(countries, ",")
	log.Println(countriesList)
	for _, v := range *dataBase.getCities() {
		log.Println(v)
		for _, n := range countriesList {
			if v.getCountry().Name == CountryId(n) {
				results[v.Name] = v
			}
		}
		for _, n := range continentList {
			if v.getContinent().Name == ContinentId(n) {
				results[v.Name] = v
			}
		}
	}
	err := json.NewEncoder(w).Encode(results)
	log.Println(err)
	log.Println(results)
}
func getCountriesByContinent(w http.ResponseWriter, req *http.Request) {
	results := make(Countries)
	fmt.Fprintf(os.Stdout, "%s %s", req.Method, req.RequestURI)
	params := mux.Vars(req)
	name := ContinentId(params[NAMEID])
	for _, v := range *dataBase.getCountries() {
		continent := v.getContinent()
		if continent.Name == name {
			results[v.Name] = v
		}
	}
	err := json.NewEncoder(w).Encode(results)
	log.Println(err)
}

const (
	// Paths
	CITIES     string = "/cities"
	COUNTRIES  string = "/countries"
	CONTINENTS string = "/continents"
	SEARCH     string = "/search"

	// Serach terms
	COUNTRIESBYCONTINENT string = COUNTRIES + CONTINENTS
	SEARCHCITIES         string = SEARCH + CITIES

	// Keywords
	NAMEID         string = "name"
	NAMEPART       string = "/{name}"
	CONTINENTSID   string = "continents"
	CONTINENTSPART string = "{continents}"
	COUNTRIESID    string = "countries"
	COUNTRIESPART  string = "{countries}"

	LISTENINGPORT string = ":8181"
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

	dataBase = NewDataBase(getMemDataBaseV1)

	router := mux.NewRouter()
	router.HandleFunc(CITIES, createCity).Methods(http.MethodPost)
	router.HandleFunc(CITIES, getCities)
	router.HandleFunc(CITIES+NAMEPART, getCity)
	router.HandleFunc(CITIES+NAMEPART, deleteCity).Methods(http.MethodDelete)

	router.HandleFunc(COUNTRIES, createCountry).Methods(http.MethodPost)
	router.HandleFunc(COUNTRIES+NAMEPART, getCountry).Methods(http.MethodGet)
	router.HandleFunc(COUNTRIES+NAMEPART, deleteCountry).Methods(http.MethodDelete)

	router.HandleFunc(CONTINENTS, createContinent).Methods(http.MethodPost)
	router.HandleFunc(CONTINENTS+NAMEPART, getContinent)
	router.HandleFunc(CONTINENTS+NAMEPART, deleteContinent).Methods(http.MethodDelete)

	router.HandleFunc(COUNTRIESBYCONTINENT+NAMEPART, getCountriesByContinent)
	router.HandleFunc(SEARCHCITIES, searchCities).Queries(CONTINENTSID, CONTINENTSPART, COUNTRIESID, COUNTRIESPART)

	// Index page
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<html><ul>")
		fmt.Fprintf(w, "<li><a href=%s>CITIES</a>\n", CITIES)
		fmt.Fprintf(w, "<li><a href=%s>COUNTRIES</a>\n", COUNTRIES)
		fmt.Fprintf(w, "<li><a href=%s>CONTINENTS</a>\n", CONTINENTS)
		fmt.Fprintln(w, "</ul></html>")
	})

	// Listen and serve
	log.Println("Listening port", LISTENINGPORT)
	log.Fatal(http.ListenAndServe(LISTENINGPORT, router))
}
