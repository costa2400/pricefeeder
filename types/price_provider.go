package types

import "github.com/NibiruChain/nibiru/x/common/asset"

// PriceProvider defines the interface for components that fetch prices from external sources.
// This is implemented by exchange-specific providers and aggregate providers that
// combine multiple price sources.
// PriceProvider must handle failures by itself.
//
//go:generate mockgen --destination mocks/price_provider.go . PriceProvider
type PriceProvider interface {
	// GetPrice fetches the price for a given asset pair.
	// Returns a Price object that includes the source, validity, and price value.
	// If a price can't be retrieved or is invalid, it should return a Price with Valid=false.
	GetPrice(asset.Pair) Price
	// Close shuts down the price provider and cleans up any resources.
	Close()
}
