package common

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"strings"
)

type K8sConfig struct {
}

func NewK8sConfig() *K8sConfig {
	return &K8sConfig{}
}

// 初始化 系统 配置
// 取默认配置文件，大家可以自己修改路径和获取方式
func (*K8sConfig) K8sRestConfig() *rest.Config {
	homeDir := strings.Replace(homedir.HomeDir(), "\\", "/", -1)
	config, err := clientcmd.BuildConfigFromFlags("", homeDir+"/.kube/config")
	if err != nil {
		log.Fatal(err)
	}
	return config
}

// 初始化client-go客户端
func (this *K8sConfig) InitClientSet() *kubernetes.Clientset {

	c, err := kubernetes.NewForConfig(this.K8sRestConfig())
	if err != nil {
		log.Fatal(err)
	}

	return c
}
