package types

// EventStream defines the interface for components that connect to the blockchain
// and listen for relevant events to trigger the price feeder's operations.
// It monitors chain activity and notifies other components when action is needed.
//
//go:generate mockgen --destination mocks/event_stream.go . EventStream
type EventStream interface {
	// VotingPeriodStarted returns a channel that signals when a new voting period begins.
	// This is when the feeder should submit new price votes to the blockchain.
	VotingPeriodStarted() <-chan VotingPeriod

	// ParamsUpdate returns a channel that provides updates to the oracle module parameters.
	// This allows the feeder to adjust to changing parameters like vote periods or accepted assets.
	ParamsUpdate() <-chan Params

	// Close shuts down the event stream and cleans up any connections to the blockchain.
	Close()
}
