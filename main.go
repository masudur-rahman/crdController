package main

import (
	clientset "github.com/masudur-rahman/crdController/pkg/client/clientset/versioned"
	appsinformers "github.com/masudur-rahman/crdController/pkg/client/informers/externalversions"
	controllers "github.com/masudur-rahman/crdController/pkg/controllers"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
	"time"
)


func main() {
	kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)

	if err !=nil {
		log.Fatal("Error building config file")
	}
	kubeclient := kubernetes.NewForConfigOrDie(config)
	appsclient := clientset.NewForConfigOrDie(config)

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeclient, time.Second*30)
	appsInformerFactory := appsinformers.NewSharedInformerFactory(appsclient, time.Second*30)

	stopCh := make(chan struct{})

	controller := controllers.NewController(kubeclient, appsclient,
		kubeInformerFactory.Apps().V1().Deployments(),
		appsInformerFactory.Controller().V1beta1().CustomDeployments())

	kubeInformerFactory.Start(stopCh)
	appsInformerFactory.Start(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		log.Println("Error running controller: %s", err.Error())
	}

}

