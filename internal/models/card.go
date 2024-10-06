package models

import (
	"errors"
	"fmt"
	"germa66/internal/utils"
	"strconv"
)

var ErrInsufficientFields = errors.New("record has insufficient fields")

func CardFields() []string {
	return []string{
		"word",
		"description",
		"backend",
	}
}

func CardFilterableFields() []string {
	return []string{
		"word",
	}
}

type Card struct {
	ID          string `json:"id"`
	Word        string `json:"word"`
	Description string `json:"description"`
	Backend     string `json:"backend"`
}

// String returns the string representation of the Card
func (c *Card) String() string {
	return fmt.Sprintf("Card: %s", c.Word)
}

// SearchFields returns the search fields of the Card
// Returns a string with the search fields values.
func (c *Card) SearchFields() string {
	return fmt.Sprintf("%s %s", c.Word, c.Description)
}

// RowToCard converts a row of string data to a Card struct.
func RowToCard(record []string, backed string, key int) (Card, error) {
	utils.LogDebug(
		fmt.Sprintf("Record: %v with key: %d and %v fields", record, key, len(record)),
	)

	requiredLength := 2

	if len(record) < requiredLength {
		return Card{}, ErrInsufficientFields
	}

	c := Card{
		ID:          strconv.Itoa(key),
		Word:        record[0],
		Description: record[1],
		Backend:     backed,
	}

	return c, nil
}
