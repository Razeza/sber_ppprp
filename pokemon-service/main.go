package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"strings"
)

type PokemonAPIResourceList struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`
}

func GetPokemon(url string) Pokemon {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var pokemon Pokemon
	err = json.NewDecoder(resp.Body).Decode(&pokemon)
	if err != nil {
		panic(err)
	}
	return pokemon
}

func GetPokemons() []Pokemon {
	const pokemonEndpoint string = "https://pokeapi.co/api/v2/pokemon/"
	resp, err := http.Get(pokemonEndpoint)
	if err != nil {
		panic(err)
	}

	var pokemonsEndpoint PokemonAPIResourceList
	err = json.NewDecoder(resp.Body).Decode(&pokemonsEndpoint)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	pokemonsCnt := 20
	pokemons := make([]Pokemon, 0, pokemonsCnt)

out:
	for {
		for _, urlMeta := range pokemonsEndpoint.Results {
			pokemon := GetPokemon(urlMeta.Url)
			pokemons = append(pokemons, pokemon)

			if len(pokemons) == pokemonsCnt {
				break out
			}
		}

		if pokemonsEndpoint.Next == "" {
			break
		}

		resp, err = http.Get(pokemonsEndpoint.Next)
		if err != nil {
			panic(err)
		}

		err = json.NewDecoder(resp.Body).Decode(&pokemonsEndpoint)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
	}

	return pokemons
}

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func getEnv(key string, defaultVal string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultVal
}

func PokemonsMux() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/get-first", func(writer http.ResponseWriter, request *http.Request) {
		pokemons := GetPokemons()

		var sb strings.Builder
		sb.WriteString("Первый покемон - ")
		sb.WriteString(pokemons[0].Name)
		sb.WriteString("\n")

		writer.Write([]byte(sb.String()))
	})

	return mux
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/pokemons/", http.StripPrefix("/pokemons", PokemonsMux()))

	host := getEnv("HOST", "")
	port := getEnv("PORT", "8083")
	http.ListenAndServe(host+":"+port, mux)
}
