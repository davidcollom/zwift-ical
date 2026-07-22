# zwift-ical

> Never miss a Zwift ride with your friends again.

Zwift events aren't part of your day-to-day calendar — so it's easy to forget about a ride with friends when your calendar is full of work and family commitments. **zwift-ical** bridges that gap by pulling upcoming events from the Zwift public API and publishing them as iCal (`.ics`) feeds that any calendar application can subscribe to and auto-update.

## 🗓️ Live calendars

All calendars are available at:

**<https://zwiftcal.beta.collom.co.uk>**

Subscribe to any `.ics` URL from that page directly in Google Calendar, Apple Calendar, Outlook, Thunderbird, or any other calendar application that supports iCal subscriptions.

## URL structure

Calendars are organised by sport and event type:

| Feed | URL pattern |
|------|-------------|
| All events | `/zwift.ics` |
| Cycling – group rides | `/cycling/rides/events.ics` |
| Cycling – races | `/cycling/races/events.ics` |
| Cycling – group workouts | `/cycling/workouts/events.ics` |
| Running – runs | `/running/runs/events.ics` |
| Running – races | `/running/races/events.ics` |
| Running – group workouts | `/running/workouts/events.ics` |
| By tag | `/<sport>/tag/<tag>.ics` |

Additional sports supported by Zwift (e.g. rowing, swimming) are automatically picked up as they appear in the API.

## How it works

```
Zwift public API → Go app → iCal files → Cloudflare Pages → Your calendar
```

1. The Go application (`cmd/zwiftcal`) calls the [Zwift upcoming events API](https://us-or-rly101.zwift.com/api/public/events/upcoming).
2. Events are grouped by sport and event type (group ride, race, workout) as well as by any tags attached to the event.
3. An `.ics` file is generated for each combination and written to `public/`.
4. A static HTML index page is generated listing all available calendar feeds.
5. The `public/` directory is deployed to [Cloudflare Pages](https://pages.cloudflare.com/) via GitHub Actions.

### Worlds / maps

Each calendar entry includes the Zwift world name in the event summary and location field:

Watopia, Richmond, London, New York, Innsbruck, Bologna TT, Yorkshire, Crit City, Makuri Islands, France, Paris, Gravel Mountain, Scotland.

## Automation

GitHub Actions keeps the feeds fresh:

| Workflow | Schedule | Purpose |
|----------|----------|---------|
| **Build And Publish** | Every hour (+ on push to `main`) | Fetches events, regenerates `.ics` files, deploys to Cloudflare Pages |
| **Regular committing** | Every 5 days | Bumps `.metadata` to keep the repo active and Cloudflare cache warm |
| **Unit Tests** | On every PR / push | Runs `go test ./...` |

> **API limitation:** The Zwift public API returns a maximum of 200 events per request and offers no pagination. The hourly schedule ensures the feeds stay as up-to-date as possible within this constraint.

## How to subscribe

### Google Calendar
1. Open Google Calendar → **Other calendars** → **From URL**.
2. Paste the `.ics` URL (e.g. `https://zwiftcal.beta.collom.co.uk/cycling/rides/events.ics`).
3. Click **Add calendar**. Google will refresh it automatically.

### Outlook
1. Go to **Add calendar** → **Subscribe from web**.
2. Paste the `.ics` URL and click **Import**.

### Apple Calendar
1. **File** → **New Calendar Subscription…**
2. Paste the `.ics` URL and click **Subscribe**.

### Thunderbird / other clients
Any application that supports iCal subscriptions (webcal:// or https://) can subscribe using the same URL.

## Development

### Prerequisites
- Go 1.24+
- Docker (optional, for container builds)

### Run locally

```bash
# Fetch events and generate the public/ directory
go run cmd/zwiftcal/main.go
```

The generated files will appear in `public/`.

### Test

```bash
go test ./...
# or via Make:
make test
```

### Coverage

```bash
make coverage          # generates coverage.html
make coverage-report   # prints per-function coverage to stdout
```

### Docker

```bash
make build             # build image
make debug             # build + drop into interactive shell
make publish           # build + push to Docker Hub
```

## Project structure

```
cmd/zwiftcal/       Main entrypoint – orchestrates fetch → generate → site
internal/
  events/           Zwift API client (fetch.go) and types (types.go)
  ical/             iCal generation (generate.go) and world map lookup (maps.go)
  site/             Static site + iCal file writer and HTML template renderer
.github/workflows/  CI/CD: build, publish, scheduled refresh, unit tests
Dockerfile          Multi-stage Docker image
Makefile            Common dev tasks
```

## Contributing

Issues and pull requests are welcome. Please run `go test ./...` before submitting a PR.

## License

See [LICENSE](LICENSE) if present, or contact the author.
