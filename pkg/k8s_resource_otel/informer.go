package k8s_resource_otel

import (
	"github.com/practice/opentelemetry-practice/pkg/common"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/klog/v2"
)

func K8sResourceInformer(c *common.ServerConfig) {
	client := common.NewK8sConfig().InitClientSet()
	fact := informers.NewSharedInformerFactoryWithOptions(client, 0, informers.WithNamespace("default"))
	podInformer := fact.Core().V1().Pods().Informer()

	podInformer.AddEventHandler(NewPodHandler(c.JaegerEndpoint))
	klog.Infof("k8s resource informer trace server start...")
	fact.Start(wait.NeverStop)
	<-wait.NeverStop
}
