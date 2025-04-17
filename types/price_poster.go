package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PricePoster defines the interface for components that submit prices to the blockchain.
// It handles the creation and submission of oracle vote transactions.
//
//go:generate mockgen --destination mocks/price_poster.go . PricePoster
type PricePoster interface {
	// SendPrices broadcasts price data to the blockchain for the current voting period.
	// It creates the necessary transactions to participate in the oracle voting protocol.
	SendPrices(vp VotingPeriod, prices []Price)

	// Whoami returns the validator address that this poster is submitting prices for.
	Whoami() sdk.ValAddress

	// Close shuts down the price poster and cleans up any resources.
	Close()
}
