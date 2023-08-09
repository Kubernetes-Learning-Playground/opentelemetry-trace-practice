package pod_otel

import (
	"fmt"
	"github.com/practice/opentelemetry-practice/pkg/common"
	"k8s.io/client-go/informers"
	//configs "k8sotel/pkg/config"
)

func WatchPod() {
	client := common.NewK8sConfig().InitClientSet()
	fact := informers.NewSharedInformerFactoryWithOptions(client, 0, informers.WithNamespace("default"))
	podInformer := fact.Core().V1().Pods().Informer()

	podInformer.AddEventHandler(NewPodHandler())
	ch := make(chan struct{})
	fmt.Println("k8s可观测开始启动...")
	podInformer.Run(ch)
}
