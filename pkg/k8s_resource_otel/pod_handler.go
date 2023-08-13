package k8s_resource_otel

import (
	"context"
	"fmt"
	"github.com/practice/opentelemetry-practice/pkg/k8s_resource_otel/helpers/lru"
	//lru "github.com/hashicorp/golang-lru/v2"
	"github.com/practice/opentelemetry-practice/pkg/opentelemetry/exporter"
	"github.com/practice/opentelemetry-practice/pkg/k8s_resource_otel/helpers/k8shelper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"log"
)

// PodCtxSet 使用lru缓存存入pod对象，
// 当update或delete时，先从缓存内获取，延续trace
var PodCtxSet *lru.Cache

func init() {
	cacheConfig := lru.NewCacheConfig(0, 12800, lru.ChangeCallbackFunc{})
	PodCtxSet = lru.NewCache(cacheConfig.LRUCacheMode(), cacheConfig)
}

type PodHandler struct {
	provider *trace.TracerProvider
}

func NewPodHandler(jaegerEndpoint string) *PodHandler {
	return &PodHandler{
		provider: exporter.NewJaegerProvider(jaegerEndpoint, exporter.ServiceInformer),
	}
}

func (p *PodHandler) OnAdd(obj interface{}, isInInitialList bool) {
	if pod, ok := obj.(*v1.Pod); ok {
		tracer := p.provider.Tracer("pods")
		newCtx, span := tracer.Start(context.Background(), pod.Spec.NodeName+" - "+pod.Name)
		defer span.End()
		carrier := propagation.MapCarrier{}
		otel.GetTextMapPropagator().Inject(newCtx, carrier) //注入
		defer func() {
			PodCtxSet.Add(pod.UID, carrier) //保存  头载体
		}()

		span.SetAttributes(
			attribute.KeyValue{
				Key:   "phase",
				Value: attribute.StringValue(string(pod.Status.Phase)),
			},
			attribute.KeyValue{
				Key:   "node",
				Value: attribute.StringValue(pod.Spec.NodeName),
			},
			attribute.KeyValue{
				Key: "creationTimestamp",
				Value: attribute.StringValue(pod.CreationTimestamp.String()),
			},
		)
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
		carrier := v.(propagation.TextMapCarrier)
		// 把trace载体信息（ex: http特定的头)注入到新ctx
		newCtx := otel.GetTextMapPropagator().Extract(context.Background(), carrier)
		tracer := p.provider.Tracer("pods")
		info := k8shelper.PrintPod(pod)

		_, span := tracer.Start(newCtx, fmt.Sprintf("%s(%s) - %s ", pod.Name, info.ContainerReady, info.Reason))

		defer span.End()
		span.SetAttributes(
			attribute.KeyValue{
				Key:   "phase",
				Value: attribute.StringValue(string(pod.Status.Phase)),
			},
			attribute.KeyValue{
				Key:   "node",
				Value: attribute.StringValue(pod.Spec.NodeName),
			},
			attribute.KeyValue{
				Key: "creationTimestamp",
				Value: attribute.StringValue(pod.CreationTimestamp.String()),
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
		carrier := v.(propagation.TextMapCarrier)
		newCtx := otel.GetTextMapPropagator().Extract(context.Background(), carrier)
		tracer := p.provider.Tracer("pods")
		info := k8shelper.PrintPod(pod)

		_, span := tracer.Start(newCtx, fmt.Sprintf("%s(%s) - %s ", pod.Name, info.ContainerReady, info.Reason))

		defer span.End()
		span.SetAttributes(
			attribute.KeyValue{
				Key:   "phase",
				Value: attribute.StringValue(string(pod.Status.Phase)),
			},
			attribute.KeyValue{
				Key:   "node",
				Value: attribute.StringValue(pod.Spec.NodeName),
			},
			attribute.KeyValue{
				Key: "creationTimestamp",
				Value: attribute.StringValue(pod.CreationTimestamp.String()),
			},
			attribute.KeyValue{
				Key: "deletionTimestamp",
				Value: attribute.StringValue(pod.DeletionTimestamp.String()),
			},
		)
	}
}

var _ cache.ResourceEventHandler = &PodHandler{}
