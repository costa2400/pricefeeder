# Pricefeeder Metrics

## Available metrics

### Basic Operational Metrics

#### `fetched_prices_total`

The total number of prices fetched from the data sources. This metric is incremented every time the price feeder fetches prices from the data sources.

**labels**:

- `source`: The data source from which the price was fetched, e.g. `Bybit`.
- `success`: The result of the fetch operation. Possible values are 'true' and 'false'.

#### `aggregate_prices_total`

The total number of times the `AggregatePriceProvider` is called to return a price. It randomly selects a source for each pair from its map of price providers. This metric is incremented every time the `AggregatePriceProvider` is called.

**labels**:

- `pair`: The pair for which the price was aggregated.
- `source`: The data source from which the price was fetched, e.g. `Bybit`.
- `success`: The result of the fetch operation. Possible values are 'true' and 'false'.

#### `prices_posted_total`

The total number of txs sent to the on-chain oracle module. This metric is incremented every time the price feeder posts a price to the on-chain oracle module.

**labels**:

- `success`: The result of the post operation. Possible values are 'true' and 'false'.

#### `price_fetch_latency_seconds`

The time it takes to fetch prices from each source in seconds. This histogram tracks latency by source and trading pair.

**labels**:

- `source`: The data source from which the price was fetched.
- `pair`: The trading pair for which the price was fetched.

#### `tx_broadcast_latency_seconds`

The time it takes to broadcast transactions in seconds. This histogram tracks latency for different transaction types.

**labels**:

- `tx_type`: The type of transaction being broadcasted.

### Data Quality Metrics

#### `price_deviation_percent`

The percentage deviation between consecutive price updates for the same pair and source. Helps identify sudden price changes that might indicate data anomalies.

**labels**:

- `pair`: The trading pair for which the price was fetched.
- `source`: The data source from which the price was fetched.

#### `data_freshness_seconds`

How old the latest price data is in seconds. Measures the time elapsed since the last successful price update.

**labels**:

- `pair`: The trading pair for which the price was fetched.
- `source`: The data source from which the price was fetched.

#### `cross_source_deviation_percent`

The percentage deviation in prices between different sources for the same pair. Helps identify inconsistencies across data sources.

**labels**:

- `pair`: The trading pair for comparison.
- `source_primary`: The primary data source for comparison.
- `source_secondary`: The secondary data source for comparison.

### System Health Metrics

#### `error_count_total`

The total number of errors by type and component. Tracks error occurrences in different parts of the system.

**labels**:

- `error_type`: The type of error that occurred (e.g., "connection", "timeout", "validation").
- `component`: The component where the error occurred (e.g., "feeder", "price_provider", "price_poster").

#### `goroutine_count`

The number of active goroutines. Useful for monitoring the overall system load and detecting potential goroutine leaks.

#### `memory_usage_bytes`

The memory usage of the application in bytes. Tracks different memory metrics for system monitoring.

**labels**:

- `type`: The type of memory metric (e.g., "alloc", "sys", "heap_alloc").

#### `connection_status`

The status of connections to external services (1 for connected, 0 for disconnected). Monitors the health of connections to external data sources and APIs.

**labels**:

- `service_type`: The type of service (e.g., "exchange", "chain", "api").
- `endpoint`: The specific endpoint being monitored.

#### `rate_limit_remaining`

The remaining rate limit for external APIs. Useful for monitoring API quota usage, especially for services with rate limits like CoinGecko.

**labels**:

- `service`: The service name (e.g., "coingecko", "binance").
- `endpoint`: The specific API endpoint.
