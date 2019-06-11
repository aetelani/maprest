package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	//cmd := exec.Command("maprest")
	//cmd.Start()
	//time.Sleep(1000 * time.Millisecond)
	//defer cmd.Process.Kill()

	url := "http://localhost:8181"
	client := &http.Client{}

	city := func(ci, co string) City { return City{Name: CityId(ci), Country: &Country{Name: CountryId(co)}} }
	var citiesAdded int = 0
	postCity := func(c City) {
		body := &bytes.Buffer{}
		json.NewEncoder(body).Encode(c)
		req, _ := http.NewRequest(http.MethodPost, url+CITIES, body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := client.Do(req)
		citiesAdded += 1
		defer resp.Body.Close()
	}
	postCity(city("Tampere", "fi"))
	postCity(city("Turku", "fi"))
	postCity(city("Rovaniemi", "fi"))
	postCity(city("Itä-Berliini", "gr"))
	postCity(city("Länsi-Berliini", "gr"))
	postCity(city("Boston", "us"))

	ciGetCities := func() (count int) {
		var results Cities
		body := &bytes.Buffer{}
		req, _ := http.NewRequest(http.MethodGet, url+CITIES, body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := client.Do(req)
		//stringbody, _ := ioutil.ReadAll(resp.Body)
		//t.Logf("Body %s", stringbody)
		//log.Println("COMES HERRTTT")
		err := json.NewDecoder(resp.Body).Decode(&results)
		log.Println(len(results), results)
		if err != nil {
			t.FailNow()
		}
		defer resp.Body.Close()
		count = len(results)
		return
	}

	ciGetCity := func(name string) (city City) {
		body := &bytes.Buffer{}
		req, _ := http.NewRequest(http.MethodGet, url+CITIES+"/"+name, body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := client.Do(req)
		//stringbody, _ := ioutil.ReadAll(resp.Body)
		//t.Logf("Body %s", stringbody)
		//log.Println("COMES HERRTTT")
		err := json.NewDecoder(resp.Body).Decode(&city)
		log.Println("Got city:", city)
		if err != nil {
			t.FailNow()
		}
		defer resp.Body.Close()
		return
	}
	ciGetCity("Tampere")

	// Added matches to readed count
	if ciGetCities() != citiesAdded {
		t.Errorf("Failed to count cities %d", citiesAdded)
		t.FailNow()
	}

	ciDeleteCity := func() {
		var result City
		body := &bytes.Buffer{}
		req, _ := http.NewRequest(http.MethodDelete, url+CITIES+"/Tampere", body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := client.Do(req)
		//stringbody, _ := ioutil.ReadAll(resp.Body)
		//t.Logf("Body %s", stringbody)
		//log.Println("COMES HERRTTT")
		err := json.NewDecoder(resp.Body).Decode(&result)
		log.Println("Delete city:", result)
		if err != nil {
			t.FailNow()
		}
		defer resp.Body.Close()
		return
	}
	ciDeleteCity()
	res := ciGetCity("Tampere")
	log.Println("After del:", res)

	postCountry := func(cs ...Country) int {
		for _, c := range cs {
			body := &bytes.Buffer{}
			json.NewEncoder(body).Encode(c)
			req, _ := http.NewRequest(http.MethodPost, url+COUNTRIES, body)
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err == nil {
				defer resp.Body.Close()
			} else {
				log.Fatalln("Failed to post country: ", resp.Status)
			}
		}
		return len(cs)
	}

	if postCountry(Country{Name: "us"}, Country{Name: "fi"}, Country{Name: "gr"}) == 0 {
		t.Failed()
	} else {
		t.Error("Success counting countries:")
	}
	ciGetCountries := func(name string) (count int) {
		body := &bytes.Buffer{}
		req, _ := http.NewRequest(http.MethodGet, url+COUNTRIES+NAMEPART+name, body)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := client.Do(req)
		result := make(Countries)
		log.Println(resp.Body)
		json.NewDecoder(resp.Body).Decode(&result)
		defer resp.Body.Close()
		log.Println(len(result), result)
		//log.Println(dataBase.getCountries())
		count = len(result)
		return
	}
	ciGetCountries("/Tampere")
}
