package main

import (
	"fmt"
	"log"
	"time"

	client "github.com/codescalersinternships/PokeAPIHTTPClient-Fatma-Ebrahim/pkg"
)

type EggGroup = client.EggGroup
type NameUrl = client.NameUrl
type Url = client.Url

func main() {
	client := client.NewClient(
		client.WithTimeout(10 * time.Second),
	)

	ability, err := client.GetPokemonAbility("1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ability.Pokemon[0])
	char,list, err := client.GetPokemonCharacteristic("1", 10, 10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("char: ",char.PossibleValues)
	fmt.Println("list: ",list.Results)

	eggGroup, err := client.GetPokemonEggGroup("", 10, 10)
	if err != nil {
		log.Fatal(err)
	}

	// Check if it's a specific egg group
	eg, ok := eggGroup.(*EggGroup)
	if ok {
		fmt.Println(eg.PokemonSpecies[0])
	} else {
		list := eggGroup.([]NameUrl)
		fmt.Println(list)
	}

}
