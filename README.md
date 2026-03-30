# Monitor Buddy

An open source monitoring tool which returns an API response showing the status code response and response time of the URLs specified in the `urls.yaml` file bundled with the application.

There is also a UI showing the responses using a Red or Green visual grid, with region filtering via a dropdown.

## Prep

Generate the `urls.yaml` file:

```bash
make prep
```

## Usage

Add your regions and sites to `urls.yaml`:

```yaml
regions:
  - name: UK
    region_id: 1
  - name: Australia
    region_id: 2

websites:
  - name: My Site
    url: https://example.com
    site_id: 1
    region_id: 1
```

Each site must have a `region_id` matching one of the defined regions.

**Start:**
```bash
make run
```

The server starts on `http://localhost:8080`.

**Stop:** Press `Ctrl+C`.

## Endpoints

| Endpoint | Description |
|---|---|
| `GET /` | Visual status grid (4 per row) with region dropdown filter |
| `GET /?region_id=1` | Visual grid filtered to a specific region |
| `GET /api/status` | JSON status for all monitored sites |
| `GET /api/status?region_id=1` | JSON status filtered to a specific region |

Each site card displays the HTTP status code and response time in milliseconds.
