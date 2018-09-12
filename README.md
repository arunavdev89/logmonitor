## Access Log Monitor

- Reads access log in W3C Log format, and emits a counter in Prometheus.
- A poller tracks/queries Prometheus for various summary stat and displays the summary.
- Also generates an alert if a certain threshold is crossed every 120s (configurable in main.go)

### How to run

- Set the file you want to monitor in `logfile` environment variable. (Have the option either changing .env or passing -e <file name> in the `docker-compose` command below,)

- Build the docker file.

```
docker build -t arunavdev89/access-log-monitor:0.2 . -f Dockerfile
````

- Run the Docker compose (which will launch a local prometheus instance)

```
docker-compose up
```

#### Test

- To see some of the alerts/summary in action run the test.go file under test/

```
go run test.go --file=<log file being monitored>
```

-- Note you may have to have docker mount/volume share permission from host to container.

### What Next?

- Add unit tests.
- Add authentication (basic) in prometheus queries.
- Add some console coloring and pretty text in reporter.go
- Make it distributed multi machine monitoring agent. Can be achieved through piping the log line to kafka (or some other distributed queue) and then reading from kafka and emitting metrics to prometheus.
- Make prometheus distributed and fault-tolerant (instead of having one instance)
- Pass prometheus metrics to Grafana for visualization/alert configuration.
- Run in k8s :)