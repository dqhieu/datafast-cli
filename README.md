# datafast

CLI for the [DataFast](https://datafa.st) analytics API. Query your website analytics, manage goals, and track payments from the terminal.

## Install

### Homebrew

```sh
brew tap dqhieu/tap
brew install datafast
```

### curl

```sh
curl -sL https://raw.githubusercontent.com/dqhieu/datafast-cli/main/install.sh | sh
```

### Go

```sh
go install github.com/dqhieu/datafast-cli@latest
```

## Setup

1. Get your API key from **Website Settings > API** in the [DataFast dashboard](https://datafa.st/dashboard)
2. Authenticate:

```sh
datafast auth login
```

## Commands

### Authentication

```sh
datafast auth login      # Save API key and website ID
datafast auth status     # Show current credentials
datafast auth logout     # Remove stored credentials
```

### Analytics

```sh
datafast analytics overview                          # Aggregate metrics
datafast analytics timeseries --interval day --start 2024-01-01 --end 2024-12-31
datafast analytics realtime                          # Current active visitors
datafast analytics realtime --watch                  # Poll every 5s
datafast analytics realtime-map                      # Real-time map data
datafast analytics pages --limit 20
datafast analytics devices
datafast analytics browsers
datafast analytics os
datafast analytics countries
datafast analytics regions
datafast analytics cities
datafast analytics referrers
datafast analytics campaigns
datafast analytics goals
datafast analytics hostnames
datafast analytics metadata
```

### Visitors

```sh
datafast visitor get <visitor_id>
```

### Goals

```sh
datafast goals create signup --visitor-id <id>
datafast goals create purchase --visitor-id <id> --param plan=pro --param source=web
datafast goals delete --name signup --start 2024-01-01 --end 2024-01-31
datafast goals delete --visitor-id <id> --force
```

### Payments

```sh
datafast payments track --amount 29.99 --currency USD --tx-id pay_123 --visitor-id <id>
datafast payments delete --tx-id pay_123 --force
datafast payments delete --visitor-id <id> --start 2024-01-01 --end 2024-12-31 --force
```

## Flags

| Flag | Description |
|------|-------------|
| `--json` | Output raw JSON |
| `--website-id` | Override configured website ID |
| `--start` | Start date (ISO 8601) |
| `--end` | End date (ISO 8601) |
| `--timezone` | Timezone (IANA format) |
| `--limit` | Limit results |
| `--fields` | Comma-separated fields |
| `--interval` | Timeseries interval: hour, day, week, month |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `DATAFAST_API_KEY` | API key (overrides config file) |
| `DATAFAST_WEBSITE_ID` | Website ID (overrides config file) |

Config is stored at `~/.config/datafast/config.json`.

## License

MIT
