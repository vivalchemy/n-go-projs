package main

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

type Mux struct {
	mux *http.ServeMux
}

type Pokemon struct {
	Id             int     `json:"id"`
	Name           string  `json:"name"`
	Height         float32 `json:"height"`
	Weight         float32 `json:"weight"`
	BaseExperience float32 `json:"base_experience"`
	Order          string  `json:"order"`
	SpecialAttack  string  `json:"special_attack"`
}

type Pokemons []Pokemon

func (p Pokemons) Len() int {
	return len(p)
}

func (p *Pokemons) Get() (*Pokemons, error) {
	return p, nil
}

func (p *Pokemons) GetById(id int) (Pokemon, error) {
	for _, pokemon := range *p {
		if pokemon.Id == id {
			return pokemon, nil
		}
	}
	return Pokemon{}, errors.New("Pokemon not found")
}

func (p *Pokemons) Set(pokemon Pokemon) {
	*p = append(*p, pokemon)
}

func (p *Pokemons) Update(pokemon Pokemon) (Pokemon, error) {
	for index, pm := range *p {
		if pm.Id == pokemon.Id {
			(*p)[index] = pokemon
			return pokemon, nil
		}
	}
	return Pokemon{}, errors.New("Pokemon not found")
}

func (p *Pokemons) Delete(id int) (Pokemon, error) {
	for index, pm := range *p {
		if pm.Id == id {
			*p = append((*p)[:index], (*p)[index+1:]...)
			return pm, nil
		}
	}
	return Pokemon{}, errors.New("Pokemon not found")
}

func NewMux() *Mux {
	return &Mux{
		mux: http.NewServeMux(),
	}
}

func (m *Mux) handle(httpMethod, path string, handler http.Handler) {
	m.mux.Handle(strings.Join([]string{httpMethod, path}, " "), handler)
}

func (m *Mux) handleFunc(httpMethod, path string, handler func(http.ResponseWriter, *http.Request)) {
	m.mux.HandleFunc(strings.Join([]string{httpMethod, path}, " "), handler)
}

func (m *Mux) HandleGet(path string, handler http.Handler) {
	m.handle(http.MethodGet, path, handler)
}

func (m *Mux) HandlePost(path string, handler http.Handler) {
	m.handle(http.MethodPost, path, handler)
}

func (m *Mux) HandleUpdate(path string, handler http.Handler) {
	m.handle(http.MethodPut, path, handler)
}

func (m *Mux) HandleDelete(path string, handler http.Handler) {
	m.handle(http.MethodDelete, path, handler)
}

func (m *Mux) StaticFileServer(path string, dir string) {
	m.HandleGet(path, http.FileServer(http.Dir(dir)))
}

func (m *Mux) HandleFuncGet(path string, handler func(http.ResponseWriter, *http.Request)) {
	m.handleFunc(http.MethodGet, path, handler)
}

func (m *Mux) HandleFuncPost(path string, handler func(http.ResponseWriter, *http.Request)) {
	m.handleFunc(http.MethodPost, path, handler)
}

func (m *Mux) HandleFuncUpdate(path string, handler func(http.ResponseWriter, *http.Request)) {
	m.handleFunc(http.MethodPut, path, handler)
}

func (m *Mux) HandleFuncDelete(path string, handler func(http.ResponseWriter, *http.Request)) {
	m.handleFunc(http.MethodDelete, path, handler)
}

var pokemons Pokemons

func getAllPokemons(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pokemons)
}

func getPokemonById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Invalid request"}`))
		return
	}

	pokemon, err := pokemons.GetById(id)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Pokemon not found"}`))
		return
	}

	json.NewEncoder(w).Encode(pokemon)
}

func addPokemon(w http.ResponseWriter, r *http.Request) {
	var pokemon Pokemon
	err := json.NewDecoder(r.Body).Decode(&pokemon)
	// log.Print("pokemon", pokemon)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Invalid request"}`))
		return
	}
	pokemon.Id = rand.Int()
	pokemons.Set(pokemon)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pokemon)
}

func updatePokemon(w http.ResponseWriter, r *http.Request) {
	var pokemon Pokemon
	err := json.NewDecoder(r.Body).Decode(&pokemon)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Invalid request"}`))
		return
	}

	pokemons.Update(pokemon)
}

func deletePokemon(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Invalid request"}`))
		return
	}
	pokemon, err := pokemons.Delete(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Pokemon not found"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pokemon)
}

func main() {
	pokemons = Pokemons{
		{
			Id:             1,
			Name:           "Bulbasaur",
			Height:         0.7,
			Weight:         6.9,
			BaseExperience: 64,
			Order:          "1",
			SpecialAttack:  "Overgrow",
		},
		{
			Id:             2,
			Name:           "Ivysaur",
			Height:         1.0,
			Weight:         13.0,
			BaseExperience: 142,
			Order:          "2",
			SpecialAttack:  "Overgrow",
		},
	}
	router := NewMux()

	router.HandleFuncGet("/", getAllPokemons)

	router.HandleFuncGet("/{id}", getPokemonById)

	router.HandleFuncPost("/", addPokemon)

	router.HandleFuncUpdate("/", updatePokemon)

	router.HandleFuncDelete("/{id}", deletePokemon)

	log.Print("Running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", router.mux))
}
