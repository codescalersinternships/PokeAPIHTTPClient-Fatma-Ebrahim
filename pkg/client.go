// This package implements a HTTP client in Go that consumes the Pokemon APIs.
// It supports two content types: plain text and JSON.
package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cenkalti/backoff/v4"
)

// Logger is used to log messages to logs.log file.
var (
	outfile, _ = os.Create("logs.log")
	logger     = log.New(outfile, "", 0)
)

type Client struct {
	client  *http.Client
	timeout time.Duration
}

type Option func(*Client)

type NameUrl struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Url struct {
	Url string `json:"url"`
}
type EffectEntry struct {
	Effect   string  `json:"effect"`
	Language NameUrl `json:"language"`
}
type EffectChanges struct {
	EffectEntries []EffectEntry `json:"effect_entries"`
	VersionGroup  NameUrl       `json:"version_group"`
}
type FlavorTextEntry struct {
	FlavorText   string  `json:"flavor_text"`
	Language     NameUrl `json:"language"`
	VersionGroup NameUrl `json:"version_group"`
}
type Name struct {
	Name     string  `json:"name"`
	Language NameUrl `json:"language"`
}
type Pokemon struct {
	IsHidden bool    `json:"is_hidden"`
	Pokemon  NameUrl `json:"pokemon"`
	Slot     int     `json:"slot"`
}
type Ability struct {
	Name              string            `json:"name"`
	IsMainSeries      bool              `json:"is_main_series"`
	Generation        map[string]string `json:"generation"`
	EffectEntries     []EffectEntry     `json:"effect_entries"`
	EffectChanges     []EffectChanges   `json:"effect_changes"`
	FlavorTextEntries []FlavorTextEntry `json:"flavor_text_entries"`
	Names             []Name            `json:"names"`
	Pokemon           []Pokemon         `json:"pokemon"`
}

type Description struct {
	Description string  `json:"description"`
	Language    NameUrl `json:"language"`
}

type Characteristics struct {
	Descriptions   []Description `json:"descriptions"`
	GeneModulo     int           `json:"gene_modulo"`
	HighestStat    NameUrl       `json:"highest_stat"`
	Id             int           `json:"id"`
	PossibleValues []int         `json:"possible_values"`
}

type EggGroup struct {
	Id             int       `json:"id"`
	Name           string    `json:"name"`
	Names          []Name    `json:"names"`
	PokemonSpecies []NameUrl `json:"pokemon_species"`
}

type NamedList struct {
	Count    int       `json:"count"`
	Next     string    `json:"next"`
	Previous string    `json:"previous"`
	Results  []NameUrl `json:"results"`
}

type UnnamedList struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []Url  `json:"results"`
}

// NewClient creates a new HTTP client with default options.
func NewClient(options ...Option) *Client {
	client := &Client{
		client:  &http.Client{},
		timeout: 30 * time.Second,
	}

	for _, opt := range options {
		opt(client)
	}

	return client
}

// WithTimeout is a functional option to set the HTTP client timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

func get(url string, duration time.Duration) (*http.Response, error) {

	var response *http.Response
	connection := func() error {
		c := http.Client{Timeout: duration}
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			logger.Println("Failed to create request", err)
			return err
		}

		response, err = c.Do(request)
		logger.Println("Sending request to server")
		if err != nil {
			logger.Println("Failed to send request", err)
			return err
		}
		logger.Println("Request sent successfully")
		return nil
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = 10 * time.Second
	err := backoff.Retry(connection, expBackoff)
	if err != nil {
		logger.Println("Failed to connect to server", err)
		return response, fmt.Errorf("error in server connection")
	}
	logger.Println("Server returned status code", response.StatusCode)

	if response.StatusCode != http.StatusOK {
		return response, fmt.Errorf("status code not OK")
	}

	return response, nil
}

func parseresponse(response *http.Response, structData interface{}) error {
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if err := json.Unmarshal(data, structData); err != nil {
		return fmt.Errorf("error in json unmarshal: %w", err)
	}
	return nil
}

func (c *Client) GetPokemonAbility(id_name string) (*Ability, error) {
	url := "https://pokeapi.co/api/v2/ability/" + id_name
	response, err := get(url, c.timeout)
	if err != nil {
		return nil, err
	}
	ability := &Ability{}
	if err := parseresponse(response, &ability); err != nil {
		return nil, err
	}
	return ability, nil
}

func (c *Client) GetPokemonCharacteristic(id_name string, offset uint, limit uint) (*Characteristics, *NamedList, error) {
	url := "https://pokeapi.co/api/v2/characteristic/"
	if id_name != "" {
		url += id_name
		response, err := get(url, c.timeout)
		if err != nil {
			return nil, nil,err
		}
		characteristics := &Characteristics{}
		if err := parseresponse(response, &characteristics); err != nil {
			return nil, nil,err
		}
		return characteristics, &NamedList{},nil
	}
	url += "?limit=" + strconv.FormatUint(uint64(limit), 10) + "&offset=" + strconv.FormatUint(uint64(offset), 10)

	response, err := get(url, c.timeout)
	if err != nil {
		return nil,nil, err
	}

	result := &NamedList{}
	if err := parseresponse(response, &result); err != nil {
		return nil,nil, err
	}

	return &Characteristics{},result, nil

}

func (c *Client) GetPokemonEggGroup(id_name string, offset uint, limit uint) (interface{}, error) {
	url := "https://pokeapi.co/api/v2/egg-group/"
	//url := "https://pokeapi.co/api/v2/characteristic/"

	if id_name != "" {
		url += id_name
		response, err := get(url, c.timeout)
		if err != nil {
			return nil, err
		}

		egggroup := &EggGroup{}
		if err := parseresponse(response, egggroup); err != nil {
			return nil, err
		}
		return egggroup, nil
	}

	url += "?limit=" + strconv.FormatUint(uint64(limit), 10) + "&offset=" + strconv.FormatUint(uint64(offset), 10)

	response, err := get(url, c.timeout)
	if err != nil {
		return nil, err
	}

	result := &NamedList{}
	if err := parseresponse(response, result); err != nil {
		return nil, err
	}

	return result.Results, nil
}
