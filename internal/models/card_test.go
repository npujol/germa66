package models_test

import (
	"testing"

	"germa66/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRowToCard(t *testing.T) {
	tests := []struct {
		name        string
		record      []string
		backend     string
		expectError bool
		errorType   error
		expected    models.Card
	}{
		{
			name:        "valid record",
			record:      []string{"hello", "hola"},
			backend:     "test-dict",
			expectError: false,
			expected: models.Card{
				Word:        "hello",
				Description: "hola",
				Backend:     "test-dict",
			},
		},
		{
			name:        "record with extra fields",
			record:      []string{"hello", "hola", "extra", "fields"},
			backend:     "test-dict",
			expectError: false,
			expected: models.Card{
				Word:        "hello",
				Description: "hola",
				Backend:     "test-dict",
			},
		},
		{
			name:        "record with whitespace",
			record:      []string{"  hello  ", "  hola  "},
			backend:     "test-dict",
			expectError: false,
			expected: models.Card{
				Word:        "hello",
				Description: "hola",
				Backend:     "test-dict",
			},
		},
		{
			name:        "insufficient fields",
			record:      []string{"hello"},
			backend:     "test-dict",
			expectError: true,
			errorType:   models.ErrInsufficientFields,
		},
		{
			name:        "empty record",
			record:      []string{},
			backend:     "test-dict",
			expectError: true,
			errorType:   models.ErrInsufficientFields,
		},
		{
			name:        "empty word",
			record:      []string{"", "hola"},
			backend:     "test-dict",
			expectError: true,
			errorType:   models.ErrEmptyWord,
		},
		{
			name:        "whitespace only word",
			record:      []string{"   ", "hola"},
			backend:     "test-dict",
			expectError: true,
			errorType:   models.ErrEmptyWord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card, err := models.RowToCard(tt.record, tt.backend)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.ErrorIs(t, err, tt.errorType)
				}
				assert.Equal(t, models.Card{}, card)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.Word, card.Word)
				assert.Equal(t, tt.expected.Description, card.Description)
				assert.Equal(t, tt.expected.Backend, card.Backend)
				assert.NotEmpty(t, card.ID, "ID should be generated")

				// Test that ID is consistent for same input
				card2, err2 := models.RowToCard(tt.record, tt.backend)
				require.NoError(t, err2)
				assert.Equal(t, card.ID, card2.ID, "ID should be consistent")
			}
		})
	}
}

func TestCard_String(t *testing.T) {
	card := models.Card{
		ID:          "test-id",
		Word:        "hello",
		Description: "hola",
		Backend:     "test-dict",
	}

	result := card.String()
	assert.Equal(t, "Card: hello", result)
}

func TestCard_SearchFields(t *testing.T) {
	card := models.Card{
		ID:          "test-id",
		Word:        "hello",
		Description: "hola",
		Backend:     "test-dict",
	}

	result := card.SearchFields()
	assert.Equal(t, "hello hola", result)
}

func TestCardFields(t *testing.T) {
	fields := models.CardFields()
	expected := []string{"word", "description", "backend"}
	assert.Equal(t, expected, fields)
}

func TestCardFilterableFields(t *testing.T) {
	fields := models.CardFilterableFields()
	expected := []string{"word", "backend"}
	assert.Equal(t, expected, fields)
}