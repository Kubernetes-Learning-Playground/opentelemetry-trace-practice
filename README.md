## opentelemetry-trace-practice

### 项目思路与设计
设计背景：
主要用于学习opentelemetry链路追踪。
1. 自建httpServer，内部嵌入trace，实现http server接口的链路追踪。
![]()
2. informer内部嵌入trace，实现pod的生命周期连路追踪。
![]()



- docker部署jaeger
```bash
docker run -d -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p 5775:5775/udp -p 6831:6831/udp -p 6832:6832/udp -p 5778:5778  -p 16686:16686 -p 14268:14268  -p 14269:14269   -p 9411:9411 jaegertracing/all-in-one:latest
```

- docker部署prometheus
```bash
docker run -itd --name=prometheus --restart=always -v /Users/zhenyu.jiang/go/src/golanglearning/new_project/opentelemetry-practice/prometheus.yml:/etc/prometheus/prometheus.yml -p 9090:9090 prom/prometheus
```