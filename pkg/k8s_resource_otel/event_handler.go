package k8s_resource_otel

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

var _ cache.ResourceEventHandler = &EventHandler{}


type EventHandler struct {
	provider *trace.TracerProvider
}

func NewEventHandler() *EventHandler {
	return &EventHandler{
		provider: GlobalJaegerProvider,
	}
}

func (e EventHandler) OnAdd(obj interface{}, isInInitialList bool) {
	if event, ok := obj.(*v1.Event); ok {
		// 获取pod引起的event，并拿到缓存中的span，并设置trace
		if event.InvolvedObject.Kind == "Pod" {
			podID := event.InvolvedObject.UID
			v, ok := PodCtxSet.Get(podID)
			if !ok {
				return
			}
			spanInfo := v.(*SpanInfo)

			tracer := e.provider.Tracer("events")
			_, evtSpan := tracer.
				Start(spanInfo.RootCtx, event.Reason)
			defer evtSpan.End()

			evtSpan.SetAttributes(
				attribute.KeyValue{
					Key:   "message",
					Value: attribute.StringValue(event.Message),
				},
			)

		}
	}
}

func (e EventHandler) OnUpdate(oldObj, newObj interface{}) {
}

func (e EventHandler) OnDelete(obj interface{}) {

}





