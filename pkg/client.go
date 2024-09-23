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
type Ability struct {
	Ability  NameUrl `json:"ability"`
	IsHidden bool    `json:"is_hidden"`
	Slot     int     `json:"slot"`
}

type Cries struct {
	Latest string `json:"latest"`
	Legacy string `json:"legacy"`
}

type GameIndex struct {
	GameIndex int     `json:"game_index"`
	Version   NameUrl `json:"version"`
}

type VersionDetail struct {
	Rarity  int     `json:"rarity"`
	Version NameUrl `json:"version"`
}

type HeldItem struct {
	Item          NameUrl         `json:"item"`
	VersionDetail []VersionDetail `json:"version_details"`
}

type VersionGroupDetail struct {
	LevelLearnedAt  int     `json:"level_learned_at"`
	MoveLearnMethod NameUrl `json:"move_learn_method"`
	VersionGroup    NameUrl `json:"version_group"`
}

type Move struct {
	Move                NameUrl              `json:"move"`
	VersionGroupDetails []VersionGroupDetail `json:"version_group_details"`
}
type DreamWorld struct {
	FrontDefault string `json:"front_default"`
	FrontFemale  string `json:"front_female"`
}

type Home struct {
	FrontDefault     string `json:"front_default"`
	FrontFemale      string `json:"front_female"`
	FrontShiny       string `json:"front_shiny"`
	FrontShinyFemale string `json:"front_shiny_female"`
}

type OfficialArtwork struct {
	FrontDefault string `json:"front_default"`
	FrontShiny   string `json:"front_shiny"`
}

type Showdown struct {
	BackDefault      string `json:"back_default"`
	BackFemale       string `json:"back_female"`
	BackShiny        string `json:"back_shiny"`
	BackShinyFemale  string `json:"back_shiny_female"`
	FrontDefault     string `json:"front_default"`
	FrontFemale      string `json:"front_female"`
	FrontShiny       string `json:"front_shiny"`
	FrontShinyFemale string `json:"front_shiny_female"`
}
type other struct {
	DreamWorld      DreamWorld      `json:"dream_world"`
	Home            Home            `json:"home"`
	OfficialArtwork OfficialArtwork `json:"official-artwork"`
	Showdown        Showdown        `json:"showdown"`
}

type Sprite struct {
	BackDefault      string `json:"back_default"`
	BackFemale       string `json:"back_female"`
	BackShiny        string `json:"back_shiny"`
	BackShinyFemale  string `json:"back_shiny_female"`
	FrontDefault     string `json:"front_default"`
	FrontFemale      string `json:"front_female"`
	FrontShiny       string `json:"front_shiny"`
	FrontShinyFemale string `json:"front_shiny_female"`
	Other            other  `json:"other"`
}

type Stat struct {
	BaseStat int     `json:"base_stat"`
	Effort   int     `json:"effort"`
	Stat     NameUrl `json:"stat"`
}

type Type struct {
	Slot int     `json:"slot"`
	Type NameUrl `json:"type"`
}

type Pokemon struct {
	Id             int         `json:"id"`
	Name           string      `json:"name"`
	BaseExperience int         `json:"base_experience"`
	Height         int         `json:"height"`
	Weight         int         `json:"weight"`
	IsDefault      bool        `json:"is_default"`
	Order          int         `json:"order"`
	Abilities      []Ability   `json:"abilities"`
	Cries          Cries       `json:"cries"`
	Forms          []NameUrl   `json:"forms"`
	GameIndices    []GameIndex `json:"game_indices"`
	Moves          []Move      `json:"moves"`
	Species        NameUrl     `json:"species"`
	Sprites        Sprite      `json:"sprites"`
	Stats          []Stat      `json:"stats"`
	Types          []Type      `json:"types"`
}

type Resources struct {
	Count    int       `json:"count"`
	Next     string    `json:"next"`
	Previous string    `json:"previous"`
	Results  []NameUrl `json:"results"`
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
		logger.Println("error reading response body")
		return fmt.Errorf("error reading response body: %w", err)
	}

	if err := json.Unmarshal(data, structData); err != nil {
		logger.Println("error in json unmarshal")
		return fmt.Errorf("error in json unmarshal: %w", err)
	}
	return nil
}

func (c *Client) GetResources(name string, limit, offset uint) (*Resources, error) {
	url := "https://pokeapi.co/api/v2/" + name + "?limit=" + strconv.FormatUint(uint64(limit), 10) + "&offset=" + strconv.FormatUint(uint64(offset), 10)

	response, err := get(url, c.timeout)
	if err != nil {
		return nil, err
	}

	result := &Resources{}
	if err := parseresponse(response, &result); err != nil {
		return nil, err
	}

	logger.Println("Resources recieved successfully")
	return result, nil
}

func (c *Client) GetPokemon(id_name string) (*Pokemon, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + id_name
	response, err := get(url, c.timeout)
	if err != nil {
		return nil, err
	}
	pokemon := &Pokemon{}
	if err := parseresponse(response, &pokemon); err != nil {
		return nil, err
	}
	logger.Printf("Pokemon %s recieved successfully", id_name)
	return pokemon, nil
}
