package models

import (
	"crypto/md5"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInsufficientFields = errors.New("record has insufficient fields")
	ErrEmptyWord          = errors.New("word cannot be empty")
)

const (
	RequiredFieldCount = 2
)

// CardFields returns the list of searchable fields for a Card
func CardFields() []string {
	return []string{
		"word",
		"description",
		"backend",
	}
}

// CardFilterableFields returns the list of filterable fields for a Card
func CardFilterableFields() []string {
	return []string{
		"word",
		"backend",
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
// It validates the input and generates a unique ID based on the word and backend.
func RowToCard(record []string, backend string) (Card, error) {
	if len(record) < RequiredFieldCount {
		return Card{}, fmt.Errorf("%w: expected at least %d fields, got %d",
			ErrInsufficientFields, RequiredFieldCount, len(record))
	}

	word := strings.TrimSpace(record[0])
	if word == "" {
		return Card{}, ErrEmptyWord
	}

	description := ""
	if len(record) > 1 {
		description = strings.TrimSpace(record[1])
	}

	// Generate unique ID using word and backend
	id := generateCardID(word, backend)

	card := Card{
		ID:          id,
		Word:        word,
		Description: description,
		Backend:     backend,
	}

	return card, nil
}

// generateCardID creates a unique identifier for a card based on word and backend
func generateCardID(word, backend string) string {
	data := fmt.Sprintf("%s:%s", strings.ToLower(word), backend)
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}
