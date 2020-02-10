package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	api "k8s.io/kubernetes/pkg/apis/core"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	rookv1alpha2 "github.com/rook/rook/pkg/client/clientset/versioned/typed/rook.io/v1alpha2"
)

var (
	kubeconfig           = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	storageclasses       = flag.String("storageclasses", "rook-ceph-block,rook-ceph-block-retain", "storageclass names to reconcile")
	rookns               = flag.String("rookcephns", "rook-ceph", "rook namespace")
	deletepod            = flag.Bool("deletepod", true, "Delete pod")
	storageclassesParsed = []string{""}
	config               *rest.Config
)

func GetClientset() (cs *kubernetes.Clientset, retErr error) {

	var err error

	if *kubeconfig == "" {
		config, err = rest.InClusterConfig()
		if err != nil {
			retErr = errors.New(err.Error())
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			retErr = errors.New(fmt.Sprintf("Error with config file %s: %s", *kubeconfig, err))
		}
	}

	cs, err = kubernetes.NewForConfig(config)
	if err != nil {
		retErr = errors.New(fmt.Sprintf("Bad config file: %s", err))
	}

	return cs, retErr
}

func main() {
	flag.Parse()
	clientset, err := GetClientset()

	if err != nil {
		log.Fatalf("%s", err)
	}

	storageclassesParsed = strings.Split(*storageclasses, ",")

	fieldSelector := fields.Set{api.PodStatusField: "Pending"}.AsSelector()

	watchlist := cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(),
		"pods",
		v1.NamespaceAll,
		fieldSelector,
	)

	_, controller := cache.NewInformer(
		watchlist,
		&v1.Pod{},
		0, //Duration is int64
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				CheckPodPvcs(obj.(*v1.Pod))
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				CheckPodPvcs(newObj.(*v1.Pod))
			},
		},
	)
	stop := make(chan struct{})
	defer close(stop)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go controller.Run(stop)
	<-sigs
}

func CheckPodPvcs(pod *v1.Pod) {
	var pvcvs *v1.PersistentVolumeClaimVolumeSource
	cs, _ := GetClientset()
	cv1 := cs.CoreV1()
	rookCs, err := rookv1alpha2.NewForConfig(config)
	rookNamespace := *rookns

	if err != nil {
		log.Fatalf("Error creating clientset for rook: %s", err)
	}

	for _, vol := range pod.Spec.Volumes {
		if vol.VolumeSource.PersistentVolumeClaim != nil {
			pvcvs = vol.VolumeSource.PersistentVolumeClaim

			pvcs, err := cv1.PersistentVolumeClaims(pod.ObjectMeta.Namespace).Get(pvcvs.ClaimName, metav1.GetOptions{})

			if err == nil {
				for _, sc := range storageclassesParsed {
					pvcsSc := *pvcs.Spec.StorageClassName
					if sc == pvcsSc {
						vol, err := rookCs.Volumes(rookNamespace).Get(pvcs.Spec.VolumeName, metav1.GetOptions{})
						if err == nil {
							if pod.ObjectMeta.Name == vol.Attachments[0].PodName {

								if vol.Attachments[0].Node != pod.Spec.NodeName {
									err := rookCs.Volumes(rookNamespace).Delete(pvcs.Spec.VolumeName, &metav1.DeleteOptions{})
									if err == nil {
										fmt.Printf("Volumes.rook.io %s deleted from %s ns!\n", pvcs.Spec.VolumeName, rookNamespace)
										if *deletepod {
											cv1.Pods(pod.ObjectMeta.Namespace).Delete(pod.ObjectMeta.Name, &metav1.DeleteOptions{})
										}
									}
								}
							}
						}

						break
					}
				}
			}

		}
	}

}
