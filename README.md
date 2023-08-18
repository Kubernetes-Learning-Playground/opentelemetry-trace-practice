## opentelemetry-trace-practice

### 项目思路与设计
设计背景：
主要用于学习opentelemetry链路追踪。
1. 自建httpServer，内部嵌入trace，实现http server接口的链路追踪。
![](https://github.com/Kubernetes-Learning-Playground/opentelemetry-trace-practice/blob/main/image/WechatIMG1298.png?raw=true)
2. informer内部嵌入trace，实现pod的生命周期连路追踪。
![](https://github.com/Kubernetes-Learning-Playground/opentelemetry-trace-practice/blob/main/image/img.png?raw=true)


接入方式主要有两种：(黄色部分属于嵌入服务中的部份，绿色部分是后端部分)
- 一种是使用各个观测后端的sdk 如下图：
![](https://github.com/Kubernetes-Learning-Playground/opentelemetry-trace-practice/blob/main/image/jaeger-prometheus.png?raw=true)
- 一种是使用opentelemetry-collector实现的sdk，使用此方法，调用者只需关注collector的sdk，然后维护对应的配置文件即可。
![](https://github.com/Kubernetes-Learning-Playground/opentelemetry-trace-practice/blob/main/image/otel.png?raw=true)

- docker部署jaeger
```bash
docker run -d -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p 5775:5775/udp -p 6831:6831/udp -p 6832:6832/udp -p 5778:5778  -p 16686:16686 -p 14268:14268  -p 14269:14269   -p 9411:9411 jaegertracing/all-in-one:latest
```

- docker部署prometheus
```bash
docker run -itd --name=prometheus --restart=always -v /Users/zhenyu.jiang/go/src/golanglearning/new_project/opentelemetry-practice/prometheus.yml:/etc/prometheus/prometheus.yml -p 9090:9090 prom/prometheus
```
