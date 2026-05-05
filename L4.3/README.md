## L4.3

![calendar banner](assets/banner.png)

<h3 align="center">A simple event calendar service written in Go, featuring a background worker for reminders via channels, asynchronous logging, and periodic cleanup of old events.</h3> <br>

<br>

## Installation
⚠️ Note: This project requires Docker Compose, regardless of how you choose to run it.  

#### 1. Run everything in containers
```bash
make
```

This will start the entire project fully containerized using Docker Compose.

#### 2. Run Calendar locally
```bash
make local
```
In this mode, only PostgreSQL is started in container via Docker Compose, while the application itself runs locally.

⚠️ Note: Local mode requires Go 1.25.1 installed on your machine.

<br>

## Configuration

### Runtime configuration

Calendar uses two configuration files, depending on the selected run mode:

[config.full.yaml](./configs/config.full.yaml) — used for the fully containerized setup

[config.dev.yaml](./configs/config.dev.yaml) — used for local development

You may optionally review and adjust the corresponding configuration file to match your preferences. The default values are suitable for most use cases.

### Environment variables and notification credentials

Calendar uses a .env file for runtime configuration. You may create your own .env file manually before running the service, or edit [.env.example](.env.example) and let it be copied automatically on startup.
If environment file does not exist, .env.example is copied to create it. If environment file already exists, it is used as-is and will not be overwritten.

⚠️ Note: Keep .env.example for local runs. Some Makefile commands rely on it and may break if it's missing.

<br>

## Shutting down

Stopping Calendar depends on how it was started:

- Local setup — press Ctrl+C to send SIGINT to the application. The service will gracefully close connections and finish any in-progress operations.  
- Full Docker setup — containers run by Docker Compose will be stopped automatically.

In both cases, to stop all services and clean up containers, run:

```bash
make down
```

⚠️ Note: In the full Docker setup, the log folder is created by the container as root and will not be removed automatically. To delete it manually, run:
```bash
sudo rm -rf <log-folder>
```

⚠️ Note: Docker Compose also creates a persistent volume for PostgreSQL data (l43_postgres_data). This volume is not removed automatically when containers are stopped. To remove it and fully reset the environment, run:
```bash
make reset
```

<br>

## Testing & Linting

Run tests and ensure code quality:

```bash
make test        # Unit tests
make lint        # Linting checks
```

<br>

## Request examples

⚠️ Note: When the service is running, a web-based UI is available at http://localhost:8080. The examples below demonstrate how to interact with the API directly using curl.

All responses are wrapped in a result field on success:

``` json
{
  "result": ...
}
```

Errors are returned in the form:

``` json
{
  "error": "..."
}
```

### Create an event

``` bash
curl -X POST http://localhost:8080/api/v1/create_event \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "date": "2028-12-04",
    "text": "Touch grass",
    "reminder": 30
  }'
```

### Response

``` json
{
  "result": {
    "event_id": "3383503d-fb71-4b8c-85bd-a914c84252a9"
  }
}
```

### Update an event

``` bash
curl -X POST http://localhost:8080/api/v1/update_event \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "event_id": "3383503d-fb71-4b8c-85bd-a914c84252a9",
    "text": "Grind leetcode",
    "new_date": "2028-12-05"
  }'
```

### Response

``` json
{
  "result": {
    "event_updated": true
  }
}
```

### Delete an event

``` bash
curl -X POST http://localhost:8080/api/v1/delete_event \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "event_id": "3383503d-fb71-4b8c-85bd-a914c84252a9"
  }'
```

### Response

``` json
{
  "result": {
    "event_deleted": true
  }
}
```

### Get events for a day

``` bash
curl "http://localhost:8080/api/v1/events_for_day?user_id=1&date=2028-12-04"
```

### Response

``` json
{
  "result": {
    "events": [
      {
        "text": "Touch grass",
        "date": "2028-12-04",
        "event_id": "3383503d-fb71-4b8c-85bd-a914c84252a9"
      }
    ]
  }
}
```

### Get events for a week

``` bash
curl "http://localhost:8080/api/v1/events_for_week?user_id=1&date=2028-12-04"
```

### Response

``` json
{
  "result": {
    "events": [
      {
        "text": "Touch grass",
        "date": "2028-12-04",
        "event_id": "3383503d-fb71-4b8c-85bd-a914c84252a9"
      }
    ]
  }
}
```

### Get events for a month

``` bash
curl "http://localhost:8080/api/v1/events_for_month?user_id=1&date=2028-12-04"
```

### Response

``` json
{
  "result": {
    "events": [
      {
        "text": "Touch grass",
        "date": "2028-12-04",
        "event_id": "3383503d-fb71-4b8c-85bd-a914c84252a9"
      }
    ]
  }
}
```

### Example of a bad request

``` bash
curl "http://localhost:8080/api/v1/events_for_day?user_id=1&date=invalid-date"
```

### Response

``` json
{
  "error": "invalid date format, expected YYYY-MM-DD"
}
```
