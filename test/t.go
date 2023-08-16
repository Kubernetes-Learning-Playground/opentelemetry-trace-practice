package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/genproto/googleapis/api/label"
	"k8s.io/client-go/informers"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	tr "go.opentelemetry.io/otel/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
)

func main() {
	// 创建 Jaeger Exporter
	exporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 创建 OpenTelemetry TracerProvider
	tracerProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			"",
			semconv.ServiceNameKey.String("test"),
			attribute.String("environment", "environment"),
			attribute.Int64("ID", 1),
			semconv.ServiceVersionKey.String("v1.20.0"),
		)),
	)
	otel.SetTracerProvider(tracerProvider)

	// 创建 Kubernetes 客户端
	kubeconfig := "/Users/zhenyu.jiang/.kube/config" // Kubernetes 配置文件的路径
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// 创建 informer
	informerFactory := informers.NewSharedInformerFactory(clientset, time.Second*30)
	podInformer := informerFactory.Core().V1().Pods().Informer()
	replicaSetInformer := informerFactory.Apps().V1().ReplicaSets().Informer()

	// 创建工作队列
	queue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "PodQueue")

	// 添加 Pod 事件处理程序
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			ctx := context.Background()
			ctx, span := otel.Tracer("kubernetes-tracing").Start(ctx, "PodEventHandler")
			defer span.End()

			// 执行与Pod相关的操作，并记录span
			handlePodEvents(ctx, pod, queue)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			pod := newObj.(*corev1.Pod)
			ctx := context.Background()
			ctx, span := otel.Tracer("kubernetes-tracing").Start(ctx, "PodEventHandler")
			defer span.End()

			// 执行与Pod相关的操作，并记录span
			handlePodEvents(ctx, pod, queue)
		},
		DeleteFunc: func(obj interface{}) {
			pod, ok := obj.(*corev1.Pod)
			if !ok {
				// 无法获取删除的 Pod 对象时的处理逻辑
				return
			}
			ctx := context.Background()
			ctx, span := otel.Tracer("kubernetes-tracing").Start(ctx, "PodEventHandler")
			defer span.End()

			// 执行与Pod相关的操作，并记录span
			handlePodEvents(ctx, pod, queue)
		},
	})

	// 添加 ReplicaSet 事件处理程序
	replicaSetInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			replicaSet := obj.(*appsv1.ReplicaSet)
			ctx := context.Background()
			ctx, span := otel.Tracer("kubernetes-tracing").Start(ctx, "ReplicaSetEventHandler")
			defer span.End()

			// 执行与ReplicaSet相关的操作，并记录span
			handleReplicaSetEvents(ctx, replicaSet)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			replicaSet := newObj.(*appsv1.ReplicaSet)
			ctx := context.Background()
			ctx, span := otel.Tracer("kubernetes-tracing").Start(ctx, "ReplicaSetEventHandler")
			defer span.End()

			// 执行与ReplicaSet相关的操作，并记录span
			handleReplicaSetEvents(ctx, replicaSet)
		},
		DeleteFunc: func(obj interface{}) {
			replicaSet, ok := obj.(*appsv1.ReplicaSet)
			if !ok {
				// 无法获取删除的 ReplicaSet 对象时的处理逻辑
				return
			}
			ctx := context.Background()
			ctx, span := otel.Tracer("kubernetes-tracing").Start(ctx, "ReplicaSetEventHandler")
			defer span.End()

			// 执行与ReplicaSet相关的操作，并记录span
			handleReplicaSetEvents(ctx, replicaSet)
		},
	})

	// 启动 informer
	stopCh := make(chan struct{})
	defer close(stopCh)
	informerFactory.Start(stopCh)

	// 等待 informer 完成同步
	cache.WaitForCacheSync(stopCh, podInformer.HasSynced, replicaSetInformer.HasSynced)

	// 启动处理队列中的事件
	go func() {
		for !queue.ShuttingDown() {
			processNextItem(queue)
		}
	}()

	// 等待程序终止信号
	<-stopCh
}

