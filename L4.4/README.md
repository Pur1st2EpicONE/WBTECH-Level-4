## L4.4

![profiler banner](assets/banner.png)

<h3 align="center">A Go-based GC and memory profiler that exposes runtime metrics in Prometheus format and ships with pprof endpoints for deeper inspection.</h3> <br>

<br>

## Overview

This project provides an HTTP service that exposes current Go runtime and garbage collector statistics at /metrics in Prometheus format.

In addition, the service includes standard pprof endpoints under /debug/pprof for profiling CPU, heap, goroutines, mutexes, and other runtime signals.

The main metrics are collected with runtime.ReadMemStats and debug.SetGCPercent.

<br>

## Monitoring configuration

The server reads its configuration from [**./config.yaml**](./config.yaml).

Monitoring is configured via files in the monitoring/ directory:

* [**monitoring/prometheus.yaml**](./monitoring/prometheus.yaml) — Prometheus scrape configuration
* [**monitoring/docker-compose.yaml**](./monitoring/docker-compose.yaml) — Prometheus and Grafana containers

Prometheus scrapes the application metrics endpoint and also monitors itself.

<br>

## How to run everything

### Option 1. Full local run + monitoring

#### Terminal 1 — start application

```bash
make
```

This builds and runs the profiler server. The terminal will be blocked while it is running.

#### Terminal 2 — start monitoring stack

```bash
make monitoring-up
```

<br>

* Metrics: [http://localhost:8080/metrics](http://localhost:8080/metrics)
* pprof: [http://localhost:8080/debug/pprof/](http://localhost:8080/debug/pprof/)
* Prometheus: [http://localhost:9090](http://localhost:9090)
* Grafana: [http://localhost:3000](http://localhost:3000)

<br>

## How to stop everything

```bash
make monitoring-down
```

To fully clean environment:

```bash
make reset
```

<br>

## Available HTTP endpoints

### Metrics endpoint

```bash
curl http://localhost:8080/metrics
```

### pprof endpoints

```bash
curl http://localhost:8080/debug/pprof/
```

<br>

### CPU profiling example

```bash
curl "http://localhost:8080/debug/pprof/profile?seconds=30" -o cpu.pprof
```

Analyze:

```bash
go tool pprof cpu.pprof
```

⚠️ Note: for more representative results, it is recommended to generate load on the server during profiling (e.g., concurrent requests or a load-testing tool). This helps produce more meaningful CPU usage samples and makes profiling results more indicative of real runtime behavior.