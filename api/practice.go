package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Pokemon struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Forms          []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	Height    int `json:"height"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Order   int    `json:"order"`
	Species struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		FrontDefault string `json:"front_default"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

func searchPokemon(app *fiber.App) {
	app.Post("/", fetchPokemon)
}

type Request struct {
	Name string `json:"name"`
}

func fetchPokemon(hello *fiber.Ctx) error {
	var req Request
	if err := hello.BodyParser(&req); err != nil {
		return hello.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
		})
	}
	if req.Name == "" {
		return hello.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name field is required",
		})
	}
	response, err := getPokemon(req.Name)
	if err != nil {
		return err
	}

	return hello.JSON(fiber.Map{
		"response": response,
	})
}

// getPokemon fetches the details of a Pokemon by its name from the PokeAPI.
// It returns a Pokemon struct and an error if any occurs during the request or unmarshalling.
func getPokemon(name string) (Pokemon, error) {
	var pokemon Pokemon
	url := "https://pokeapi.co/api/v2/pokemon/" + name
	response, err := http.Get(url)
	if err != nil {
		return pokemon, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return pokemon, fmt.Errorf("failed to get pokemon: %s", response.Status)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		// fmt.Println(pokemon)
	}
	err = json.Unmarshal(body, &pokemon)
	fmt.Println(pokemon)
	return pokemon, err
}
