# PokeAPIHTTPClient-Fatma-Ebrahim

This package implements an HTTP client in Go that consumes the pokemon API https://pokeapi.co/
## Installation

To install the client package, run the following command:

```shell
git clone https://github.com/codescalersinternships/PokeAPIHTTPClient-Fatma-Ebrahim.git
```
To install the needed dependencies:

```shell
go mod download
```

## Usage

Here's an example of how to use the `Client` package:

In a `main.go` file:

```go
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
	fmt.Println(res.Results)

	pokemon, err := client.GetPokemon("1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pokemon.GameIndices[0])

}
```

In the terminal, run the command:

```shell

go run main.go
```

Now you can see the response of the request Pokemon or Resource in the terminal!
