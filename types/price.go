package types

import (
	"time"

	"github.com/NibiruChain/nibiru/x/common/asset"
	"github.com/NibiruChain/nibiru/x/common/set"
	"github.com/rs/zerolog"
)

const (
	// PriceTimeout defines how long a price is considered valid after being fetched
	PriceTimeout = 15 * time.Second
)

// RawPrice represents a price fetched from an exchange with its timestamp
type RawPrice struct {
	Price      float64
	UpdateTime time.Time
}

// Price defines the processed price data that will be submitted to the blockchain.
// It includes metadata like the source and validity status.
type Price struct {
	// Pair defines the symbol we're posting prices for.
	Pair asset.Pair
	// Price defines the symbol's price.
	Price float64
	// SourceName defines the source which is providing the prices.
	SourceName string
	// Valid reports whether the price is valid or not.
	// If not valid then an abstain vote will be posted.
	// Computed from the update time.
	Valid bool
}

// FetchPricesFunc is the function type used to fetch updated prices from an exchange.
// Each price source implements this function to query their specific API.
// The symbols passed are the symbols we require prices for.
// The returned map must map symbol to its float64 price, or an error.
// If there's a failure in updating only one price then the map can be returned
// without the provided symbol.
type FetchPricesFunc func(symbols set.Set[Symbol], logger zerolog.Logger) (map[Symbol]float64, error)
