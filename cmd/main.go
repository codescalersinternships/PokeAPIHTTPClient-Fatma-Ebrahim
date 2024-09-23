package main

import (
	"fmt"
	"log"
	"time"

	client "github.com/codescalersinternships/PokeAPIHTTPClient-Fatma-Ebrahim/pkg"
)

func main() {
	client := client.NewClient(
		client.WithTimeout(10 * time.Second),
	)

	res, err := client.GetResources("contest-effect", 10, 10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Results[0].Url)

	pokemon, err := client.GetPokemon("1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pokemon.GameIndices[0])

}
