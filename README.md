# opentelemetry-trace-practice

```bash
docker run -d -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p 5775:5775/udp -p 6831:6831/udp -p 6832:6832/udp -p 5778:5778  -p 16686:16686 -p 14268:14268  -p 14269:14269   -p 9411:9411 jaegertracing/all-in-one:latest
```

```bash
docker run -itd --name=prometheus --restart=always -v /Users/zhenyu.jiang/go/src/golanglearning/new_project/opentelemetry-practice/prometheus.yml:/etc/prometheus/prometheus.yml -p 9090:9090 prom/prometheus
```