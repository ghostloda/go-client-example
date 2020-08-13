package main

import (
	"flag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

func main() {
	//ToCreatedDeployment()
	//yamlCreateDeployment()
	getPod()
}

func newClient() *kubernetes.Clientset {
	var kubeconfig *string
	//if home := homeDir(); home != "" {
	//	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	//} else {
	//	kubeconfig = flag.String("kubeconfig", "D:\\work\\codePro\\k8s-api\\kubeconfig", "absolute path to the kubeconfig file")
	//}
	//flag.Parse()

	kubeconfig = flag.String("kubeconfig", "D:\\work\\codePro\\k8s-api\\kubeconfig", "absolute path to the kubeconfig file")

	// use the current context in kubeconfig

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset

}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
