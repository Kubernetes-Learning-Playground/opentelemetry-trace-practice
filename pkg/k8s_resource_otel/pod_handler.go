package k8s_resource_otel

import (
	"context"
	"fmt"
	"github.com/practice/opentelemetry-practice/pkg/k8s_resource_otel/helpers/k8shelper"
	"github.com/practice/opentelemetry-practice/pkg/k8s_resource_otel/helpers/lru"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"log"
)

// PodCtxSet 使用lru缓存存入pod对象，
// 当update或delete时，先从缓存内获取，延续trace
var PodCtxSet *lru.Cache

// SpanInfo 贯穿整个链路的Span
type SpanInfo struct {
	// RootCtx 根context，理解为最上层trace需要传递的context
	RootCtx context.Context
	// Ctx 子context，第二层级context
	Ctx     context.Context
	// Carrier 载体
	Carrier propagation.TextMapCarrier
}

func init() {
	cacheConfig := lru.NewCacheConfig(0, 12800, lru.ChangeCallbackFunc{})
	PodCtxSet = lru.NewCache(cacheConfig.LRUCacheMode(), cacheConfig)
}

type PodHandler struct {
	provider *trace.TracerProvider
}

func NewPodHandler(jaegerEndpoint string) *PodHandler {
	return &PodHandler{
		//provider: exporter.NewJaegerProvider(jaegerEndpoint, exporter.ServiceInformer),
	}
}

func (p *PodHandler) OnAdd(obj interface{}, isInInitialList bool) {
	if pod, ok := obj.(*v1.Pod); ok {
		tracer := p.provider.Tracer("pods")
		// 初始化 rootCtx podLifeCtx
		rootCtx, rootSpan := tracer.Start(context.Background(), fmt.Sprintf("pod-%s/%s", pod.Name, pod.Namespace))
		podLifeCtx, _ := tracer.Start(rootCtx, "pod-lifecycle")

		carrier := propagation.MapCarrier{}
		otel.GetTextMapPropagator().Inject(podLifeCtx, carrier) // 注入
		defer func() {
			// 保存信息
			PodCtxSet.Add(pod.UID, &SpanInfo{
				RootCtx:     rootCtx,
				Ctx: podLifeCtx,
				Carrier: carrier,
			})
		}()

		// 最外层trace需要记录的信息字段
		rootSpan.SetName(fmt.Sprintf("%s - %s", pod.Spec.NodeName, pod.Name))
		rootSpan.SetAttributes(
			attribute.KeyValue{
				Key:   "node",
				Value: attribute.StringValue(pod.Spec.NodeName),
			},
			attribute.KeyValue{
				Key:   "creationTimestamp",
				Value: attribute.StringValue(pod.CreationTimestamp.String()),
			},
		)
		if len(pod.OwnerReferences) != 0 {
			rootSpan.SetAttributes(attribute.KeyValue{
				Key: "ownerReference",
				Value: attribute.StringValue(fmt.Sprintf("name: %s kind: %s", pod.OwnerReferences[0].Name, pod.OwnerReferences[0].Kind)),
			})
		}
	}
}

func (p *PodHandler) OnUpdate(oldObj, newObj interface{}) {
	if pod, ok := newObj.(*v1.Pod); ok {
		// 从缓存获取
		v, ok := PodCtxSet.Get(pod.UID)
		if !ok {
			log.Println("not found carrier:", pod.Name)
			return
		}
		spanInfo := v.(*SpanInfo)
		// 把trace载体信息（ex: http特定的头)注入到新ctx
		newCtx := otel.GetTextMapPropagator().Extract(context.Background(), spanInfo.Carrier)
		tracer := p.provider.Tracer("pods")
		info := k8shelper.PrintPod(pod)

		// 处理完成or异常情况
		childSpan := oteltrace.SpanFromContext(spanInfo.Ctx)
		if childSpan.IsRecording() {
			if info.Reason == "Completed" || info.Reason == "Error" {
				childSpan.SetName(fmt.Sprintf("%s - %s(%s) ", pod.Spec.NodeName, pod.Name, info.ContainerReady))
				childSpan.End()
			}
		}

		// 基于Ctx链路的trace继续跟踪
		_, span := tracer.Start(newCtx, fmt.Sprintf("%s(%s) - %s", pod.Name, info.ContainerReady, info.Reason))

		defer span.End()


		if info.Reason == "Error" {
			span.RecordError(fmt.Errorf("pod error"))
		}

		// 记录需要的字段
		span.SetAttributes(
			attribute.KeyValue{
				Key:   "node",
				Value: attribute.StringValue(pod.Spec.NodeName),
			},
			attribute.KeyValue{
				Key:   "phase",
				Value: attribute.StringValue(string(pod.Status.Phase)),
			},
			attribute.KeyValue{
				Key:   "creationTimestamp",
				Value: attribute.StringValue(pod.CreationTimestamp.String()),
			},
			attribute.KeyValue{
				Key:   "eventMessage",
				Value: attribute.StringValue(pod.Status.Message),
			},
		)
	}
}

func (p *PodHandler) OnDelete(obj interface{}) {
	if pod, ok := obj.(*v1.Pod); ok {

		v, ok := PodCtxSet.Get(pod.UID)
		if !ok {
			log.Println("not found carrier:", pod.Name)
			return
		}
		spanInfo := v.(*SpanInfo)

		// 当删除操作时，需要结束trace追踪
		parentSpan := oteltrace.SpanFromContext(spanInfo.RootCtx)
		childSpan := oteltrace.SpanFromContext(spanInfo.Ctx)

		parentSpan.SetStatus(codes.Unset, "pod deleted")
		parentSpan.SetName(fmt.Sprintf("%s - %s(deleted)", pod.Spec.NodeName, pod.Name))

		defer childSpan.End()
		defer parentSpan.End()


		childSpan.SetAttributes(
			attribute.KeyValue{
				Key:   "node",
				Value: attribute.StringValue(pod.Spec.NodeName),
			},
			attribute.KeyValue{
				Key:   "creationTimestamp",
				Value: attribute.StringValue(pod.CreationTimestamp.String()),
			},
			attribute.KeyValue{
				Key:   "deletionTimestamp",
				Value: attribute.StringValue(pod.DeletionTimestamp.String()),
			},
		)
	}
}

var _ cache.ResourceEventHandler = &PodHandler{}
