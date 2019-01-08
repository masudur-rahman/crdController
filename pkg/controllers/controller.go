package controller

import (
	"fmt"
	crdv1beta1 "github.com/masudur-rahman/crdController/pkg/apis/controller.crd/v1beta1"
	clientset "github.com/masudur-rahman/crdController/pkg/client/clientset/versioned"
	informers "github.com/masudur-rahman/crdController/pkg/client/informers/externalversions/controller.crd/v1beta1"
	listers "github.com/masudur-rahman/crdController/pkg/client/listers/controller.crd/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"log"
	"time"
)

const controllerAgentName = "crdController"

type Controller struct {
	kubeclientset kubernetes.Interface
	appsclientset clientset.Interface

	deploymentsLister	appslisters.DeploymentLister
	deploymentsSynced	cache.InformerSynced
	customLister 		listers.CustomDeploymentLister
	customSynced		cache.InformerSynced

	workqueue	workqueue.RateLimitingInterface
}


func NewController(kubeclientset kubernetes.Interface, appsclientset clientset.Interface, deploymentInformer appsinformers.DeploymentInformer, customInformer informers.CustomDeploymentInformer) *Controller {
	log.Println("Creating NewController")
	controller := &Controller{
		kubeclientset: 		kubeclientset,
		appsclientset: 		appsclientset,
		deploymentsLister: 	deploymentInformer.Lister(),
		deploymentsSynced: 	deploymentInformer.Informer().HasSynced,
		customLister: 		customInformer.Lister(),
		customSynced: 		customInformer.Informer().HasSynced,

		workqueue: 			workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "CustomDeployments"),
	}

	customInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueCustomDeploy,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueCustomDeploy(new)
		},
	})
	deploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newDepl := new.(*appsv1.Deployment)
			oldDepl := old.(*appsv1.Deployment)
			//log.Println("Deployment Informer")
			if newDepl.ResourceVersion == oldDepl.ResourceVersion {
				return
			}

			controller.handleObject(new)
		},
	})
	log.Println("Controller created")
	return controller
}

func (c *Controller) Run(threadiness int, stopCh <- chan struct{}) error {

	log.Println("Controller running")


	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	log.Println("Starting CustomDeploy Controller")

	log.Println("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.deploymentsSynced, c.customSynced); !ok {
		//log.Println("something")
		return fmt.Errorf("Failed to wait for caches to sync")
	}

	log.Println("Starting workers")

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	log.Println("Started workers")
	<- stopCh
	log.Println("Shutting down workers")
	return nil
}

func (c *Controller) runWorker() {
	log.Println("Funciton runWorker()")
	var x = 1
	for c.processNextWorkItem() {
		log.Println("processnextWorkItem", x)
		x += 1
	}
}

func (c *Controller) processNextWorkItem() bool {

	log.Println("Processing next work Item")

	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)
		var key string
		var ok bool

		if key, ok = obj.(string); !ok {
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("Expected string in workqu3e but got %v\n", obj))
			return nil
		}

		if err := c.syncHandler(key); err != nil {
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("Error syncing '%s' : %s, requesting", key, err.Error())
		}

		c.workqueue.Forget(obj)
		log.Println("Successfully synced")
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

func (c *Controller) syncHandler(key string) error {

	log.Println("SyncHandler")

	namespace, name, err := cache.SplitMetaNamespaceKey(key)

	if err != nil {
		utilruntime.HandleError(fmt.Errorf("Invalid reource key: %s", key))
	}

	customdeploy, err := c.customLister.CustomDeployments(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("customdeploy '%s' in workqueue no longler exists", key))
			return nil
		}
		return err
	}

	deploymentName := customdeploy.Spec.DeploymentName
	if deploymentName == "" {
		utilruntime.HandleError(fmt.Errorf("%s: deployment name must be specified", key))
	}
	if err != nil {
		return err
	}

	deployment, err := c.deploymentsLister.Deployments(customdeploy.Namespace).Get(deploymentName)
	if errors.IsNotFound(err) {
		deployment, err = c.kubeclientset.AppsV1().Deployments(customdeploy.Namespace).Create(newDeployment(customdeploy))
	}
	if err != nil {
		return err
	}

	if customdeploy.Spec.Replicas != nil && *customdeploy.Spec.Replicas != *deployment.Spec.Replicas {
		log.Printf("CustomDeploy %s replicas : %d, deployment replicas: %d\n", name, *customdeploy.Spec.Replicas, *deployment.Spec.Replicas)
		deployment, err = c.kubeclientset.Apps().Deployments(customdeploy.Namespace).Update(newDeployment(customdeploy))
	}
	if err != nil {
		return err
	}

	err = c.updateCustomDeployStatus(customdeploy, deployment)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) updateCustomDeployStatus(customdeploy *crdv1beta1.CustomDeployment, deployment *appsv1.Deployment) error {
	log.Println("Updating customdeployStatus")
	customCopy := customdeploy.DeepCopy()
	customCopy.Status.AvailableReplicas = deployment.Status.AvailableReplicas

	_, err := c.appsclientset.ControllerV1beta1().CustomDeployments(customdeploy.Namespace).Update(customCopy)
	return err
}

func (c *Controller) enqueueCustomDeploy(obj interface{}){
	var key string
	var err error 

	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)

	log.Println("CustomDeployment enqueued")
}

func (c *Controller) handleObject(obj interface{}) {
	log.Println("handleObject")

	
	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("Error decoding Object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("Error decoding tombstone, invalid type"))
		}
		log.Println("handleObject -1")
		log.Println("Processing object: %s", object.GetName())
	}
	log.Println("handleObject -2")
	log.Printf("Processing object: %s\n", object.GetName())

	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		if ownerRef.Kind != "CustomDeploy" {
			return
		}
		customdeploy, err := c.customLister.CustomDeployments(object.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			log.Printf("ignoring orphaned object '%s' of customdeploy '%s'\n", object.GetSelfLink(), ownerRef.Name)
		}

		log.Println("Re-enqueuing CustomDeployment")

		c.enqueueCustomDeploy(customdeploy)
		return
	}

}

func newDeployment(customdeploy *crdv1beta1.CustomDeployment) *appsv1.Deployment {
	log.Println("New Deployment")
	labels := map[string]string{
		"app"		: "appscode",
		"controller": customdeploy.Name,
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: 		customdeploy.Spec.DeploymentName,
			Namespace:	customdeploy.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(customdeploy, schema.GroupVersionKind{
					Group: 		crdv1beta1.SchemeGroupVersion.Group,
					Version:	crdv1beta1.SchemeGroupVersion.Version,
					Kind: 		"CustomDeployment",
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: customdeploy.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "appscode",
							Image: "masudjuly02/appscodeserver",
						},
					},
				},
			},
		},

	}
}
