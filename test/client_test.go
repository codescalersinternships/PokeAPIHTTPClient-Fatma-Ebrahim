package test

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
	"time"

	client "github.com/codescalersinternships/PokeAPIHTTPClient-Fatma-Ebrahim/pkg"
)

func TestGetPokemon(t *testing.T) {
	t.Run("Test GetPokemon function", func(t *testing.T) {
		c := client.NewClient(
			client.WithTimeout(10 * time.Second),
		)
		pokemon, err := c.GetPokemon("1")
		if err != nil {
			log.Fatal(err)
		}
		if pokemon.Id != 1 {
			t.Errorf("expected 1 got %v", pokemon.Id)
		}
		if pokemon.Name != "bulbasaur" {
			t.Errorf("expected bulbasaur got %v", pokemon.Name)
		}

		if pokemon.Abilities[0].Ability.Url != "https://pokeapi.co/api/v2/ability/65/" {
			t.Errorf("expected %q got %q", "https://pokeapi.co/api/v2/ability/65/", pokemon.Abilities[0].Ability.Url)
		}
		if pokemon.Cries.Latest != "https://raw.githubusercontent.com/PokeAPI/cries/main/cries/pokemon/latest/1.ogg" {
			t.Errorf("expected %q got %q", "https://raw.githubusercontent.com/PokeAPI/cries/main/cries/pokemon/latest/1.ogg", pokemon.Cries.Latest)
		}
		if pokemon.Weight != 69 {
			t.Errorf("expected 69 got %v", pokemon.Weight)
		}

		if pokemon.Types[0].Type.Url != "https://pokeapi.co/api/v2/type/12/" {
			t.Errorf("expected %q got %q", "https://pokeapi.co/api/v2/type/12/", pokemon.Types[0].Type.Url)
		}
		if pokemon.Sprites.BackDefault != "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/back/1.png" {
			t.Errorf("expected %q got %q", "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/back/1.png", pokemon.Sprites.BackDefault)
		}
		if pokemon.Sprites.Other.DreamWorld.FrontDefault != "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/dream-world/1.svg" {
			t.Errorf("expected %q got %q", "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/dream-world/1.svg", pokemon.Sprites.Other.DreamWorld.FrontDefault)
		}
	})

}

func TestGetResources(t *testing.T) {
	t.Run("Test GetResources function", func(t *testing.T) {
		c := client.NewClient(
			client.WithTimeout(10 * time.Second),
		)
		got, err := c.GetResources("pokemon", 3, 0)
		if err != nil {
			log.Fatal(err)
		}

		data := []byte(`{
	  "count": 1302,
	  "next": "https://pokeapi.co/api/v2/pokemon/?offset=3&limit=3",
	  "previous": null,
	  "results": [
		{
		  "name": "bulbasaur",
		  "url": "https://pokeapi.co/api/v2/pokemon/1/"
		},
		{
		  "name": "ivysaur",
		  "url": "https://pokeapi.co/api/v2/pokemon/2/"
		},
		{
		  "name": "venusaur",
		  "url": "https://pokeapi.co/api/v2/pokemon/3/"
		}
	  ]
	}`)

		want := &client.Resources{}
		err = json.Unmarshal(data, want)
		if err != nil {
			log.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v got %v", want, got)
		}

	})

}
