# Monitor Buddy 

An open source monitoring tool which returns an API response showing the status code response of the urls specified in the urls.yaml file bundled with the application. 

There is also an UI showing the responses using a Red or Green visual grid.

## Prep 

Generate the `urls.yaml` file: 

```bash
make prep
```

## Usage

Add the sites you want to monitor to `urls.yaml`:

```yaml
websites:
  - name: My Site
    url: https://example.com
    site_id: 1
```

**Start:**
```bash
make run
```

The server starts on `http://localhost:8080`.

**Stop:** Press `Ctrl+C`.

## Endpoints

| Endpoint | Description |
|---|---|
| `GET /` | Visual Red/Green status grid |
| `GET /api/status` | JSON status for all monitored sites |