//func handlePodEvents(ctx context.Context, pod *corev1.Pod) {
//	// 创建子span
//	ctx, span := otel.Tracer("kubernetes-tracing").Start(
//		ctx,
//		"HandlePodEvents",
//	)
//	span.SetAttributes(
//		attribute.KeyValue{
//			Key:   "pod.name",
//			Value: attribute.StringValue(string(pod.Name)),
//		},
//		attribute.KeyValue{
//			Key:   "pod.namespace",
//			Value: attribute.StringValue(string(pod.Namespace)),
//		},
//	)
//	defer span.End()
//
//	// 在子span中执行与Pod相关的操作，例如记录日志、执行其他函数等
//	fmt.Println("Handling Pod Event:", pod.Name)
//	// ...
//}

// 在 handlePodEvents 函数中
//func handlePodEvents(ctx context.Context, pod *corev1.Pod, queue workqueue.RateLimitingInterface) {
//	// 创建 span
//	ctx, span := otel.Tracer("kubernetes-tracing").Start(ctx, "handlePodEvents")
//	defer span.End()
//
//	// 追踪 Pod 的操作
//	span.SetAttributes(
//		attribute.KeyValue{
//			Key:   "pod.name",
//			Value: attribute.StringValue(string(pod.Name)),
//		},
//		attribute.KeyValue{
//			Key:   "pod.namespace",
//			Value: attribute.StringValue(string(pod.Namespace)),
//		},
//	)
//
//	// 判断是否关联到 ReplicaSet
//	if isPodRelatedToReplicaSet(pod) {
//		// 创建子 span，并与父 span 关联
//		cc, childSpan := otel.Tracer("kubernetes-tracing").Start(ctx, "handlePodEvents-ReplicaSet")
//		defer childSpan.End()
//
//		// 将 Pod 的 TraceContext 传递给子 span
//		// 将 Pod 的 TraceContext 传递给子 span
//		carrier := propagation.MapCarrier{}
//		otel.GetTextMapPropagator().Inject(cc, carrier)
//
//
//		// 在子 span 中执行与 Pod 相关的操作...
//		// 使用 childSpan 记录相关的标签和属性
//		childSpan.SetAttributes(
//			attribute.KeyValue{
//				Key:   "pod.name",
//				Value: attribute.StringValue(string(pod.Name)),
//			},
//			attribute.KeyValue{
//				Key:   "pod.namespace",
//				Value: attribute.StringValue(string(pod.Namespace)),
//			},
//		)
//
//		// 执行与 Pod 相关的操作...
//	} else {
//		// 在主 span 中执行与 Pod 相关的操作...
//	}
//
//	// 添加 Pod 操作到工作队列，以便进一步处理
//	queue.Add(pod)
//}

func handlePodEvents(ctx context.Context, pod *corev1.Pod) {
	// 创建 span
	ctx, span := otel.Tracer("kubernetes-tracing").Start(ctx, "handlePodEvents")
	defer span.End()

	// 追踪 Pod 的操作
	span.SetAttributes(
		label.String("pod.name", pod.Name),
		label.String("pod.namespace", pod.Namespace),
		label.String("pod.phase", string(pod.Status.Phase)),
	)

	// 判断是否关联到 ReplicaSet
	if isPodRelatedToReplicaSet(pod) {
		// 从 ReplicaSet 的标签中提取唯一标识符
		rsLabels := pod.OwnerReferences[0].Labels
		rsName := rsLabels["app.kubernetes.io/name"]
		rsNamespace := pod.Namespace

		// 根据 ReplicaSet 的唯一标识符找到已存在的 Span
		rsSpan := findReplicaSetSpan(ctx, rsName, rsNamespace)
		if rsSpan != nil {
			// 在当前 Pod 的 Span 中添加链接（link），关联到 ReplicaSet 的 Span
			span.AddEvent("Link to ReplicaSet", trace.WithAttributes(label.String("replicaset.name", rsName)))
			span.AddLink(trace.Link{SpanContext: rsSpan.SpanContext(), Attributes: []label.KeyValue{
				label.String("replicaset.name", rsName),
				label.String("replicaset.namespace", rsNamespace),
			}})
		}
	}

	// 执行与 Pod 相关的操作...
}

