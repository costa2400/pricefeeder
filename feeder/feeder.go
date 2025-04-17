package feeder

import (
	"fmt"
	"time"

	"github.com/NibiruChain/pricefeeder/types"
	"github.com/rs/zerolog"
)

var (
	// InitTimeout defines how long to wait for initial parameters before failing
	InitTimeout = 15 * time.Second
)

// Feeder is the core component that coordinates price fetching and submission.
// It connects to the blockchain, listens for events that trigger price updates,
// fetches prices from exchanges, and submits them to the blockchain.
type Feeder struct {
	logger zerolog.Logger

	stop chan struct{} // Signal to stop the feeder
	done chan struct{} // Signal that the feeder has stopped

	params types.Params // Oracle module parameters

	eventStream   types.EventStream   // Connects to the blockchain and receives events
	pricePoster   types.PricePoster   // Submits price votes to the blockchain
	priceProvider types.PriceProvider // Fetches prices from exchanges
}

// NewFeeder creates a new price feeder instance with provided dependencies.
func NewFeeder(eventStream types.EventStream, priceProvider types.PriceProvider, pricePoster types.PricePoster, logger zerolog.Logger) *Feeder {
	f := &Feeder{
		logger:        logger,
		stop:          make(chan struct{}),
		done:          make(chan struct{}),
		params:        types.Params{},
		eventStream:   eventStream,
		pricePoster:   pricePoster,
		priceProvider: priceProvider,
	}

	return f
}

// Run starts the feeder's main loop that listens for events and processes them.
func (f *Feeder) Run() {
	f.initParamsOrDie()

	go f.loop()
}

// initParamsOrDie gets the initial params from the event stream or panics if the timeout is exceeded.
func (f *Feeder) initParamsOrDie() {
	select {
	case initParams := <-f.eventStream.ParamsUpdate():
		f.handleParamsUpdate(initParams)
	case <-time.After(InitTimeout):
		panic("init timeout deadline exceeded")
	}
}

// loop is the main event processing loop. It handles three types of events:
// 1. Stop signals to shutdown the feeder
// 2. Parameter updates from the blockchain
// 3. New voting periods that trigger price submissions
func (f *Feeder) loop() {
	defer f.close()

	for {
		select {
		case <-f.stop:
			f.logger.Debug().Msg("stop signal received")
			return
		case params := <-f.eventStream.ParamsUpdate():
			f.logger.Info().Interface("changes", params).Msg("params changed")
			f.handleParamsUpdate(params)
		case vp := <-f.eventStream.VotingPeriodStarted():
			f.logger.Info().Interface("voting-period", vp).Msg("new voting period")
			f.handleVotingPeriod(vp)
		}
	}
}

// close properly shuts down all feeder components.
func (f *Feeder) close() {
	f.eventStream.Close()
	f.pricePoster.Close()
	f.priceProvider.Close()
	close(f.done)
}

// handleParamsUpdate processes parameter updates from the blockchain.
func (f *Feeder) handleParamsUpdate(params types.Params) {
	f.params = params
}

// handleVotingPeriod is triggered when a new voting period starts.
// It fetches prices for all configured pairs and submits them to the blockchain.
func (f *Feeder) handleVotingPeriod(vp types.VotingPeriod) {
	// gather prices
	prices := make([]types.Price, len(f.params.Pairs))
	for i, p := range f.params.Pairs {
		price := f.priceProvider.GetPrice(p)
		if !price.Valid {
			f.logger.Err(fmt.Errorf("no valid price")).Str("asset", p.String()).Str("source", price.SourceName)
			price.Price = 0
		}
		prices[i] = price
	}

	// send prices
	f.pricePoster.SendPrices(vp, prices)
}

// Close gracefully shuts down the feeder.
func (f *Feeder) Close() {
	close(f.stop)
	<-f.done
}
