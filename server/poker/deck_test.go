package poker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeckGetNextCard(t *testing.T) {
	cardCount := 0
	deck := NewDeck()
	cardMap := make(map[Card]int)

	for {
		cardCount++
		if card, err := deck.GetNextCard(); err != nil {
			assert.Equal(t, DeckSize+1, cardCount, "cardCount should be equal to DeckSize+1")
		} else {
			assert.True(t, cardCount <= DeckSize, "cardCount should be less than or equal to DeckSize")
			cardMap[*card]++
		}

		if cardCount == DeckSize+1 {
			break
		}

		for _, count := range cardMap {
			assert.Equal(t, 1, count)
		}
	}
}
