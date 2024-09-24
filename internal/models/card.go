package models

import (
	"errors"
	"fmt"
	"germa66/internal/utils"
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
func RowToCard(record []string, backed string) (Card, error) {
	requiredLength := 2
	if len(record) < requiredLength {
		return Card{}, ErrInsufficientFields
	}
	utils.LogInfo(record)
	c := Card{
		ID:          record[0],
		Word:        record[0],
		Description: record[1],
		Backend:     backed,
	}
	utils.LogInfo(c.String())
	return c, nil
}
