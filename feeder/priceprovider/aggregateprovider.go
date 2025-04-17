package priceprovider

import (
	"encoding/json"

	"github.com/NibiruChain/nibiru/x/common/asset"
	"github.com/NibiruChain/pricefeeder/metrics"
	"github.com/NibiruChain/pricefeeder/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog"
)

var _ types.PriceProvider = (*AggregatePriceProvider)(nil)

// AggregatePriceProvider combines multiple price providers into one.
// It gets prices from multiple exchanges and returns the first valid price
// it finds for each trading pair.
type AggregatePriceProvider struct {
	logger    zerolog.Logger
	providers map[int]types.PriceProvider // we use a map here to provide random ranging (since golang's map range is unordered)
}

// NewAggregatePriceProvider creates an AggregatePriceProvider that manages
// multiple exchange-specific price providers. It configures connections to
// each exchange based on the provided configuration.
func NewAggregatePriceProvider(
	sourcesToPairSymbolMap map[string]map[asset.Pair]types.Symbol,
	sourceConfigMap map[string]json.RawMessage,
	logger zerolog.Logger,
) types.PriceProvider {
	providers := make(map[int]types.PriceProvider, len(sourcesToPairSymbolMap))
	i := 0
	for sourceName, pairToSymbolMap := range sourcesToPairSymbolMap {
		providers[i] = NewPriceProvider(sourceName, pairToSymbolMap, sourceConfigMap[sourceName], logger)
		i++
	}

	return AggregatePriceProvider{
		logger:    logger.With().Str("component", "aggregate-price-provider").Logger(),
		providers: providers,
	}
}

// Prometheus metric for tracking price aggregation success/failure
var aggregatePriceProvider = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: metrics.PrometheusNamespace,
	Name:      "aggregate_prices_total",
	Help:      "The total number prices provided by the aggregate price provider, by pair, source, and success status",
}, []string{"pair", "source", "success"})

// GetPrice fetches the first available and correct price from the wrapped PriceProviders.
// It iterates through the available providers in a randomized order until it finds
// a valid price. If no valid price is found, it returns an invalid price.
func (a AggregatePriceProvider) GetPrice(pair asset.Pair) types.Price {
	// iterate randomly, if we find a valid price, we return it
	// otherwise we go onto the next PriceProvider to ask for prices.
	for _, p := range a.providers {
		price := p.GetPrice(pair)
		if price.Valid {
			aggregatePriceProvider.WithLabelValues(pair.String(), price.SourceName, "true").Inc()
			return price
		}
	}

	// if we reach here no valid symbols were found
	a.logger.Warn().Str("pair", pair.String()).Msg("no valid price found")
	aggregatePriceProvider.WithLabelValues(pair.String(), "missing", "false").Inc()
	return types.Price{
		SourceName: "missing",
		Pair:       pair,
		Price:      0,
		Valid:      false,
	}
}

// Close properly shuts down all underlying price providers.
func (a AggregatePriceProvider) Close() {
	for _, p := range a.providers {
		p.Close()
	}
}