// 判断 Pod 是否关联到 ReplicaSet
func isPodRelatedToReplicaSet(pod *corev1.Pod) bool {
	for _, owner := range pod.OwnerReferences {
		if owner.Kind == "ReplicaSet" {
			return true
		}
	}
	return false
}

// 根据 ReplicaSet 的唯一标识符找到已存在的 Span
func findReplicaSetSpan(ctx context.Context, rsName, rsNamespace string) trace.Span {
	// 创建一个新的 SpanContext，用于匹配 ReplicaSet 的 Span
	rsSpanContext := tr.SpanContext{
		TraceID: jaeger.TraceID{},
		SpanID:  jaeger.SpanID{},
		// 设置其他所需的 SpanContext 属性
	}

	// 使用 OpenTelemetry 的 SpanProcessor 或 SpanExporter 等机制，根据 ReplicaSet 的唯一标识符找到已存在的 Span
	// 这里只是一个示例，需要根据实际情况进行实现
	// 您可以使用 OpenTelemetry 的 SpanProcessor 或 SpanExporter 等机制来存储和检索 Span
	// 或者使用 OpenTelemetry 的 Context API 来存储和检索 Span
	// 这里假设已找到 ReplicaSet 的 Span
	// 您可以根据实际情况进行修改和调整
	tr.w
	_, rsSpan := otel.Tracer("kubernetes-tracing").Start(ctx, "findReplicaSetSpan", tr.WithSpanContext(rsSpanContext))

	return rsSpan
}

// 判断 Pod 是否关联到 ReplicaSet
//func isPodRelatedToReplicaSet(pod *corev1.Pod) bool {
//	for _, owner := range pod.OwnerReferences {
//		if owner.Kind == "ReplicaSet" {
//			return true
//		}
//	}
//	return false
//}

func handleReplicaSetEvents(ctx context.Context, replicaSet *appsv1.ReplicaSet) {
	// 创建子span
	ctx, span := otel.Tracer("kubernetes-tracing").Start(
		ctx,
		"HandleReplicaSetEvents",
	)
	span.SetAttributes(
		attribute.KeyValue{
			Key:   "replicaset.name",
			Value: attribute.StringValue(string(replicaSet.Name)),
		},
		attribute.KeyValue{
			Key:   "replicaset.namespace",
			Value: attribute.StringValue(string(replicaSet.Namespace)),
		},
	)
	defer span.End()

	// 在子span中执行与ReplicaSet相关的操作，例如记录日志、执行其他函数等
	fmt.Println("Handling ReplicaSet Event:", replicaSet.Name)
	// ...
}

func processNextItem(queue workqueue.RateLimitingInterface) {
	// 处理队列中的下一个事件
	obj, shutdown := queue.Get()
	if shutdown {
		return
	}

	err := func(obj interface{}) error {
		defer queue.Done(obj)

		var key string
		var ok bool
		if key, ok = obj.(string); !ok {
			queue.Forget(obj)
			return fmt.Errorf("expected string but got %#v", obj)
		}

		// 处理事件
		err := func() error {
			// 处理事件的逻辑
			return nil
		}()

		if err != nil {
			queue.AddRateLimited(obj)
			return fmt.Errorf("error processing item with key %q: %w", key, err)
		}

		queue.Forget(obj)
		return nil
	}(obj)

	if err != nil {
		log.Println(err)
	}

	queue.Forget(obj)
}
