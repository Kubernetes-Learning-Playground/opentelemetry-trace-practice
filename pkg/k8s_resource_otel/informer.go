package k8s_resource_otel

import (
	"github.com/practice/opentelemetry-practice/pkg/common"
	"github.com/practice/opentelemetry-practice/pkg/opentelemetry/exporter"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/klog/v2"
	"os"
	"os/signal"
)

func K8sResourceInformer(c *common.ServerConfig) {
	client := common.NewK8sConfig().InitClientSet()
	fact := informers.NewSharedInformerFactoryWithOptions(client, 0, informers.WithNamespace("default"))
	GlobalJaegerProvider = exporter.NewJaegerProvider(c.JaegerEndpoint, exporter.ServiceInformer)
	podInformer := fact.Core().V1().Pods().Informer()
	podInformer.AddEventHandler(NewPodHandler())

	eventInformer := fact.Core().V1().Events().Informer()
	eventInformer.AddEventHandler(NewEventHandler())

	klog.Infof("k8s resource informer trace server start...")

	// 启动shareInformer
	fact.Start(wait.NeverStop)

	notifyCh := make(chan os.Signal, 1)
	signal.Notify(notifyCh, os.Interrupt, os.Kill)
	<-notifyCh
}
