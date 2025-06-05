package poker

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

var AceOfClubs = Card{Rank: Ace, Suit: Clubs}
var AceOfDiamonds = Card{Rank: Ace, Suit: Diamonds}
var AceOfHearts = Card{Rank: Ace, Suit: Hearts}
var AceOfSpades = Card{Rank: Ace, Suit: Spades}

func TestCardOutput(t *testing.T) {
	assert.Equal(t, "Ace of Clubs", AceOfClubs.String(), "Ace of Clubs should be Ace of Clubs")
	assert.Equal(t, "A♣", AceOfClubs.Symbol(), "Ace of Clubs should be A♣")

	assert.Equal(t, "Ace of Diamonds", AceOfDiamonds.String(), "Ace of Diamonds should be Ace of Diamonds")
	assert.Equal(t, "A♦", AceOfDiamonds.Symbol(), "Ace of Diamonds should be A♦")

	assert.Equal(t, "Ace of Hearts", AceOfHearts.String(), "Ace of Hearts should be Ace of Hearts")
	assert.Equal(t, "A♥", AceOfHearts.Symbol(), "Ace of Hearts should be A♥")

	assert.Equal(t, "Ace of Spades", AceOfSpades.String(), "Ace of Spades should be Ace of Spades")
	assert.Equal(t, "A♠", AceOfSpades.Symbol(), "Ace of Spades should be A♠")
}

func TestCardSort(t *testing.T) {
	cards := []Card{
		{Rank: Ace, Suit: Clubs},
		{Rank: Ten, Suit: Hearts},
		{Rank: Queen, Suit: Hearts},
		{Rank: Two, Suit: Spades},
		{Rank: Queen, Suit: Spades},
	}

	expectedRanks := []CardRank{Two, Ten, Queen, Queen, Ace}
	sort.Sort(ByCard(cards))

	for i, card := range cards {
		assert.Equal(t, expectedRanks[i], card.Rank, "Card %d should be %s", i, expectedRanks[i])
	}
}
